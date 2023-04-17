package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type siteRecoveryReplicationFabricsResolver struct{}

func (siteRecoveryReplicationFabricsResolver) ResourceTypes() []string {
	return []string{"azurerm_site_recovery_services_vault_hyperv_site", "azurerm_site_recovery_fabric"}
}

func (siteRecoveryReplicationFabricsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSiteRecoveryReplicationFabricsClient(resourceGroupId.SubscriptionId, resourceGroupId.Name, id.Names()[0])
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	prop := resp.Fabric.Properties
	if prop == nil {
		return "", fmt.Errorf("unexpected nil prop in response")
	}
	switch prop.CustomDetails.(type) {
	case *armrecoveryservicessiterecovery.AzureFabricSpecificDetails:
		return "azurerm_site_recovery_fabric", nil
	case *armrecoveryservicessiterecovery.HyperVSiteDetails:
		return "azurerm_site_recovery_services_vault_hyperv_site", nil
	default:
		return "", fmt.Errorf("unknown site recovery replication fabric detail type: %T", prop.CustomDetails)
	}
}
