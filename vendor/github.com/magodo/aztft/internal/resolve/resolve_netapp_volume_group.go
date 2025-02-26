package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type netappVolumeGroupResolver struct{}

func (netappVolumeGroupResolver) ResourceTypes() []string {
	return []string{
		"azurerm_netapp_volume_group_oracle",
		"azurerm_netapp_volume_group_sap_hana",
	}
}

func (netappVolumeGroupResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetAppVolumeGroupClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.VolumeGroupDetails.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}
	gmetadata := props.GroupMetaData
	if gmetadata == nil {
		return "", fmt.Errorf("unexpected nil groupMetaData in response")
	}
	appType := gmetadata.ApplicationType
	if appType == nil {
		return "", fmt.Errorf("unexpected nil applicationType in response")
	}

	switch strings.ToUpper(string(*appType)) {
	case "ORACLE":
		return "azurerm_netapp_volume_group_oracle", nil
	case "SAP-HANA":
		return "azurerm_netapp_volume_group_sap_hana", nil
	}
	return "", fmt.Errorf("unknown application type: %s", *appType)
}
