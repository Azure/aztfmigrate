package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type appServiceSitesResolver struct{}

func (appServiceSitesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_logic_app_standard",
		"azurerm_linux_function_app",
		"azurerm_windows_function_app",
		"azurerm_linux_web_app",
		"azurerm_windows_web_app",
	}
}

func (appServiceSitesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAppServiceWebAppsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	// The value of kind for different resource are listed below:
	//
	// azurerm_logic_app_standard	: functionapp,workflowapp or functionapp,linux,container,workflowapp
	// azurerm_linux_function_app	: functionapp,linux
	// azurerm_windows_function_app	: functionapp
	// azurerm_linux_web_app		: app,linux
	// azurerm_windows_web_app		: app,container,windows

	kinds := strings.Split(*kind, ",")
	m := map[string]bool{}
	for _, k := range kinds {
		m[strings.ToLower(k)] = true
	}

	if m["workflowapp"] && m["functionapp"] {
		return "azurerm_logic_app_standard", nil
	}

	if m["functionapp"] {
		if m["linux"] {
			return "azurerm_linux_function_app", nil
		}
		return "azurerm_windows_function_app", nil
	}

	if m["app"] {
		if m["linux"] {
			return "azurerm_linux_web_app", nil
		}
		return "azurerm_windows_web_app", nil
	}

	return "", fmt.Errorf("unknown kind: %s", *kind)
}
