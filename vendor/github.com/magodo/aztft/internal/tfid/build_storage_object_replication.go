package tfid

import (
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildStorageObjectReplication(b *client.ClientBuilder, id armid.ResourceId, _ string) (string, error) {
	// This is not supported as in the response body of the source policy only contains the destination policy's storage account name.
	// In order to get the destination policy id, we'll have to query the storage account resource id by name, which will hit the Azure Resource list API bug.
	// Therefore, there is no good way to implement this at this moment.
	return "", fmt.Errorf("this is not supported yet")
}
