package loader

import "github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm/types"

type HardcodeDependencyLoader struct {
}

func (h HardcodeDependencyLoader) Load() ([]types.Dependency, error) {
	return []types.Dependency{
		{
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.DBforMariaDB/servers",
			ExampleConfiguration: `
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mariadb_server" "example" {
  name                = "example-mariadb-server"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login          = "mariadbadmin"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_2"
  storage_mb = 5120
  version    = "10.2"

  auto_grow_enabled             = true
  backup_retention_days         = 7
  geo_redundant_backup_enabled  = false
  public_network_access_enabled = true
  ssl_enforcement_enabled       = true
}
`,
			ResourceType:     "azurerm_mariadb_server",
			ReferredProperty: "id",
		},
		{
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.DBforMySQL/flexibleServers",
			ExampleConfiguration: `
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-fs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "psqladmin"
  administrator_password = "H@Sh1CoR3!"
  backup_retention_days  = 7
  sku_name               = "GP_Standard_D2ds_v4"
}
`,
			ResourceType:     "azurerm_mysql_flexible_server",
			ReferredProperty: "id",
		},
		{
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.DBforPostgreSQL/flexibleServers",
			ExampleConfiguration: `
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  administrator_login    = "psqladmin"
  administrator_password = "H@Sh1CoR3!"

  storage_mb = 32768

  sku_name   = "GP_Standard_D4s_v3"

}
`,
			ResourceType:     "azurerm_postgresql_flexible_server",
			ReferredProperty: "id",
		},
		{
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.Sql/servers",
			ExampleConfiguration: `
resource "azurerm_resource_group" "example" {
  name     = "database-rg"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server" "example" {
  name                         = "mssqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.example.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }

  tags = {
    environment = "production"
  }
}
`,
			ResourceType:     "azurerm_mssql_server",
			ReferredProperty: "id",
		},
		{
			Pattern: "/subscriptions/resourceGroups/providers/Microsoft.Synapse/workspaces",
			ExampleConfiguration: `
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  tags = {
    Env = "production"
  }
}
`,
			ResourceType:     "azurerm_synapse_workspace",
			ReferredProperty: "id",
		},
		// override those customer managed key resource which is not a real resource
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_log_analytics_cluster_customer_managed_key",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_eventhub_namespace_customer_managed_key",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_cognitive_account_customer_managed_key",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_storage_account_customer_managed_key",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_kusto_cluster_customer_managed_key",
			ReferredProperty:     "id",
		},
		// override role assignment
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_synapse_role_assignment",
			ReferredProperty:     "id",
		},
		// override all kinds of associations
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_network_interface_nat_rule_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_subnet_route_table_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_network_interface_application_security_group_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_virtual_desktop_workspace_application_group_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_nat_gateway_public_ip_prefix_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_management_group_subscription_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_network_interface_application_gateway_backend_address_pool_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_subnet_nat_gateway_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_subnet_network_security_group_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_network_interface_backend_address_pool_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_network_interface_security_group_association",
			ReferredProperty:     "id",
		},
		{
			Pattern:              "",
			ExampleConfiguration: "",
			ResourceType:         "azurerm_nat_gateway_public_ip_association",
			ReferredProperty:     "id",
		},
	}, nil
}

var _ DependencyLoader = HardcodeDependencyLoader{}
