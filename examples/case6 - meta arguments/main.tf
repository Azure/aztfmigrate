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
  name     = "henglu211220-resource-group"
  location = "west europe"
}

resource "azapi_resource" "test" {
  name                   = "henglu1"
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
    create_before_destroy = true
    prevent_destroy       = true
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}


resource "azapi_resource" "test1" {
  name                   = "anotherhenglu1"
  parent_id              = azurerm_resource_group.test.id
  type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

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
    create_before_destroy = true
    prevent_destroy       = true
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}