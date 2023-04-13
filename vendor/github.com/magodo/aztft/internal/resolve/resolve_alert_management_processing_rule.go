package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertsmanagement/armalertsmanagement"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type alertsManagementProcessingRulesResolver struct{}

func (alertsManagementProcessingRulesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_monitor_alert_processing_rule_suppression",
		"azurerm_monitor_alert_processing_rule_action_group",
	}
}

func (alertsManagementProcessingRulesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAlertsManagementProcessingRulesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.GetByName(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}

	// Below logic is not ideal, as the API version used here is not the same as what was used in the provider, and the two APIs are not compatible.
	// The reason why we don't use the same API as provider is the version of the alertsmanagement that used the same version depends on a non-compatible version of azcore.
	// The check below is derived from my local test by using the provider to provision both resources and use the SDK API version to GET.
	props := resp.AlertProcessingRule.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	actions := props.Actions
	if len(actions) != 1 {
		return "", fmt.Errorf("expect 1 action, got=%d", len(actions))
	}

	switch actions[0].(type) {
	case *armalertsmanagement.AddActionGroups:
		return "azurerm_monitor_alert_processing_rule_action_group", nil
	case *armalertsmanagement.RemoveAllActionGroups:
		return "azurerm_monitor_alert_processing_rule_suppression", nil
	default:
		return "", fmt.Errorf("unknown action type: %T", actions[0])
	}
}
