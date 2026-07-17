package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v9"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type kubernetesClusterResolver struct{}

func (kubernetesClusterResolver) ResourceTypes() []string {
	return []string{
		"azurerm_kubernetes_automatic_cluster",
		"azurerm_kubernetes_cluster",
	}
}

func (kubernetesClusterResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewKubernetesClusterClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	sku := resp.ManagedCluster.SKU
	if sku == nil {
		return "", fmt.Errorf("unexpected nil sku in response")
	}
	if sku.Name == nil {
		return "", fmt.Errorf("unexpected nil sku.Name in response")
	}

	switch strings.ToLower(string(*sku.Name)) {
	case strings.ToLower(string(armcontainerservice.ManagedClusterSKUNameAutomatic)):
		return "azurerm_kubernetes_automatic_cluster", nil
	case strings.ToLower(string(armcontainerservice.ManagedClusterSKUNameBase)):
		return "azurerm_kubernetes_cluster", nil
	default:
		return "", fmt.Errorf("unknown sku name: %s", *sku.Name)
	}
}
