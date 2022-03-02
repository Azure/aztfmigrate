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
  name     = "henglu-resource-group"
  location = "west europe"
}

data "azurerm_client_config" "current" {
}

variable "accounts" {
  type = map(any)
  default = {
    "item1" = {
      name = "test1"
      sku  = "Basic"
    }
    "item2" = {
      name = "test2"
      sku  = "Basic"
    }
  }
}

resource "azurerm_automation_account" "test" {
  location            = azurerm_resource_group.test.location
  name                = each.value.name
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  for_each = {
    item1 = {
      name = "henglutest1"
    }
    item2 = {
      name = "henglutest2"
    }
  }
}

// some comment
output "sku1" {
  value = azurerm_automation_account.test["item1"].sku_name
}

output "sku2" {
  value = azurerm_automation_account.test["item2"].sku_name
}
