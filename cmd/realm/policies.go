// Copyright © 2023 SECO Mind Srl
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

package realm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/astarte-platform/astartectl/utils"
	"github.com/spf13/cobra"
)

// policiesCmd represents the policies command
var policiesCmd = &cobra.Command{
	Use:     "policies",
	Short:   "Manage policies",
	Long:    `List, show, install or delete trigger delivery policies in your realm.`,
	Aliases: []string{"policy"},
}

var policiesListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List trigger delivery policies",
	Long:    `List the name of trigger delivery policies installed in the realm.`,
	Example: `  astartectl realm-management policies list`,
	RunE:    policiesListF,
	Aliases: []string{"ls"},
}

var policiesShowCmd = &cobra.Command{
	Use:     "show <policy_name>",
	Short:   "Show trigger delivery policy",
	Long:    `Shows a trigger delivery policy installed in the realm.`,
	Example: `  astartectl realm-management policies show my_policy`,
	Args:    cobra.ExactArgs(1),
	RunE:    policiesShowF,
}

var policiesInstallCmd = &cobra.Command{
	Use:   "install <policy_file>",
	Short: "Install trigger delivery policy",
	Long: `Install the given trigger delivery policies in the realm.
<policy_file> must be a path to a JSON file containing a valid Astarte trigger delivery policy.`,
	Example: `  astartectl realm-management policies install my_policy.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    policiesInstallF,
}

var policiesDeleteCmd = &cobra.Command{
	Use:     "delete <policy_name>",
	Short:   "Delete trigger delivery policy",
	Long:    `Deletes the specified trigger delivery policies from the realm.`,
	Example: `  astartectl realm-management policies delete my_policy`,
	Args:    cobra.ExactArgs(1),
	RunE:    policiesDeleteF,
	Aliases: []string{"del"},
}

func init() {
	RealmManagementCmd.AddCommand(policiesCmd)

	policiesCmd.AddCommand(
		policiesListCmd,
		policiesShowCmd,
		policiesInstallCmd,
		policiesDeleteCmd,
	)
}

func policiesListF(command *cobra.Command, args []string) error {
	policiesCall, err := astarteAPIClient.ListTriggerDeliveryPolicies(realm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	utils.MaybeCurlAndExit(policiesCall, astarteAPIClient)

	policiesRes, err := policiesCall.Run(astarteAPIClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rawPolicies, _ := policiesRes.Parse()
	policies, _ := rawPolicies.([]string)
	fmt.Println(policies)
	return nil
}

func policiesShowF(command *cobra.Command, args []string) error {
	policyName := args[0]

	getPolicyCall, err := astarteAPIClient.GetTriggerDeliveryPolicy(realm, policyName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	utils.MaybeCurlAndExit(getPolicyCall, astarteAPIClient)

	getPolicyRes, err := getPolicyCall.Run(astarteAPIClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	policyDefinition, _ := getPolicyRes.Parse()
	respJSON, _ := json.MarshalIndent(policyDefinition, "", "  ")
	fmt.Println(string(respJSON))

	return nil
}

func policiesInstallF(command *cobra.Command, args []string) error {
	policyFile, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	var policyBody map[string]interface{}
	err = json.Unmarshal(policyFile, &policyBody)
	if err != nil {
		return err
	}

	installPolicyCall, err := astarteAPIClient.InstallTriggerDeliveryPolicy(realm, policyBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	utils.MaybeCurlAndExit(installPolicyCall, astarteAPIClient)

	installPolicyRes, err := installPolicyCall.Run(astarteAPIClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, _ = installPolicyRes.Parse()

	fmt.Println("ok")
	return nil
}

func policiesDeleteF(command *cobra.Command, args []string) error {
	policyName := args[0]
	deletePolicyCall, err := astarteAPIClient.DeleteTriggerDeliveryPolicy(realm, policyName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	utils.MaybeCurlAndExit(deletePolicyCall, astarteAPIClient)

	deletePolicyRes, err := deletePolicyCall.Run(astarteAPIClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, _ = deletePolicyRes.Parse()

	fmt.Println("ok")
	return nil
}
