package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentscripts/v2"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type deploymentScriptsResolver struct{}

func (deploymentScriptsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_resource_deployment_script_azure_cli",
		"azurerm_resource_deployment_script_azure_power_shell",
	}
}

func (deploymentScriptsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDeploymentScriptsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.DeploymentScriptClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}

	switch model.(type) {
	case *armdeploymentscripts.AzureCliScript:
		return "azurerm_resource_deployment_script_azure_cli", nil
	case *armdeploymentscripts.AzurePowerShellScript:
		return "azurerm_resource_deployment_script_azure_power_shell", nil
	default:
		return "", fmt.Errorf("unknown deployment scripts type: %T", model)
	}
}
