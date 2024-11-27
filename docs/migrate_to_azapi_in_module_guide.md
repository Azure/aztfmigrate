# Guide to migrate AzureRM resources to AzAPi in Module

## Introduction

This guide is intended to help you migrate your existing AzureRM resources to AzApi in Terraform Module. 


## Prerequisites

- Terraform 1.8.0 or later
- Aztfmigrate 2.1.0 or later

## Migration Steps

This guide will walk you through the steps to migrate your existing AzureRM resources to AzApi in Terraform Module. 
It will use [terraform-azurerm-avm-ptn-aks-dev](https://github.com/Azure/terraform-azurerm-avm-ptn-aks-dev) as an example.
The target is to migrate the `azurerm_management_lock` to `azapi_resource` in the module.

### Step 1: Deploy the resources in root module

The migration flow only works with deployed resources. And it only supports the resources in the root module.

Please add the necessary terraform code to deploy the resources in the root module. In the `terraform-azurerm-avm-ptn-aks-dev` example, you can add the following code to deploy the `azurerm_management_lock` resource.

```hcl
// provider.tf

provider "azurerm" {
  features {}
  
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

```

```hcl
// terraform.tfvars

location            = "eastus"
resource_group_name = "heng-aks" // name of an existing resource group
name                = "henglutest"
lock = {    // enable the lock
  kind = "CanNotDelete"
  name = "mylock"
}
```

After adding the code, run `terraform init` and `terraform apply` to deploy the resources.

**Note**: Please use the Azure CLI to authenticate with Azure before running the `terraform` commands.

### Step 2: Migrate the resources in the module

Run `aztfmigrate plan -to azapi` to list all resources that can be migrated from `azurerm` provider to `azapi` provider.

```bash
$ aztfmigrate plan -to azapi

2024/11/27 09:35:55 [INFO] target provider: azapi
2024/11/27 09:35:55 [INFO] initializing terraform...
2024/11/27 09:35:56 [INFO] running terraform plan...
2024/11/27 09:36:28 [INFO]

The tool will perform the following actions:

The following resources will be migrated:
        azurerm_management_lock.this will be replaced with azapi_resource.lock_this
        azurerm_role_assignment.acr will be replaced with azapi_resource.roleAssignment_acr
        azurerm_user_assigned_identity.aks will be replaced with azapi_resource.userAssignedIdentity_aks
        azurerm_container_registry.this will be replaced with azapi_resource.registry_this
        azurerm_kubernetes_cluster.this will be replaced with azapi_resource.managedCluster_this

The following resources can't be migrated:

The following resources will be ignored in migration:

```

Create an ignore file `aztfmigrate.ignore` to include the resources that you don't want to migrate. For this example, you can add the following content to the `aztfmigrate.ignore` file.

```
// aztfmigrate.ignore
azurerm_role_assignment.acr
azurerm_user_assigned_identity.aks
azurerm_container_registry.this
azurerm_kubernetes_cluster.this
```

Run `aztfmigrate plan -to azapi` again to list the resources that can be migrated. It will ignore the resources listed in the `aztfmigrate.ignore` file.

```bash
$ aztfmigrate plan -to azapi

2024/11/27 09:40:59 [INFO] target provider: azapi
2024/11/27 09:40:59 [INFO] initializing terraform...
2024/11/27 09:40:59 [INFO] running terraform plan...
2024/11/27 09:41:25 [INFO]

The tool will perform the following actions:

The following resources will be migrated:
        azurerm_management_lock.this will be replaced with azapi_resource.lock_this

The following resources can't be migrated:

The following resources will be ignored in migration:
        azurerm_role_assignment.acr
        azurerm_user_assigned_identity.aks
        azurerm_container_registry.this
        azurerm_kubernetes_cluster.this
```

Run `aztfmigrate migrate -to azapi` to migrate the resources from `azurerm` provider to `azapi` provider.

```bash
$ aztfmigrate migrate -to azapi

2024/11/27 09:42:14 [INFO] working directory: /Users/luheng/go/src/github.com/Azure/terraform-azurerm-avm-ptn-aks-dev
2024/11/27 09:42:14 [INFO] initializing terraform...
2024/11/27 09:42:14 [INFO] running terraform plan...
2024/11/27 09:42:30 [INFO]

The tool will perform the following actions:

The following resources will be migrated:
        azurerm_management_lock.this will be replaced with azapi_resource.lock_this

The following resources can't be migrated:

The following resources will be ignored in migration:
        azurerm_role_assignment.acr
        azurerm_user_assigned_identity.aks
        azurerm_container_registry.this
        azurerm_kubernetes_cluster.this

2024/11/27 09:42:30 [INFO] generating import config...
2024/11/27 09:42:30 [INFO] migrating resources...
2024/11/27 09:42:30 [INFO] generating new config for resource azurerm_management_lock.this...
2024/11/27 09:42:30 [INFO] generating config...
2024/11/27 09:42:33 [INFO] updating config...
2024/11/27 09:42:33 [INFO] removing azurerm_management_lock.this from config
2024/11/27 09:42:33 [INFO] replacing references with migrated resource...
```

Below is the migrated code block.

```hcl

# # required AVM resources interfaces
# resource "azurerm_management_lock" "this" {
#   count = var.lock != null ? 1 : 0
# 
#   lock_level = var.lock.kind
#   name       = coalesce(var.lock.name, "lock-${var.lock.kind}")
#   scope      = azurerm_kubernetes_cluster.this.id
#   notes      = var.lock.kind == "CanNotDelete" ? "Cannot delete the resource or its child resources." : "Cannot delete or modify the resource or its child resources."
# }
moved {
  from = azurerm_management_lock.this
  to   = azapi_resource.lock_this
}

resource "azapi_resource" "lock_this" {
  count     = 1
  type      = "Microsoft.Authorization/locks@2020-05-01"
  parent_id = azurerm_kubernetes_cluster.this.id
  name      = "mylock"
  body = {
    properties = {
      level = "CanNotDelete"
      notes = "Cannot delete the resource or its child resources."
    }
  }
  ignore_casing             = false
  ignore_missing_property   = true
  schema_validation_enabled = true
}
```

The generated code may not be perfect, you may need to update rest of the code manually. For example, the reference to the original resource block is not updated to the new resource block. You can update the reference manually.

Now if you run `terraform plan` command, you will see below output:

```bash

Terraform will perform the following actions:

  # azurerm_management_lock.this[0] has moved to azapi_resource.lock_this[0]
    resource "azapi_resource" "lock_this" {
        id                        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/heng-aks/providers/Microsoft.ContainerService/managedClusters/aks-henglutest/providers/Microsoft.Authorization/locks/mylock"
        name                      = "mylock"
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 0 to change, 0 to destroy.
```

If you want to migrate other resources, you can repeat the above steps after running `terraform apply` to apply the changes.

### Step 3: Update the migrated code

The migrated code is using literal values. You need to update the code to use expressions and variables like the original code.

```hcl

resource "azapi_resource" "lock_this" {
  count     = var.lock != null ? 1 : 0
  type      = "Microsoft.Authorization/locks@2020-05-01"
  parent_id = azurerm_kubernetes_cluster.this.id
  name      = coalesce(var.lock.name, "lock-${var.lock.kind}")
  body = {
    properties = {
      level = var.lock.kind
      notes =var.lock.kind == "CanNotDelete" ? "Cannot delete the resource or its child resources." : "Cannot delete or modify the resource or its child resources."
    }
  }
  ignore_casing             = false
  ignore_missing_property   = true
  schema_validation_enabled = true
}
```

It is also possible to rename the `azapi_resource` label, for example:

```hcl

moved {
  from = azurerm_management_lock.this
  to   = azapi_resource.cluster_lock // previously azapi_resource.lock_this
}

resource "azapi_resource" "cluster_lock" { // previously azapi_resource.lock_this
  ...
}
```

After updating the code, if you run `terraform plan` command, you will see below output:

```bash

Terraform will perform the following actions:

  # azurerm_management_lock.this[0] has moved to azapi_resource.cluster_lock[0]
    resource "azapi_resource" "cluster_lock" {
        id                        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/heng-aks/providers/Microsoft.ContainerService/managedClusters/aks-henglutest/providers/Microsoft.Authorization/locks/mylock"
        name                      = "mylock"
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 0 to change, 0 to destroy.
```

### Step 4: Clean up

After updating the migrated code, you can destroy the deployed resources by running `terraform destroy` command.

Then you can remove the original code block from the module.

```hcl
// main.tf

// Below code block is the original code, you can remove it after the migration
# resource "azurerm_management_lock" "this" {
#   count = var.lock != null ? 1 : 0
# 
#   lock_level = var.lock.kind
#   name       = coalesce(var.lock.name, "lock-${var.lock.kind}")
#   scope      = azurerm_kubernetes_cluster.this.id
#   notes      = var.lock.kind == "CanNotDelete" ? "Cannot delete the resource or its child resources." : "Cannot delete or modify the resource or its child resources."
# }
```

Also, you can remove the `provider.tf`, `terraform.tfvars` and `aztfmigrate.ignore` files from the module.

## Frequently Asked Questions

### Why do I need to deploy the resources in the root module?

The migration flow generates new azapi resources based on the existing azurerm resources. Currently, the migration flow only supports the resources in the root module and will improve in the future.


### Why do I need to authenticate with Azure CLI?

The migration flow needs to authenticate with Azure to access the existing resources and generate new azapi resource code. The Azure CLI authentication is easy to use in multiple terraform used in the migration flow. It will improve in the future.

### How long does the migration take?

The first migration may take a longer time(1~2 minutes, depending on the network speed) because it needs to download the necessary providers. The subsequent migrations will be faster(about 15s).

### How about migrating multiple resources?

The migration flow supports migrating multiple resources at the same time. You can select multiple resource blocks and click on the light bulb icon to migrate them.


### How about migrating resources from azapi to azurerm?

The `azurerm` provider doesn't support `moved` block yet. The migration will support migrating resources from `azapi` to `azurerm` in the future.