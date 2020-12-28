package registry

import (
	"os"

	"github.com/getlantern/deepcopy"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
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
	for _, r := range GlobalRegistry {
		// Skip non enabled registries
		if !r.Enabled {
			continue
		}

		tempDir := r.Local + "_temp"
		err := util.CreateDirPathIfNotExists(tempDir)
		if err != nil {
			log.Error("Error while creating dir path: ", tempDir)
			// TODO how to escape
			os.Exit(1)
		}

		err = util.RemoveDirContents(tempDir)
		if err != nil {
			log.Error("Error while removing contents at dir path: ", tempDir)
			// TODO how to escape
			os.Exit(1)
		}

		err = util.GitPullHead(r.Link, r.Branch, tempDir)
		if err != nil {
			log.Error("Error while pulling releases information to dir path: ", tempDir, " Registry: ", r)
			// TODO how to escape
			os.Exit(1)
		}

		ok, err := c.registryPreSanity(tempDir)
		if err != nil {
			log.Error("Prerun Sanity resulted in error: ", tempDir, " Registry: ", r)
			// TODO how to escape
			os.Exit(1)
		}

		if !ok {
			log.Warning("Upstream registry did not pass registry pre sanity tests. Reverting to older registry!. Registry: ", r)
			continue
		}

		err = util.MoveDir(tempDir, r.Local)
		if err != nil {
			log.Error("Error while moving directory: ", tempDir, " to directory ", r.Local)
			// TODO how to escape
			os.Exit(1)
		}

		err = util.RemoveDirPathIfExists(tempDir)
		if err != nil {
			log.Error("Error while removing dir: ", tempDir)
			// TODO how to escape
			os.Exit(1)
		}

		log.Debug("Successfully pulled registry: ", r)
	}
	return nil
}

func (c *RegistryConfig) registryPreSanity(dirPath string) (bool, error) {
	log.Warning("Yet to implement")
	return true, nil
}
