package tfid

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageShareFile(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	shareId, err := buildStorageShare(b, id.Parent(), spec)
	if err != nil {
		return "", err
	}
	uri, err := url.Parse(shareId)
	if err != nil {
		return "", fmt.Errorf("parsing uri %s: %v", shareId, err)
	}
	path := id.Names()[3]
	path = strings.ReplaceAll(path, ":", "/")
	uri = uri.JoinPath(path)
	return uri.String(), nil
}
