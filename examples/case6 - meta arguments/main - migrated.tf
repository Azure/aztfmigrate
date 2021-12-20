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
  name     = "henglu211220-resource-group"
  location = "west europe"
}

resource "azurerm_automation_account" "test" {
  location            = azurerm_resource_group.test.location
  name                = "henglu1"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  depends_on          = [azurerm_resource_group.test]
  lifecycle {
    create_before_destroy = true
    prevent_destroy       = true
  }
  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}

resource "azurerm_automation_account" "test1" {
  location            = azurerm_resource_group.test.location
  name                = "anotherhenglu1"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  depends_on          = [azurerm_resource_group.test, azurerm_automation_account.test]
  lifecycle {
    create_before_destroy = true
    prevent_destroy       = true
  }
  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}

