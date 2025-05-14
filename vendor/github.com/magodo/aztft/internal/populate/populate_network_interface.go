package populate

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

func populateNetworkInterface(b *client.ClientBuilder, id armid.ResourceId) ([]armid.ResourceId, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkInterfacesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Interface.Properties
	if props == nil {
		return nil, nil
	}

	nsgAssociations, err := networkInterfacePopulateNSGAssociation(id, props)
	if err != nil {
		return nil, fmt.Errorf("populating for NSG associations: %v", err)
	}

	bapAssociations, err := networkInterfacePopulateIpConfigAssociations(id, props)
	if err != nil {
		return nil, fmt.Errorf("populating for Application Gateway BAP associations: %v", err)
	}

	var result []armid.ResourceId
	result = append(result, nsgAssociations...)
	result = append(result, bapAssociations...)
	return result, nil
}

func networkInterfacePopulateNSGAssociation(id armid.ResourceId, props *armnetwork.InterfacePropertiesFormat) ([]armid.ResourceId, error) {
	nsgProp := props.NetworkSecurityGroup
	if nsgProp == nil {
		return nil, nil
	}

	nsgId := nsgProp.ID
	if nsgId == nil {
	}

	nsgAzureId, err := armid.ParseResourceId(*nsgId)
	if err != nil {
		return nil, fmt.Errorf("parsing resource id %q: %v", *nsgId, err)
	}
	azureId := id.Clone().(*armid.ScopedResourceId)
	azureId.AttrTypes = append(azureId.AttrTypes, "networkSecurityGroups")
	azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(nsgAzureId.String())))

	return []armid.ResourceId{azureId}, nil
}

func networkInterfacePopulateIpConfigAssociations(id armid.ResourceId, props *armnetwork.InterfacePropertiesFormat) ([]armid.ResourceId, error) {
	var result []armid.ResourceId
	for _, ipConfig := range props.IPConfigurations {
		if ipConfig == nil {
			continue
		}

		ipConfigProps := ipConfig.Properties
		if ipConfigProps == nil {
			continue
		}

		if ipConfig.ID == nil {
			continue
		}

		ipConfigId, err := armid.ParseResourceId(*ipConfig.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", *ipConfig.ID, err)
		}

		for _, bap := range ipConfigProps.ApplicationGatewayBackendAddressPools {
			if bap == nil || bap.ID == nil {
				continue
			}
			azureId, err := networkInterfacePopulateIpConfigApplicationGatewayBackendAddressPoolAssociation(ipConfigId, *bap.ID)
			if err != nil {
				return nil, err
			}
			result = append(result, azureId)
		}

		for _, asg := range ipConfigProps.ApplicationSecurityGroups {
			if asg == nil || asg.ID == nil {
				continue
			}
			azureId, err := networkInterfacePopulateIpConfigApplicationSecurityGroupAssociation(ipConfigId, *asg.ID)
			if err != nil {
				return nil, err
			}
			result = append(result, azureId)
		}

		for _, natRule := range ipConfigProps.LoadBalancerInboundNatRules {
			if natRule == nil || natRule.ID == nil {
				continue
			}
			azureId, err := networkInterfacePopulateIpConfigLoadBalancerNatRuleAssociation(ipConfigId, *natRule.ID)
			if err != nil {
				return nil, err
			}
			result = append(result, azureId)
		}

		for _, bap := range ipConfigProps.LoadBalancerBackendAddressPools {
			if bap == nil || bap.ID == nil {
				continue
			}
			azureId, err := networkInterfacePopulateIpConfigLoadBalancerBackendAddressPoolAssociation(ipConfigId, *bap.ID)
			if err != nil {
				return nil, err
			}
			result = append(result, azureId)
		}
	}
	return result, nil
}

func networkInterfacePopulateIpConfigApplicationGatewayBackendAddressPoolAssociation(ipConfigId armid.ResourceId, id string) (armid.ResourceId, error) {
	bapId, err := armid.ParseResourceId(id)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %v", id, err)
	}
	azureId := ipConfigId.Clone().(*armid.ScopedResourceId)
	azureId.AttrTypes = append(azureId.AttrTypes, "applicationGatewayBackendAddressPools")
	azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(bapId.String())))
	return azureId, nil
}

func networkInterfacePopulateIpConfigApplicationSecurityGroupAssociation(ipConfigId armid.ResourceId, id string) (armid.ResourceId, error) {
	asgId, err := armid.ParseResourceId(id)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %v", id, err)
	}
	azureId := ipConfigId.Clone().(*armid.ScopedResourceId)
	azureId.AttrTypes = append(azureId.AttrTypes, "applicationSecurityGroups")
	azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(asgId.String())))
	return azureId, nil
}

func networkInterfacePopulateIpConfigLoadBalancerNatRuleAssociation(ipConfigId armid.ResourceId, id string) (armid.ResourceId, error) {
	natRuleId, err := armid.ParseResourceId(id)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %v", id, err)
	}
	azureId := ipConfigId.Clone().(*armid.ScopedResourceId)
	azureId.AttrTypes = append(azureId.AttrTypes, "loadBalancerInboundNatRules")
	azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(natRuleId.String())))
	return azureId, nil
}

func networkInterfacePopulateIpConfigLoadBalancerBackendAddressPoolAssociation(ipConfigId armid.ResourceId, id string) (armid.ResourceId, error) {
	bapId, err := armid.ParseResourceId(id)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %v", id, err)
	}
	azureId := ipConfigId.Clone().(*armid.ScopedResourceId)
	azureId.AttrTypes = append(azureId.AttrTypes, "loadBalancerBackendAddressPools")
	azureId.AttrNames = append(azureId.AttrNames, base64.StdEncoding.EncodeToString([]byte(bapId.String())))
	return azureId, nil
}
