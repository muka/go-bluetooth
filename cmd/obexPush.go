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
	obex_push_example "github.com/muka/go-bluetooth/examples/obex_push"
	"github.com/spf13/cobra"
)

// obexPushCmd represents the obexPush command
var obexPushCmd = &cobra.Command{
	Use:   "obex-push",
	Short: "Obex push example",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		adapterID, err := cmd.Flags().GetString("adapterID")
		if err != nil {
			fail(err)
		}

		if len(args) < 2 {
			failArgs([]string{"target_address", "file_path"})
		}

		fail(obex_push_example.Run(args[0], args[1], adapterID))
	},
}

func init() {
	rootCmd.AddCommand(obexPushCmd)
}
