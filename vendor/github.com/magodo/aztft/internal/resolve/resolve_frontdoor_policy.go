package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type frontdoorPoliciesResolver struct{}

func (frontdoorPoliciesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_cdn_frontdoor_firewall_policy",
		"azurerm_frontdoor_firewall_policy",
	}
}

func (frontdoorPoliciesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewFrontdoorPoliciesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	sku := resp.WebApplicationFirewallPolicy.SKU
	if sku == nil {
		return "", fmt.Errorf("unexpected nil sku in response")
	}
	skuName := sku.Name
	if skuName == nil {
		return "", fmt.Errorf("unexpected nil sku name in response")
	}
	switch *skuName {
	case armfrontdoor.SKUNameClassicAzureFrontDoor:
		return "azurerm_frontdoor_firewall_policy", nil
	case armfrontdoor.SKUNameStandardAzureFrontDoor,
		armfrontdoor.SKUNamePremiumAzureFrontDoor:
		return "azurerm_cdn_frontdoor_firewall_policy", nil
	default:
		return "", fmt.Errorf("unknown frontdoor firewall policy sku name %s", *skuName)
	}
}
