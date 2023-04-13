package tfid

import (
	"context"
	"fmt"
	"net/url"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageDfs(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStorageAccountsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.GetProperties(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Account.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	endpoints := props.PrimaryEndpoints
	if endpoints == nil {
		return "", fmt.Errorf("unexpected nil properties.primaryEndpoints in response")
	}
	dfsEndpoint := endpoints.Dfs
	if dfsEndpoint == nil {
		return "", fmt.Errorf("unexpected nil properties.primaryEndpoints.dfs in response")
	}
	uri, err := url.Parse(*dfsEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %v", *dfsEndpoint, err)
	}
	uri = uri.JoinPath(id.Names()[1])
	return uri.String(), nil
}
