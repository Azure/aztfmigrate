package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v4"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type machineLearningOutboundRulesResolver struct{}

func (machineLearningOutboundRulesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_machine_learning_workspace_network_outbound_rule_service_tag",
		"azurerm_machine_learning_workspace_network_outbound_rule_fqdn",
		"azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint",
	}
}

func (machineLearningOutboundRulesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewMachineLearningOutboundRulesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.OutboundRuleBasicResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armmachinelearning.FqdnOutboundRule:
		return "azurerm_machine_learning_workspace_network_outbound_rule_fqdn", nil
	case *armmachinelearning.PrivateEndpointOutboundRule:
		return "azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint", nil
	case *armmachinelearning.ServiceTagOutboundRule:
		return "azurerm_machine_learning_workspace_network_outbound_rule_service_tag", nil
	default:
		return "", fmt.Errorf("unknown outbound rule resource type %T", props)
	}
}
