package cmd_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/cmd"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/tf"
)

func TestPlan_basic(t *testing.T) {
	planTestCase(t, basic(), []string{"azapi_resource.test2", "azapi_patch_resource.test"}, false)
}

func TestPlan_foreach(t *testing.T) {
	planTestCase(t, foreach(), []string{"azapi_resource.test"}, false)
}

func TestPlan_nestedBlock(t *testing.T) {
	planTestCase(t, nestedBlock(), []string{"azapi_resource.test"}, false)
}

func TestPlan_count(t *testing.T) {
	planTestCase(t, count(), []string{"azapi_resource.test"}, false)
}

func TestPlan_nestedBlockPatch(t *testing.T) {
	planTestCase(t, nestedBlockPatch(), []string{"azapi_patch_resource.test"}, false)
}

func TestPlan_metaArguments(t *testing.T) {
	planTestCase(t, metaArguments(), []string{"azapi_resource.test1"}, false)
}

func TestPlan_strictMode(t *testing.T) {
	planTestCase(t, basic(), []string{}, true)
}

func planTestCase(t *testing.T, content string, expectMigratedAddresses []string, strictMode bool) {
	if len(os.Getenv("TF_ACC")) == 0 {
		t.Skipf("Set `TF_ACC=true` to enable this test")
	}
	dir := tempDir(t)
	filename := filepath.Join(dir, "main.tf")
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	terraform, err := tf.NewTerraform(dir, false, false)
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
	planCommand := cmd.PlanCommand{Ui: ui, Strict: strictMode}
	resources, patchResources := planCommand.Plan(terraform, true)

	expectSet := make(map[string]bool)
	for _, value := range expectMigratedAddresses {
		expectSet[value] = true
	}
	for _, r := range resources {
		if !expectSet[r.OldAddress(nil)] {
			t.Fatalf("expect %s not migrated, but got it migrated", r.OldAddress(nil))
		}
	}
	for _, r := range patchResources {
		if !expectSet[r.OldAddress()] {
			t.Fatalf("expect %s not migrated, but got it migrated", r.OldAddress())
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
