package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordGuildEmoji() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordGuildEmojiCreate,
		Read:   resourceDiscordGuildEmojiRead,
		Update: resourceDiscordGuildEmojiUpdate,
		Delete: resourceDiscordGuildEmojiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"guild_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Description: "128x128px base64 image data",
				Required:    true,
			},
			"roles": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"animated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"require_colons": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// =====================================================================================================================
// Discord Emoji CRUD Operations
// =====================================================================================================================

func resourceDiscordGuildEmojiCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	emoji, err := s.GuildEmojiCreate(
		d.Get("guild_id").(string),
		d.Get("name").(string),
		d.Get("image").(string),
		d.Get("roles").([]string),
	)
	if err != nil {
		return err
	}

	d.SetId(emoji.ID)
	d.Set("required_colons", emoji.RequireColons)
	d.Set("animated", emoji.Animated)
	d.Set("managed", emoji.Managed)

	return nil
}

func resourceDiscordGuildEmojiRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	guild := d.Get("guild_id").(string)

	// TODO(tobbbles): Contribute back to upstream discordgo and use library function
	resp, err := s.RequestWithBucketID(
		http.MethodGet,
		fmt.Sprint("%sguilds/%s/emojis/%s", discordgo.EndpointAPI, guild, d.Id()),
		nil,
		discordgo.EndpointGuildRoles(guild),
	)
	if err != nil {
		return err
	}

	var emoji *discordgo.Emoji
	if err := json.Unmarshal(resp, emoji); err != nil {
		return err
	}

	d.SetId(emoji.ID)
	d.Set("name", emoji.Name)
	d.Set("managed", emoji.Managed)
	d.Set("roles", emoji.Roles)
	d.Set("animated", emoji.Animated)
	d.Set("require_colons", emoji.RequireColons)

	return nil
}

func resourceDiscordGuildEmojiUpdate(d *schema.ResourceData, meta interface{}) error {
	_, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	return nil
}

func resourceDiscordGuildEmojiDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	guild := d.Get("guild_id").(string)

	log.Printf("[DEBUG] Delete Discord guild %s emoji %s", guild, d.Id())

	return s.GuildEmojiDelete(guild, d.Id())
}
