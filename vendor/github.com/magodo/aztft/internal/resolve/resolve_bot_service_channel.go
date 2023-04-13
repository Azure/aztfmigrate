package resolve

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/botservice/armbotservice"
	"github.com/magodo/armid"
	"github.com/magodo/aztft/internal/client"
)

type botServiceChannelsResolver struct{}

func (botServiceChannelsResolver) ResourceTypes() []string {
	return []string{
		"azurerm_bot_channel_directline",
		"azurerm_bot_channel_sms",
		"azurerm_bot_channel_line",
		"azurerm_bot_channel_alexa",
		"azurerm_bot_channel_direct_line_speech",
		"azurerm_bot_channel_slack",
		"azurerm_bot_channel_facebook",
		"azurerm_bot_channel_email",
		"azurerm_bot_channel_ms_teams",
		"azurerm_bot_channel_web_chat",
	}
}

func (botServiceChannelsResolver) Resolve(b *client.ClientBuilder, id armid.ResourceId) (string, error) {
	resourceGroupId := id.RootScope().(*armid.ResourceGroup)
	client, err := b.NewBotServiceChannelsClient(resourceGroupId.SubscriptionId)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(context.Background(), resourceGroupId.Name, id.Names()[0], id.Names()[1], nil)
	if err != nil {
		return "", fmt.Errorf("retrieving %q: %v", id, err)
	}
	props := resp.BotChannel.Properties
	if props == nil {
		return "", fmt.Errorf("unexpected nil properties in response")
	}

	switch props.(type) {
	case *armbotservice.DirectLineChannel:
		return "azurerm_bot_channel_directline", nil
	case *armbotservice.SmsChannel:
		return "azurerm_bot_channel_sms", nil
	case *armbotservice.LineChannel:
		return "azurerm_bot_channel_line", nil
	case *armbotservice.AlexaChannel:
		return "azurerm_bot_channel_alexa", nil
	case *armbotservice.DirectLineSpeechChannel:
		return "azurerm_bot_channel_direct_line_speech", nil
	case *armbotservice.SlackChannel:
		return "azurerm_bot_channel_slack", nil
	case *armbotservice.FacebookChannel:
		return "azurerm_bot_channel_facebook", nil
	case *armbotservice.EmailChannel:
		return "azurerm_bot_channel_email", nil
	case *armbotservice.MsTeamsChannel:
		return "azurerm_bot_channel_ms_teams", nil
	case *armbotservice.WebChatChannel:
		return "azurerm_bot_channel_web_chat", nil
	default:
		return "", fmt.Errorf("unknown bot channel: %T", props)
	}
}
