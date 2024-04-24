package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildApiManagementApi(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewApiManagementApiClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	rev := props.APIRevision
	if rev == nil {
		return "", fmt.Errorf("unexpected nil properties.APIRevision in response")
	}
	if err := id.Normalize(spec); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s;rev=%s", id.String(), *rev), nil
}
