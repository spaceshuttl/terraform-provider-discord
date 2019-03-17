package discord

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

func dataSourceDiscordGuildRole() *schema.Resource {
	return &schema.Resource{
		Read:   dataDiscordGuildRoleRead,

		Schema: map[string]*schema.Schema{
			"guild_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mentionable": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hoist": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"color": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"position": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataDiscordGuildRoleRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	guild := d.Get("guild_id").(string)

	body, err := s.RequestWithBucketID(http.MethodGet, discordgo.EndpointGuildRole(guild, d.Id()), nil, discordgo.EndpointGuildRole(guild, ""))
	if err != nil {
		return err
	}

	var role *discordgo.Role
	if err := json.Unmarshal(body, role); err != nil {
		return err
	}

	d.SetId(role.ID)
	d.Set("name", role.Name)
	d.Set("managed", role.Managed)
	d.Set("mentionable", role.Mentionable)
	d.Set("hoist", role.Hoist)
	d.Set("color", role.Color)
	d.Set("position", role.Position)
	d.Set("permissions", role.Permissions)

	return nil
}
