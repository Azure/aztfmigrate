package loader

import "github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm/types"

type DependencyLoader interface {
	Load() ([]types.Dependency, error)
}
