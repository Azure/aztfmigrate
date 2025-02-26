package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type hybridkubernetesConnectedClusterResolver struct{}

func (hybridkubernetesConnectedClusterResolver) ResourceTypes() []string {
	return []string{"azurerm_arc_kubernetes_cluster", "azurerm_arc_kubernetes_provisioned_cluster"}
}

func (hybridkubernetesConnectedClusterResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewHybridKubernetesConnectedClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}

	// TODO: The Azure/azure-sdk-for-go uses the API version: 2021-10-01, which has no "Kind" defined.
	_ = resp

	// kind := resp.Kind
	// if kind == nil {
	// 	return "", fmt.Errorf("unexpected nil kind in response")
	// }
	return "azurerm_arc_kubernetes_cluster", nil
}
