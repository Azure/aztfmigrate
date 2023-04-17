package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type securityInsightsSecurityMLAnalyticsSettingsResolver struct{}

func (securityInsightsSecurityMLAnalyticsSettingsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_sentinel_alert_rule_anomaly_built_in",
		"azurerm_sentinel_alert_rule_anomaly_duplicate",
	}
}

func (securityInsightsSecurityMLAnalyticsSettingsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSecurityInsightsSecurityMLAnalyticsSettingsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.ParentScope().Names()[0], id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.SecurityMLAnalyticsSettingClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}

	switch model := model.(type) {
	case *armsecurityinsights.AnomalySecurityMLAnalyticsSettings:
		// TODO: figure out how to resolve azurerm_sentinel_alert_rule_anomaly_{built_in|duplicate}
		return "azurerm_sentinel_alert_rule_anomaly_built_in", nil
	default:
		return "", fmt.Errorf("unknown security ML analytics setting type: %T", model)
	}
}
