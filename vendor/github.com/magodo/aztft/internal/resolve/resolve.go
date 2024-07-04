package resolve

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type resolver interface {
	Resolve(*client.ClientBuilder, armid.ResourceId) (string, error)
	ResourceTypes() []string
}

var Resolvers = map[string]map[string]resolver{
	"/MICROSOFT.COMPUTE/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virtualMachinesResolver{},
	},
	"/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virtualMachineScaleSetsResolver{},
	},
	"/MICROSOFT.DEVTESTLAB/LABS/VIRTUALMACHINES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": devTestVirtualMachinesResolver{},
	},
	"/MICROSOFT.APIMANAGEMENT/SERVICE/IDENTITYPROVIDERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": apiManagementIdentitiesResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/BACKUPPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": recoveryServicesBackupProtectionPoliciesResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/BACKUPFABRICS/PROTECTIONCONTAINERS/PROTECTEDITEMS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": recoveryServicesBackupProtectedItemsResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/REPLICATIONFABRICS/REPLICATIONPROTECTIONCONTAINERS/REPLICATIONPROTECTEDITEMS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": recoveryServicesReplicationProtectedItemsResolver{},
	},
	"/MICROSOFT.DATAPROTECTION/BACKUPVAULTS/BACKUPPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataProtectionBackupPoliciesResolver{},
	},
	"/MICROSOFT.DATAPROTECTION/BACKUPVAULTS/BACKUPINSTANCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataProtectionBackupInstancesResolver{},
	},
	"/MICROSOFT.SYNAPSE/WORKSPACES/INTEGRATIONRUNTIMES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": synapseIntegrationRuntimesResolver{},
	},
	"/MICROSOFT.DIGITALTWINS/DIGITALTWINSINSTANCES/ENDPOINTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": digitalTwinsEndpointsResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/TRIGGERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryTriggersResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/DATASETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryDatasetsResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/DATAFLOWS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryDataFlowsResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/LINKEDSERVICES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryLinkedServicesResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/INTEGRATIONRUNTIMES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryIntegrationRuntimesResolver{},
	},
	"/MICROSOFT.DATAFACTORY/FACTORIES/CREDENTIALS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": dataFactoryCredentialsResolver{},
	},
	"/MICROSOFT.KUSTO/CLUSTERS/DATABASES/DATACONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": kustoDataConnectionsResolver{},
	},
	"/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/COMPUTES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": machineLearningComputesResolver{},
	},
	"/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/DATASTORES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": machineLearningDataStoresResolver{},
	},
	"/MICROSOFT.TIMESERIESINSIGHTS/ENVIRONMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": timeSeriesInsightsEnvironmentResolver{},
	},
	"/MICROSOFT.TIMESERIESINSIGHTS/ENVIRONMENTS/EVENTSOURCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": timeSeriesInsightsEventSourcesResolver{},
	},
	"/MICROSOFT.STORAGECACHE/CACHES/STORAGETARGETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": storageCacheTargetsResolver{},
	},
	"/MICROSOFT.AUTOMATION/AUTOMATIONACCOUNTS/CONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": automationConnectionsResolver{},
	},
	"/MICROSOFT.AUTOMATION/AUTOMATIONACCOUNTS/VARIABLES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": automationVariablesResolver{},
	},
	"/MICROSOFT.BOTSERVICE/BOTSERVICES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": botServiceBotsResolver{},
	},
	"/MICROSOFT.BOTSERVICE/BOTSERVICES/CHANNELS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": botServiceChannelsResolver{},
	},
	"/MICROSOFT.SECURITYINSIGHTS/DATACONNECTORS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": securityInsightsDataConnectorsResolver{},
	},
	"/MICROSOFT.SECURITYINSIGHTS/ALERTRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": securityInsightsAlertRulesResolver{},
	},
	"/MICROSOFT.SECURITYINSIGHTS/SECURITYMLANALYTICSSETTINGS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES": securityInsightsSecurityMLAnalyticsSettingsResolver{},
	},
	"/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES/DATASOURCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": operationalInsightsDataSourcesResolver{},
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APPS/BINDINGS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appPlatformBindingsResolver{},
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APPS/DEPLOYMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appPlatformDeploymentsResolver{},
	},
	"/MICROSOFT.DATASHARE/ACCOUNTS/SHARES/DATASETS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": datashareDatasetsResolver{},
	},
	"/MICROSOFT.HDINSIGHT/CLUSTERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": hdInsightClustersResolver{},
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/INPUTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": streamAnalyticsInputsResolver{},
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/OUTPUTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": streamAnalyticsOutputsResolver{},
	},
	"/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/FUNCTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": streamAnalyticsFunctionsResolver{},
	},
	"/MICROSOFT.INSIGHTS/SCHEDULEDQUERYRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": monitorScheduledQueryRulesResolver{},
	},
	"/MICROSOFT.CDN/PROFILES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": cdnProfilesResolver{},
	},
	"/MICROSOFT.WEB/CERTIFICATES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceCertificatesResolver{},
	},
	"/MICROSOFT.WEB/SITES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceSitesResolver{},
	},
	"/MICROSOFT.WEB/SITES/SLOTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceSiteSlotsResolver{},
	},
	"/MICROSOFT.WEB/SITES/HYBRIDCONNECTIONNAMESPACES/RELAYS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceSiteHybridConnectionsResolver{},
	},
	"/MICROSOFT.WEB/HOSTINGENVIRONMENTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": appServiceEnvironemntsResolver{},
	},
	"/MICROSOFT.ALERTSMANAGEMENT/ACTIONRULES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": alertsManagementProcessingRulesResolver{},
	},
	"/MICROSOFT.NETWORK/VIRTUALHUBS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virtualHubsResolver{},
	},
	"/MICROSOFT.NETWORK/VIRTUALHUBS/BGPCONNECTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virtualHubBgpConnectionsResolver{},
	},
	"/MICROSOFT.NETWORK/FRONTDOORWEBAPPLICATIONFIREWALLPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": frontdoorPoliciesResolver{},
	},
	"/MICROSOFT.NETWORK/NETWORKWATCHERS/PACKETCAPTURES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": networkPacketCaptureResolver{},
	},
	"/MICROSOFT.RESOURCES/DEPLOYMENTSCRIPTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": deploymentScriptsResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/REPLICATIONPOLICIES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": siteRecoveryReplicationPoliciesResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/REPLICATIONFABRICS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": siteRecoveryReplicationFabricsResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/REPLICATIONFABRICS/REPLICATIONPROTECTIONCONTAINERS/REPLICATIONPROTECTIONCONTAINERMAPPINGS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": siteRecoveryReplicationProtectionContainerMappingResolver{},
	},
	"/MICROSOFT.RECOVERYSERVICES/VAULTS/REPLICATIONFABRICS/REPLICATIONNETWORKS/REPLICATIONNETWORKMAPPINGS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": siteRecoveryReplicationNetworkMappingResolver{},
	},
	"/MICROSOFT.STORAGEMOVER/STORAGEMOVERS/ENDPOINTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": storageMoverEndpointsResolver{},
	},
	"/MICROSOFT.COSTMANAGEMENT/SCHEDULEDACTIONS": {
		"/SUBSCRIPTIONS": costmanagementScheduleActionsResolver{},
	},
	"/MICROSOFT.INSIGHTS/WEBTESTS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": applicationInsightsWebTestsResolver{},
	},
	"/MICROSOFT.LOGIC/WORKFLOWS/ACTIONS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": logicAppAction{},
	},
	"/MICROSOFT.LOGIC/WORKFLOWS/TRIGGERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": logicAppTrigger{},
	},
	"/PALOALTONETWORKS.CLOUDNGFW/FIREWALLS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": paloalToNetworkFirewall{},
	},
	"/MICROSOFT.SERVICELINKER/LINKERS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.WEB/SITES": serviceConnectorAppServiceResolver{},
	},
	"/MICROSOFT.APPPLATFORM/SPRING/APMS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": springApmsResolver{},
	},
	"/MICROSOFT.WORKLOADS/SAPVIRTUALINSTANCES": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": sapVirtualInstancesResolver{},
	},
	"/MICROSOFT.COMPUTE/VIRTUALMACHINES/DATADISKS": {
		"/SUBSCRIPTIONS/RESOURCEGROUPS": virutalMachineDataDiskResolver{},
	},
}

type ResolveError struct {
	ResourceId armid.ResourceId
	Err        error
}

func (e ResolveError) Error() string {
	return e.ResourceId.String() + ": " + e.Err.Error()
}

func (e *ResolveError) Unwrap() error { return e.Err }

func getResolver(id armid.ResourceId) (resolver, bool) {
	routeKey := strings.ToUpper(id.RouteScopeString())
	var parentScopeKey string
	if id.ParentScope() != nil {
		parentScopeKey = strings.ToUpper(id.ParentScope().ScopeString())
	}
	m, ok := Resolvers[routeKey]
	if !ok {
		return nil, false
	}
	resolver, ok := m[parentScopeKey]
	return resolver, ok
}

func NeedsAPI(id armid.ResourceId) bool {
	_, ok := getResolver(id)
	return ok
}

// Resolve resolves a given resource id via Azure API to disambiguate and return a single matched TF resource type.
func Resolve(id armid.ResourceId, cred azcore.TokenCredential, clientOpt arm.ClientOptions) (string, error) {
	// Ensure the API client can be built.
	b := &client.ClientBuilder{Cred: cred, ClientOpt: clientOpt}

	resolver, ok := getResolver(id)
	if !ok {
		return "", ResolveError{ResourceId: id, Err: fmt.Errorf("no resolver found for %q", id)}
	}
	rt, err := resolver.Resolve(b, id)
	if err != nil {
		return "", ResolveError{ResourceId: id, Err: fmt.Errorf("resolving %q: %v", id, err)}
	}
	return rt, nil
}
