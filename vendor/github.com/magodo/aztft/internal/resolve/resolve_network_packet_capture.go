package resolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type networkPacketCaptureResolver struct{}

func (networkPacketCaptureResolver) ResourceTypes() []string {
	return []string{"azurerm_virtual_machine_scale_set_packet_capture", "azurerm_virtual_machine_packet_capture"}
}

func (networkPacketCaptureResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewNetworkPacketCapturesClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.PacketCaptureResult.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil property in response")
	}
	targetId := props.Target
	if targetId == nil {
		return "", fmt.Errorf("unexpected nil target id in response")
	}

	tid, err := armid.ParseResourceId(*targetId)
	if err != nil {
		return "", fmt.Errorf("parsing target id %q: %v", *targetId, err)
	}

	if len(tid.Types()) != 1 {
		return "", fmt.Errorf("un-supported resource types for this target id: %v", tid.Types())
	}

	switch rt := strings.ToUpper(tid.Types()[0]); rt {
	case "VIRTUALMACHINESCALESETS":
		return "azurerm_virtual_machine_scale_set_packet_capture", nil
	case "VIRTUALMACHINES":
		return "azurerm_virtual_machine_packet_capture", nil
	default:
		return "", fmt.Errorf("unknown resource type: %s", rt)
	}
}
