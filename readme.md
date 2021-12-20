# azurerm-restapi to azurerm

## Introduction
This tool is used to migrate resources from terraform `azurerm-restapi` provider to `azurerm` provider.

## How to use it?
1. Clone this repo to local.
2. `go install` under project directory.
3. Run `azurerm-restapi-to-azurerm.exe` under your terraform working directory, 
   it will migrate all resources from `azurerm-restapi` provider to `azurerm` provider, 
   both terraform configuration and state.
   
## Examples
There're some examples to show the migration results.
1. [case1 - basic](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case1%20-%20basic)
2. [case2 - for_each](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case2%20-%20for_each)
3. [case3 - nested block](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case3%20-%20nested%20block)
4. [case4 - count](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case4%20-%20count)
5. [case5 - nested block patch](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case5%20-%20nested%20block%20patch)
6. [case6 - meta argument](https://github.com/ms-henglu/azurerm-restapi-to-azurerm/tree/master/examples/case6%20-%20meta%20argument)
   
## Features
- [x] Support resource `azurerm-restapi_resource` migration
- [x] Support resource `azurerm-restapi_patch_resource` migration
- [x] Support meta-argument `for_each`
- [x] Support meta-argument `count`
- [x] Support meta-argument `depends_on`, `lifecycle` and `provisioner`
- [x] Support dependency injection in array and primitive value.
- [x] Support dependency injection in Map and other complicated struct value.
- [ ] Support migration based on `azurerm` provider's property coverage
- [ ] Support data source `azurerm-restapi_resource` migration.

## Known limitations
1. References to local variables can't be migrated.
2. Usage of `dynamic` can't be migrated.
