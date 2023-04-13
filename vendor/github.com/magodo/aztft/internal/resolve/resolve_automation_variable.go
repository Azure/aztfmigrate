package resolve

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type automationVariablesResolver struct{}

func (automationVariablesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_automation_variable_datetime",
		"azurerm_automation_variable_string",
		"azurerm_automation_variable_bool",
		"azurerm_automation_variable_int",
	}
}

func (automationVariablesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAutomationVariableClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Variable.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	value := props.Value
	if value == nil {
		return "", fmt.Errorf("unexpected nil properties.value in response")
	}

	// Referenced from: https://github.com/hashicorp/terraform-provider-azurerm/blob/a053df86e2d9790bf0d99a3283f979c8a944d3f5/internal/services/automation/automation_variable.go#L37
	datePattern := regexp.MustCompile(`"\\/Date\((-?[0-9]+)\)\\/"`)
	matches := datePattern.FindStringSubmatch(*value)
	if len(matches) == 2 && matches[0] == *value {
		if _, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			return "azurerm_automation_variable_datetime", nil
		}
	}
	if _, err = strconv.Unquote(*value); err == nil {
		return "azurerm_automation_variable_string", nil
	}
	if _, err = strconv.ParseInt(*value, 10, 32); err == nil {
		return "azurerm_automation_variable_int", nil
	}
	if _, err = strconv.ParseBool(*value); err == nil {
		return "azurerm_automation_variable_bool", nil
	}
	return "", fmt.Errorf("can't resolve resource type from value: %q", *value)
}
