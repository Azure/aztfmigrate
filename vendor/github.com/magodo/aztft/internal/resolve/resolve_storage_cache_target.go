package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type storageCacheTargetsResolver struct{}

func (storageCacheTargetsResolver) ResourceTypes() []string {
	return []string{"azurerm_hpc_cache_blob_nfs_target", "azurerm_hpc_cache_blob_target", "azurerm_hpc_cache_nfs_target"}
}

func (storageCacheTargetsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStorageCacheTargetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.StorageTarget.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	tt := props.TargetType
	if tt == nil {
		return "", fmt.Errorf("unexpected nil targetType in response")
	}

	switch *tt {
	case armstoragecache.StorageTargetTypeBlobNfs:
		return "azurerm_hpc_cache_blob_nfs_target", nil
	case armstoragecache.StorageTargetTypeClfs:
		return "azurerm_hpc_cache_blob_target", nil
	case armstoragecache.StorageTargetTypeNfs3:
		return "azurerm_hpc_cache_nfs_target", nil
	default:
		return "", fmt.Errorf("unknown resource type: %s", *tt)
	}
}
