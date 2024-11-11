package cmd_test

import (
	"fmt"
	"github.com/Azure/aztfmigrate/azurerm"
	"github.com/Azure/aztfmigrate/cmd"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/aztfmigrate/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrate_basic(t *testing.T) {
	migrateTestCase(t, basic(), "azurerm")
}

func TestMigrate_foreach(t *testing.T) {
	migrateTestCase(t, foreach(), "azurerm")
}

func TestMigrate_nestedBlock(t *testing.T) {
	migrateTestCase(t, nestedBlock(), "azurerm")
}

func TestMigrate_count(t *testing.T) {
	migrateTestCase(t, count(), "azurerm")
}

func TestMigrate_nestedBlockUpdate(t *testing.T) {
	migrateTestCase(t, nestedBlockUpdate(), "azurerm")
}

func TestMigrate_metaArguments(t *testing.T) {
	migrateTestCase(t, metaArguments(), "azurerm")
}

func TestMigrate_basicAzureRM(t *testing.T) {
	migrateTestCase(t, basicAzureRM(), "azapi")
}

func TestMigrate_foreachAzureRM(t *testing.T) {
	migrateTestCase(t, foreachAzureRM(), "azapi")
}

func TestMigrate_countAzureRM(t *testing.T) {
	migrateTestCase(t, countAzureRM(), "azapi")
}

func TestMigrate_metaArgumentsAzureRM(t *testing.T) {
	migrateTestCase(t, metaArgumentsAzureRM(), "azapi")
}

func migrateTestCase(t *testing.T, content string, targetProvider string, ignore ...string) {
	if len(os.Getenv("TF_ACC")) == 0 {
		t.Skipf("Set `TF_ACC=true` to enable this test")
	}
	dir := tempDir(t)
	filename := filepath.Join(dir, "main.tf")
	err := os.WriteFile(filename, []byte(`
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
`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dir: %s", dir)
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

	err = os.WriteFile(filename, []byte(content), 0600)
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
	p, err := terraform.Plan(nil)
	if err != nil {
		log.Fatal(err)
	}
	allResources := types.ListResourcesFromPlan(p)
	for index, r := range allResources {
		switch v := r.(type) {
		case *types.AzapiResource:
			resourceId := v.Instances[0].ResourceId
			resourceTypes, _, err := azurerm.GetAzureRMResourceType(resourceId)
			if err != nil {
				t.Fatal(err)
			}
			v.ResourceType = resourceTypes[0]
			allResources[index] = r
		case *types.AzapiUpdateResource:
			resourceId := v.Id
			resourceTypes, _, err := azurerm.GetAzureRMResourceType(resourceId)
			if err != nil {
				t.Fatal(err)
			}
			v.ResourceType = resourceTypes[0]
			allResources[index] = r
		}
	}

	migrateResources := make([]types.AzureResource, 0)
	for _, r := range allResources {
		if r.TargetProvider() == targetProvider {
			migrateResources = append(migrateResources, r)
		}
	}

	migrateCommand := cmd.MigrateCommand{Ui: ui, TargetProvider: targetProvider}
	migrateCommand.MigrateResources(terraform, migrateResources)

	// check generic resources are migrated
	config, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("migration result: \n%s", string(config))
	file, diag := hclwrite.ParseConfig(config, filename, hcl.InitialPos)
	if diag != nil && diag.HasErrors() {
		t.Fatal(diag.Error())
	}
	if file == nil {
		t.Fatal("expect a valid file, but got nil")
	}
	migratedSet := make(map[string]bool)
	for _, r := range migrateResources {
		migratedSet[r.OldAddress(nil)] = true
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
	plan, err := terraform.Plan(nil)
	if err != nil {
		t.Fatal(err)
	}
	if changes := GetChanges(plan); len(changes) != 0 {
		t.Fatalf("expect no plan-diff, but got %v", changes)
	}

}

func tempDir(t *testing.T) string {
	tmpDir := filepath.Join(os.TempDir(), "aztfmigrate", t.Name())
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
	return fmt.Sprintf("acctest%d", rand.Intn(100000))
}

func template() string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azapi" {
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

resource "azapi_resource" "test" {
  name                   = var.AutomationName
  parent_id              = azurerm_resource_group.test.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  response_export_values = ["name", "identity", "properties.sku"]

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = local.AutomationSku
      }
    }
  }
}

resource "azapi_resource" "test2" {
  name      = "${var.AutomationName}another"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  location  = azurerm_resource_group.test.location
  body = {
    properties = {
      sku = {
        name = azapi_resource.test.output.properties.sku.name
      }
    }
  }
}

resource "azurerm_automation_account" "test1" {
  location            = "westeurope"
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  resource_id            = azurerm_automation_account.test1.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  response_export_values = ["properties.sku"]
  body = {
    tags = {
      key = var.Label
    }
  }
}

output "accountName" {
  value = azapi_resource.test.output.name
}

output "patchAccountSKU" {
  value = azapi_update_resource.test.output.properties.sku.name
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


resource "azapi_resource" "test" {
  name      = "henglu${each.value.name}"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = each.value.sku
      }
    }
  }

  for_each = var.accounts
}
`, template(), randomResourceName(), randomResourceName())
}

func nestedBlock() string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                            = "%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "GRS"
  allow_nested_items_to_be_public = false
}

variable "description" {
  type    = string
  default = "this is my desc"
}

variable "defName" {
  type    = string
  default = "def1"
}

resource "azapi_resource" "test" {
  name      = "%s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Network/serviceEndpointPolicies@2020-11-01"

  body = {
    location = "westeurope"
    tags     = {}
    properties = {
      serviceEndpointPolicyDefinitions = [
        {
          name = var.defName
          properties = {
            service     = "Microsoft.Storage"
            description = var.description
            serviceResources = [
              azurerm_storage_account.test.id,
              azurerm_resource_group.test.id
            ]
          }
        }
      ]
    }
  }
}
`, template(), randomResourceName(), randomResourceName())
}

func count() string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  name      = "%s${count.index}"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  location  = azurerm_resource_group.test.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

  count = 2
}
`, template(), randomResourceName())
}

func nestedBlockUpdate() string {
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

resource "azapi_update_resource" "test" {
  resource_id = azurerm_container_registry.test.id
  type        = "Microsoft.ContainerRegistry/registries@2019-05-01"
  body = {
    properties = {
      networkRuleSet = {
        defaultAction = "Deny"
        ipRules = [
          {
            action = var.action
            value  = "7.7.7.7"
          },
          {
            action = var.action
            value  = "2.2.2.2"
          }
        ]
      }
    }
  }
}
`, template(), randomResourceName())
}

func metaArguments() string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  name                   = "%s"
  parent_id              = azurerm_resource_group.test.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  response_export_values = ["name", "identity", "properties.sku"]

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

  depends_on = [azurerm_resource_group.test]

  lifecycle {
    create_before_destroy = false
    prevent_destroy       = false
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}


resource "azapi_resource" "test1" {
  name      = "%s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

  depends_on = [azurerm_resource_group.test, azapi_resource.test]

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

func basicAzureRM() string {
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

resource "azurerm_automation_account" "automationAccount" {
  name                = var.AutomationName
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Basic"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_automation_account" "automationAccount2" {
  name                = "${var.AutomationName}another"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = azurerm_automation_account.automationAccount.sku_name
  identity {
    type = "SystemAssigned"
  }
}

output "accountName" {
  value = azurerm_automation_account.automationAccount.name
}


`, template(), randomResourceName())
}

func foreachAzureRM() string {
	return fmt.Sprintf(`
%s
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

resource "azurerm_automation_account" "test" {
  name                = "henglu${each.value.name}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = each.value.sku
  identity {
    type = "SystemAssigned"
  }

  for_each = var.accounts
}
`, template(), randomResourceName(), randomResourceName())
}

func countAzureRM() string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "test" {
  name                = "%s${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Basic"
  identity {
    type = "SystemAssigned"
  }

  count = 2
}


`, template(), randomResourceName())
}

func metaArgumentsAzureRM() string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "automationAccount" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Basic"
  identity {
    type = "SystemAssigned"
  }

  depends_on = [azurerm_resource_group.test]

  lifecycle {
    create_before_destroy = false
    prevent_destroy       = false
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}
`, template(), randomResourceName())
}
