package registry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/jedib0t/go-pretty/table"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Upstream configuration
type RegistryConfig []types.Registry

var GlobalRegistry RegistryConfig

func SetupGlobalRegistry(reg RegistryConfig) {
	GlobalRegistry = make(RegistryConfig, len(reg))
	for i := 0; i < len(reg); i++ {
		deepcopy.Copy(&GlobalRegistry[i], &reg[i])
	}
}

// Download upstream VCS to Home Clone location
func (c *RegistryConfig) Sync() error {
	type WorkerResult struct {
		Registry  types.Registry
		Completed bool
		Error     error
	}

	workerChan := make(chan WorkerResult, len(*c))

	start := time.Now()
	for _, r := range *c {
		go func(wc chan WorkerResult, r types.Registry) {
			// Skip non enabled registries
			if !r.Enabled {
				wc <- WorkerResult{Registry: r, Completed: true, Error: nil}
				return
			}

			tempDir := r.Local + "_temp"
			err := util.CreateDirPathIfNotExists(tempDir)
			if err != nil {
				log.Error("Error while creating dir path: ", tempDir, " ", err)
				wc <- WorkerResult{Registry: r, Completed: false, Error: err}
				return
			}

			err = util.RemoveDirContents(tempDir)
			if err != nil {
				log.Error("Error while removing contents at dir path: ", tempDir, " ", err)
				wc <- WorkerResult{Registry: r, Completed: false, Error: err}
				return
			}

			err = util.GitPullHead(r.Link, r.Branch, tempDir)
			if err != nil {
				log.Error("Error while pulling releases information to dir path: ", tempDir, " Registry: ", r, " ", err)
				wc <- WorkerResult{Registry: r, Completed: false, Error: err}
				return
			}

			ok, err := c.registryPreSanity(tempDir)
			if err != nil {
				log.Error("Prerun Sanity resulted in error: ", tempDir, " Registry: ", r, " ", err)
				wc <- WorkerResult{Registry: r, Completed: false, Error: err}
				return
			}

			if !ok {
				log.Warning("Upstream registry did not pass registry pre sanity tests. Reverting to older registry!. Registry: ", r)
				wc <- WorkerResult{Registry: r, Completed: true, Error: errors.New("Registry not updated, not able to fetch from upstream")}
				return
			}

			err = util.MoveDir(tempDir, r.Local)
			if err != nil {
				log.Error("Error while moving directory: ", tempDir, " to directory ", r.Local, " ", err)
				wc <- WorkerResult{Registry: r, Completed: false, Error: err}
				return
			}

			err = util.RemoveDirPathIfExists(tempDir)
			if err != nil {
				log.Error("Error while removing dir: ", tempDir, " ", err)
				wc <- WorkerResult{Registry: r, Completed: false, Error: err}
				return
			}

			wc <- WorkerResult{Registry: r, Completed: true, Error: nil}
			return
		}(workerChan, r)
	}

	for i := 0; i < len(*c); i++ {
		work := <-workerChan
		if work.Completed {
			if work.Error != nil {
				log.Warning("Registry ", work.Registry, " completed with error ", work.Error)
			}
		} else {
			log.Error("Registry ", work.Registry, " failed due to error: ", work.Error, ". Aborting application.")
			os.Exit(1)
		}
	}
	elapsed := time.Since(start)

	log.Debug("Remote registeries pulled in ", elapsed)
	return nil
}

// TODO add sanity checks here
func (c *RegistryConfig) registryPreSanity(dirPath string) (bool, error) {
	return true, nil
}

func (c *RegistryConfig) GetVersions(project string, subscriptions []string, currentVersion string, updatePolicy string, runtime string) ([]ProjectVersion, error) {
	if !util.IsValidUpdatePolicy(updatePolicy) {
		return []ProjectVersion{}, errors.New("Unknown update policy: " + updatePolicy)
	}

	var grsMap = make(map[string]types.Registry)
	for _, v := range *c {
		grsMap[v.Name] = v
	}

	for _, s := range subscriptions {
		if sub, ok := grsMap[s]; !sub.Enabled || !ok {
			return []ProjectVersion{}, errors.New("Registry does not support subscription to " + s)
		}
	}

	var projectVersions []ProjectVersion

	for _, s := range subscriptions {
		sub, _ := grsMap[s]
		releaseFile := sub.Local + "/projects/" + project + "/releases.json"
		if _, err := os.Stat(releaseFile); os.IsNotExist(err) {
			return projectVersions, errors.New("Cannot find " + releaseFile)
		}
		file, _ := ioutil.ReadFile(releaseFile)
		releasesJson := types.ReleaseJSON{}
		err := json.Unmarshal([]byte(file), &releasesJson)
		if err != nil {
			return projectVersions, err
		}
		switch releasesJson.JSONVersion {
		case 1:
			versions, err := c.decodeReleasesJsonVersion1(releasesJson.Data, s, runtime, currentVersion, updatePolicy)
			if err != nil {
				return projectVersions, err
			}
			projectVersions = append(projectVersions, versions...)
		default:
			return projectVersions, errors.New("Cannot decode releases json with JSON version: " + strconv.Itoa(releasesJson.JSONVersion))
		}
	}

	sort.Slice(projectVersions, func(i, j int) bool {
		return projectVersions[i].ReleaseTime.After(projectVersions[j].ReleaseTime)
	})

	return projectVersions, nil
}

func (c *RegistryConfig) decodeReleasesJsonVersion1(data interface{}, subscription string, runtime string, currentVersion string, updatePolicy string) ([]ProjectVersion, error) {
	var isPublic bool = (subscription == "public")
	var isFirstRun bool = (currentVersion == "0.0.0")

	currMaj, currMin, currPatch, currSub, currBuild, err := util.DecodeVersionString(currentVersion)
	if err != nil {
		return []ProjectVersion{}, err
	}

	var versions []ProjectVersion
	var missedVersions []string
	// TODO more error checking
	for MajVer, MajVerData := range data.(map[string]interface{}) {
		for MinVer, MinVerData := range MajVerData.(map[string]interface{}) {
			for PatchVer, PatchVerData := range MinVerData.(map[string]interface{}) {
				for Build, BuildData := range PatchVerData.(map[string]interface{}) {

					maj, err1 := strconv.Atoi(MajVer)
					min, err2 := strconv.Atoi(MinVer)
					patch, err3 := strconv.Atoi(PatchVer)
					build, err4 := strconv.Atoi(Build)

					if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
						continue
					}

					var fullVersion = MajVer + "." + MinVer + "." + PatchVer
					if !isPublic {
						fullVersion = fullVersion + "-" + subscription + "." + Build
					}

					if !isFirstRun && !util.CanUseVersion(maj, min, patch, subscription, build,
						currMaj, currMin, currPatch, currSub, currBuild,
						updatePolicy) {
						if isPublic && util.IsHigherVersion(maj, min, patch,
							currMaj, currMin, currPatch, currSub) {
							missedVersions = append(missedVersions, fullVersion)
						}
						continue
					}

					bundles, bundlesok := BuildData.(map[string]interface{})["bundles"]
					if !bundlesok {
						continue
					}
					runtimeData, runtimeok := bundles.(map[string]interface{})[runtime]
					if !runtimeok {
						continue
					}
					runnerId, runnerIdok := runtimeData.(map[string]interface{})["runner"]
					if !runnerIdok {
						continue
					}
					runnerData, runnerDataok := runtimeData.(map[string]interface{})["data"]
					if !runnerDataok {
						continue
					}
					var r, ok = BuildData.(map[string]interface{})["time"]
					if !ok {
						continue
					}
					var reltime, err = time.Parse(time.RFC822Z, r.(string))
					if err != nil {
						continue
					}
					var desc, ok3 = BuildData.(map[string]interface{})["description"]
					if !ok3 {
						continue
					}
					var version = ProjectVersion{
						ReleaseType: subscription,
						Version:     fullVersion,
						Description: desc.(string),
						ReleaseTime: reltime,
						RunnerId:    runnerId.(string),
						RunnerData:  runnerData,
					}
					versions = append(versions, version)
				}
			}
		}
	}
	if len(missedVersions) > 0 {
		var updatesMissedString = ""
		for _, v := range missedVersions {
			updatesMissedString = updatesMissedString + "[version \"" + v + "\"] "
		}
		log.Warning("It seems that you are missing vital updates on public releases channel." +
			" You are not able to upgrade to them due to update policy applied on the project." +
			" Here are versions you are missing: " + updatesMissedString)

	}
	return versions, nil
}

func (c *RegistryConfig) PrettyPrintProjectVersions(versions []ProjectVersion) {
	t := util.GetTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Type", "Version", "Time", "Description", "Runner"})
	for _, v := range versions {
		var releaseType = v.ReleaseType
		if releaseType == "rtw" {
			releaseType = "release"
		}
		t.AppendRow(table.Row{releaseType, v.Version, v.ReleaseTime, v.Description, v.RunnerId})
	}
	// terminalColorCapability, err := exec.Command("tput", "colors").Output()
	// if err == nil && strings.TrimSpace(string(terminalColorCapability)) == "256" && isatty.IsTerminal(os.Stdout.Fd()) {
	// 	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	// }
	t.Render()
}

func (c *RegistryConfig) GetVersionToRun(projectName string) (ProjectVersion, error) {
	var proj types.Project
	err := viper.UnmarshalKey(projectName, &proj)
	if err != nil {
		return ProjectVersion{}, err
	}
	versions, err := c.GetVersions(projectName, proj.Subscription, proj.CurrentVersion, proj.UpdatePolicy, proj.Runtime)
	if err != nil {
		return ProjectVersion{}, err
	}

	if len(versions) > 0 {
		return versions[0], nil
	}
	return ProjectVersion{}, errors.New("No version found for running the project " + projectName)
}

type ProjectVersion struct {
	ReleaseType string
	Version     string
	Description string
	ReleaseTime time.Time
	RunnerId    string
	RunnerData  interface{}
}
