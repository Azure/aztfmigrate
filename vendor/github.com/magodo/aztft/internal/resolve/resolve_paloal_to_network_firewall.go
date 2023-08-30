package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type paloalToNetworkFirewall struct{}

func (paloalToNetworkFirewall) ResourceTypes() []string {
	return []string{
		"azurerm_palo_alto_next_generation_firewall_virtual_network_panorama",
		"azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama",
		"azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack",
		"azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack",
	}
}

func (paloalToNetworkFirewall) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewPaloalToNetworkFirewallsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.FirewallResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	if props.NetworkProfile == nil {
		return "", fmt.Errorf("unexpected nil networkProfile in response")
	}
	if props.NetworkProfile.NetworkType == nil {
		return "", fmt.Errorf("unexpected nil networkProfile.networkType in response")
	}

	networkType := *props.NetworkProfile.NetworkType

	if props.IsPanoramaManaged != nil && *props.IsPanoramaManaged == armpanngfw.BooleanEnumTRUE {
		switch networkType {
		case armpanngfw.NetworkTypeVNET:
			return "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama", nil
		case armpanngfw.NetworkTypeVWAN:
			return "azurerm_palo_alto_next_generation_firewall_vhub_panorama", nil
		default:
			return "", fmt.Errorf("unknown network type: %s", networkType)
		}
	} else {
		switch networkType {
		case armpanngfw.NetworkTypeVNET:
			return "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack", nil
		case armpanngfw.NetworkTypeVWAN:
			return "azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack", nil
		default:
			return "", fmt.Errorf("unknown network type: %s", networkType)
		}
	}
}
