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
  name     = "hl1214-resource-group"
  location = "west europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "hl1214storageacct"
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

resource "azurerm_subnet_service_endpoint_storage_policy" "test" {
  location            = azurerm_resource_group.test.location
  name                = "henglu-policy"
  resource_group_name = azurerm_resource_group.test.name
  definition {
    description       = var.description
    name              = var.defName
    service_resources = [azurerm_resource_group.test.id, azurerm_storage_account.test.id]
  }
}

