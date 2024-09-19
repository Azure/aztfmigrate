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

# resource "azapi_resource" "test" {
#   name      = "henglu${count.index}"
#   parent_id = azurerm_resource_group.test.id
#   type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
#   location  = azurerm_resource_group.test.location
#   body = {
#     properties = {
#       sku = {
#         name = "Basic"
#       }
#     }
#   }
# 
#   count = 2
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
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu-resource-group/providers/Microsoft.Automation/automationAccounts/henglu0" = 0
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/henglu-resource-group/providers/Microsoft.Automation/automationAccounts/henglu1" = 1
  }
  id = each.key
  to = azurerm_automation_account.test[each.value]
}

resource "azurerm_automation_account" "test" {
  name                = "henglu${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  location            = azurerm_resource_group.test.location
  count               = 2
}
