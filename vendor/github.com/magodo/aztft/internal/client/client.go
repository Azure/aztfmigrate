package client

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertsmanagement/armalertsmanagement"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datashare/armdatashare"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/desktopvirtualization/armdesktopvirtualization"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/digitaltwins/armdigitaltwins"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/domainservices/armdomainservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsight/armhdinsight"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iothub/armiothub"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/logic/armlogic"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentscripts"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagepool/armstoragepool"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloads/armworkloads"
)

type ClientBuilder struct {
	Cred      azcore.TokenCredential
	ClientOpt arm.ClientOptions
}

func (b *ClientBuilder) NewVirtualMachinesClient(subscriptionId string) (*armcompute.VirtualMachinesClient, error) {
	return armcompute.NewVirtualMachinesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewVirtualMachineScaleSetsClient(subscriptionId string) (*armcompute.VirtualMachineScaleSetsClient, error) {
	return armcompute.NewVirtualMachineScaleSetsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDevTestVirtualMachinesClient(subscriptionId string) (*armdevtestlabs.VirtualMachinesClient, error) {
	return armdevtestlabs.NewVirtualMachinesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewRecoveryservicesBackupProtectedItemsClient(subscriptionId string) (*armrecoveryservicesbackup.ProtectedItemsClient, error) {
	return armrecoveryservicesbackup.NewProtectedItemsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewRecoveryServicesBackupProtectionPoliciesClient(subscriptionId string) (*armrecoveryservicesbackup.ProtectionPoliciesClient, error) {
	return armrecoveryservicesbackup.NewProtectionPoliciesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataProtectionBackupPoliciesClient(subscriptionId string) (*armdataprotection.BackupPoliciesClient, error) {
	return armdataprotection.NewBackupPoliciesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataProtectionBackupInstancesClient(subscriptionId string) (*armdataprotection.BackupInstancesClient, error) {
	return armdataprotection.NewBackupInstancesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSynapseIntegrationRuntimesClient(subscriptionId string) (*armsynapse.IntegrationRuntimesClient, error) {
	return armsynapse.NewIntegrationRuntimesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDigitalTwinsEndpointsClient(subscriptionId string) (*armdigitaltwins.EndpointClient, error) {
	return armdigitaltwins.NewEndpointClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryTriggersClient(subscriptionId string) (*armdatafactory.TriggersClient, error) {
	return armdatafactory.NewTriggersClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryDatasetsClient(subscriptionId string) (*armdatafactory.DatasetsClient, error) {
	return armdatafactory.NewDatasetsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryDataFlowsClient(subscriptionId string) (*armdatafactory.DataFlowsClient, error) {
	return armdatafactory.NewDataFlowsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryLinkedServicesClient(subscriptionId string) (*armdatafactory.LinkedServicesClient, error) {
	return armdatafactory.NewLinkedServicesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryIntegrationRuntimesClient(subscriptionId string) (*armdatafactory.IntegrationRuntimesClient, error) {
	return armdatafactory.NewIntegrationRuntimesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDataFactoryCredentialsClient(subscriptionId string) (*armdatafactory.CredentialOperationsClient, error) {
	return armdatafactory.NewCredentialOperationsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewKustoDataConnectionsClient(subscriptionId string) (*armkusto.DataConnectionsClient, error) {
	return armkusto.NewDataConnectionsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewMachineLearningComputeClient(subscriptionId string) (*armmachinelearning.ComputeClient, error) {
	return armmachinelearning.NewComputeClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewMachineLearningDataStoreClient(subscriptionId string) (*armmachinelearning.DatastoresClient, error) {
	return armmachinelearning.NewDatastoresClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}
func (b *ClientBuilder) NewTimeSeriesInsightEnvironmentsClient(subscriptionId string) (*armtimeseriesinsights.EnvironmentsClient, error) {
	return armtimeseriesinsights.NewEnvironmentsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewTimeSeriesInsightEventSourcesClient(subscriptionId string) (*armtimeseriesinsights.EventSourcesClient, error) {
	return armtimeseriesinsights.NewEventSourcesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStorageCacheTargetsClient(subscriptionId string) (*armstoragecache.StorageTargetsClient, error) {
	return armstoragecache.NewStorageTargetsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAutomationConnectionClient(subscriptionId string) (*armautomation.ConnectionClient, error) {
	return armautomation.NewConnectionClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAutomationVariableClient(subscriptionId string) (*armautomation.VariableClient, error) {
	return armautomation.NewVariableClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAutomationJobScheduleClient(subscriptionId string) (*armautomation.JobScheduleClient, error) {
	return armautomation.NewJobScheduleClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewBotServiceBotsClient(subscriptionId string) (*armbotservice.BotsClient, error) {
	return armbotservice.NewBotsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewBotServiceChannelsClient(subscriptionId string) (*armbotservice.ChannelsClient, error) {
	return armbotservice.NewChannelsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSecurityInsightsDataConnectorsClient(subscriptionId string) (*armsecurityinsights.DataConnectorsClient, error) {
	return armsecurityinsights.NewDataConnectorsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSecurityInsightsAlertRulesClient(subscriptionId string) (*armsecurityinsights.AlertRulesClient, error) {
	return armsecurityinsights.NewAlertRulesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSecurityInsightsSecurityMLAnalyticsSettingsClient(subscriptionId string) (*armsecurityinsights.SecurityMLAnalyticsSettingsClient, error) {
	return armsecurityinsights.NewSecurityMLAnalyticsSettingsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewOperationalInsightsDataSourcesClient(subscriptionId string) (*armoperationalinsights.DataSourcesClient, error) {
	return armoperationalinsights.NewDataSourcesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAppPlatformBindingsClient(subscriptionId string) (*armappplatform.BindingsClient, error) {
	return armappplatform.NewBindingsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAppPlatformDeploymentsClient(subscriptionId string) (*armappplatform.DeploymentsClient, error) {
	return armappplatform.NewDeploymentsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDatashareDatasetsClient(subscriptionId string) (*armdatashare.DataSetsClient, error) {
	return armdatashare.NewDataSetsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewHDInsightClustersClient(subscriptionId string) (*armhdinsight.ClustersClient, error) {
	return armhdinsight.NewClustersClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStreamAnalyticsInputsClient(subscriptionId string) (*armstreamanalytics.InputsClient, error) {
	return armstreamanalytics.NewInputsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStreamAnalyticsOutputsClient(subscriptionId string) (*armstreamanalytics.OutputsClient, error) {
	return armstreamanalytics.NewOutputsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStreamAnalyticsFunctionsClient(subscriptionId string) (*armstreamanalytics.FunctionsClient, error) {
	return armstreamanalytics.NewFunctionsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewMonitorScheduledQueryRulesClient(subscriptionId string) (*armmonitor.ScheduledQueryRulesClient, error) {
	return armmonitor.NewScheduledQueryRulesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewCdnProfilesClient(subscriptionId string) (*armcdn.ProfilesClient, error) {
	return armcdn.NewProfilesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAppServiceCertificatesClient(subscriptionId string) (*armappservice.CertificatesClient, error) {
	return armappservice.NewCertificatesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAppServiceWebAppsClient(subscriptionId string) (*armappservice.WebAppsClient, error) {
	return armappservice.NewWebAppsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAppServiceEnvironmentsClient(subscriptionId string) (*armappservice.EnvironmentsClient, error) {
	return armappservice.NewEnvironmentsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewAlertsManagementProcessingRulesClient(subscriptionId string) (*armalertsmanagement.AlertProcessingRulesClient, error) {
	return armalertsmanagement.NewAlertProcessingRulesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDomainServiceClient(subscriptionId string) (*armdomainservices.Client, error) {
	return armdomainservices.NewClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStorageObjectReplicationPoliciesClient(subscriptionId string) (*armstorage.ObjectReplicationPoliciesClient, error) {
	return armstorage.NewObjectReplicationPoliciesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStorageFileSharesClient(subscriptionId string) (*armstorage.FileSharesClient, error) {
	return armstorage.NewFileSharesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStorageAccountsClient(subscriptionId string) (*armstorage.AccountsClient, error) {
	return armstorage.NewAccountsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewKeyVaultVaultsClient(subscriptionId string) (*armkeyvault.VaultsClient, error) {
	return armkeyvault.NewVaultsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewKeyVaultKeysClient(subscriptionId string) (*armkeyvault.KeysClient, error) {
	return armkeyvault.NewKeysClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewKeyVaultSecretsClient(subscriptionId string) (*armkeyvault.SecretsClient, error) {
	return armkeyvault.NewSecretsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkVirtualHubsClient(subscriptionId string) (*armnetwork.VirtualHubsClient, error) {
	return armnetwork.NewVirtualHubsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkVirtualHubBgpConnectionClient(subscriptionId string) (*armnetwork.VirtualHubBgpConnectionClient, error) {
	return armnetwork.NewVirtualHubBgpConnectionClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkInterfacesClient(subscriptionId string) (*armnetwork.InterfacesClient, error) {
	return armnetwork.NewInterfacesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkSubnetsClient(subscriptionId string) (*armnetwork.SubnetsClient, error) {
	return armnetwork.NewSubnetsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkNatGatewaysClient(subscriptionId string) (*armnetwork.NatGatewaysClient, error) {
	return armnetwork.NewNatGatewaysClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkPacketCapturesClient(subscriptionId string) (*armnetwork.PacketCapturesClient, error) {
	return armnetwork.NewPacketCapturesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkManagementDeploymentStatusClient(subscriptionId string) (*armnetwork.ManagerDeploymentStatusClient, error) {
	return armnetwork.NewManagerDeploymentStatusClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetworkLoadBalancersClient(subscriptionId string) (*armnetwork.LoadBalancersClient, error) {
	return armnetwork.NewLoadBalancersClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewFrontdoorPoliciesClient(subscriptionId string) (*armfrontdoor.PoliciesClient, error) {
	return armfrontdoor.NewPoliciesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDesktopVirtualizationWorkspacesClient(subscriptionId string) (*armdesktopvirtualization.WorkspacesClient, error) {
	return armdesktopvirtualization.NewWorkspacesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStoragePoolDiskPoolsClient(subscriptionId string) (*armstoragepool.DiskPoolsClient, error) {
	return armstoragepool.NewDiskPoolsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStoragePoolIscsiTargetsClient(subscriptionId string) (*armstoragepool.IscsiTargetsClient, error) {
	return armstoragepool.NewIscsiTargetsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewDeploymentScriptsClient(subscriptionId string) (*armdeploymentscripts.Client, error) {
	return armdeploymentscripts.NewClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSiteRecoveryReplicationPoliciesClient(subscriptionId, resourceGroupName, vaultName string) (*armrecoveryservicessiterecovery.ReplicationPoliciesClient, error) {
	return armrecoveryservicessiterecovery.NewReplicationPoliciesClient(
		vaultName,
		resourceGroupName,
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSiteRecoveryReplicationFabricsClient(subscriptionId, resourceGroupName, vaultName string) (*armrecoveryservicessiterecovery.ReplicationFabricsClient, error) {
	return armrecoveryservicessiterecovery.NewReplicationFabricsClient(
		vaultName,
		resourceGroupName,
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSiteRecoveryReplicationProtectedItemsClient(subscriptionId, resourceGroupName, vaultName string) (*armrecoveryservicessiterecovery.ReplicationProtectedItemsClient, error) {
	return armrecoveryservicessiterecovery.NewReplicationProtectedItemsClient(
		vaultName,
		resourceGroupName,
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSiteRecoveryReplicationProtectionContainerMappingsClient(subscriptionId, resourceGroupName, vaultName string) (*armrecoveryservicessiterecovery.ReplicationProtectionContainerMappingsClient, error) {
	return armrecoveryservicessiterecovery.NewReplicationProtectionContainerMappingsClient(
		vaultName,
		resourceGroupName,
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewSiteRecoveryReplicationNetworkMappingsClient(subscriptionId, resourceGroupName, vaultName string) (*armrecoveryservicessiterecovery.ReplicationNetworkMappingsClient, error) {
	return armrecoveryservicessiterecovery.NewReplicationNetworkMappingsClient(
		vaultName,
		resourceGroupName,
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewStorageMoverEndpointsClient(subscriptionId string) (*armstoragemover.EndpointsClient, error) {
	return armstoragemover.NewEndpointsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewCostManagementScheduledActionsClient() (*armcostmanagement.ScheduledActionsClient, error) {
	return armcostmanagement.NewScheduledActionsClient(
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewApplicationInsightsWebTestsClient(subscriptionId string) (*armapplicationinsights.WebTestsClient, error) {
	return armapplicationinsights.NewWebTestsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewLogicWorkflowsClient(subscriptionId string) (*armlogic.WorkflowsClient, error) {
	return armlogic.NewWorkflowsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewPaloalToNetworkFirewallsClient(subscriptionId string) (*armpanngfw.FirewallsClient, error) {
	return armpanngfw.NewFirewallsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewIothubsClient(subscriptionId string) (*armiothub.ResourceClient, error) {
	return armiothub.NewResourceClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewApiManagementApiClient(subscriptionId string) (*armapimanagement.APIClient, error) {
	return armapimanagement.NewAPIClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewNetAppAccountClient(subscriptionId string) (*armnetapp.AccountsClient, error) {
	return armnetapp.NewAccountsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewWorkloadSAPVirtualInstanceClient(subscriptionId string) (*armworkloads.SAPVirtualInstancesClient, error) {
	return armworkloads.NewSAPVirtualInstancesClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}

func (b *ClientBuilder) NewContainerAppEnvironmentsClient(subscriptionId string) (*armappcontainers.ManagedEnvironmentsClient, error) {
	return armappcontainers.NewManagedEnvironmentsClient(
		subscriptionId,
		b.Cred,
		&b.ClientOpt,
	)
}
