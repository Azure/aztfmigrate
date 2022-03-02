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


resource "azapi_resource" "test" {
  name        = "henglu${each.value.name}"
  parent_id   = azurerm_resource_group.test.id
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