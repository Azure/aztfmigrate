package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type costmanagementScheduleActionsResolver struct{}

func (costmanagementScheduleActionsResolver) ResourceTypes() []string {
	return []string{"azurerm_cost_anomaly_alert", "azurerm_cost_management_scheduled_action"}
}

func (costmanagementScheduleActionsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	client, err := b.NewCostManagementScheduledActionsClient()
	if err != nil {
		return "", err
	}
	resp, err := client.GetByScope(context.Background(), strings.TrimPrefix(id.ParentScope().String(), "/"), id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.ScheduledAction.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	switch *kind {
	case armcostmanagement.ScheduledActionKindEmail:
		return "azurerm_cost_management_scheduled_action", nil
	case armcostmanagement.ScheduledActionKindInsightAlert:
		return "azurerm_cost_anomaly_alert", nil
	default:
		return "", fmt.Errorf("unknown costmanagement scheduled action kind: %s", *kind)
	}
}
