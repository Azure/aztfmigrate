package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type siteRecoveryReplicationProtectionContainerMappingResolver struct{}

func (siteRecoveryReplicationProtectionContainerMappingResolver) ResourceTypes() []string {
	return []string{
		"azurerm_site_recovery_hyperv_replication_policy_association",
		"azurerm_site_recovery_protection_container_mapping",

		// There is no SDk defined proverSpecificDetails type for this type, so just ingoring it for now
		"azurerm_site_recovery_vmware_replication_policy_association",
	}
}

func (siteRecoveryReplicationProtectionContainerMappingResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSiteRecoveryReplicationProtectionContainerMappingsClient(resourceGroupId.SubscriptionId, resourceGroupId.Name, id.Names()[0])
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[1], id.Names()[2], id.Names()[3], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	prop := resp.ProtectionContainerMapping.Properties
	if prop == nil {
		return "", fmt.Errorf("unexpected nil prop in response")
	}
	switch prop.ProviderSpecificDetails.(type) {
	case *armrecoveryservicessiterecovery.ProtectionContainerMappingProviderSpecificDetails:
		return "azurerm_site_recovery_hyperv_replication_policy_association", nil
	case *armrecoveryservicessiterecovery.A2AProtectionContainerMappingDetails:
		return "azurerm_site_recovery_protection_container_mapping", nil
	default:
		return "", fmt.Errorf("unsupported site recovery replication container mapping detail type: %T", prop.ProviderSpecificDetails)
	}
}
