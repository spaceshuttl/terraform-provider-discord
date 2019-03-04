package discord

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/bwmarrin/discordgo"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Discord Bot token for authenticating with Discord's API",
				DefaultFunc: schema.EnvDefaultFunc("DISCORD_TOKEN", nil),
			},
		},

		//DataSourcesMap: map[string]*schema.Resource{
		//	"discord_guild":      dataSourceDiscordGuild(),
		//	"discord_guild_role": dataSourceDiscordGuildRole(),
		//},

		ResourcesMap: map[string]*schema.Resource{
			"discord_guild":         resourceDiscordGuild(),
			"discord_guild_role":    resourceDiscordGuildRole(),
			"discord_guild_emoji":   resourceDiscordGuildEmoji(),
			"discord_guild_channel": resourceDiscordGuildChannel(),
		},

		ConfigureFunc: discordProviderConfigure,
	}
}

func discordProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	return discordgo.New("Bot " + d.Get("token").(string))
}

var (
	ErrClientNotConfigured = errors.New("discord client not properly configured")
)
