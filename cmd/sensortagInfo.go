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
	sensortag_info_example "github.com/muka/go-bluetooth/examples/sensortag_info"
	"github.com/spf13/cobra"
)

// sensortagInfoCmd represents the sensortagInfo command
var sensortagInfoCmd = &cobra.Command{
	Use:   "sensortag-info",
	Short: "Retrieve TI SensorTag sensors informations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		adapterID, err := cmd.Flags().GetString("adapterID")
		if err != nil {
			fail(err)
		}

		if len(args) < 1 {
			failArgs([]string{"sensortag_address"})
		}

		fail(sensortag_info_example.Run(args[0], adapterID))
	},
}

func init() {
	rootCmd.AddCommand(sensortagInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sensortagInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sensortagInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
