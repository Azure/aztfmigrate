package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type dataProtectionBackupPoliciesResolver struct{}

func (dataProtectionBackupPoliciesResolver) ResourceTypes() []string {
	return []string{"azurerm_data_protection_backup_policy_postgresql", "azurerm_data_protection_backup_policy_disk", "azurerm_data_protection_backup_policy_blob_storage"}
}

func (dataProtectionBackupPoliciesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataProtectionBackupPoliciesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.BaseBackupPolicyResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	policy, ok := props.(*armdataprotection.BackupPolicy)
	if !ok {
		return "", fmt.Errorf("unknown type of the property: %T", props)
	}
	if len(policy.DatasourceTypes) != 1 {
		return "", fmt.Errorf("provider only support backup policy that has exactly one datasourceType specified, got=%d", len(policy.DatasourceTypes))
	}
	pdt := policy.DatasourceTypes[0]
	if pdt == nil {
		return "", fmt.Errorf("unexpected nil datasource type")
	}
	switch strings.ToUpper(*pdt) {
	case "MICROSOFT.DBFORPOSTGRESQL/SERVERS/DATABASES":
		return "azurerm_data_protection_backup_policy_postgresql", nil
	case "MICROSOFT.COMPUTE/DISKS":
		return "azurerm_data_protection_backup_policy_disk", nil
	case "MICROSOFT.STORAGE/STORAGEACCOUNTS/BLOBSERVICES":
		return "azurerm_data_protection_backup_policy_blob_storage", nil
	default:
		return "", fmt.Errorf("unknown data source type: %s", *pdt)
	}
}
