package tfid

import (
	"context"
	"fmt"
	"net/url"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageShare(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
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
	fileEndpoint := endpoints.File
	if fileEndpoint == nil {
		return "", fmt.Errorf("unexpected nil properties.primaryEndpoints.file in response")
	}
	uri, err := url.Parse(*fileEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %v", *fileEndpoint, err)
	}
	uri = uri.JoinPath(id.Names()[2])
	return uri.String(), nil
}
