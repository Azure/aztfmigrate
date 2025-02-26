package populate

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

// populateFunc populates the hypothetic azure resource ids that represent the property like resources of the specified resource.
type populateFunc func(*client.ClientBuilder, armid.ResourceId) ([]armid.ResourceId, error)

var populaters = map[string]populateFunc{
	"azurerm_linux_virtual_machine":     populateVirtualMachine,
	"azurerm_windows_virtual_machine":   populateVirtualMachine,
	"azurerm_network_interface":         populateNetworkInterface,
	"azurerm_virtual_desktop_workspace": populateVirtualDesktopWorkspace,
	"azurerm_nat_gateway":               populateNatGateway,
	"azurerm_subnet":                    populateSubnet,
	"azurerm_logic_app_workflow":        populateLogicAppWorkflow,
	"azurerm_iothub":                    populateIotHub,
	"azurerm_netapp_account":            populateNetAppAccount,
	"azurerm_lb":                        populateLoadBalancer,
	"azurerm_container_app_environment": populateContainerAppEnv,
	"azurerm_mssql_job":                 populateMssqlJob,
}

func NeedsAPI(rt string) bool {
	_, ok := populaters[rt]
	return ok
}

func Populate(id armid.ResourceId, rt string, cred azcore.TokenCredential, clientOpt arm.ClientOptions) ([]armid.ResourceId, error) {
	populater, ok := populaters[rt]
	if !ok {
		return nil, nil
	}

	b := &client.ClientBuilder{
		Cred:      cred,
		ClientOpt: clientOpt,
	}

	return populater(b, id)
}
