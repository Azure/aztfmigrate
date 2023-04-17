package populate

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateDiskPoolIscsiTarget(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStoragePoolIscsiTargetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.IscsiTarget.Properties
	if props == nil {
		return nil, nil
	}

	var result []armid.ResourceId

	for _, lun := range props.Luns {
		if lun == nil {
			continue
		}
		if lun.ManagedDiskAzureResourceID == nil {
			continue
		}
		diskId, err := armid.ParseResourceId(*lun.ManagedDiskAzureResourceID)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", *lun.ManagedDiskAzureResourceID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "disks")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(diskId.String())))

		result = append(result, azureId)
	}

	return result, nil
}
