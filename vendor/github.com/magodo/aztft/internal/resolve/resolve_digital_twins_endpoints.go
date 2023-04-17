package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/digitaltwins/armdigitaltwins"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type digitalTwinsEndpointsResolver struct{}

func (digitalTwinsEndpointsResolver) ResourceTypes() []string {
	return []string{"azurerm_digital_twins_endpoint_eventgrid", "azurerm_digital_twins_endpoint_eventhub", "azurerm_digital_twins_endpoint_servicebus"}
}

func (digitalTwinsEndpointsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDigitalTwinsEndpointsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.EndpointResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armdigitaltwins.EventGrid:
		return "azurerm_digital_twins_endpoint_eventgrid", nil
	case *armdigitaltwins.EventHub:
		return "azurerm_digital_twins_endpoint_eventhub", nil
	case *armdigitaltwins.ServiceBus:
		return "azurerm_digital_twins_endpoint_servicebus", nil
	default:
		return "", fmt.Errorf("unknown endpoint type %T", props)
	}
}
