package azurerm

import (
	"github.com/Azure/azapi2azurerm/azurerm/loader"
	"github.com/Azure/azapi2azurerm/azurerm/types"
	"github.com/Azure/azapi2azurerm/helper"
)

var deps = make([]types.Dependency, 0)

func init() {
	mappingJsonLoader := loader.MappingJsonDependencyLoader{}
	hardcodeLoader := loader.HardcodeDependencyLoader{}
	deps = make([]types.Dependency, 0)
	depsMap := make(map[string]types.Dependency)
	if temp, err := mappingJsonLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType+"."+dep.ReferredProperty] = dep
		}
	}
	if temp, err := hardcodeLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType+"."+dep.ReferredProperty] = dep
		}
	}
	for _, dep := range depsMap {
		if dep.ReferredProperty == "id" {
			deps = append(deps, dep)
		}
	}
}

func GetAzureRMResourceType(id string) []string {
	resourceTypes := make([]string, 0)
	for _, dep := range deps {
		if helper.IsValueMatchPattern(id, dep.Pattern) {
			resourceTypes = append(resourceTypes, dep.ResourceType)
		}
	}
	return resourceTypes
}
