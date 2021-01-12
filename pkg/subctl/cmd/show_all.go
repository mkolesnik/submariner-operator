/*
© 2021 Red Hat, Inc. and others

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
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/submariner-io/submariner-operator/apis/submariner/v1alpha1"
	"k8s.io/client-go/rest"
)

var showAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Show information related to a submariner cluster",
	Long: `This command shows information related to a submariner cluster:
 networks, endpoints, gateways, connections and component versions.`,
	Run: showAll,
}

func init() {
	showCmd.AddCommand(showAllCmd)
}

func showAll(cmd *cobra.Command, args []string) {
	printInfoForClustersUsing(allInfoPrinter)
}

func allInfoPrinter(config *rest.Config, submariner *v1alpha1.Submariner) {
	fmt.Println("Showing Network details")
	showNetworkSingleCluster(config)
	fmt.Println("")
	fmt.Println("\nShowing Endpoint details")
	endpointsPrinter(config, submariner)
	fmt.Println("\nShowing Connection details")
	connectionsPrinter(config, submariner)
	fmt.Println("\nShowing Gateway details")
	gatewaysPrinter(config, submariner)
	fmt.Println("\nShowing version details")
	clusterVersionsPrinter(config, submariner)
}
