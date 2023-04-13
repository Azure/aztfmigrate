package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type cdnProfilesResolver struct{}

func (cdnProfilesResolver) ResourceTypes() []string {
	return []string{"azurerm_cdn_frontdoor_profile", "azurerm_cdn_profile"}
}

func (cdnProfilesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewCdnProfilesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	sku := resp.Profile.SKU
	if sku == nil {
		return "", fmt.Errorf("unexpected nil properties.sku in response")
	}
	skuName := sku.Name
	if skuName == nil {
		return "", fmt.Errorf("unexpected nil properties.sku.name in response")
	}
	switch *skuName {
	case armcdn.SKUNamePremiumAzureFrontDoor,
		armcdn.SKUNameStandardAzureFrontDoor:
		return "azurerm_cdn_frontdoor_profile", nil
	case armcdn.SKUNameStandardAkamai,
		armcdn.SKUNameStandardChinaCdn,
		armcdn.SKUNameStandardVerizon,
		armcdn.SKUNameStandardMicrosoft,
		armcdn.SKUNamePremiumVerizon:
		return "azurerm_cdn_profile", nil
	default:
		return "", fmt.Errorf("unknown sku name %s", *skuName)
	}
}
