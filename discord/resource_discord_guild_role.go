package discord

import (
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

func resourceDiscordGuildRoleCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	guild := d.Get("guild_id").(string)

	role, err := s.GuildRoleCreate(guild)
	if err != nil {
		return err
	}

	role, err = s.GuildRoleEdit(guild, role.ID, d.Get("name").(string), d.Get("color").(int),
		d.Get("hoist").(bool), d.Get("permissions").(int), d.Get("mentionable").(bool))
	if err != nil {
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
		d.Get("hoist").(bool), d.Get("permissions").(int), d.Get("mentionable").(bool))
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
		d.Id(),
	)
}
