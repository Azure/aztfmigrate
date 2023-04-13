package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type timeSeriesInsightsEventSourcesResolver struct{}

func (timeSeriesInsightsEventSourcesResolver) ResourceTypes() []string {
	return []string{"azurerm_iot_time_series_insights_event_source_iothub", "azurerm_iot_time_series_insights_event_source_eventhub"}
}

func (timeSeriesInsightsEventSourcesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewTimeSeriesInsightEventSourcesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.EventSourceResourceClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch model.(type) {
	case *armtimeseriesinsights.IoTHubEventSourceResource:
		return "azurerm_iot_time_series_insights_event_source_iothub", nil
	case *armtimeseriesinsights.EventHubEventSourceResource:
		return "azurerm_iot_time_series_insights_event_source_eventhub", nil
	default:
		return "", fmt.Errorf("unknown environment type %T", model)
	}
}
