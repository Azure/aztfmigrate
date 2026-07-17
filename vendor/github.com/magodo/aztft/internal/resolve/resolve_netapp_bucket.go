package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type netappBucketResolver struct{}

func (netappBucketResolver) ResourceTypes() []string {
	return []string{
		"azurerm_netapp_volume_bucket",
		"azurerm_netapp_volume_bucket_with_server",
	}
}

func (netappBucketResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetAppVolumeBucketClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], id.Names()[2], id.Names()[3], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Bucket.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}
	server := props.Server
	if server != nil {
		return "azurerm_netapp_volume_bucket_with_server", nil
	} else {
		return "azurerm_netapp_volume_bucket", nil
	}
}
