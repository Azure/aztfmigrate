# aztfmigrate

## Introduction
This tool is used to migrate resources between terraform `azapi` provider and `azurerm` provider.


## Command Usage
```
PS C:\Users\henglu\go\src\github.com\Azure\aztfmigrate> aztfmigrate.exe            
Usage: aztfmigrate [--version] [--help] <command> [<args>]

Available commands are:
    migrate    Migrate azapi resources to azurerm resources in current working directory
    plan       Show terraform resources which can be migrated to azurerm or azapi resources in current working directory
    version    Displays the version of the migration tool
```

1. Run `aztfmigrate plan -target-provider=azurerm` under your terraform working directory, 
   it will list all resources that can be migrated from `azapi` provider to `azurerm` provider.
   The Terraform addresses listed in file `aztfmigrate.ignore` will be ignored during migration.
```
2022/01/25 14:34:46 [INFO] searching azapi_resource & azapi_update_resource...
2022/01/25 14:34:55 [INFO]

The tool will perform the following actions:

The following resources will be migrated:
azapi_resource.test2 will be replaced with azurerm_automation_account
azapi_update_resource.test will be replaced with azurerm_automation_account

The following resources can't be migrated:
azapi_resource.test: input properties not supported: [], output properties not supported: [identity.principalId, identity.type, identity.tenantId]

The following resources will be ignored in migration:
   ```
2. Run `aztfmigrate migrate -target-provider=azurerm` under your terraform working directory, 
   it will migrate above resources from `azapi` provider to `azurerm` provider, 
   both terraform configuration and state.
   The Terraform addresses listed in file `aztfmigrate.ignore` will be ignored during migration.
   
## Examples
There're some examples to show the migration results.
1. [case1 - basic](https://github.com/Azure/aztfmigrate/tree/master/examples/case1%20-%20basic)
2. [case2 - for_each](https://github.com/Azure/aztfmigrate/tree/master/examples/case2%20-%20for_each)
3. [case3 - nested block](https://github.com/Azure/aztfmigrate/tree/master/examples/case3%20-%20nested%20block)
4. [case4 - count](https://github.com/Azure/aztfmigrate/tree/master/examples/case4%20-%20count)
5. [case5 - nested block patch](https://github.com/Azure/aztfmigrate/tree/master/examples/case5%20-%20nested%20block%20patch)
6. [case6 - meta argument](https://github.com/Azure/aztfmigrate/tree/master/examples/case6%20-%20meta%20arguments)
7. [case7 - ignore](https://github.com/Azure/aztfmigrate/tree/master/examples/case7%20-%20ignore)


## Install

### From Release

Precompiled binaries and Window MSI are available at [Releases](https://github.com/Azure/aztfmigrate/releases).

For Mac OS users, you need to run the following command to remove the quarantine flag.
```bash
xattr -d com.apple.quarantine aztfmigrate 
```

### From Package Manager

#### dnf (Linux)

Supported versions:

- RHEL 8 (amd64, arm64)
- RHEL 9 (amd64, arm64)

1. Import the Microsoft repository key:

    ```
    rpm --import https://packages.microsoft.com/keys/microsoft.asc
    ```

2. Add `packages-microsoft-com-prod` repository:

    ```
    ver=8 # or 9
    dnf install -y https://packages.microsoft.com/config/rhel/${ver}/packages-microsoft-prod.rpm
    ```

3. Install:

    ```
    dnf install aztfmigrate
    ```

#### apt (Linux)

Supported versions:

- Ubuntu 20.04 (amd64, arm64)
- Ubuntu 22.04 (amd64, arm64)

1. Import the Microsoft repository key:

    ```
    curl -sSL https://packages.microsoft.com/keys/microsoft.asc > /etc/apt/trusted.gpg.d/microsoft.asc
    ```

2. Add `packages-microsoft-com-prod` repository:

    ```
    ver=20.04 # or 22.04
    apt-add-repository https://packages.microsoft.com/ubuntu/${ver}/prod
    ```

3. Install:

    ```
    apt-get install aztfmigrate
    ```

#### AUR (Linux)

```bash
yay -S aztfmigrate
```
   
## Features
- [x] Support resource `azapi_resource` migration
- [x] Support resource `azapi_update_resource` migration
- [x] Support meta-argument `for_each`
- [x] Support meta-argument `count`
- [x] Support meta-argument `depends_on`, `lifecycle` and `provisioner`
- [x] Support dependency injection in array and primitive value.
- [x] Support dependency injection in Map and other complicated struct value.
- [x] Support user input when there're multiple/none `azurerm` resource match for the resource id
- [x] Support migration based on `azurerm` provider's property coverage
- [x] Support ignore terraform addresses listed in file `aztfmigrate.ignore`
- [ ] Support data source `azapi_resource` migration.

## Known limitations
1. References to local variables can't be migrated.
2. Usage of `dynamic` can't be migrated.
3. Update resource used to manage CMK can't be migrated.

## Credits

We wish to thank HashiCorp for the use of some MPLv2-licensed code from their open source project [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk).
