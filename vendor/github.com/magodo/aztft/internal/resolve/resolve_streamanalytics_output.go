package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type streamAnalyticsOutputsResolver struct{}

func (streamAnalyticsOutputsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_stream_analytics_output_servicebus_topic",
		"azurerm_stream_analytics_output_blob",
		"azurerm_stream_analytics_output_mssql",
		"azurerm_stream_analytics_output_table",
		"azurerm_stream_analytics_output_cosmosdb",
		"azurerm_stream_analytics_output_servicebus_queue",
		"azurerm_stream_analytics_output_eventhub",
		"azurerm_stream_analytics_output_powerbi",
		"azurerm_stream_analytics_output_synapse",
		"azurerm_stream_analytics_output_function",
	}
}

func (streamAnalyticsOutputsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStreamAnalyticsOutputsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Output.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	ds := props.Datasource
	if ds == nil {
		return "", fmt.Errorf("unexpected nil properties.datasource in response")
	}
	switch ds.(type) {
	case *armstreamanalytics.ServiceBusTopicOutputDataSource:
		return "azurerm_stream_analytics_output_servicebus_topic", nil
	case *armstreamanalytics.BlobOutputDataSource:
		return "azurerm_stream_analytics_output_blob", nil
	case *armstreamanalytics.AzureSQLDatabaseOutputDataSource:
		return "azurerm_stream_analytics_output_mssql", nil
	case *armstreamanalytics.AzureTableOutputDataSource:
		return "azurerm_stream_analytics_output_table", nil
	case *armstreamanalytics.DocumentDbOutputDataSource:
		return "azurerm_stream_analytics_output_cosmosdb", nil
	case *armstreamanalytics.ServiceBusQueueOutputDataSource:
		return "azurerm_stream_analytics_output_servicebus_queue", nil
	case *armstreamanalytics.EventHubOutputDataSource:
		return "azurerm_stream_analytics_output_eventhub", nil
	case *armstreamanalytics.PowerBIOutputDataSource:
		return "azurerm_stream_analytics_output_powerbi", nil
	case *armstreamanalytics.AzureSynapseOutputDataSource:
		return "azurerm_stream_analytics_output_synapse", nil
	case *armstreamanalytics.AzureFunctionOutputDataSource:
		return "azurerm_stream_analytics_output_function", nil
	default:
		return "", fmt.Errorf("unknown output data source type: %T", ds)
	}
}
