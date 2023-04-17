package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type streamAnalyticsInputsResolver struct{}

func (streamAnalyticsInputsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_stream_analytics_stream_input_eventhub_v2",
		"azurerm_stream_analytics_stream_input_eventhub",
		"azurerm_stream_analytics_stream_input_blob",
		"azurerm_stream_analytics_stream_input_iothub",
		"azurerm_stream_analytics_reference_input_mssql",
		"azurerm_stream_analytics_reference_input_blob",
	}
}

func (streamAnalyticsInputsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewStreamAnalyticsInputsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Input.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	switch props := props.(type) {
	case *armstreamanalytics.StreamInputProperties:
		ds := props.Datasource
		if ds == nil {
			return "", fmt.Errorf("unexpected nil properties.datasource in response")
		}
		switch ds := ds.(type) {
		case *armstreamanalytics.EventHubStreamInputDataSource:
			if ds.Type == nil {
				return "", fmt.Errorf("unexpected nil properties.datasource.type in response")
			}
			switch strings.ToUpper(*ds.Type) {
			case "MICROSOFT.SERVICEBUS/EVENTHUB":
				return "azurerm_stream_analytics_stream_input_eventhub", nil
			case "MICROSOFT.EVENTHUB/EVENTHUB":
				return "azurerm_stream_analytics_stream_input_eventhub_v2", nil
			default:
				return "", fmt.Errorf("unknown properties.datasource.type: %s", *ds.Type)
			}
		case *armstreamanalytics.BlobStreamInputDataSource:
			return "azurerm_stream_analytics_stream_input_blob", nil
		case *armstreamanalytics.IoTHubStreamInputDataSource:
			return "azurerm_stream_analytics_stream_input_iothub", nil
		default:
			return "", fmt.Errorf("unknown input property data source type: %T", ds)
		}
	case *armstreamanalytics.ReferenceInputProperties:
		ds := props.Datasource
		if ds == nil {
			return "", fmt.Errorf("unexpected nil properties.datasource in response")
		}
		switch ds.(type) {
		case *armstreamanalytics.AzureSQLReferenceInputDataSource:
			return "azurerm_stream_analytics_reference_input_mssql", nil
		case *armstreamanalytics.BlobReferenceInputDataSource:
			return "azurerm_stream_analytics_reference_input_blob", nil
		default:
			return "", fmt.Errorf("unknown input property data source type: %T", ds)
		}

	default:
		return "", fmt.Errorf("unknown input property type: %T", props)
	}
}
