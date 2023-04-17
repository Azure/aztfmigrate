package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type dataFactoryIntegrationRuntimesResolver struct{}

func (dataFactoryIntegrationRuntimesResolver) ResourceTypes() []string {
	return []string{"azurerm_data_factory_integration_runtime_azure_ssis", "azurerm_data_factory_integration_runtime_azure", "azurerm_data_factory_integration_runtime_self_hosted"}
}

func (dataFactoryIntegrationRuntimesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataFactoryIntegrationRuntimesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.IntegrationRuntimeResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props := props.(type) {
	case *armdatafactory.ManagedIntegrationRuntime:
		tp := props.TypeProperties
		if tp == nil {
			return "", fmt.Errorf("unexpected nil properties.typeProperties in response")
		}
		if tp.SsisProperties != nil {
			return "azurerm_data_factory_integration_runtime_azure_ssis", nil
		}
		return "azurerm_data_factory_integration_runtime_azure", nil
	case *armdatafactory.SelfHostedIntegrationRuntime:
		return "azurerm_data_factory_integration_runtime_self_hosted", nil
	default:
		return "", fmt.Errorf("unknown integration runtime type %T", props)
	}
}
