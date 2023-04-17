package tfid

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildKeyVaultCertificate(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	// We use the key client here as the certificate is a data plane only resource, which is a combination of both a key and secret, with the same name.
	client, err := b.NewKeyVaultKeysClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Key.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	uri := props.KeyURIWithVersion
	if uri == nil {
		return "", fmt.Errorf("unexpected nil properties.keyUriWithVersion in response")
	}
	keyUrl, err := url.Parse(*uri)
	if err != nil {
		return "", fmt.Errorf("failed to parse uri %s: %v", *uri, err)
	}
	segs := strings.Split(keyUrl.Path, "/")
	segs[1] = "certificates"
	keyUrl.Path = strings.Join(segs, "/")
	return keyUrl.String(), nil
}
