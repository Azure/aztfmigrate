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
  name     = "hl1214-resource-group"
  location = "west europe"
}

data "azurerm_client_config" "current" {
}

variable "AutomationName" {
  type    = string
  default = "hl1214"
}

variable "Label" {
  type    = string
  default = "value"
}

locals {
  AutomationSku = "Basic"
}

resource "azurerm_automation_account" "test" {
  location            = azurerm_resource_group.test.location
  name                = var.AutomationName
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_account" "test2" {
  location            = azurerm_resource_group.test.location
  name                = "hl1214another"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_automation_account.test.sku_name
}

resource "azurerm_automation_account" "test1" {
  location            = "westeurope"
  name                = "hl1214-2"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  tags = {
    key = "value"
  }
}

output "accountName" {
  value = azurerm_automation_account.test.name
}

output "patchAccountSKU" {
  value = azurerm_automation_account.test1.sku_name
}

