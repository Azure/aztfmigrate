package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateContainerAppEnv(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewContainerAppEnvironmentsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.ManagedEnvironment.Properties
	if props == nil {
		return nil, nil
	}

	var result []armid.ResourceId

	if cfg := props.CustomDomainConfiguration; cfg != nil {
		if cfg.DNSSuffix != nil {
			cid := id.Clone().(*armid.ScopedResourceId)
			cid.AttrTypes = append(cid.AttrTypes, "customDomains")
			cid.AttrNames = append(cid.AttrTypes, "default")
			result = append(result, cid)
		}
	}

	return result, nil
}
