package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type dataFactoryDataFlowsResolver struct{}

func (dataFactoryDataFlowsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_data_factory_data_flow",
		"azurerm_data_factory_flowlet_data_flow",
	}
}

func (dataFactoryDataFlowsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataFactoryDataFlowsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.DataFlowResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armdatafactory.Flowlet:
		return "azurerm_data_factory_flowlet_data_flow", nil
	case *armdatafactory.MappingDataFlow:
		return "azurerm_data_factory_data_flow", nil
	default:
		return "", fmt.Errorf("unknown data flow type %T", props)
	}
}
