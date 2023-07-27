package resolve

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type logicAppAction struct{}

func (logicAppAction) ResourceTypes() []string {
	return []string{"azurerm_logic_app_action_custom", "azurerm_logic_app_action_http"}
}

func (logicAppAction) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
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
		Actions map[string]struct {
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

	if len(def.Actions) == 0 {
		return "", fmt.Errorf("unexpected nil actions")
	}

	action, ok := def.Actions[id.Names()[1]]
	if !ok {
		return "", fmt.Errorf("can't find action with name %s", id.Names()[1])
	}
	switch strings.ToLower(action.Type) {
	case "http":
		return "azurerm_logic_app_action_http", nil
	default:
		return "azurerm_logic_app_action_custom", nil
	}
}
