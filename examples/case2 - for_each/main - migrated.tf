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
      name = "acctest3505"
      sku  = "Basic"
    }
    "item2" = {
      name = "acctest62"
      sku  = "Basic"
    }
  }
}

# resource "azapi_resource" "test" {
#   name      = "henglu${each.value.name}"
#   parent_id = azurerm_resource_group.test.id
#   type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
# 
#   location = azurerm_resource_group.test.location
#   identity {
#     type = "SystemAssigned"
#   }
# 
#   body = {
#     properties = {
#       sku = {
#         name = each.value.sku
#       }
#     }
#   }
# 
#   for_each = var.accounts
# }
# 
removed {
  from = azapi_resource.test
  lifecycle {
    destroy = false
  }
}

import {
  for_each = {
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest4094/providers/Microsoft.Automation/automationAccounts/hengluacctest3505" = "item1"
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest4094/providers/Microsoft.Automation/automationAccounts/hengluacctest62"   = "item2"
  }
  id = each.key
  to = azurerm_automation_account.test[each.value]
}

resource "azurerm_automation_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  location            = azurerm_resource_group.test.location
  name                = each.value.name
  identity {
    type = "SystemAssigned"
  }
  for_each = {
    item1 = {
      name = "hengluacctest3505"
    }
    item2 = {
      name = "hengluacctest62"
    }
  }
}
