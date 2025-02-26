package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateMssqlJob(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSqlJobsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], id.Names()[2], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Job.Properties
	if props == nil {
		return nil, nil
	}

	if props.Schedule == nil {
		return nil, nil
	}

	sid := id.Clone().(*armid.ScopedResourceId)
	sid.AttrTypes = append(sid.AttrTypes, "schedules")
	sid.AttrNames = append(sid.AttrTypes, "default")
	return []armid.ResourceId{sid}, nil
}
