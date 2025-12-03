package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/oracledatabase/armoracledatabase"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type oracleAutonomousDatabaseResolver struct{}

func (oracleAutonomousDatabaseResolver) ResourceTypes() []string {
	return []string{
		"azurerm_oracle_autonomous_database",
		"azurerm_oracle_autonomous_database_clone_from_database",
		"azurerm_oracle_autonomous_database_clone_from_backup",
	}
}

func (oracleAutonomousDatabaseResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewOracleAutonomousDatabaseClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}
	switch props.(type) {
	case *armoracledatabase.AutonomousDatabaseProperties:
		return "azurerm_oracle_autonomous_database", nil
	case *armoracledatabase.AutonomousDatabaseCloneProperties:
		return "azurerm_oracle_autonomous_database_clone_from_database", nil
	case *armoracledatabase.AutonomousDatabaseFromBackupTimestampProperties:
		return "azurerm_oracle_autonomous_database_clone_from_backup", nil
	default:
		return "", fmt.Errorf("unknown database properties type: %T", props)
	}
}
