package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type virtualHubsResolver struct{}

func (virtualHubsResolver) ResourceTypes() []string {
	return []string{"azurerm_route_server", "azurerm_virtual_hub"}
}

func (virtualHubsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkVirtualHubsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.VirtualHub.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	vwan := props.VirtualWan

	if vwan == nil {
		return "azurerm_route_server", nil
	}
	return "azurerm_virtual_hub", nil
}
