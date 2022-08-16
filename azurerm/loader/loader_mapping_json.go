package loader

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/Azure/azapi2azurerm/azurerm/types"
)

type MappingJsonDependencyLoader struct {
	MappingJsonFilepath string
}

//go:embed mappings.json
var mappingsJson string

func (m MappingJsonDependencyLoader) Load() ([]types.Dependency, error) {
	var mappings []types.Mapping
	var data []byte
	var err error
	if len(m.MappingJsonFilepath) > 0 {
		data, err = os.ReadFile(m.MappingJsonFilepath)
		if err != nil {
			return []types.Dependency{}, err
		}
	} else {
		data = []byte(mappingsJson)
	}
	err = json.Unmarshal(data, &mappings)
	if err != nil {
		return []types.Dependency{}, err
	}
	deps := make([]types.Dependency, 0)
	for _, mapping := range mappings {
		deps = append(deps, types.Dependency{
			Pattern:              mapping.IdPattern,
			ExampleConfiguration: mapping.ExampleConfiguration,
			ResourceType:         mapping.ResourceType,
			ReferredProperty:     "id",
		})
	}

	return deps, nil
}

var _ DependencyLoader = MappingJsonDependencyLoader{}
