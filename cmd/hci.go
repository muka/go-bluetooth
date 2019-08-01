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
	hci_updown_example "github.com/muka/go-bluetooth/examples/hci_updown"
	"github.com/spf13/cobra"
)

// hciCmd represents the hci command
var hciCmd = &cobra.Command{
	Use:   "hci",
	Short: "Example usage of the HCI interface",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		adapterID, err := cmd.Flags().GetString("adapterID")
		if err != nil {
			fail(err)
		}
		fail(hci_updown_example.Run(adapterID))
	},
}

func init() {
	rootCmd.AddCommand(hciCmd)
}
