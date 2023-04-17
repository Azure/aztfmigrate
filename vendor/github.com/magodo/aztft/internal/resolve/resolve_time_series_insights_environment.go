package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type timeSeriesInsightsEnvironmentResolver struct{}

func (timeSeriesInsightsEnvironmentResolver) ResourceTypes() []string {
	return []string{"azurerm_iot_time_series_insights_standard_environment", "azurerm_iot_time_series_insights_gen2_environment"}
}

func (timeSeriesInsightsEnvironmentResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewTimeSeriesInsightEnvironmentsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.EnvironmentResourceClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch model.(type) {
	case *armtimeseriesinsights.Gen1EnvironmentResource:
		return "azurerm_iot_time_series_insights_standard_environment", nil
	case *armtimeseriesinsights.Gen2EnvironmentResource:
		return "azurerm_iot_time_series_insights_gen2_environment", nil
	default:
		return "", fmt.Errorf("unknown environment type %T", model)
	}
}
