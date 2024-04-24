package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type springApmsResolver struct{}

func (springApmsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_spring_cloud_dynatrace_application_performance_monitoring",
		"azurerm_spring_cloud_application_insights_application_performance_monitoring",
		"azurerm_spring_cloud_new_relic_application_performance_monitoring",
		"azurerm_spring_cloud_elastic_application_performance_monitoring",
		"azurerm_spring_cloud_app_dynamics_application_performance_monitoring",
	}
}

func (springApmsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	c, err := b.NewRawClient()
	if err != nil {
		return "", err
	}
	resp, err := c.Get(context.Background(), id.String(), "2023-11-01-preview")
	if err != nil {
		return "", err
	}
	m, ok := resp.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("GET on %q: response is not a map: %T", id, resp)
	}
	props, ok := m["properties"]
	if !ok {
		return "", fmt.Errorf("response of GET on %q has no `properties`", id)
	}
	m, ok = props.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("GET on %q: response.properties is not a map: %T", id, props)
	}
	typ, ok := m["type"].(string)
	if !ok {
		return "", fmt.Errorf("GET on %q: response.properties.type is not a string: %T", id, m["type"])
	}
	switch typ {
	case "ElasticAPM":
		return "azurerm_spring_cloud_elastic_application_performance_monitoring", nil
	case "Dynatrace":
		return "azurerm_spring_cloud_dynatrace_application_performance_monitoring", nil
	case "NewRelic":
		return "azurerm_spring_cloud_new_relic_application_performance_monitoring", nil
	case "ApplicationInsights":
		return "azurerm_spring_cloud_application_insights_application_performance_monitoring", nil
	case "AppDynamics":
		return "azurerm_spring_cloud_app_dynamics_application_performance_monitoring", nil
	default:
		return "", fmt.Errorf("unknown spring APM type: %s", typ)
	}
}
