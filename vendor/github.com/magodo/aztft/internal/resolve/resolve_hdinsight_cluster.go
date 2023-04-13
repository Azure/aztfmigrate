package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type hdInsightClustersResolver struct{}

func (hdInsightClustersResolver) ResourceTypes() []string {
	return []string{
		"azurerm_hdinsight_kafka_cluster",
		"azurerm_hdinsight_hadoop_cluster",
		"azurerm_hdinsight_spark_cluster",
		"azurerm_hdinsight_hbase_cluster",
		"azurerm_hdinsight_interactive_query_cluster",
	}
}

func (hdInsightClustersResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewHDInsightClustersClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Cluster.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	definition := props.ClusterDefinition
	if definition == nil {
		return "", fmt.Errorf("unexpected nil properties.clusterDefinition in response")
	}
	kind := definition.Kind
	if kind == nil {
		return "", fmt.Errorf("unexpected nil properties.clusterDefinition.kind in response")
	}

	switch strings.ToUpper(*kind) {
	case "KAFKA":
		return "azurerm_hdinsight_kafka_cluster", nil
	case "HADOOP":
		return "azurerm_hdinsight_hadoop_cluster", nil
	case "SPARK":
		return "azurerm_hdinsight_spark_cluster", nil
	case "HBASE":
		return "azurerm_hdinsight_hbase_cluster", nil
	case "INTERACTIVEHIVE":
		return "azurerm_hdinsight_interactive_query_cluster", nil
	default:
		return "", fmt.Errorf("unknown cluster kind: %s", *kind)
	}
}
