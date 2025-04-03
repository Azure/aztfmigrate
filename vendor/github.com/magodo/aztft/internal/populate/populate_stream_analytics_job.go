package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateStreamAnalyticsJob(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStreamAnalyticsJobsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.StreamingJob.Properties
	if props == nil {
		return nil, nil
	}

	if props.JobStorageAccount == nil {
		return nil, nil
	}

	if props.JobStorageAccount.AccountName == nil {
		return nil, nil
	}

	sid := id.Clone().(*armid.ScopedResourceId)
	sid.AttrTypes = append(sid.AttrTypes, "storageAccounts")
	sid.AttrNames = append(sid.AttrTypes, *props.JobStorageAccount.AccountName)
	return []armid.ResourceId{sid}, nil
}
