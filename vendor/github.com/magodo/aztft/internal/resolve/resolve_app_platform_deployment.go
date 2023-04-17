package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appplatform/armappplatform"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type appPlatformDeploymentsResolver struct{}

func (appPlatformDeploymentsResolver) ResourceTypes() []string {
	return []string{"azurerm_spring_cloud_build_deployment", "azurerm_spring_cloud_java_deployment", "azurerm_spring_cloud_container_deployment"}
}

func (appPlatformDeploymentsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAppPlatformDeploymentsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], id.Names()[2], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.DeploymentResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	source := props.Source
	if source == nil {
		return "", fmt.Errorf("unexpected nil properties.source in response")
	}

	switch source.(type) {
	case *armappplatform.BuildResultUserSourceInfo:
		return "azurerm_spring_cloud_build_deployment", nil
	case *armappplatform.JarUploadedUserSourceInfo:
		return "azurerm_spring_cloud_java_deployment", nil
	case *armappplatform.CustomContainerUserSourceInfo:
		return "azurerm_spring_cloud_container_deployment", nil
	default:
		return "", fmt.Errorf("unknown spring cloud deployment source type: %T", source)
	}
}
