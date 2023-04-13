package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type appServiceSiteSlotsResolver struct{}

func (appServiceSiteSlotsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_linux_function_app_slot",
		"azurerm_windows_function_app_slot",
		"azurerm_linux_web_app_slot",
		"azurerm_windows_web_app_slot",
	}
}

func (appServiceSiteSlotsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAppServiceWebAppsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.GetSlot(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	// The value of kind for different resource are listed below:
	//
	// azurerm_windows_function_app_slot: functionapp
	// azurerm_linux_function_app_slot	: functionapp,linux
	// azurerm_windows_web_app_slot		: app
	// azurerm_linux_web_app_slot		: app,linux

	kinds := strings.Split(*kind, ",")
	m := map[string]bool{}
	for _, k := range kinds {
		m[strings.ToLower(k)] = true
	}

	if m["functionapp"] {
		if m["linux"] {
			return "azurerm_linux_function_app_slot", nil
		}
		return "azurerm_windows_function_app_slot", nil
	}

	if m["app"] {
		if m["linux"] {
			return "azurerm_linux_web_app_slot", nil
		}
		return "azurerm_windows_web_app_slot", nil
	}

	return "", fmt.Errorf("unknown kind: %s", *kind)
}
