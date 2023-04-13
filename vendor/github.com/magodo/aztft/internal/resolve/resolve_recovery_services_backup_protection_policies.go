package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type recoveryServicesBackupProtectionPoliciesResolver struct{}

func (recoveryServicesBackupProtectionPoliciesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_backup_policy_vm",
		"azurerm_backup_policy_vm_workload",
		"azurerm_backup_policy_file_share",
	}
}

func (recoveryServicesBackupProtectionPoliciesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewRecoveryServicesBackupProtectionPoliciesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[0], resourceGroupId.Name, id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.ProtectionPolicyResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy:
		return "azurerm_backup_policy_vm", nil
	case *armrecoveryservicesbackup.AzureFileShareProtectionPolicy:
		return "azurerm_backup_policy_file_share", nil
	case *armrecoveryservicesbackup.AzureVMWorkloadProtectionPolicy:
		return "azurerm_backup_policy_vm_workload", nil
	default:
		return "", fmt.Errorf("unknown policy type: %T", props)
	}
}
