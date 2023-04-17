package resolve

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type appServiceSiteHybridConnectionsResolver struct{}

func (appServiceSiteHybridConnectionsResolver) ResourceTypes() []string {
	return []string{"azurerm_web_app_hybrid_connection", "azurerm_function_app_hybrid_connection"}
}

func (appServiceSiteHybridConnectionsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	// Resolve the resource type by resolving its parent resource, i.e. the sites.
	rt, err := appServiceSitesResolver{}.Resolve(b, id.Parent().Parent())
	if err != nil {
		return "", err
	}

	switch rt {
	case "azurerm_windows_web_app", "azurerm_linux_web_app":
		return "azurerm_web_app_hybrid_connection", nil
	case "azurerm_windows_function_app", "azurerm_linux_function_app":
		return "azurerm_function_app_hybrid_connection", nil
	default:
		return "", fmt.Errorf("unknown parent resource type: %s", rt)
	}
}
