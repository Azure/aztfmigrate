package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type cognitiveAccountConnectionResolver struct{}

func (cognitiveAccountConnectionResolver) ResourceTypes() []string {
	return []string{
		"azurerm_cognitive_account_connection_api_key",
		"azurerm_cognitive_account_connection_account_key",
		"azurerm_cognitive_account_connection_entra_id",
		"azurerm_cognitive_account_connection_account_managed_identity",
		"azurerm_cognitive_account_connection_custom_keys",
	}
}

func (cognitiveAccountConnectionResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewCognitiveServiceAccountConnectionsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}

	switch props := props.(type) {
	case *armcognitiveservices.APIKeyAuthConnectionProperties:
		return "azurerm_cognitive_account_connection_api_key", nil
	case *armcognitiveservices.AccountKeyAuthTypeConnectionProperties:
		return "azurerm_cognitive_account_connection_account_key", nil
	case *armcognitiveservices.OAuth2AuthTypeConnectionProperties:
		return "azurerm_cognitive_account_connection_entra_id", nil
	case *armcognitiveservices.ManagedIdentityAuthTypeConnectionProperties:
		return "azurerm_cognitive_account_connection_account_managed_identity", nil
	case *armcognitiveservices.CustomKeysConnectionProperties:
		return "azurerm_cognitive_account_connection_custom_keys", nil
	default:
		return "", fmt.Errorf("unknown connection type: %T", props)
	}
}
