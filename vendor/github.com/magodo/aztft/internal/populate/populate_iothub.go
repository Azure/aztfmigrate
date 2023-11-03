package populate

import (
	"context"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateIotHub(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewIothubsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Description.Properties
	if props == nil {
		return nil, nil
	}

	if props.Routing == nil || props.Routing.Endpoints == nil {
		return nil, nil
	}

	endpoints := *props.Routing.Endpoints

	var result []armid.ResourceId

	for _, ep := range endpoints.EventHubs {
		if ep == nil {
			continue
		}
		if ep.Name == nil {
			continue
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "endpointsEventhub")
		azureId.AttrNames = append(azureId.AttrNames, *ep.Name)
		result = append(result, azureId)
	}
	for _, ep := range endpoints.ServiceBusQueues {
		if ep == nil {
			continue
		}
		if ep.Name == nil {
			continue
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "endpointsServicebusQueue")
		azureId.AttrNames = append(azureId.AttrNames, *ep.Name)
		result = append(result, azureId)
	}
	for _, ep := range endpoints.ServiceBusTopics {
		if ep == nil {
			continue
		}
		if ep.Name == nil {
			continue
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "endpointsServicebusTopic")
		azureId.AttrNames = append(azureId.AttrNames, *ep.Name)
		result = append(result, azureId)
	}
	for _, ep := range endpoints.StorageContainers {
		if ep == nil {
			continue
		}
		if ep.Name == nil {
			continue
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "endpointsStorageContainer")
		azureId.AttrNames = append(azureId.AttrNames, *ep.Name)
		result = append(result, azureId)
	}

	return result, nil
}
