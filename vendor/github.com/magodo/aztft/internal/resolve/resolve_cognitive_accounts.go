package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type cognitiveAccountsResolver struct{}

func (cognitiveAccountsResolver) ResourceTypes() []string {
	return []string{"azurerm_cognitive_account", "azurerm_ai_services"}
}

func (cognitiveAccountsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewCognitiveServiceAccountsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.Account.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}

	if strings.EqualFold(*kind, "AIServices") {
		return "azurerm_ai_services", nil
	}
	return "azurerm_cognitive_account", nil
}
