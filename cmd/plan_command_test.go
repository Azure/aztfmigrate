package cmd_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/aztfmigrate/cmd"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/mitchellh/cli"
)

func TestPlan_basic(t *testing.T) {
	planTestCase(t, basic(), []string{"azapi_resource.test", "azapi_resource.test2", "azapi_update_resource.test"}, false, "azurerm")
}

func TestPlan_foreach(t *testing.T) {
	planTestCase(t, foreach(), []string{"azapi_resource.test"}, false, "azurerm")
}

func TestPlan_nestedBlock(t *testing.T) {
	planTestCase(t, nestedBlock(), []string{"azapi_resource.test"}, false, "azurerm")
}

func TestPlan_count(t *testing.T) {
	planTestCase(t, count(), []string{"azapi_resource.test"}, false, "azurerm")
}

func TestPlan_nestedBlockUpdate(t *testing.T) {
	planTestCase(t, nestedBlockUpdate(), []string{"azapi_update_resource.test"}, false, "azurerm")
}

func TestPlan_metaArguments(t *testing.T) {
	planTestCase(t, metaArguments(), []string{"azapi_resource.test", "azapi_resource.test1"}, false, "azurerm")
}

func TestPlan_strictMode(t *testing.T) {
	planTestCase(t, basic(), []string{}, true, "azurerm")
}

func planTestCase(t *testing.T, content string, expectMigratedAddresses []string, strictMode bool, targetProvider string) {
	if len(os.Getenv("TF_ACC")) == 0 {
		t.Skipf("Set `TF_ACC=true` to enable this test")
	}
	dir := tempDir(t)
	filename := filepath.Join(dir, "main.tf")
	err := os.WriteFile(filename, []byte(content), 0600)
	if err != nil {
		t.Fatal(err)
	}

	terraform, err := tf.NewTerraform(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = terraform.Destroy()
		if err != nil {
			t.Fatalf("destroy config: %+v", err)
		}
	})

	err = terraform.Apply()
	if err != nil {
		t.Fatalf("apply config: %+v", err)
	}

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			Reader:      os.Stdin,
			ErrorWriter: os.Stderr,
		},
	}
	planCommand := cmd.PlanCommand{Ui: ui, Strict: strictMode, TargetProvider: targetProvider}
	resources := planCommand.Plan(terraform, true)

	expectSet := make(map[string]bool)
	for _, value := range expectMigratedAddresses {
		expectSet[value] = true
	}

	if len(resources) != len(expectMigratedAddresses) {
		t.Errorf("expect %d resources to be migrated, but got %d resources migrated", len(expectMigratedAddresses), len(resources))
	}

	for _, r := range resources {
		if !expectSet[r.OldAddress(nil)] {
			t.Fatalf("expect %s not migrated, but got it migrated", r.OldAddress(nil))
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
