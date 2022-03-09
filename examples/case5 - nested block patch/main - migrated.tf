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

resource "azurerm_container_registry" "test" {
  name                = "henglutest"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false

  network_rule_set = [
    {
      default_action = "Deny"
      ip_rule = [
        {
          action   = var.action
          ip_range = "7.7.7.7/32"
        },
        {
          action   = var.action
          ip_range = "2.2.2.2/32"
        }
      ]
      virtual_network = []
    }
  ]
}

variable "action" {
  type    = string
  default = "Allow"
}

