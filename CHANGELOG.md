## v2.7.0
Target azurerm version: v4.44.0

## v2.6.0
Target azurerm version: v4.39.0

## v2.5.0
Target azurerm version: v4.34.0

## v2.4.0
Target azurerm version: v4.28.0

ENHANCEMENTS:
- Support specifying the provider version used in the migration in the `terraform` block.

## v2.3.0
Target azurerm version: v4.20.0

BUG FIXES:
- Fix a bug that some azurerm resources could not be migrated correctly because the azurerm ID is not the same as the azure resource ID.

## v2.2.0
Target azurerm version: v4.14.0

## v2.1.0

FEATURES:
- Support `-var-file` option to specify the path to the terraform variable file.
- Support migrating resources from `azurerm` provider to `azapi` provider.
- When migrating resources from `azurerm` provider to `azapi` provider, it will generate `moved` block to modify the terraform state.

## v2.0.0-beta

Target azurerm version: v4.0.0

FEATURES:
- The new migration flow uses `import` and `removed` block instead of importing resources and removing resources from terraform state directly.
- Support `working-dir` flag to specify the working directory

## v1.15.0
Target azurerm version: v3.114.0

## v1.14.0
Target azurerm version: v3.110.0

## v1.13.0
Target azurerm version: v3.99.0
Target azapi version: v1.13.0

## v1.12.0
Target azurerm version: v3.83.0

## v1.11.0
Target azurerm version: v3.83.0

## v1.10.0
Target azurerm version: v3.79.0

## v1.9.0
Target azurerm version: v3.71.0

## v1.8.0
Target azurerm version: v3.66.0

## v1.7.0
Target azurerm version: v3.61.0

ENHANCEMENTS:
- Refactor: use `tfadd` to generate config from state

BUG FIXES:
- Fix import with `for_each` statement

## v1.6.0
Target azurerm version: v3.55.0

ENHANCEMENTS:
- Refactor: use aztft to get resource type & upgrade to go 1.19

BUG FIXES:
- Fix import with `count` statement

## v1.5.0
Target azurerm version: v3.50.0

## v1.4.0
Target azurerm version: v3.45.0

## v1.3.0
Target azurerm version: v3.41.0

## v1.2.0
Target azurerm version: v3.37.0

## v1.1.0
Target azurerm version: v3.31.0

## v1.0.0
Target azurerm version: v3.24.0

## v0.6.0
Target azurerm version: v3.22.0

## v0.5.0
Target azurerm version: v3.18.0

FEATURES:
- Refresh state after migrating update resources.

## v0.4.0
Target azurerm version: v3.11.0

## v0.3.0
Target azurerm version: v3.1.0

## v0.2.0
Target azurerm version: v3.0.2

## v0.1.0
Target azurerm version: v2.99.0

FEATURES:
- Support resource `azapi_resource` migration
- Support resource `azapi_update_resource` migration
- Support meta-argument `for_each`
- Support meta-argument `count`
- Support meta-argument `depends_on`, `lifecycle` and `provisioner`
- Support dependency injection in array and primitive value.
- Support dependency injection in Map and other complicated struct value.
- Support user input when there are multiple/none `azurerm` resource match for the resource id
- Support migration based on `azurerm` provider's property coverage
- Support ignore terraform addresses listed in file `azapi2azurerm.ignore`
