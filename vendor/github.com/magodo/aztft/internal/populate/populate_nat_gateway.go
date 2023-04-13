package populate

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateNatGateway(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkNatGatewaysClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.NatGateway.Properties
	if props == nil {
		return nil, nil
	}

	pipAssociations, err := natGatewayPopulatePublicIpAssociation(id, props)
	if err != nil {
		return nil, fmt.Errorf("populating for public ip associations: %v", err)
	}
	pipPrefixAssociations, err := natGatewayPopulatePublicIpPrefixAssociation(id, props)
	if err != nil {
		return nil, fmt.Errorf("populating for public ip prefix associations: %v", err)
	}

	var result []armid.ResourceId
	result = append(result, pipAssociations...)
	result = append(result, pipPrefixAssociations...)

	return result, nil
}

func natGatewayPopulatePublicIpAssociation(id armid.ResourceId, props *armnetwork.NatGatewayPropertiesFormat) ([]armid.ResourceId, error) {
	var result []armid.ResourceId

	for _, pip := range props.PublicIPAddresses {
		if pip == nil {
			continue
		}
		if pip.ID == nil {
			continue
		}
		pipId, err := armid.ParseResourceId(*pip.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", *pip.ID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "publicIPAddresses")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(pipId.String())))

		result = append(result, azureId)
	}

	return result, nil
}

func natGatewayPopulatePublicIpPrefixAssociation(id armid.ResourceId, props *armnetwork.NatGatewayPropertiesFormat) ([]armid.ResourceId, error) {
	var result []armid.ResourceId

	for _, prefix := range props.PublicIPPrefixes {
		if prefix == nil {
			continue
		}
		if prefix.ID == nil {
			continue
		}
		prefixId, err := armid.ParseResourceId(*prefix.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing resource id %q: %v", *prefix.ID, err)
		}
		azureId := id.Clone().(*armid.ScopedResourceId)
		azureId.AttrTypes = append(azureId.AttrTypes, "publicIPPrefixes")
		azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(prefixId.String())))

		result = append(result, azureId)
	}

	return result, nil
}
