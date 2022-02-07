package cmd_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/cmd"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/tf"
)

func TestMigrate_basic(t *testing.T) {
	migrateTestCase(t, basic())
}

func TestMigrate_foreach(t *testing.T) {
	migrateTestCase(t, foreach())
}

func TestMigrate_nestedBlock(t *testing.T) {
	migrateTestCase(t, nestedBlock())
}

func TestMigrate_count(t *testing.T) {
	migrateTestCase(t, count())
}

func TestMigrate_nestedBlockPatch(t *testing.T) {
	migrateTestCase(t, nestedBlockPatch())
}

func TestMigrate_metaArguments(t *testing.T) {
	migrateTestCase(t, metaArguments())
}

func migrateTestCase(t *testing.T, content string, ignore ...string) {
	if len(os.Getenv("TF_ACC")) == 0 {
		t.Skipf("Set `TF_ACC=true` to enable this test")
	}
	dir := tempDir(t)
	filename := filepath.Join(dir, "main.tf")
	err := ioutil.WriteFile(filename, []byte(`
terraform {
  required_providers {
    azurerm = {
      version = ">= 2.92.0"
    }
  }
}
provider "azurerm" {
  features {}
}
`), 0644)
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

	_ = terraform.Init()

	err = ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
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
	planCommand := cmd.PlanCommand{Ui: ui}
	resources, patchResources := planCommand.Plan(terraform, false)
	migrateCommand := cmd.MigrateCommand{Ui: ui}
	migrateCommand.MigrateGenericResource(terraform, resources)
	migrateCommand.MigrateGenericPatchResource(terraform, patchResources)

	// check generic resources are migrated
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	file, diag := hclwrite.ParseConfig(config, filename, hcl.InitialPos)
	if diag != nil && diag.HasErrors() {
		t.Fatal(diag.Error())
	}
	if file == nil {
		t.Fatal("expect a valid file, but got nil")
	}
	migratedSet := make(map[string]bool)
	for _, r := range resources {
		migratedSet[r.OldAddress(nil)] = true
	}
	for _, r := range patchResources {
		migratedSet[r.OldAddress()] = true
	}
	ignoreSet := make(map[string]bool)
	for _, r := range ignore {
		ignoreSet[r] = true
	}
	for _, block := range file.Body().Blocks() {
		if block.Type() != "resource" {
			continue
		}
		if len(block.Labels()) != 2 {
			continue
		}
		address := strings.Join(block.Labels(), ".")
		if migratedSet[address] {
			t.Fatalf("expect %s to be migrated, but still exist in config, config: \n%s", address, string(config))
		}
	}

	// check no plan-diff
	plan, err := terraform.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if changes := tf.GetChanges(plan); len(changes) != 0 {
		t.Fatalf("expect no plan-diff, but got %v", changes)
	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func tempDir(t *testing.T) string {
	tmpDir := filepath.Join(os.TempDir(), "azurerm-restapi-to-azurerm", t.Name())
	err := os.MkdirAll(tmpDir, 0o755)
	if err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
	})
	return tmpDir
}

func randomResourceName() string {
	return fmt.Sprintf("acctest%d", rand.Intn(10000))
}

func template() string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    azurerm-restapi = {
      source = "Azure/azurerm-restapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azurerm-restapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "west europe"
}
`, randomResourceName())
}

func basic() string {
	return fmt.Sprintf(`
%s
data "azurerm_client_config" "current" {
}

variable "AutomationName" {
  type    = string
  default = "%s"
}

variable "Label" {
  type    = string
  default = "value"
}

locals {
  AutomationSku = "Basic"
}

resource "azurerm-restapi_resource" "test" {
  name                   = var.AutomationName
  parent_id              = azurerm_resource_group.test.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  response_export_values = ["name", "identity", "properties.sku"]

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = local.AutomationSku
      }
    }
  })
}

resource "azurerm-restapi_resource" "test2" {
  name        = "${var.AutomationName}another"
  parent_id   = azurerm_resource_group.test.id
  type        = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  location    = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = jsondecode(azurerm-restapi_resource.test.output).properties.sku.name
      }
    }
  })
}

resource "azurerm_automation_account" "test1" {
  location            = "westeurope"
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm-restapi_patch_resource" "test" {
  resource_id            = azurerm_automation_account.test1.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  response_export_values = ["properties.sku"]
  body = jsonencode({
    tags = {
      key = var.Label
    }
  })
}

output "accountName" {
  value = jsondecode(azurerm-restapi_resource.test.output).name
}

output "patchAccountSKU" {
  value = jsondecode(azurerm-restapi_patch_resource.test.output).properties.sku.name
}
`, template(), randomResourceName(), randomResourceName())
}

func foreach() string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}


variable "accounts" {
  type = map(any)
  default = {
    "item1" = {
      name = "%s"
      sku  = "Basic"
    }
    "item2" = {
      name = "%s"
      sku  = "Basic"
    }
  }
}


resource "azurerm-restapi_resource" "test" {
  name        = "henglu${each.value.name}"
  parent_id   = azurerm_resource_group.test.id
  type        = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = each.value.sku
      }
    }
  })

  for_each = var.accounts
}
`, template(), randomResourceName(), randomResourceName())
}

func nestedBlock() string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

variable "description" {
  type    = string
  default = "this is my desc"
}

variable "defName" {
  type    = string
  default = "def1"
}

resource "azurerm-restapi_resource" "test" {
  name        = "%s"
  parent_id   = azurerm_resource_group.test.id
  type        = "Microsoft.Network/serviceEndpointPolicies@2020-11-01"

  body = <<BODY
{
    "location": "westeurope",
    "tags": {},
    "properties": {
        "serviceEndpointPolicyDefinitions": [
            {
                "name": "${var.defName}",
                "properties": {
                    "service": "Microsoft.Storage",
                    "description": "${var.description}",
                    "serviceResources": [
                        "${azurerm_storage_account.test.id}",
                        "${azurerm_resource_group.test.id}"
                    ]
                }
            }
        ]
    }
}
  BODY
}
`, template(), randomResourceName(), randomResourceName())
}

func count() string {
	return fmt.Sprintf(`
%s

resource "azurerm-restapi_resource" "test" {
  name        = "%s${count.index}"
  parent_id   = azurerm_resource_group.test.id
  type        = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  location    = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  count = 2
}
`, template(), randomResourceName())
}

func nestedBlockPatch() string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false

  network_rule_set = [{
    default_action = "Deny"
    ip_rule = [{
      action   = "Allow"
      ip_range = "1.1.1.1/32"
      }, {
      action   = "Allow"
      ip_range = "8.8.8.8/32"
    }]
    virtual_network = []
  }]
}

variable "action" {
  type    = string
  default = "Allow"
}

resource "azurerm-restapi_patch_resource" "test" {
  resource_id = azurerm_container_registry.test.id
  type        = "Microsoft.ContainerRegistry/registries@2019-05-01"
  body        = <<BODY
{
    "properties": {
        "networkRuleSet": {
            "defaultAction": "Deny",
            "ipRules": [
                {
                    "action": "${var.action}",
                    "value": "7.7.7.7"
                },
                {
                    "action": "${var.action}",
                    "value": "2.2.2.2"
                }
            ]
        }
    }
}
    BODY
}
`, template(), randomResourceName())
}

func metaArguments() string {
	return fmt.Sprintf(`
%s

resource "azurerm-restapi_resource" "test" {
  name                   = "%s"
  parent_id              = azurerm_resource_group.test.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  response_export_values = ["name", "identity", "properties.sku"]

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  depends_on = [azurerm_resource_group.test]

  lifecycle {
    create_before_destroy = false
    prevent_destroy       = false
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}


resource "azurerm-restapi_resource" "test1" {
  name                   = "%s"
  parent_id              = azurerm_resource_group.test.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  depends_on = [azurerm_resource_group.test, azurerm-restapi_resource.test]

  lifecycle {
    create_before_destroy = false
    prevent_destroy       = false
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}
`, template(), randomResourceName(), randomResourceName())
}
