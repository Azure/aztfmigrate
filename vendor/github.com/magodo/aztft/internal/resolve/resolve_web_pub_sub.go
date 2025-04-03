package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type webPubSubResolver struct{}

func (webPubSubResolver) ResourceTypes() []string {
	return []string{"azurerm_web_pubsub_socketio", "azurerm_web_pubsub"}
}

func (webPubSubResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewWebPubSubsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.ResourceInfo.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}

	switch *kind {
	case armwebpubsub.ServiceKindSocketIO:
		return "azurerm_web_pubsub_socketio", nil
	default:
		return "azurerm_web_pubsub", nil
	}
}
