package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type automationConnectionsResolver struct{}

func (automationConnectionsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_automation_connection_service_principal",
		"azurerm_automation_connection_certificate",
		"azurerm_automation_connection_classic_certificate",
		"azurerm_automation_connection",
	}
}

func (automationConnectionsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAutomationConnectionClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Connection.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	connType := props.ConnectionType
	if connType == nil {
		return "", fmt.Errorf("unexpected nil properties.connectionType in response")
	}
	connTypeName := connType.Name
	if connTypeName == nil {
		return "", fmt.Errorf("unexpected nil property.connectionType.name in response")
	}

	switch *connTypeName {
	case "AzureServicePrincipal":
		return "azurerm_automation_connection_service_principal", nil
	case "Azure":
		return "azurerm_automation_connection_certificate", nil
	case "AzureClassicCertificate":
		return "azurerm_automation_connection_classic_certificate", nil
	default:
		return "azurerm_automation_connection", nil
	}
}
