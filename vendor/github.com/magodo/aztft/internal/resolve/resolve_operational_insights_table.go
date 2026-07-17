package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v3"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type operationalInsightsTableResolver struct{}

func (operationalInsightsTableResolver) ResourceTypes() []string {
	return []string{
		"azurerm_log_analytics_workspace_table_custom_log", "azurerm_log_analytics_workspace_table_microsoft",
	}
}

func (operationalInsightsTableResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewOperationalInsightsTablesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}

	schema := props.Schema
	if schema == nil {
		return "", fmt.Errorf("unexpected nil properties.schema in response")
	}

	tableType := schema.TableType
	if tableType == nil {
		return "", fmt.Errorf("unexpected nil properties.schema.tableType in response")
	}

	switch *tableType {
	case armoperationalinsights.TableTypeEnumCustomLog:
		return "azurerm_log_analytics_workspace_table_custom_log", nil
	case armoperationalinsights.TableTypeEnumMicrosoft:
		return "azurerm_log_analytics_workspace_table_microsoft", nil
	default:
		return "", fmt.Errorf("uknown table type: %v", *tableType)
	}
}
