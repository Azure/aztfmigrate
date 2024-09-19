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

data "azurerm_client_config" "current" {
}

variable "AutomationName" {
  type    = string
  default = "hl1214"
}

variable "Label" {
  type    = string
  default = "value"
}

locals {
  AutomationSku = "Basic"
}

# resource "azapi_resource" "test" {
#   name                   = var.AutomationName
#   parent_id              = azurerm_resource_group.test.id
#   type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
#   response_export_values = ["name", "identity", "properties.sku"]
# 
#   location = azurerm_resource_group.test.location
#   identity {
#     type = "SystemAssigned"
#   }
# 
#   body = {
#     properties = {
#       sku = {
#         name = local.AutomationSku
#       }
#     }
#   }
# }
# 

removed {
  from = azapi_resource.test
  lifecycle {
    destroy = false
  }
}

import {
  id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hl1214-resource-group/providers/Microsoft.Automation/automationAccounts/hl1214"
  to = azurerm_automation_account.test
}

resource "azurerm_automation_account" "test" {
  location            = azurerm_resource_group.test.location
  name                = var.AutomationName
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  identity {
    type = "SystemAssigned"
  }
}

# resource "azapi_resource" "test2" {
#   name      = "${var.AutomationName}another"
#   parent_id = azurerm_resource_group.test.id
#   type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
#   location  = azurerm_resource_group.test.location
#   body = {
#     properties = {
#       sku = {
#         name = azapi_resource.test.output.properties.sku.name
#       }
#     }
#   }
# }
# 
removed {
  from = azapi_resource.test2
  lifecycle {
    destroy = false
  }
}

import {
  id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest3348/providers/Microsoft.Automation/automationAccounts/acctest4893another"
  to = azurerm_automation_account.test2
}

resource "azurerm_automation_account" "test2" {
  location            = azurerm_resource_group.test.location
  name                = "hl1214another"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = azurerm_automation_account.test.sku_name
}

resource "azurerm_automation_account" "test1" {
  location            = "westeurope"
  name                = "hl1214-2"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
  tags = {
    key = var.Label
  }
}

# resource "azapi_update_resource" "test" {
#   resource_id            = azurerm_automation_account.test1.id
#   type                   = "Microsoft.Automation/automationAccounts@2020-01-13-preview"
#   response_export_values = ["properties.sku"]
#   body = {
#     tags = {
#       key = var.Label
#     }
#   }
# }
# 
removed {
  from = azapi_update_resource.test
  lifecycle {
    destroy = false
  }
}

output "accountName" {
  value = azurerm_automation_account.test.name
}

output "patchAccountSKU" {
  value = azurerm_automation_account.test1.sku_name
}
