package tfid

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func buildActiveDirectoryDomainService(b *client.ClientBuilder, id armid.ResourceId, spec string) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDomainServiceClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.DomainService.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	if len(props.ReplicaSets) == 0 {
		return "", fmt.Errorf("unexpected 0 properties.replicaSets in response")
	}
	initReplicaSet := props.ReplicaSets[0]
	if initReplicaSet == nil {
		return "", fmt.Errorf("unexpected nil properties.replicaSets[0] in response")
	}
	initReplicaSetId := initReplicaSet.ReplicaSetID
	if initReplicaSetId == nil {
		return "", fmt.Errorf("unexpected nil properties.replicaSets[0].replicaSetId in response")
	}
	rid := id.(*armid.ScopedResourceId)
	rid.AttrTypes = append(rid.AttrTypes, "initialReplicaSetId")
	rid.AttrNames = append(rid.AttrNames, *initReplicaSetId)
	if err := id.Normalize(spec); err != nil {
		return "", fmt.Errorf("normalizing id %q with import spec %q: %v", id.String(), spec, err)
	}
	return id.String(), nil
}
