package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v3"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type machineLearningDataStoresResolver struct{}

func (machineLearningDataStoresResolver) ResourceTypes() []string {
	return []string{
		"azurerm_machine_learning_datastore_fileshare",
		"azurerm_machine_learning_datastore_blobstorage",
		"azurerm_machine_learning_datastore_datalake_gen2",
	}
}

func (machineLearningDataStoresResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewMachineLearningDataStoreClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Datastore.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props.(type) {
	case *armmachinelearning.AzureBlobDatastore:
		return "azurerm_machine_learning_datastore_blobstorage", nil
	case *armmachinelearning.AzureFileDatastore:
		return "azurerm_machine_learning_datastore_fileshare", nil
	case *armmachinelearning.AzureDataLakeGen2Datastore:
		return "azurerm_machine_learning_datastore_datalake_gen2", nil
	default:
		return "", fmt.Errorf("unknown data store resource type %T", props)
	}
}
