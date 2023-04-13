package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateVirtualMachine(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewVirtualMachinesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.VirtualMachine.Properties
	if props == nil {
		return nil, nil
	}
	storageProfile := props.StorageProfile
	if storageProfile == nil {
		return nil, nil
	}

	var mdiskIds []string
	for _, disk := range storageProfile.DataDisks {
		if disk == nil {
			continue
		}
		mdisk := disk.ManagedDisk
		if mdisk == nil {
			continue
		}
		if mdisk.ID == nil {
			continue
		}
		mdiskIds = append(mdiskIds, *mdisk.ID)
	}

	var result []armid.ResourceId
	for _, mdiskId := range mdiskIds {
		mdiskAzureId, err := armid.ParseResourceId(mdiskId)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", mdiskId, err)
		}
		diskName := mdiskAzureId.Names()[0]

		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "dataDisks")
		azureId.AttrNames = append(azureId.AttrNames, diskName)

		result = append(result, azureId)
	}
	return result, nil
}
