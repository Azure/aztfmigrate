package coverage_test

import (
	"reflect"
	"testing"

	"github.com/Azure/aztfmigrate/azurerm/coverage"
)

func Test_GetCoverage(t *testing.T) {
	testcases := []struct {
		ApiVersion string
		IdPattern  string
		Properties []string
		Covered    []string
		Uncovered  []string
	}{
		{
			ApiVersion: "2015-10-31",
			IdPattern:  "/subscriptions/resourceGroups/providers/Microsoft.Automation/automationAccounts",
			Properties: []string{"location", "tags"},
			Covered:    []string{"location", "tags"},
			Uncovered:  []string{},
		},

		{
			ApiVersion: "2015-10-31",
			IdPattern:  "/subscriptions/resourceGroups/providers/Microsoft.Automation/automationAccounts",
			Properties: []string{"publicNetworkAccess", "tags"},
			Covered:    []string{"tags"},
			Uncovered:  []string{"publicNetworkAccess"},
		},

		{
			ApiVersion: "2015-10-31",
			IdPattern:  "/subscriptions/resourceGroups/providers/Microsoft.Automation/automationAccounts",
			Properties: []string{"properties.sku.name", "tags"},
			Covered:    []string{"properties.sku.name", "tags"},
			Uncovered:  []string{},
		},

		{
			ApiVersion: "2021-07-01",
			IdPattern:  "/subscriptions/resourceGroups/providers/Microsoft.MACHINELEARNINGSERVICES/WORKSPACES/COMPUTES",
			Properties: []string{"properties.computeLocation", "tags"},
			Covered:    []string{"properties.computeLocation", "tags"},
			Uncovered:  []string{},
		},
	}

	for _, testcase := range testcases {
		c, uc := coverage.GetPutCoverage(testcase.Properties, testcase.IdPattern)
		if !reflect.DeepEqual(testcase.Covered, c) || !reflect.DeepEqual(testcase.Uncovered, uc) {
			t.Errorf("expect %v and %v but got %v and %v, testcase: %#v", testcase.Covered, testcase.Uncovered, c, uc, testcase)
		}
	}
}
