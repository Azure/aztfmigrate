package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type virtualMachinesResolver struct{}

func (virtualMachinesResolver) ResourceTypes() []string {
	return []string{"azurerm_linux_virtual_machine", "azurerm_windows_virtual_machine", "azurerm_virtual_machine"}
}

func (virtualMachinesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
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

	if props.OSProfile == nil {
		// Per: https://github.com/hashicorp/terraform-provider-azurerm/blob/c8d1a23b143360eaf5ee371840cc4d5ee286eddc/internal/services/compute/virtual_machine_import.go#L45-L48
		return "azurerm_virtual_machine", nil
	}

	storageProfile := props.StorageProfile
	if storageProfile == nil {
		return "", fmt.Errorf("unexpected nil storage profile in response")
	}

	osDisk := storageProfile.OSDisk
	if osDisk == nil {
		return "", fmt.Errorf("unexpected nil OS Disk in storage profile")
	}

	if osDisk.Vhd != nil {
		// Per: https://github.com/hashicorp/terraform-provider-azurerm/blob/c8d1a23b143360eaf5ee371840cc4d5ee286eddc/internal/services/compute/virtual_machine_import.go#L36-L38
		return "azurerm_virtual_machine", nil
	}

	osType := osDisk.OSType
	if osType == nil {
		return "", fmt.Errorf("unexpected nil OS Type in OS Disk")
	}

	switch *osType {
	case armcompute.OperatingSystemTypesLinux:
		return "azurerm_linux_virtual_machine", nil
	case armcompute.OperatingSystemTypesWindows:
		return "azurerm_windows_virtual_machine", nil
	default:
		return "", fmt.Errorf("Unknown OS Type: %s", *osType)
	}
}
