package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type storageMoverEndpointsResolver struct{}

func (storageMoverEndpointsResolver) ResourceTypes() []string {
	return []string{"azurerm_storage_mover_source_endpoint", "azurerm_storage_mover_target_endpoint"}
}

func (storageMoverEndpointsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStorageMoverEndpointsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Endpoint.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armstoragemover.NfsMountEndpointProperties:
		return "azurerm_storage_mover_source_endpoint", nil
	case *armstoragemover.AzureStorageBlobContainerEndpointProperties:
		return "azurerm_storage_mover_target_endpoint", nil
	default:
		return "", fmt.Errorf("unknown storage mover endpoint type: %T", props)
	}
}
