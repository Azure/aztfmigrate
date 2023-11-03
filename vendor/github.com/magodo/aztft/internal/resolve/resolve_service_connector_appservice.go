package resolve

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type serviceConnectorAppServiceResolver struct{}

func (serviceConnectorAppServiceResolver) ResourceTypes() []string {
	return []string{
		"azurerm_function_app_connection",
		"azurerm_app_service_connection",
	}
}

func (serviceConnectorAppServiceResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	rt, err := appServiceSitesResolver{}.Resolve(b, id.ParentScope())
	if err != nil {
		return "", err
	}

	switch rt {
	case "azurerm_logic_app_standard",
		"azurerm_linux_function_app",
		"azurerm_windows_function_app":
		return "azurerm_function_app_connection", nil
	case "azurerm_linux_web_app",
		"azurerm_windows_web_app":
		return "azurerm_app_service_connection", nil
	}
	return "", fmt.Errorf("unknown app service site resource type: %s", rt)
}
