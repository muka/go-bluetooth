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
	service_example "github.com/muka/go-bluetooth/examples/service"
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A service / client example",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		adapterID, err := cmd.Flags().GetString("adapterID")
		if err != nil {
			fail(err)
		}

		if len(args) < 1 {
			failArgs([]string{"mode [server|client]"})
		}

		if args[0] == "client" {
			if len(args) < 2 {
				failArgs([]string{
					"please specify the adapter HW address that expose the service (eg. using hciconfig)",
				})
			}
		} else {
			args = append(args, "")
		}

		fail(service_example.Run(adapterID, args[0], args[1]))
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

}
