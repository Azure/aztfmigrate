package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type logicAppTrigger struct{}

func (logicAppTrigger) ResourceTypes() []string {
	return []string{"azurerm_logic_app_trigger_recurrence", "azurerm_logic_app_trigger_custom", "azurerm_logic_app_trigger_http_request"}
}

func (logicAppTrigger) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewLogicWorkflowsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Workflow.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}

	type Def struct {
		Triggers map[string]struct {
			Type string
		}
	}

	rb, err := json.Marshal(props.Definition)
	if err != nil {
		return "", fmt.Errorf("marshaling definition: %v", err)
	}
	var def Def
	if err := json.Unmarshal(rb, &def); err != nil {
		return "", fmt.Errorf("unmarshaling definition: %v", err)
	}

	if len(def.Triggers) == 0 {
		return "", fmt.Errorf("unexpected nil triggers")
	}

	trigger, ok := def.Triggers[id.Names()[1]]
	if !ok {
		return "", fmt.Errorf("can't find trigger with name %s", id.Names()[1])
	}
	switch strings.ToLower(trigger.Type) {
	case "request":
		return "azurerm_logic_app_trigger_http_request", nil
	case "recurrence":
		return "azurerm_logic_app_trigger_recurrence", nil
	default:
		return "azurerm_logic_app_trigger_custom", nil
	}
}
