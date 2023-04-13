package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type securityInsightsAlertRulesResolver struct{}

func (securityInsightsAlertRulesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_sentinel_alert_rule_nrt",
		"azurerm_sentinel_alert_rule_fusion",
		"azurerm_sentinel_alert_rule_machine_learning_behavior_analytics",
		"azurerm_sentinel_alert_rule_ms_security_incident",
		"azurerm_sentinel_alert_rule_scheduled",
		"azurerm_sentinel_alert_rule_threat_intelligence",
	}
}

func (securityInsightsAlertRulesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewSecurityInsightsAlertRulesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.ParentScope().Names()[0], id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.AlertRuleClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}

	switch model.(type) {
	case *armsecurityinsights.NrtAlertRule:
		return "azurerm_sentinel_alert_rule_nrt", nil
	case *armsecurityinsights.FusionAlertRule:
		return "azurerm_sentinel_alert_rule_fusion", nil
	case *armsecurityinsights.MLBehaviorAnalyticsAlertRule:
		return "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics", nil
	case *armsecurityinsights.MicrosoftSecurityIncidentCreationAlertRule:
		return "azurerm_sentinel_alert_rule_ms_security_incident", nil
	case *armsecurityinsights.ScheduledAlertRule:
		return "azurerm_sentinel_alert_rule_scheduled", nil
	case *armsecurityinsights.ThreatIntelligenceAlertRule:
		return "azurerm_sentinel_alert_rule_threat_intelligence", nil
	default:
		return "", fmt.Errorf("unknown alert rule type: %T", model)
	}
}
