package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type machineLearningWorkspaceResolver struct{}

func (machineLearningWorkspaceResolver) ResourceTypes() []string {
	return []string{
		"azurerm_machine_learning_workspace",
		"azurerm_ai_foundry_project",
		"azurerm_ai_foundry",
	}
}

func (machineLearningWorkspaceResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewMachineLearningWorkspaceClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	kind := resp.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil kind in response")
	}
	switch strings.ToUpper(string(*kind)) {
	case "DEFAULT", "FEATURESTORE":
		return "azurerm_machine_learning_workspace", nil
	case "PROJECT":
		return "azurerm_ai_foundry_project", nil
	case "HUB":
		return "azurerm_ai_foundry", nil
	}
	return "", fmt.Errorf("unknown kind: %s", *kind)
}
