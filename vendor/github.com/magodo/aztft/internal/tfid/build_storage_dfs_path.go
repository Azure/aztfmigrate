package tfid

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageDfsPath(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	dfsId, err := buildStorageDfs(b, id.Parent(), spec)
	if err != nil {
		return "", err
	}
	uri, err := url.Parse(dfsId)
	if err != nil {
		return "", fmt.Errorf("parsing uri %s: %v", dfsId, err)
	}
	path := id.Names()[2]
	path = strings.ReplaceAll(path, ":", "/")
	uri = uri.JoinPath(path)
	return uri.String(), nil
}
