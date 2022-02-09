# azurerm-restapi to azurerm

## Introduction
This tool is used to migrate resources from terraform `azurerm-restapi` provider to `azurerm` provider.

## How to setup?
1. Clone this repo to local.
2. `go install` under project directory.
   
## Command Usage
```
PS C:\Users\henglu\go\src\github.com\ms-henglu\azurerm-restapi-to-azurerm> azurerm-restapi-to-azurerm.exe            
Usage: azurerm-restapi-to-azurerm [--version] [--help] <command> [<args>]

Available commands are:
    migrate    Migrate azurerm-restapi resources to azurerm resources in current working directory
    plan       Show azurerm-restapi resources which can migrate to azurerm resources in current working directory
    version    Displays the version of the migration tool
```

1. Run `azurerm-restapi-to-azurerm plan` under your terraform working directory, 
   it will list all resources that can be migrated from `azurerm-restapi` provider to `azurerm` provider.
   The Terraform addresses listed in file `azurerm-restapi-to-azurerm.ignore` will be ignored during migration.
```
2022/01/25 14:34:46 [INFO] searching azurerm-restapi_resource & azurerm-restapi_patch_resource...
2022/01/25 14:34:55 [INFO]

The tool will perform the following actions:

The following resources will be migrated:
azurerm-restapi_resource.test2 will be replaced with azurerm_automation_account
azurerm-restapi_patch_resource.test will be replaced with azurerm_automation_account

The following resources can't be migrated:
azurerm-restapi_resource.test: input properties not supported: [], output properties not supported: [identity.principalId, identity.type, identity.tenantId]

The following resources will be ignored in migration:
   ```
2. Run `azurerm-restapi-to-azurerm migrate` under your terraform working directory, 
   it will migrate above resources from `azurerm-restapi` provider to `azurerm` provider, 
   both terraform configuration and state.
   The Terraform addresses listed in file `azurerm-restapi-to-azurerm.ignore` will be ignored during migration.
   
## Examples
There're some examples to show the migration results.
1. [case1 - basic](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case1%20-%20basic)
2. [case2 - for_each](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case2%20-%20for_each)
3. [case3 - nested block](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case3%20-%20nested%20block)
4. [case4 - count](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case4%20-%20count)
5. [case5 - nested block patch](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case5%20-%20nested%20block%20patch)
6. [case6 - meta argument](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case6%20-%20meta%20arguments)
7. [case7 - ignore](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case7%20-%20ignore)
   
## Features
- [x] Support resource `azurerm-restapi_resource` migration
- [x] Support resource `azurerm-restapi_patch_resource` migration
- [x] Support meta-argument `for_each`
- [x] Support meta-argument `count`
- [x] Support meta-argument `depends_on`, `lifecycle` and `provisioner`
- [x] Support dependency injection in array and primitive value.
- [x] Support dependency injection in Map and other complicated struct value.
- [x] Support user input when there're multiple/none `azurerm` resource match for the resource id
- [x] Support migration based on `azurerm` provider's property coverage
- [x] Support ignore terraform addresses listed in file `azurerm-restapi-to-azurerm.ignore`
- [ ] Support data source `azurerm-restapi_resource` migration.

## Known limitations
1. References to local variables can't be migrated.
2. Usage of `dynamic` can't be migrated.
3. Patch resource used to manage CMK can't be migrated.
