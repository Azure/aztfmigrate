package populate

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateSubnet(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkSubnetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Subnet.Properties
	if props == nil {
		return nil, nil
	}

	var result []armid.ResourceId

	if props.RouteTable != nil && props.RouteTable.ID != nil {
		routeTableId, err := armid.ParseResourceId(*props.RouteTable.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", *props.RouteTable.ID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "routeTables")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(routeTableId.String())))
		result = append(result, azureId)
	}

	if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.ID != nil {
		nsgId, err := armid.ParseResourceId(*props.NetworkSecurityGroup.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", *props.NetworkSecurityGroup.ID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "networkSecurityGroups")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(nsgId.String())))
		result = append(result, azureId)
	}

	if props.NatGateway != nil && props.NatGateway.ID != nil {
		natGwId, err := armid.ParseResourceId(*props.NatGateway.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", *props.NatGateway.ID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "natGateways")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(natGwId.String())))
		result = append(result, azureId)
	}

	return result, nil
}
