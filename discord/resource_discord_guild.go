package discord

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"net/http"
)

func resourceDiscordGuild() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordGuildCreate,
		Read:   resourceDiscordGuildRead,
		Update: resourceDiscordGuildUpdate,
		Delete: resourceDiscordGuildDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the guild (2-100 characters)",
				ValidateFunc: func(interface{}, string) ([]string, []error) {
					return nil, nil
				},
				Sensitive: false,
				Required:  true,
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Voice region ID",
				Optional:    true,
			},
			"icon": &schema.Schema{
				Type:        schema.TypeString,
				Description: "base64 128x128 jpeg image for the guild icon",
				Optional:    true,
			},
			"verification_level": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Verification level",
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 4),
				Optional:     true,
			},
			"default_message_notifications": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Default message notification level",
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1),
				Optional:     true,
			},
			"explicit_content_filter": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "explicit content filter level",
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 2),
				Optional:     true,
			},
			"afk_channel_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the voice AFK channel",
				Computed:    true,
				Optional:    true,
			},
			"embed_channel_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the channel welcome messages go to",
				Computed:    true,
				Optional:    true,
			},
			"owner_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the server owner",
				Computed:    true,
			},
			"splash": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Hash of the guild's splash",
				Optional:    true,
			},
			"afk_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Timeout, in seconds, before a user is considered AFK",
				Optional:    true,
				Default:     300,
			},
			"member_count": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The number of members in the server",
				Computed:    true,
			},
			"embed_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether the guild has embedding enabled.",
				Computed:    true,
			},
			"large": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether the guild is considered large.",
				Computed:    true,
			},
			//"roles": &schema.Schema{
			//	Type:     schema.TypeList,
			//	Elem:     resourceDiscordGuildRole(),
			//	Optional: true,
			//},
			//"channels": &schema.Schema{
			//	Type:     schema.TypeList,
			//	Elem:     resourceDiscordGuildChannel(),
			//	Optional: true,
			//},
			//"emojis": &schema.Schema{
			//	Type:     schema.TypeList,
			//	Elem:     resourceDiscordGuildEmoji(),
			//	Optional: true,
			//},
		},
	}
}

// =====================================================================================================================
// Discord Guild CRUD Operations
// =====================================================================================================================

type GuildCreate struct {
	Name                        string `json:"name"`
	Region                      string `json:"region"`
	Icon                        string `json:"icon,omitempty"`
	VerificationLevel           int    `json:"verification_level,omitempty"`
	DefaultMessageNotifications int    `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       int    `json:"explicit_content_filter,omitempty"`
}

// resourceDiscordGuildCreate
func resourceDiscordGuildCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	gc := &GuildCreate{
		Name:                        d.Get("name").(string),
		Region:                      d.Get("region").(string),
		Icon:                        d.Get("icon").(string),
		VerificationLevel:           d.Get("verification_level").(int),
		DefaultMessageNotifications: d.Get("default_message_notifications").(int),
		ExplicitContentFilter:       d.Get("explicit_content_filter").(int),
	}

	resp, err := s.RequestWithBucketID(http.MethodPost, discordgo.EndpointGuildCreate, gc, discordgo.EndpointGuildCreate)
	if err != nil {
		return err
	}

	guild := &discordgo.Guild{}
	if err := json.Unmarshal(resp, guild); err != nil {
		return err
	}

	d.SetId(guild.ID)

	// Lets clean up the default channels as they don't really make in this setup
	channels, err := s.GuildChannels(guild.ID)
	for i := 0; i < len(channels); i++ {
		_, _ = s.ChannelDelete(channels[i].ID)
	}

	return resourceDiscordGuildRead(d, meta)
}

// resourceDiscordGuildRead
func resourceDiscordGuildRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	g, err := s.Guild(d.Id())
	if err != nil {
		return err
	}

	fmt.Println(g.AfkChannelID)

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

	return nil
}

// resourceDiscordGuildUpdate
func resourceDiscordGuildUpdate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	var vl discordgo.VerificationLevel

	switch d.Get("verification_level").(int) {
	case 0:
		vl = discordgo.VerificationLevelNone
		break
	case 1:
		vl = discordgo.VerificationLevelLow
		break
	case 2:
		vl = discordgo.VerificationLevelMedium
		break
	case 3:
		vl = discordgo.VerificationLevelHigh
	}

	data := discordgo.GuildParams{
		Name:                        d.Get("name").(string),
		Region:                      d.Get("region").(string),
		VerificationLevel:           &vl,
		DefaultMessageNotifications: d.Get("default_message_notifications").(int),
		AfkChannelID:                d.Get("afk_channel_id").(string),
		AfkTimeout:                  d.Get("afk_timeout").(int),
		Icon:                        d.Get("icon").(string),
		OwnerID:                     d.Get("owner_id").(string),
		Splash:                      d.Get("splash").(string),
	}

	_, err := s.GuildEdit(d.Id(), data)
	if err != nil {
		return err
	}

	return resourceDiscordGuildRead(d, meta)
}

// resourceDiscordGuildDelete
func resourceDiscordGuildDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	log.Printf("[DEBUG] Delete Discord guild %s", d.Id())

	_, err := s.RequestWithBucketID("DELETE", discordgo.EndpointGuild(d.Id()), nil, discordgo.EndpointGuild(d.Id()))
	return err
}
