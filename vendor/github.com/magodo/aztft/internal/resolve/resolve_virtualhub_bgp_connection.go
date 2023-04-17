package resolve

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type virtualHubBgpConnectionsResolver struct{}

func (virtualHubBgpConnectionsResolver) ResourceTypes() []string {
	return []string{"azurerm_route_server_bgp_connection", "azurerm_virtual_hub_bgp_connection"}
}

func (virtualHubBgpConnectionsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	// The two connections are actually the same, disambiguate them via their parent resource
	t, err := virtualHubsResolver{}.Resolve(b, id)
	if err != nil {
		return "", err
	}
	switch t {
	case "azurerm_route_server":
		return "azurerm_route_server_bgp_connection", nil
	case "azurerm_virtual_hub":
		return "azurerm_virtual_hub_bgp_connection", nil
	}
	return "", fmt.Errorf("unknown parent resource type: %s", t)
}
