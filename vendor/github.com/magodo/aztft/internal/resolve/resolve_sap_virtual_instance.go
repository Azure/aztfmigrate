package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloads/armworkloads"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type sapVirtualInstancesResolver struct{}

func (sapVirtualInstancesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_workloads_sap_single_node_virtual_instance",
		"azurerm_workloads_sap_three_tier_virtual_instance",
		"azurerm_workloads_sap_discovery_virtual_instance",
	}
}

func (sapVirtualInstancesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewWorkloadSAPVirtualInstanceClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.SAPVirtualInstance.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	configRaw := props.Configuration
	if configRaw == nil {
		return "", fmt.Errorf("unexpected nil Configuration in response")
	}
	switch config := configRaw.(type) {
	case *armworkloads.DiscoveryConfiguration:
		return "azurerm_workloads_sap_discovery_virtual_instance", nil

	case *armworkloads.DeploymentWithOSConfiguration:
		infraConfigRaw := config.InfrastructureConfiguration
		if infraConfigRaw == nil {
			return "", fmt.Errorf("unexpected nil Configuration.InfrastructureConfiguration in response")
		}
		switch infraConfigRaw.(type) {
		case *armworkloads.SingleServerConfiguration:
			return "azurerm_workloads_sap_single_node_virtual_instance", nil
		case *armworkloads.ThreeTierConfiguration:
			return "azurerm_workloads_sap_three_tier_virtual_instance", nil
		default:
			return "", fmt.Errorf("unexpected Configuration.InfrastructureConfiguration type in response, got=%T", infraConfigRaw)
		}
	default:
		return "", fmt.Errorf("unexpected Configuration type in response, got=%T", configRaw)
	}
}
