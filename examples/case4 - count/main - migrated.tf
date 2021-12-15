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
  name     = "henglu-resource-group"
  location = "west europe"
}

resource "azurerm_automation_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  location            = azurerm_resource_group.test.location
  name                = "henglu${count.index}"
  count               = 2
}

