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


resource "azurerm-restapi_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.Automation/automationAccounts/henglu${each.value.name}"
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