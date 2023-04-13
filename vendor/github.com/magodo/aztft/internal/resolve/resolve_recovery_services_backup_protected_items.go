package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type recoveryServicesBackupProtectedItemsResolver struct{}

func (recoveryServicesBackupProtectedItemsResolver) ResourceTypes() []string {
	return []string{"azurerm_backup_protected_vm", "azurerm_backup_protected_file_share"}
}

func (recoveryServicesBackupProtectedItemsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewRecoveryservicesBackupProtectedItemsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), id.Names()[0], resourceGroupId.Name, id.Names()[1], id.Names()[2], id.Names()[3], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.ProtectedItemResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem:
		return "azurerm_backup_protected_vm", nil
	case *armrecoveryservicesbackup.AzureFileshareProtectedItem:
		return "azurerm_backup_protected_file_share", nil
	default:
		return "", fmt.Errorf("unknown protected item type: %T", props)
	}
}
