package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/marlinprotocol/ctl2/version"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

type ConfigConfig struct {
	HomeDir string
}

func (cfg *ConfigConfig) Initialise() error {
	// TODO add links
	log.Info("marlinctl will ask you a few set of inputs to correctly setup your config file. For most of the following options, the given defaults are what you need. In case you need help setting up marlinctl, refer init logs at <insert link here>")

	// Get marlindir config
	marlinDir, err := cfg.getMarlinDir()
	if err != nil {
		return err
	}

	// Get upstream repo config
	upstream, err := cfg.getUpstreamRepo()
	if err != nil {
		return err
	}

	// Get platform and runtime config
	platform, runtime, err := cfg.getPlatformAndRuntime()
	if err != nil {
		return err
	}

	data := struct {
		ConfigVersion, Platform, Runtime, Upstream string
	}{version.CfgVersion, platform, runtime, upstream}

	// Write configuration
	log.Info("Writing configuration at dir: ", marlinDir)
	err = os.MkdirAll(marlinDir, os.ModePerm)
	if err != nil {
		return err
	}

	cfgTemplate, err := template.New("cfgTemplate").Parse(version.CfgTemplate)
	if err != nil {
		return err
	}

	cfgFile, err := os.Create(marlinDir + "/marlinctl_config.yaml")
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	err = cfgTemplate.Execute(cfgFile, data)
	if err != nil {
		return err
	}
	log.Info("Configuration written successfully to disk")

	return nil
}

func (cfg *ConfigConfig) RemoveLocalConfig() error {
	err := os.Remove(cfg.HomeDir + "/marlinctl_config.yaml")
	return err
}

func (cfg *ConfigConfig) getMarlinDir() (string, error) {
	validate := func(input string) error {
		directoryPathRegex := `^/|(/[\w-]+)+$`

		validDirectoryPath, err := regexp.Match(directoryPathRegex, []byte(input))
		if err != nil {
			return err
		}

		if !validDirectoryPath {
			return errors.New("Invalid directory path")
		}

		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	prompt := promptui.Prompt{
		Label:       "Marlinctl directory",
		Default:     home + "/.marlinctl",
		Validate:    validate,
		AllowEdit:   true,
		HideEntered: true,
	}

	result, err := prompt.Run()
	time.Sleep(300 * time.Millisecond) // Added because sometime terminal does not clear and makes the experience inconsistent

	if err != nil {
		return "", err
	}

	return result, nil
}

func (cfg *ConfigConfig) getUpstreamRepo() (string, error) {
	validate := func(input string) error {
		githubRemoteRepoRegex := `((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`

		validGithubRepo, err := regexp.Match(githubRemoteRepoRegex, []byte(input))
		if err != nil {
			return err
		}

		if !validGithubRepo {
			return errors.New("Invalid upstream git repository")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:       "Upstream releases repository",
		Default:     "https://github.com/marlinprotocol/releases.git",
		Validate:    validate,
		AllowEdit:   true,
		HideEntered: true,
	}

	result, err := prompt.Run()
	time.Sleep(300 * time.Millisecond) // Added because sometime terminal does not clear and makes the experience inconsistent

	if err != nil {
		return "", err
	}

	return result, nil
}

func (cfg *ConfigConfig) getPlatformAndRuntime() (string, string, error) {
	platform := runtime.GOOS + "-" + runtime.GOARCH

	supportedPlatformsRuntimeCombinations := []string{"linux-amd64.supervisor", "linux-amd64.systemd"}

	defaultOption := -1
	for i := 0; i < len(supportedPlatformsRuntimeCombinations); i++ {
		if match, err := regexp.Match(platform, []byte(supportedPlatformsRuntimeCombinations[i])); err == nil && match {
			defaultOption = i
			break
		}
	}

	if defaultOption == -1 {
		log.Warn("We could not find a relevant platform-runtime combination for your platform: " + platform +
			". There is a possibility that marlinctl does not support your platform as of yet." +
			"You are free to still choose one of the options, however it is recommended that you contact dev team regarding support for your platform.")
		defaultOption = 0
	} else {
		log.Info("Identified system runtime as: ", platform, ". Select one of the options of format: "+platform+".<runtime_you_want>")
	}

	prompt := promptui.Select{
		Label:        "Platform and runtime",
		Size:         50,
		CursorPos:    defaultOption,
		Items:        supportedPlatformsRuntimeCombinations,
		HideSelected: true,
	}

	_, result, err := prompt.Run()
	time.Sleep(300 * time.Millisecond) // Added because sometime terminal does not clear and makes the experience inconsistent

	if err != nil {
		return "", "", err
	}
	combination := strings.Split(result, ".")

	return combination[0], combination[1], nil
}
