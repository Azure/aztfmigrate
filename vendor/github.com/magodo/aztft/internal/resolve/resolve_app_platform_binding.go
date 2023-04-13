package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type appPlatformBindingsResolver struct{}

func (appPlatformBindingsResolver) ResourceTypes() []string {
	return []string{"azurerm_spring_cloud_app_cosmosdb_association", "azurerm_spring_cloud_app_redis_association", "azurerm_spring_cloud_app_mysql_association"}
}

func (appPlatformBindingsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAppPlatformBindingsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], id.Names()[2], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.BindingResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	params := props.BindingParameters
	if params == nil {
		return "", fmt.Errorf("unexpected nil properties.bindingParams in response")
	}

	switch {
	case params["apiType"] != nil:
		return "azurerm_spring_cloud_app_cosmosdb_association", nil
	case params["useSsl"] != nil:
		return "azurerm_spring_cloud_app_redis_association", nil
	case params["databaseName"] != nil && params["username"] != nil:
		return "azurerm_spring_cloud_app_mysql_association", nil
	default:
		return "", fmt.Errorf("unknown spring binding type")
	}
}
