package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateRedisEnterpriseDatabase(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewRedisEnterpriseDatabaseClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], "default", nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return nil, nil
	}

	if props.GeoReplication == nil || props.GeoReplication.LinkedDatabases == nil || len(props.GeoReplication.LinkedDatabases) == 0 {
		return nil, nil
	}

	eid := id.Clone().(*armid.ScopedResourceId)
	eid.AttrTypes = append(eid.AttrTypes, "replications")
	eid.AttrNames = append(eid.AttrNames, "default")
	return []armid.ResourceId{eid}, nil
}
