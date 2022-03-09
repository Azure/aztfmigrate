terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "example-resource-group"
  location = "west europe"
}

data "azurerm_client_config" "current" {
}

variable "AutomationName" {
  type    = string
  default = "henglu1"
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
  response_export_values = ["name", "properties.sku"]

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

resource "azapi_resource" "test2" {
  name        = "${var.AutomationName}another"
  parent_id   = azurerm_resource_group.test.id
  type        = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
  location    = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = jsondecode(azapi_resource.test.output).properties.sku.name
      }
    }
  })
}

resource "azurerm_automation_account" "test1" {
  location            = "westeurope"
  name                = "henglu2"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_patch_resource" "test" {
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
  value = jsondecode(azapi_resource.test.output).name
}

output "patchAccountSKU" {
  value = jsondecode(azapi_patch_resource.test.output).properties.sku.name
}
