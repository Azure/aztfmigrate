package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateLogicAppWorkflow(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewLogicWorkflowsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Workflow.Properties
	if props == nil {
		return nil, nil
	}

	if props.Definition == nil {
		return nil, nil
	}

	def := props.Definition.(map[string]interface{})

	var result []armid.ResourceId

	if actionsRaw, ok := def["actions"]; ok {
		actions := actionsRaw.(map[string]interface{})
		for name := range actions {
			azureId := id.Clone().(*armid.ScopedResourceId)
			azureId.AttrTypes = append(azureId.AttrTypes, "actions")
			azureId.AttrNames = append(azureId.AttrNames, name)
			result = append(result, azureId)
		}
	}

	if triggersRaw, ok := def["triggers"]; ok {
		triggers := triggersRaw.(map[string]interface{})
		for name := range triggers {
			azureId := id.Clone().(*armid.ScopedResourceId)
			azureId.AttrTypes = append(azureId.AttrTypes, "triggers")
			azureId.AttrNames = append(azureId.AttrNames, name)
			result = append(result, azureId)
		}
	}

	return result, nil
}
