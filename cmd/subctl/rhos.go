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

//nolint:dupl // Similar code in aws.go, azure.go, rhos.go, but not duplicate
package subctl

import (
	"github.com/spf13/cobra"
	"github.com/submariner-io/admiral/pkg/reporter"
	"github.com/submariner-io/subctl/internal/cli"
	"github.com/submariner-io/subctl/internal/exit"
	"github.com/submariner-io/subctl/pkg/cloud/cleanup"
	"github.com/submariner-io/subctl/pkg/cloud/prepare"
	"github.com/submariner-io/subctl/pkg/cloud/rhos"
	"github.com/submariner-io/subctl/pkg/cluster"
)

var (
	rhosConfig rhos.Config

	rhosPrepareCmd = &cobra.Command{
		Use:     "rhos",
		Short:   "Prepare an OpenShift RHOS cloud",
		Long:    "This command prepares an OpenShift installer-provisioned infrastructure (IPI) on RHOS cloud for Submariner installation.",
		PreRunE: checkRHOSFlags,
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnError(cloudRestConfigProducer.RunOnSelectedContext(
				func(clusterInfo *cluster.Info, namespace string, status reporter.Interface) error {
					return prepare.RHOS( //nolint:wrapcheck // Not needed.
						clusterInfo, &cloudOptions.ports, &rhosConfig, cloudOptions.useLoadBalancer, status)
				}, cli.NewReporter()))
		},
	}

	rhosCleanupCmd = &cobra.Command{
		Use:   "rhos",
		Short: "Clean up an RHOS cloud",
		Long: "This command cleans up an OpenShift installer-provisioned infrastructure (IPI) on RHOS-based" +
			" cloud after Submariner uninstallation.",
		PreRunE: checkRHOSFlags,
		Run: func(cmd *cobra.Command, args []string) {
			exit.OnError(cloudRestConfigProducer.RunOnSelectedContext(
				func(clusterInfo *cluster.Info, namespace string, status reporter.Interface) error {
					return cleanup.RHOS(clusterInfo, &rhosConfig, status) //nolint:wrapcheck // No need to wrap errors here.
				}, cli.NewReporter()))
		},
	}
)

func init() {
	addGeneralRHOSFlags := func(command *cobra.Command) {
		command.Flags().StringVar(&rhosConfig.InfraID, infraIDFlag, "", "OpenStack infra ID")
		command.Flags().StringVar(&rhosConfig.Region, regionFlag, "", "OpenStack region")
		command.Flags().StringVar(&rhosConfig.ProjectID, projectIDFlag, "", "OpenStack project ID")
		command.Flags().StringVar(&rhosConfig.OcpMetadataFile, "ocp-metadata", "",
			"OCP metadata.json file (or the directory containing it) from which to read the RHOS infra ID "+
				"and region from (takes precedence over the specific flags)")
		command.Flags().StringVar(&rhosConfig.CloudEntry, cloudEntryFlag, "", "Specific cloud configuration to use from the clouds.yaml")
	}

	addGeneralRHOSFlags(rhosPrepareCmd)
	rhosPrepareCmd.Flags().IntVar(&rhosConfig.Gateways, "gateways", defaultNumGateways,
		"Number of gateways to deploy")
	rhosPrepareCmd.Flags().StringVar(&rhosConfig.GWInstanceType, "gateway-instance", "PnTAE.CPU_4_Memory_8192_Disk_50",
		"Type of gateway instance machine")
	rhosPrepareCmd.Flags().BoolVar(&rhosConfig.DedicatedGateway, "dedicated-gateway", true,
		"Whether a dedicated gateway node has to be deployed")

	cloudPrepareCmd.AddCommand(rhosPrepareCmd)

	addGeneralRHOSFlags(rhosCleanupCmd)
	cloudCleanupCmd.AddCommand(rhosCleanupCmd)
}

func checkRHOSFlags(cmd *cobra.Command, args []string) error {
	if rhosConfig.OcpMetadataFile == "" {
		expectFlag(infraIDFlag, rhosConfig.InfraID)
		expectFlag(regionFlag, rhosConfig.Region)
		expectFlag(projectIDFlag, rhosConfig.ProjectID)
	}

	return nil
}
