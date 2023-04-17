package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildKeyVaultSecret(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewKeyVaultSecretsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Secret.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	uri := props.SecretURIWithVersion
	if uri == nil {
		return "", fmt.Errorf("unexpected nil properties.secretUriWithVersion in response")
	}
	return *uri, nil
}
