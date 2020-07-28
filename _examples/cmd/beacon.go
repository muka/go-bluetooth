//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	beacon_example "github.com/muka/go-bluetooth/examples/beacon"
	"github.com/spf13/cobra"
)

// beaconCmd represents the beacon command
var beaconCmd = &cobra.Command{
	Use:   "beacon",
	Short: "Advertising example",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		adapterID, err := cmd.Flags().GetString("adapterID")
		if err != nil {
			fail(err)
		}

		if len(args) == 0 {
			failArg("type: ibeacon or eddystone")
		}

		var beaconType string = "URL"
		if len(args) == 2 {
			beaconType = args[1]
		}

		fail(beacon_example.Run(args[0], beaconType, adapterID))

	},
}

func init() {
	rootCmd.AddCommand(beaconCmd)
}
