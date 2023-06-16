package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type siteRecoveryReplicationNetworkMappingResolver struct{}

func (siteRecoveryReplicationNetworkMappingResolver) ResourceTypes() []string {
	return []string{"azurerm_site_recovery_hyperv_network_mapping", "azurerm_site_recovery_network_mapping"}
}

func (siteRecoveryReplicationNetworkMappingResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSiteRecoveryReplicationNetworkMappingsClient(resourceGroupId.SubscriptionId, resourceGroupId.Name, id.Names()[0])
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[1], id.Names()[2], id.Names()[3], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	prop := resp.NetworkMapping.Properties
	if prop == nil {
		return "", fmt.Errorf("unexpected nil prop in response")
	}
	switch prop.FabricSpecificSettings.(type) {
	case *armrecoveryservicessiterecovery.AzureToAzureNetworkMappingSettings:
		return "azurerm_site_recovery_network_mapping", nil
	case *armrecoveryservicessiterecovery.VmmToAzureNetworkMappingSettings:
		return "azurerm_site_recovery_hyperv_network_mapping", nil
	default:
		return "", fmt.Errorf("unsupported site recovery replication network mapping detail type: %T", prop.FabricSpecificSettings)
	}
}
