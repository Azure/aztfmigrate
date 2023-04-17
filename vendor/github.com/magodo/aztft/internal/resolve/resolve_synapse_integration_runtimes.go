package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type synapseIntegrationRuntimesResolver struct{}

func (synapseIntegrationRuntimesResolver) ResourceTypes() []string {
	return []string{"azurerm_synapse_integration_runtime_azure", "azurerm_synapse_integration_runtime_self_hosted"}
}

func (synapseIntegrationRuntimesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSynapseIntegrationRuntimesClient(resourceGroupId.SubscriptionId)
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
	switch props.(type) {
	case *armsynapse.ManagedIntegrationRuntime:
		return "azurerm_synapse_integration_runtime_azure", nil
	case *armsynapse.SelfHostedIntegrationRuntime:
		return "azurerm_synapse_integration_runtime_self_hosted", nil
	default:
		return "", fmt.Errorf("unknown integration runtime type: %T", props)
	}
}
