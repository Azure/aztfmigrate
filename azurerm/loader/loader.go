package loader

import "github.com/Azure/azapi2azurerm/azurerm/types"

type DependencyLoader interface {
	Load() ([]types.Dependency, error)
}
