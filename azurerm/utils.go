package azurerm

import (
	"github.com/magodo/aztft/aztft"
)

func GetAzureRMResourceType(id string) ([]string, bool, error) {
	resourceTypes := make([]string, 0)
	types, exact, err := aztft.QueryType(id, nil)
	if err != nil {
		return nil, false, err
	}
	for _, t := range types {
		resourceTypes = append(resourceTypes, t.TFType)
	}
	return resourceTypes, exact, err
}
