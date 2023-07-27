package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type applicationInsightsWebTestsResolver struct{}

func (applicationInsightsWebTestsResolver) ResourceTypes() []string {
	return []string{"azurerm_application_insights_web_test", "azurerm_application_insights_standard_web_test"}
}

func (applicationInsightsWebTestsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewApplicationInsightsWebTestsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.WebTest.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	switch *kind {
	case armapplicationinsights.WebTestKindPing,
		armapplicationinsights.WebTestKindMultistep:
		return "azurerm_application_insights_web_test", nil
	default:
		// Actually, the kind is "standard", but the SDK is not updated to support it yet
		return "azurerm_application_insights_standard_web_test", nil
	}
}
