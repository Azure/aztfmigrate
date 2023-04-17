package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type dataFactoryDatasetsResolver struct{}

func (dataFactoryDatasetsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_data_factory_dataset_postgresql",
		"azurerm_data_factory_dataset_snowflake",
		"azurerm_data_factory_dataset_parquet",
		"azurerm_data_factory_custom_dataset",
		"azurerm_data_factory_dataset_json",
		"azurerm_data_factory_dataset_azure_blob",
		"azurerm_data_factory_dataset_delimited_text",
		"azurerm_data_factory_dataset_cosmosdb_sqlapi",
		"azurerm_data_factory_dataset_sql_server_table",
		"azurerm_data_factory_dataset_http",
		"azurerm_data_factory_dataset_binary",
		"azurerm_data_factory_dataset_mysql",
	}
}

func (dataFactoryDatasetsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDataFactoryDatasetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.DatasetResource.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armdatafactory.AzurePostgreSQLTableDataset:
		return "azurerm_data_factory_dataset_postgresql", nil
	case *armdatafactory.SnowflakeDataset:
		return "azurerm_data_factory_dataset_snowflake", nil
	case *armdatafactory.ParquetDataset:
		return "azurerm_data_factory_dataset_parquet", nil
	case *armdatafactory.CustomDataset:
		return "azurerm_data_factory_custom_dataset", nil
	case *armdatafactory.JSONDataset:
		return "azurerm_data_factory_dataset_json", nil
	case *armdatafactory.AzureBlobDataset:
		return "azurerm_data_factory_dataset_azure_blob", nil
	case *armdatafactory.DelimitedTextDataset:
		return "azurerm_data_factory_dataset_delimited_text", nil
	case *armdatafactory.DocumentDbCollectionDataset:
		return "azurerm_data_factory_dataset_cosmosdb_sqlapi", nil
	case *armdatafactory.SQLServerTableDataset:
		return "azurerm_data_factory_dataset_sql_server_table", nil
	case *armdatafactory.HTTPDataset:
		return "azurerm_data_factory_dataset_http", nil
	case *armdatafactory.BinaryDataset:
		return "azurerm_data_factory_dataset_binary", nil
	case *armdatafactory.MySQLTableDataset:
		return "azurerm_data_factory_dataset_mysql", nil
	default:
		return "", fmt.Errorf("unknown dataset type %T", props)
	}
}
