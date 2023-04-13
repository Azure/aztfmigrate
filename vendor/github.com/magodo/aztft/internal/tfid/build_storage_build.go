package tfid

import (
	"fmt"
	"net/url"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageBlob(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	containerId, err := buildStorageContainer(b, id.Parent(), spec)
	if err != nil {
		return "", err
	}
	uri, err := url.Parse(containerId)
	if err != nil {
		return "", fmt.Errorf("parsing uri %s: %v", containerId, err)
	}
	uri = uri.JoinPath(id.Names()[3])
	return uri.String(), nil
}
