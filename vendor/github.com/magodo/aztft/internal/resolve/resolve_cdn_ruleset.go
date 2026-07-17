package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type cdnRuleSetsResolver struct{}

func (cdnRuleSetsResolver) ResourceTypes() []string {
	return []string{"azurerm_cdn_frontdoor_rule_set", "azurerm_cdn_frontdoor_batch_rule_set"}
}

func (cdnRuleSetsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewCdnRuleSetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}

	// TODO: There is no props.BatchMode in the SDK (even in v3).
	return "azurerm_cdn_frontdoor_rule_set", nil
}
