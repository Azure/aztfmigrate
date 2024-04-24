package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateNetAppAccount(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetAppAccountClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return nil, nil
	}

	// The "azurerm_netapp_account_encryption" is only for key source Microsoft.KeyVault
	if props.Encryption == nil || props.Encryption.KeySource == nil || *props.Encryption.KeySource == "Microsoft.NetApp" {
		return nil, nil
	}

	eid := id.Clone().(*armid.ScopedResourceId)
	eid.AttrTypes = append(eid.AttrTypes, "encryptions")
	eid.AttrNames = append(eid.AttrNames, "enc1")
	return []armid.ResourceId{eid}, nil
}
