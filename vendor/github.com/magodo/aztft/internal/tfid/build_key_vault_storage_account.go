package tfid

import (
	"context"
	"fmt"
	"net/url"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildKeyVaultStorageAccount(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewKeyVaultVaultsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Vault.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	puri := props.VaultURI
	if puri == nil {
		return "", fmt.Errorf("unexpected nil properties.vaultUri in response")
	}
	uri, err := url.Parse(*puri)
	if err != nil {
		return "", fmt.Errorf("parsing uri %s: %v", *puri, err)
	}
	uri.Path = "/storage/" + id.Names()[1]
	return uri.String(), nil
}
