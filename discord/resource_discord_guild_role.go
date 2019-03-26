package discord

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordGuildRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordGuildRoleCreate,
		Read:   resourceDiscordGuildRoleRead,
		Update: resourceDiscordGuildRoleUpdate,
		Delete: resourceDiscordGuildRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"guild_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Associated Guild ID for the role",
				Required:    true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"color": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  0,
				Optional: true,
			},
			"hoist": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"mentionable": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// =====================================================================================================================
// Discord Guild CRUD Operations
// =====================================================================================================================

type RoleCreate struct {
	Name        string `json:"name"`
	Permissions int    `json:"permissions"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Mentionable bool   `json:"mentionable"`
}

func resourceDiscordGuildRoleCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	guild := d.Get("guild_id").(string)

	rc := &RoleCreate{
		Name:        d.Get("name").(string),
		Permissions: d.Get("permissions").(int),
		Color:       d.Get("color").(int),
		Hoist:       d.Get("hoist").(bool),
		Mentionable: d.Get("mentionable").(bool),
	}

	resp, err := s.RequestWithBucketID("POST", discordgo.EndpointGuildRoles(guild), rc, discordgo.EndpointGuildRoles(guild))
	if err != nil {
		return err
	}

	// Marshal results into schema
	var role *discordgo.Role
	if err := json.Unmarshal(resp, role); err != nil {
		return err
	}

	d.SetId(role.ID)

	return resourceDiscordGuildRoleRead(d, meta)
}

func resourceDiscordGuildRoleRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	roles, err := s.GuildRoles(d.Get("guild_id").(string))
	if err != nil {
		return err
	}

	// Match the role by ID
	for _, role := range roles {
		if role.ID != d.Id() {
			continue
		}

		d.Set("name", role.Name)
		d.Set("permissions", role.Permissions)
		d.Set("color", role.Color)
		d.Set("hoist", role.Hoist)
		d.Set("managed", role.Managed)
		d.Set("mentionable", role.Mentionable)
		d.Set("position", role.Position)

		break
	}

	return nil
}

func resourceDiscordGuildRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	_, err := s.GuildRoleEdit(
		d.Get("guild_id").(string), d.Id(), d.Get("name").(string), d.Get("color").(int),
		d.Get("hoist").(bool), d.Get("permissions").(int), d.Get("mention").(bool))
	if err != nil {
		return err
	}

	return resourceDiscordGuildRoleRead(d, meta)
}

func resourceDiscordGuildRoleDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	return s.GuildRoleDelete(
		d.Get("guild_id").(string),
		d.Get("id").(string),
	)
}
