// Copyright © 2019 luca capra
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
	agent_example "github.com/muka/go-bluetooth/examples/agent"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "A bluez Agent1 example",
	Long:  `An example of agent interaction to exchange a passkey during pairing`,
	Run: func(cmd *cobra.Command, args []string) {

		// id, err := cmd.Flags().GetString("id")
		// if err != nil {
		// 	fail(err)
		// }

		adapterID, err := cmd.Flags().GetString("adapterID")
		if err != nil {
			fail(err)
		}

		if len(args) == 0 {
			failArg("Device mac")
		}
		id := args[0]

		fail(agent_example.Run(id, adapterID))
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)

	agentCmd.Flags().String("id", "", "Help message for toggle")
}
