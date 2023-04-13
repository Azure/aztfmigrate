package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicessiterecovery"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type siteRecoveryReplicationPoliciesResolver struct{}

func (siteRecoveryReplicationPoliciesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_site_recovery_hyperv_replication_policy",
		"azurerm_site_recovery_replication_policy",
		"azurerm_site_recovery_vmware_replication_policy",
	}
}

func (siteRecoveryReplicationPoliciesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSiteRecoveryReplicationPoliciesClient(resourceGroupId.SubscriptionId, resourceGroupId.Name, id.Names()[0])
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	prop := resp.Policy.Properties
	if prop == nil {
		return "", fmt.Errorf("unexpected nil prop in response")
	}
	switch prop.ProviderSpecificDetails.(type) {
	case *armrecoveryservicessiterecovery.HyperVReplicaAzurePolicyDetails:
		return "azurerm_site_recovery_hyperv_replication_policy", nil
	case *armrecoveryservicessiterecovery.A2APolicyDetails:
		return "azurerm_site_recovery_replication_policy", nil
	case *armrecoveryservicessiterecovery.VmwareCbtPolicyDetails:
		return "azurerm_site_recovery_vmware_replication_policy", nil
	default:
		return "", fmt.Errorf("unknown site recovery replication policy detail type: %T", prop.ProviderSpecificDetails)
	}
}
