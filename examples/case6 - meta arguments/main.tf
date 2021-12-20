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

resource "azurerm-restapi_resource" "test" {
  resource_id            = "${azurerm_resource_group.test.id}/providers/Microsoft.Automation/automationAccounts/henglu1"
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
    create_before_destroy = true
    prevent_destroy       = true
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}


resource "azurerm-restapi_resource" "test1" {
  resource_id            = "${azurerm_resource_group.test.id}/providers/Microsoft.Automation/automationAccounts/anotherhenglu1"
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
    create_before_destroy = true
    prevent_destroy       = true
  }

  provisioner "local-exec" {
    command = "echo the resource id is ${self.id}"
  }
}