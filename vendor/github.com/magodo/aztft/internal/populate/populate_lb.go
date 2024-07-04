package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateLoadBalancer(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkLoadBalancersClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.LoadBalancer.Properties
	if props == nil {
		return nil, nil
	}

	var result []armid.ResourceId

	for _, rule := range props.LoadBalancingRules {
		if rule == nil || rule.ID == nil {
			continue
		}
		id, err := armid.ParseResourceId(*rule.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", *rule.ID, err)
		}
		result = append(result, id)
	}

	for _, probe := range props.Probes {
		if probe == nil || probe.ID == nil {
			continue
		}
		id, err := armid.ParseResourceId(*probe.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", *probe.ID, err)
		}
		result = append(result, id)
	}

	return result, nil
}
