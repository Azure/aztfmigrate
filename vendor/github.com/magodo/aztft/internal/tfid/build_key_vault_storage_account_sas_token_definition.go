package tfid

import (
	"fmt"
	"net/url"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildKeyVaultStorageAccountSasTokenDefinition(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	storageId, err := buildKeyVaultStorageAccount(b, id.Parent(), spec)
	if err != nil {
		return "", err
	}
	uri, err := url.Parse(storageId)
	if err != nil {
		return "", fmt.Errorf("parsing uri %s: %v", storageId, err)
	}
	uri = uri.JoinPath("sas", id.Names()[2])
	return uri.String(), nil
}
