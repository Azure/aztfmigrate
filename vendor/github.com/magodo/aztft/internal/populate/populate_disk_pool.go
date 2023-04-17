package populate

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagepool/armstoragepool"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateDiskPool(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStoragePoolDiskPoolsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.DiskPool.Properties
	if props == nil {
		return nil, nil
	}

	managedDiskAssociation, err := diskPoolPopulateManagedDiskAssociation(id, props)
	if err != nil {
		return nil, fmt.Errorf("populating for managed disk associations: %v", err)
	}

	var result []armid.ResourceId
	result = append(result, managedDiskAssociation...)

	return result, nil
}

func diskPoolPopulateManagedDiskAssociation(id armid.ResourceId, props *armstoragepool.DiskPoolProperties) ([]armid.ResourceId, error) {
	var result []armid.ResourceId

	for _, disk := range props.Disks {
		if disk == nil {
			continue
		}
		if disk.ID == nil {
			continue
		}
		diskId, err := armid.ParseResourceId(*disk.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", *disk.ID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "disks")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(diskId.String())))

		result = append(result, azureId)
	}

	return result, nil
}
