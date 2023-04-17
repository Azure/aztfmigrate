package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type virtualMachineScaleSetsResolver struct{}

func (virtualMachineScaleSetsResolver) ResourceTypes() []string {
	return []string{"azurerm_orchestrated_virtual_machine_scale_set", "azurerm_linux_virtual_machine_scale_set", "azurerm_windows_virtual_machine_scale_set"}
}

func (virtualMachineScaleSetsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewVirtualMachineScaleSetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.VirtualMachineScaleSet.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	// If the VMSS is created with orchestration mode "Uniform" (i.e. either linux/windows vmss), the orchestrationMode is not returned in the GET response body.
	if orchMode := props.OrchestrationMode; orchMode != nil && *orchMode == armcompute.OrchestrationModeFlexible {
		return "azurerm_orchestrated_virtual_machine_scale_set", nil
	}

	profile := props.VirtualMachineProfile
	if profile == nil {
		return "", fmt.Errorf("unexpected nil virtualMachineProfile in response")
	}
	osProfile := profile.OSProfile
	if osProfile == nil {
		return "", fmt.Errorf("unexpected nil virtualMachineProfile.osProfile in response")
	}
	switch {
	case osProfile.LinuxConfiguration != nil:
		return "azurerm_linux_virtual_machine_scale_set", nil
	case osProfile.WindowsConfiguration != nil:
		return "azurerm_windows_virtual_machine_scale_set", nil
	default:
		return "", fmt.Errorf("both windowsConfiguration and linuxConfiguration in OS profile is null")
	}
}
