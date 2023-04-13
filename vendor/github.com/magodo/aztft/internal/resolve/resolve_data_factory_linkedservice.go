package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type dataFactoryLinkedServicesResolver struct{}

func (dataFactoryLinkedServicesResolver) ResourceTypes() []string {
	return []string{
		"azurerm_data_factory_linked_service_azure_sql_database",
		"azurerm_data_factory_linked_service_cosmosdb",
		"azurerm_data_factory_linked_service_azure_table_storage",
		"azurerm_data_factory_linked_service_web",
		"azurerm_data_factory_linked_service_kusto",
		"azurerm_data_factory_linked_service_azure_file_storage",
		"azurerm_data_factory_linked_service_azure_search",
		"azurerm_data_factory_linked_service_azure_databricks",
		"azurerm_data_factory_linked_service_key_vault",
		"azurerm_data_factory_linked_service_postgresql",
		"azurerm_data_factory_linked_service_mysql",
		"azurerm_data_factory_linked_service_data_lake_storage_gen2",
		"azurerm_data_factory_linked_service_sftp",
		"azurerm_data_factory_linked_service_cosmosdb_mongoapi",
		"azurerm_data_factory_linked_service_azure_function",
		"azurerm_data_factory_linked_service_synapse",
		"azurerm_data_factory_linked_service_snowflake",
		"azurerm_data_factory_linked_service_odbc",
		"azurerm_data_factory_linked_service_azure_blob_storage",
		"azurerm_data_factory_linked_service_odata",
		"azurerm_data_factory_linked_service_sql_server",
		"azurerm_data_factory_linked_custom_service",
	}
}

func (dataFactoryLinkedServicesResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataFactoryLinkedServicesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.LinkedServiceResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armdatafactory.AzureSQLDatabaseLinkedService:
		return "azurerm_data_factory_linked_service_azure_sql_database", nil
	case *armdatafactory.CosmosDbLinkedService:
		return "azurerm_data_factory_linked_service_cosmosdb", nil
	case *armdatafactory.AzureTableStorageLinkedService:
		return "azurerm_data_factory_linked_service_azure_table_storage", nil
	case *armdatafactory.WebLinkedService:
		return "azurerm_data_factory_linked_service_web", nil
	case *armdatafactory.AzureDataExplorerLinkedService:
		return "azurerm_data_factory_linked_service_kusto", nil
	case *armdatafactory.AzureFileStorageLinkedService:
		return "azurerm_data_factory_linked_service_azure_file_storage", nil
	case *armdatafactory.AzureSearchLinkedService:
		return "azurerm_data_factory_linked_service_azure_search", nil
	case *armdatafactory.AzureDatabricksLinkedService:
		return "azurerm_data_factory_linked_service_azure_databricks", nil
	case *armdatafactory.AzureKeyVaultLinkedService:
		return "azurerm_data_factory_linked_service_key_vault", nil
	case *armdatafactory.PostgreSQLLinkedService:
		return "azurerm_data_factory_linked_service_postgresql", nil
	case *armdatafactory.MySQLLinkedService:
		return "azurerm_data_factory_linked_service_mysql", nil
	case *armdatafactory.AzureBlobFSLinkedService:
		return "azurerm_data_factory_linked_service_data_lake_storage_gen2", nil
	case *armdatafactory.SftpServerLinkedService:
		return "azurerm_data_factory_linked_service_sftp", nil
	case *armdatafactory.CosmosDbMongoDbAPILinkedService:
		return "azurerm_data_factory_linked_service_cosmosdb_mongoapi", nil
	case *armdatafactory.AzureFunctionLinkedService:
		return "azurerm_data_factory_linked_service_azure_function", nil
	case *armdatafactory.AzureSQLDWLinkedService:
		return "azurerm_data_factory_linked_service_synapse", nil
	case *armdatafactory.SnowflakeLinkedService:
		return "azurerm_data_factory_linked_service_snowflake", nil
	case *armdatafactory.OdbcLinkedService:
		return "azurerm_data_factory_linked_service_odbc", nil
	case *armdatafactory.AzureBlobStorageLinkedService:
		return "azurerm_data_factory_linked_service_azure_blob_storage", nil
	case *armdatafactory.ODataLinkedService:
		return "azurerm_data_factory_linked_service_odata", nil
	case *armdatafactory.SQLServerLinkedService:
		return "azurerm_data_factory_linked_service_sql_server", nil
	default:
		// By default, we return the custom service resource, as it is a general resource that supports all kinds of linked service
		return "azurerm_data_factory_linked_custom_service", nil
	}
}
