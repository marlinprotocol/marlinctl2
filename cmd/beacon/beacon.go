/*
Copyright © 2020 MARLIN TEAM <info@marlin.pro>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package beacon

import (
	"github.com/marlinprotocol/ctl2/cmd/beacon/actions"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var BeaconCmd = &cobra.Command{
	Use:   "beacon",
	Short: "Marlin Beacon",
	Long:  `Marlin Beacon`,
}

func init() {
	BeaconCmd.AddCommand(actions.CreateCmd)
	BeaconCmd.AddCommand(actions.StatusCmd)
	BeaconCmd.AddCommand(actions.DestroyCmd)
	BeaconCmd.AddCommand(actions.VersionsCmd)
	BeaconCmd.AddCommand(actions.LogsCmd)
}
