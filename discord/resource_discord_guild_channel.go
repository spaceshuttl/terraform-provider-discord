package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceDiscordGuildChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordGuildChannelCreate,
		Read:   resourceDiscordGuildChannelRead,
		Update: resourceDiscordGuildChannelUpdate,
		Delete: resourceDiscordGuildChannelDelete,
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
			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_message_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsfw": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"icon": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"position": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bitrate": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			// TODO: Implement Recipients
			//"recipients": &schema.Schema{
			//	Type: schema.TypeList,
			//},
			"permission_overwrites": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"role", "member"}, false),
						},
						"allow": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"deny": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
			"user_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"parent_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// =====================================================================================================================
// Discord Channel CRUD Operations
// =====================================================================================================================

func resourceDiscordGuildChannelCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	var overwrites []*discordgo.PermissionOverwrite

	if po, ok := d.GetOk("permission_overwrites"); ok {
		permissionOverwrites := po.([]interface{})

		for _, p := range permissionOverwrites {
			overwrite := p.(map[string]interface{})
			overwrites = append(overwrites, &discordgo.PermissionOverwrite{
				ID:    overwrite["id"].(string),
				Type:  overwrite["type"].(string),
				Allow: overwrite["allow"].(int),
				Deny:  overwrite["deny"].(int),
			})
		}
	}

	data := discordgo.GuildChannelCreateData{
		Name:                 d.Get("name").(string),
		Type:                 discordgo.ChannelType(d.Get("type").(int)),
		Topic:                d.Get("topic").(string),
		Bitrate:              d.Get("bitrate").(int),
		UserLimit:            d.Get("user_limit").(int),
		ParentID:             d.Get("parent_id").(string),
		NSFW:                 d.Get("nsfw").(bool),
		PermissionOverwrites: overwrites,
	}

	channel, err := s.GuildChannelCreateComplex(
		d.Get("guild_id").(string),
		data)
	if err != nil {
		return err
	}

	d.SetId(channel.ID)
	d.Set("name", channel.Name)
	d.Set("type", channel.Type)
	d.Set("topic", channel.Topic)
	d.Set("bitrate", channel.Bitrate)
	d.Set("user_limit", channel.UserLimit)
	d.Set("parent_id", channel.ParentID)
	d.Set("nsfw", channel.NSFW)

	return nil
}

func resourceDiscordGuildChannelRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	channel, err := s.Channel(d.Id())
	if err != nil {
		return err
	}

	d.SetId(channel.ID)

	return nil
}

func resourceDiscordGuildChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	_, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	return nil
}

func resourceDiscordGuildChannelDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}
	guild := d.Get("guild_id").(string)

	log.Printf("[DEBUG] Delete Discord guild %s channel %s", guild, d.Id())

	if _, err := s.ChannelDelete(d.Id()); err != nil {
		return err
	}

	return nil
}
