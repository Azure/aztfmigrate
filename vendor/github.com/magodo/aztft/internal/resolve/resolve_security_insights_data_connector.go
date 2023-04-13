package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type securityInsightsDataConnectorsResolver struct{}

func (securityInsightsDataConnectorsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_sentinel_data_connector_microsoft_cloud_app_security",
		"azurerm_sentinel_data_connector_azure_active_directory",
		"azurerm_sentinel_data_connector_office_365",
		"azurerm_sentinel_data_connector_office_atp",
		"azurerm_sentinel_data_connector_threat_intelligence",
		"azurerm_sentinel_data_connector_aws_s3",
		"azurerm_sentinel_data_connector_aws_cloud_trail",
		"azurerm_sentinel_data_connector_azure_security_center",
		"azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection",
		"azurerm_sentinel_data_connector_azure_advanced_threat_protection",
		"azurerm_sentinel_data_connector_office_365_project",
		"azurerm_sentinel_data_connector_dynamics_365",
		"azurerm_sentinel_data_connector_iot",
		"azurerm_sentinel_data_connector_office_irm",
		"azurerm_sentinel_data_connector_office_power_bi",
		"azurerm_sentinel_data_connector_microsoft_threat_protection",
		"azurerm_sentinel_data_connector_threat_intelligence_taxii",
		"azurerm_sentinel_data_connector_microsoft_threat_intelligence",
	}
}

func (securityInsightsDataConnectorsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSecurityInsightsDataConnectorsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.ParentScope().Names()[0], id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.DataConnectorClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}

	switch model.(type) {
	case *armsecurityinsights.MCASDataConnector:
		return "azurerm_sentinel_data_connector_microsoft_cloud_app_security", nil
	case *armsecurityinsights.AADDataConnector:
		return "azurerm_sentinel_data_connector_azure_active_directory", nil
	case *armsecurityinsights.OfficeDataConnector:
		return "azurerm_sentinel_data_connector_office_365", nil
	case *armsecurityinsights.OfficeIRMDataConnector:
		return "azurerm_sentinel_data_connector_office_irm", nil
	case *armsecurityinsights.OfficePowerBIDataConnector:
		return "azurerm_sentinel_data_connector_office_power_bi", nil
	case *armsecurityinsights.Office365ProjectDataConnector:
		return "azurerm_sentinel_data_connector_office_365_project", nil
	case *armsecurityinsights.OfficeATPDataConnector:
		return "azurerm_sentinel_data_connector_office_atp", nil
	case *armsecurityinsights.TIDataConnector:
		return "azurerm_sentinel_data_connector_threat_intelligence", nil
	case *armsecurityinsights.AwsS3DataConnector:
		return "azurerm_sentinel_data_connector_aws_s3", nil
	case *armsecurityinsights.AwsCloudTrailDataConnector:
		return "azurerm_sentinel_data_connector_aws_cloud_trail", nil
	case *armsecurityinsights.ASCDataConnector:
		return "azurerm_sentinel_data_connector_azure_security_center", nil
	case *armsecurityinsights.MDATPDataConnector:
		return "azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection", nil
	case *armsecurityinsights.AATPDataConnector:
		return "azurerm_sentinel_data_connector_azure_advanced_threat_protection", nil
	case *armsecurityinsights.Dynamics365DataConnector:
		return "azurerm_sentinel_data_connector_dynamics_365", nil
	case *armsecurityinsights.IoTDataConnector:
		return "azurerm_sentinel_data_connector_iot", nil
	case *armsecurityinsights.MTPDataConnector:
		return "azurerm_sentinel_data_connector_microsoft_threat_protection", nil
	case *armsecurityinsights.TiTaxiiDataConnector:
		return "azurerm_sentinel_data_connector_threat_intelligence_taxii", nil
	case *armsecurityinsights.MSTIDataConnector:
		return "azurerm_sentinel_data_connector_microsoft_threat_intelligence", nil
	default:
		return "", fmt.Errorf("unknown data connector type: %T", model)
	}
}
