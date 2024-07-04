package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type virutalMachineDataDiskResolver struct{}

func (virutalMachineDataDiskResolver) ResourceTypes() []string {
	return []string{
		"azurerm_virtual_machine_implicit_data_disk_from_source",
		"azurerm_virtual_machine_data_disk_attachment",
	}
}

func (virutalMachineDataDiskResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewVirtualMachinesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}

	props := resp.VirtualMachine.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	profile := props.StorageProfile
	if profile == nil {
		return "", fmt.Errorf("unexpected nil storageProfile")
	}

	diskName := id.Names()[1]
	for _, disk := range profile.DataDisks {
		if disk == nil {
			continue
		}
		if disk.Name == nil {
			continue
		}
		if *disk.Name != diskName {
			continue
		}
		createOpt := disk.CreateOption
		if createOpt == nil {
			return "", fmt.Errorf("unexpected nil storageProfile.dataDisks.*.createOption")
		}
		switch *createOpt {
		case armcompute.DiskCreateOptionTypesEmpty, armcompute.DiskCreateOptionTypesAttach:
			return "azurerm_virtual_machine_data_disk_attachment", nil
		case armcompute.DiskCreateOptionTypesCopy:
			return "azurerm_virtual_machine_implicit_data_disk_from_source", nil
		default:
			return "", fmt.Errorf("unexpected storageProfile.dataDisks.*.createOption: %v", *createOpt)
		}
	}
	return "", fmt.Errorf("data disk named %q not found", diskName)
}
