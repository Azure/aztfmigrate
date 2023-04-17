package tfid

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageTableEntity(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
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
	tableEndpoint := endpoints.Table
	if tableEndpoint == nil {
		return "", fmt.Errorf("unexpected nil properties.primaryEndpoints.table in response")
	}
	baseUri := strings.TrimSuffix(*tableEndpoint, "/")
	return fmt.Sprintf("%s/%s(PartitionKey='%s',RowKey='%s')", baseUri, id.Names()[2], id.Names()[3], id.Names()[4]), nil
}
