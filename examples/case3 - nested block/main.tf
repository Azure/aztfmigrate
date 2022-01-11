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
  name     = "hl1214-resource-group"
  location = "west europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "hl1214storageacct"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

variable "description" {
  type    = string
  default = "this is my desc"
}

variable "defName" {
  type    = string
  default = "def1"
}

resource "azurerm-restapi_resource" "test" {
  name        = "henglu-policy"
  parent_id   = azurerm_resource_group.test.id
  type        = "Microsoft.Network/serviceEndpointPolicies@2020-11-01"

  body = <<BODY
{
    "location": "westeurope",
    "tags": {},
    "properties": {
        "serviceEndpointPolicyDefinitions": [
            {
                "name": "${var.defName}",
                "properties": {
                    "service": "Microsoft.Storage",
                    "description": "${var.description}",
                    "serviceResources": [
                        "${azurerm_storage_account.test.id}",
                        "${azurerm_resource_group.test.id}"
                    ]
                }
            }
        ]
    }
}
  BODY
}