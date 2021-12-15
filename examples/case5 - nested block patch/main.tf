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

resource "azurerm_container_registry" "test" {
  name                = "henglutest"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false

  network_rule_set = [{
    default_action = "Deny"
    ip_rule = [{
      action   = "Allow"
      ip_range = "1.1.1.1/32"
      }, {
      action   = "Allow"
      ip_range = "8.8.8.8/32"
    }]
    virtual_network = []
  }]
}

variable "action" {
  type    = string
  default = "Allow"
}

resource "azurerm-restapi_patch_resource" "test" {
  resource_id = azurerm_container_registry.test.id
  type        = "Microsoft.ContainerRegistry/registries@2019-05-01"
  body        = <<BODY
{
    "properties": {
        "networkRuleSet": {
            "defaultAction": "Deny",
            "ipRules": [
                {
                    "action": "${var.action}",
                    "value": "7.7.7.7"
                },
                {
                    "action": "${var.action}",
                    "value": "2.2.2.2"
                }
            ]
        }
    }
}
    BODY
}