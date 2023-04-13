package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datashare/armdatashare"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type datashareDatasetsResolver struct{}

func (datashareDatasetsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_data_share_dataset_kusto_cluster",
		"azurerm_data_share_dataset_data_lake_gen2",
		"azurerm_data_share_dataset_kusto_database",
		"azurerm_data_share_dataset_blob_storage",
	}
}

func (datashareDatasetsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewDatashareDatasetsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], id.Names()[2], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	model := resp.DataSetClassification
	if model == nil {
		return "", fmt.Errorf("unexpected nil model in response")
	}
	switch model.(type) {
	case *armdatashare.KustoClusterDataSet:
		return "azurerm_data_share_dataset_kusto_cluster", nil
	case *armdatashare.ADLSGen2FileDataSet:
		return "azurerm_data_share_dataset_data_lake_gen2", nil
	case *armdatashare.KustoDatabaseDataSet:
		return "azurerm_data_share_dataset_kusto_database", nil
	case *armdatashare.BlobDataSet:
		return "azurerm_data_share_dataset_blob_storage", nil
	default:
		return "", fmt.Errorf("unknown dataset type: %T", model)
	}
}
