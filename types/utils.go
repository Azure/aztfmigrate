package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func getApiVersion(value interface{}) string {
	if valueMap, ok := value.(map[string]interface{}); ok && valueMap["type"] != nil {
		if typeValue, ok := valueMap["type"].(string); ok {
			if parts := strings.Split(typeValue, "@"); len(parts) == 2 {
				return parts[1]
			}
		}
	}
	return ""
}

func getId(value interface{}) string {
	if valueMap, ok := value.(map[string]interface{}); ok && valueMap["id"] != nil {
		if resourceId, ok := valueMap["id"].(string); ok {
			return resourceId
		}
	}
	return ""
}

func getOutputsForAddress(address string, refValueMap map[string]interface{}) []Output {
	res := make([]Output, 0)
	for key, value := range refValueMap {
		if strings.HasPrefix(key, fmt.Sprintf("%s.output.", address)) {
			res = append(res, Output{
				OldName: key,
				Value:   value,
			})
		}
	}
	return res
}

func getReferencesForAddress(address string, p *tfjson.Plan, refValueMap map[string]interface{}) []Reference {
	res := make([]Reference, 0)
	for _, r := range p.Config.RootModule.Resources {
		if r.Address == address {
			for _, expression := range r.Expressions {
				res = append(res, listReferences(expression)...)
			}
			temp := make([]Reference, 0)
			for _, ref := range res {
				// if it refers to some resource's id before, after migration, it will refer to its name now
				// TODO: use regex
				if len(strings.Split(ref.Name, ".")) == 3 && strings.HasSuffix(ref.Name, "id") {
					temp = append(temp, Reference{
						Name: ref.Name[0:strings.LastIndex(ref.Name, ".")] + ".name",
					})
					if strings.HasPrefix(ref.Name, "azurerm_resource_group") {
						temp = append(temp, Reference{
							Name: ref.Name[0:strings.LastIndex(ref.Name, ".")] + ".location",
						})
					}
				}
				if len(strings.Split(ref.Name, ".")) == 3 && strings.HasSuffix(ref.Name, "name") {
					temp = append(temp, Reference{
						Name: ref.Name[0:strings.LastIndex(ref.Name, ".")] + ".id",
					})
				}
			}
			res = append(res, temp...)
			break
		}
	}

	refSet := make(map[string]Reference)
	for i := range res {
		refSet[res[i].Name] = res[i]
	}
	res = make([]Reference, 0)
	for _, ref := range refSet {
		ref.Value = refValueMap[ref.Name]
		res = append(res, ref)
	}
	return res
}

func listReferences(expression *tfjson.Expression) []Reference {
	res := make([]Reference, 0)
	for _, ref := range expression.References {
		res = append(res, Reference{
			Name: ref,
		})
	}
	for _, block := range expression.NestedBlocks {
		for _, exp := range block {
			res = append(res, listReferences(exp)...)
		}
	}
	return res
}

func getRefValueMap(p *tfjson.Plan) map[string]interface{} {
	refValueMap := make(map[string]interface{})
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Change.Before == nil {
			continue
		}
		prefix := resourceChange.Address
		if strings.HasPrefix(prefix, "azapi") {
			if beforeMap, ok := resourceChange.Change.Before.(map[string]interface{}); ok && beforeMap["output"] != nil {
				if output, ok := beforeMap["output"].(string); ok {
					var outputObj interface{}
					if err := json.Unmarshal([]byte(output), &outputObj); err == nil {
						propValueMap := getPropValueMap(outputObj, fmt.Sprintf("jsondecode(%s.output)", prefix))
						for key, value := range propValueMap {
							refValueMap[key] = value
						}
					} else {
						if outputObj := beforeMap["output"]; outputObj != nil {
							propValueMap := getPropValueMap(outputObj, fmt.Sprintf("%s.output", prefix))
							for key, value := range propValueMap {
								refValueMap[key] = value
							}
						}
					}
				}
			}
		}
		propValueMap := getPropValueMap(resourceChange.Change.Before, prefix)
		for key, value := range propValueMap {
			refValueMap[key] = value
		}
	}
	for key, variable := range p.Variables {
		refValueMap["var."+key] = variable.Value
	}
	return refValueMap
}

func getPropValueMap(input interface{}, prefix string) map[string]interface{} {
	res := make(map[string]interface{})
	if input == nil {
		return res
	}
	switch cur := input.(type) {
	case map[string]interface{}:
		for key, value := range cur {
			propValueMap := getPropValueMap(value, prefix+"."+key)
			for k, v := range propValueMap {
				res[k] = v
			}
		}
	case []interface{}:
		for index, value := range cur {
			propValueMap := getPropValueMap(value, fmt.Sprintf("%s.%d", prefix, index))
			for k, v := range propValueMap {
				res[k] = v
			}
		}
	default:
		res[prefix] = cur
	}
	return res
}

func getInputProperties(address string, p *tfjson.Plan) []string {
	for _, resourceChange := range p.ResourceChanges {
		if resourceChange == nil || resourceChange.Change == nil || resourceChange.Change.Before == nil ||
			(resourceChange.Address != address && !strings.HasPrefix(resourceChange.Address, address+"[")) {
			continue
		}

		stateMap, ok := resourceChange.Change.Before.(map[string]interface{})
		if !ok {
			return nil
		}

		props := make([]string, 0)
		if stateMap["tags"] != nil {
			if tags, ok := stateMap["tags"].(map[string]interface{}); ok && len(tags) > 0 {
				props = append(props, "tags")
			}
		}
		if stateMap["identity"] != nil {
			if identities, ok := stateMap["identity"].([]interface{}); ok && len(identities) > 0 {
				if identity, ok := identities[0].(map[string]interface{}); ok {
					if identity["type"] != nil {
						if identityType, ok := identity["type"].(string); ok && len(identityType) > 0 {
							props = append(props, "identity.type")
						}
						if identityIds, ok := identity["identity_ids"].([]interface{}); ok && len(identityIds) > 0 {
							props = append(props, "identity.userAssignedIdentities")
						}
					}
				}
			}
		}
		if stateMap["body"] != nil {
			if body, ok := stateMap["body"].(string); ok {
				var bodyObj interface{}
				if err := json.Unmarshal([]byte(body), &bodyObj); err == nil {
					propValueMap := getPropValueMap(bodyObj, "")
					propSet := make(map[string]bool)
					for key := range propValueMap {
						key = strings.TrimPrefix(key, ".")
						if strings.HasPrefix(key, "tags") {
							key = "tags"
						}
						propSet[key] = true
					}
					for key := range propSet {
						key = removeIndexOfProp(key)
						props = append(props, key)
					}
				}
			} else {
				if bodyObj := stateMap["body"]; bodyObj != nil {
					propValueMap := getPropValueMap(bodyObj, "")
					propSet := make(map[string]bool)
					for key := range propValueMap {
						key = strings.TrimPrefix(key, ".")
						if strings.HasPrefix(key, "tags") {
							key = "tags"
						}
						propSet[key] = true
					}
					for key := range propSet {
						key = removeIndexOfProp(key)
						props = append(props, key)
					}
				}
			}
		}

		return props
	}
	return nil
}

func removeIndexOfProp(prop string) string {
	parts := strings.Split(prop, ".")
	res := make([]string, 0)
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err == nil {
			continue
		}
		res = append(res, part)
	}
	return strings.Join(res, ".")
}
