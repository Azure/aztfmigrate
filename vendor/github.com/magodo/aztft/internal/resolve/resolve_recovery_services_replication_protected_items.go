package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type recoveryServicesReplicationProtectedItemsResolver struct{}

func (recoveryServicesReplicationProtectedItemsResolver) ResourceTypes() []string {
	return []string{"azurerm_site_recovery_replicated_vm", "azurerm_site_recovery_vmware_replicated_vm"}
}

func (recoveryServicesReplicationProtectedItemsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSiteRecoveryReplicationProtectedItemsClient(resourceGroupId.SubscriptionId, resourceGroupId.Name, id.Names()[0])
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[1], id.Names()[2], id.Names()[3], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	details := props.ProviderSpecificDetails
	if details == nil {
		return "", fmt.Errorf("unexpected nil property.providerSpecificDetails in response")
	}
	settings := details.GetReplicationProviderSpecificSettings()
	if settings == nil {
		return "", fmt.Errorf("unexpected nil property.providerSpecificDetails.settings() in response")
	}
	typ := settings.InstanceType
	if typ == nil {
		return "", fmt.Errorf("unexpected nil property.providerSpecificDetails.instanceType in response")
	}

	switch *typ {
	case "A2ACrossClusterMigration":
		return "azurerm_site_recovery_replicated_vm", nil
	case "InMageRcm":
		return "azurerm_site_recovery_vmware_replicated_vm", nil
	default:
		return "", fmt.Errorf("unknown replication protected items type: %s", *typ)
	}
}
