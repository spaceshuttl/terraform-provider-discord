package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDiscordGuild() *schema.Resource {
	return &schema.Resource{
		Read:   dataDiscordGuildRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"icon": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"afk_channel_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"embed_channel_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"splash": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"afk_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"member_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"verification_level": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"embed_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"large": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_message_notifications": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     dataSourceDiscordGuildRole(),
				Computed: true,
			},
		},
	}
}

// resourceDiscordGuildRead
func dataDiscordGuildRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	g, err := s.Guild(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", g.ID)
	d.Set("name", g.Name)
	d.Set("icon", g.Icon)
	d.Set("region", g.Region)
	d.Set("afk_channel_id", g.AfkChannelID)
	d.Set("embed_channel_id", g.EmbedChannelID)
	d.Set("owner_id", g.OwnerID)
	d.Set("joined_at", g.JoinedAt)
	d.Set("splash", g.Splash)
	d.Set("afk_timeout", g.AfkTimeout)
	d.Set("member_count", g.MemberCount)
	d.Set("verification_level", g.VerificationLevel)
	d.Set("embed_enabled", g.EmbedEnabled)
	d.Set("large", g.Large)
	d.Set("default_message_notifications", g.DefaultMessageNotifications)

	// Flattened lists
	d.Set("roles", flattenRoles(g.Roles))
	d.Set("emojis", flattenEmojis(g.Emojis))
	d.Set("channels", flattenChannels(g.Channels))
	d.Set("members", flattenMembers(g.Members))
	return nil
}

func dataDiscordGuildDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	_, err := s.GuildDelete(
		d.Get("guild_id").(string),
	)
	return err
}
