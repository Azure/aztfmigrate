package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type virtualMachinesResolver struct{}

func (virtualMachinesResolver) ResourceTypes() []string {
	return []string{"azurerm_linux_virtual_machine", "azurerm_windows_virtual_machine"}
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
	osProfile := props.OSProfile
	if osProfile == nil {
		return "", fmt.Errorf("unexpected nil OS profile in response")
	}

	switch {
	case osProfile.LinuxConfiguration != nil:
		return "azurerm_linux_virtual_machine", nil
	case osProfile.WindowsConfiguration != nil:
		return "azurerm_windows_virtual_machine", nil
	default:
		return "", fmt.Errorf("both windowsConfiguration and linuxConfiguration in OS profile is null")
	}
}
