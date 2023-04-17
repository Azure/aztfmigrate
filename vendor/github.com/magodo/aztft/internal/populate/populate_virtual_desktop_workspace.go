package populate

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateVirtualDesktopWorkspace(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDesktopVirtualizationWorkspacesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Workspace.Properties
	if props == nil {
		return nil, nil
	}
	var applicationGroupIds []string
	for _, id := range props.ApplicationGroupReferences {
		if id == nil {
			continue
		}
		applicationGroupIds = append(applicationGroupIds, *id)
	}

	var result []armid.ResourceId
	for _, applicationGroupId := range applicationGroupIds {
		applicationGroupAzureId, err := armid.ParseResourceId(applicationGroupId)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", applicationGroupId, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "applicationGroups")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(applicationGroupAzureId.String())))

		result = append(result, azureId)
	}
	return result, nil
}
