package types

import (
	"fmt"
	"log"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func ListResourcesFromPlan(p *tfjson.Plan) []AzureResource {
	resources := make([]AzureResource, 0)
	if p == nil {
		return resources
	}

	idMap := make(map[string]*tfjson.ResourceChange)
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.ProviderName != "registry.terraform.io/hashicorp/azurerm" {
			continue
		}
		idMap[getId(resourceChange.Change.Before)] = resourceChange
	}

	azapiResourceMap := make(map[string]*AzapiResource)
	azapiUpdateResources := make([]AzapiUpdateResource, 0)
	azurermResourceMap := make(map[string]*AzurermResource)
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil {
			continue
		}

		switch resourceChange.Type {
		case "azapi_resource":
			address := fmt.Sprintf("%s.%s", resourceChange.Type, resourceChange.Name)
			if azapiResourceMap[address] == nil {
				azapiResourceMap[address] = &AzapiResource{
					Label:        resourceChange.Name,
					ResourceType: "",
					Instances:    make([]Instance, 0),
				}
			}

			azapiResourceMap[address].Instances = append(azapiResourceMap[address].Instances, Instance{
				Index:      resourceChange.Index,
				ResourceId: getId(resourceChange.Change.Before),
				ApiVersion: getApiVersion(resourceChange.Change.Before),
			})

		case "azapi_update_resource":
			resourceId := getId(resourceChange.Change.Before)
			if idMap[resourceId] == nil {
				log.Printf("[WARN] resource azapi_update_resource.%s's target is not in the same terraform working directory", resourceChange.Name)
				continue
			}
			rc := idMap[resourceId]
			azapiUpdateResources = append(azapiUpdateResources, AzapiUpdateResource{
				OldLabel:     resourceChange.Name,
				Label:        rc.Name,
				Id:           resourceId,
				ApiVersion:   getApiVersion(resourceChange.Change.Before),
				ResourceType: rc.Type,
				Change:       rc.Change,
			})
		default:
			if strings.HasPrefix(resourceChange.Type, "azurerm") {
				address := fmt.Sprintf("%s.%s", resourceChange.Type, resourceChange.Name)
				id := getId(resourceChange.Change.Before)
				if azurermResourceMap[address] == nil {
					azurermResourceMap[address] = &AzurermResource{
						OldResourceType: resourceChange.Type,
						OldLabel:        resourceChange.Name,
						NewResourceType: "azapi_resource",
						NewLabel:        NewLabel(id, resourceChange.Name),
						Instances:       make([]Instance, 0),
					}
				}

				azurermResourceMap[address].Instances = append(azurermResourceMap[address].Instances, Instance{
					Index:      resourceChange.Index,
					ResourceId: id,
				})
			}
		}
	}

	refValueMap := getRefValueMap(p)
	azapiResources := make([]AzapiResource, 0)
	for _, v := range azapiResourceMap {
		azapiResources = append(azapiResources, *v)
	}

	azurermResources := make([]AzurermResource, 0)
	for _, resourceChange := range azurermResourceMap {
		azurermResources = append(azurermResources, *resourceChange)
	}

	for index, resource := range azapiResources {
		azapiResources[index].References = getReferencesForAddress(resource.OldAddress(nil), p, refValueMap)
		azapiResources[index].InputProperties = getInputProperties(resource.OldAddress(nil), p)
	}
	for i, resource := range azapiResources {
		outputPropSet := make(map[string]bool)
		for j, instance := range resource.Instances {
			azapiResources[i].Instances[j].Outputs = getOutputsForAddress(resource.OldAddress(instance.Index), refValueMap)
			for _, output := range azapiResources[i].Instances[j].Outputs {
				prop := strings.TrimPrefix(output.OldName, fmt.Sprintf("%s.output.", resource.OldAddress(instance.Index)))
				if strings.HasPrefix(prop, "identity.userAssignedIdentities") {
					prop = "identity.userAssignedIdentities"
				}
				outputPropSet[prop] = true
			}
		}
		azapiResources[i].OutputProperties = make([]string, 0)
		for key := range outputPropSet {
			azapiResources[i].OutputProperties = append(azapiResources[i].OutputProperties, key)
		}
	}

	for index, resource := range azapiUpdateResources {
		azapiUpdateResources[index].References = getReferencesForAddress(resource.OldAddress(nil), p, refValueMap)
		azapiUpdateResources[index].InputProperties = getInputProperties(resource.OldAddress(nil), p)
	}
	for i, resource := range azapiUpdateResources {
		azapiUpdateResources[i].outputs = getOutputsForAddress(resource.OldAddress(nil), refValueMap)
		azapiUpdateResources[i].OutputProperties = make([]string, 0)
		for _, output := range azapiUpdateResources[i].outputs {
			azapiUpdateResources[i].OutputProperties = append(azapiUpdateResources[i].OutputProperties, strings.TrimPrefix(output.OldName, fmt.Sprintf("%s.output.", resource.OldAddress(nil))))
		}
	}

	for index, resource := range azurermResources {
		azurermResources[index].References = getReferencesForAddress(resource.OldAddress(nil), p, refValueMap)
	}

	for _, resource := range azapiResources {
		resources = append(resources, &resource)
	}
	for _, resource := range azapiUpdateResources {
		resources = append(resources, &resource)
	}
	for _, resource := range azurermResources {
		resources = append(resources, &resource)
	}

	return resources
}
