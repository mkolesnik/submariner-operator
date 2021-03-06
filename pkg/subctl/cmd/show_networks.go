/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

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
	submarinerclientset "github.com/submariner-io/submariner-operator/pkg/client/clientset/versioned"
	"github.com/submariner-io/submariner-operator/pkg/discovery/network"
	"github.com/submariner-io/submariner-operator/pkg/subctl/cmd/utils"
	"github.com/submariner-io/submariner-operator/pkg/subctl/cmd/utils/restconfig"
	"k8s.io/client-go/rest"
)

// showNetworksCmd represents the show networks command
var showNetworksCmd = &cobra.Command{
	Use:   "networks",
	Short: "Get information on your cluster related to submariner",
	Long: `This command shows the status of submariner in your cluster,
and the relevant network details from your cluster.`,
	PreRunE: checkVersionMismatch,
	Run:     showNetwork,
}

func init() {
	showCmd.AddCommand(showNetworksCmd)
}

func showNetwork(cmd *cobra.Command, args []string) {
	configs, err := restconfig.ForClusters(kubeConfig, kubeContexts)
	utils.ExitOnError("Error getting REST config for cluster", err)

	for _, item := range configs {
		fmt.Println()
		fmt.Printf("Showing network details for cluster %q:\n", item.ClusterName)
		showNetworkSingleCluster(item.Config)
	}
}

func showNetworkSingleCluster(config *rest.Config) {
	submariner := getSubmarinerResource(config)

	var clusterNetwork *network.ClusterNetwork
	var msg string
	if submariner != nil {
		msg = "    Discovered network details via Submariner:"
		clusterNetwork = &network.ClusterNetwork{
			PodCIDRs:      []string{submariner.Status.ClusterCIDR},
			ServiceCIDRs:  []string{submariner.Status.ServiceCIDR},
			NetworkPlugin: submariner.Status.NetworkPlugin,
			GlobalCIDR:    submariner.Status.GlobalCIDR,
		}
	} else {
		msg = "    Discovered network details"
		dynClient, clientSet, err := restconfig.Clients(config)
		utils.ExitOnError("Error creating clients for cluster", err)

		submarinerClient, err := submarinerclientset.NewForConfig(config)
		utils.ExitOnError("Unable to get the Submariner client", err)

		clusterNetwork, err = network.Discover(dynClient, clientSet, submarinerClient, OperatorNamespace)
		utils.ExitOnError("There was an error discovering network details for this cluster", err)
	}
	if clusterNetwork != nil {
		fmt.Println(msg)
	}
	clusterNetwork.Show()
}
