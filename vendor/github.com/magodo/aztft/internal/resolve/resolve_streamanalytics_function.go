package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type streamAnalyticsFunctionsResolver struct{}

func (streamAnalyticsFunctionsResolver) ResourceTypes() []string {
	return []string{"azurerm_stream_analytics_function_javascript_uda", "azurerm_stream_analytics_function_javascript_udf"}
}

func (streamAnalyticsFunctionsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStreamAnalyticsFunctionsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Function.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props := props.(type) {
	case *armstreamanalytics.AggregateFunctionProperties:
		return "azurerm_stream_analytics_function_javascript_uda", nil
	case *armstreamanalytics.FunctionProperties:
		return "azurerm_stream_analytics_function_javascript_udf", nil
	default:
		return "", fmt.Errorf("unknown input property type: %T", props)
	}
}
