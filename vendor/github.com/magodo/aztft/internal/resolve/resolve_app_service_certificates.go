package resolve

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type appServiceCertificatesResolver struct{}

func (appServiceCertificatesResolver) ResourceTypes() []string {
	return []string{"azurerm_app_service_certificate", "azurerm_app_service_managed_certificate"}
}

func (appServiceCertificatesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewAppServiceCertificatesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.AppCertificate.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	if props.ServerFarmID == nil {
		return "azurerm_app_service_certificate", nil
	}
	return "azurerm_app_service_managed_certificate", nil
}
