package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v7"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type dataFactoryCredentialsResolver struct{}

func (dataFactoryCredentialsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_data_factory_credential_user_managed_identity",
		"azurerm_data_factory_credential_service_principal",
	}
}

func (dataFactoryCredentialsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataFactoryCredentialsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.CredentialResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armdatafactory.ManagedIdentityCredential:
		return "azurerm_data_factory_credential_user_managed_identity", nil
	case *armdatafactory.ServicePrincipalCredential:
		return "azurerm_data_factory_credential_service_principal", nil
	default:
		return "", fmt.Errorf("unknown data flow type %T", props)
	}
}
