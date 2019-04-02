package discord

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordGuildMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordGuildMemberCreate,
		Read:   resourceDiscordGuildMemberRead,
		Update: resourceDiscordGuildMemberUpdate,
		Delete: resourceDiscordGuildMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"guild_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"member_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

// =====================================================================================================================
// Discord Guild CRUD Operations
// =====================================================================================================================

// resourceDiscordGuildCreate
func resourceDiscordGuildMemberCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("member_id").(string))

	return resourceDiscordGuildMemberUpdate(d, meta)
}

// resourceDiscordGuildRead
func resourceDiscordGuildMemberRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	m, err := s.GuildMember(d.Get("guild_id").(string), d.Id())
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read member: %s", err))
	}

	d.Set("roles", m.Roles)

	return nil
}

// resourceDiscordGuildUpdate
func resourceDiscordGuildMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	var roles []string

	roleList := d.Get("roles").([]interface{})

	for _, r := range roleList {
		roles = append(roles, r.(string))
	}

	err := s.GuildMemberEdit(d.Get("guild_id").(string), d.Id(), roles)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to update member: %s", err))
	}

	return resourceDiscordGuildMemberRead(d, meta)
}

// resourceDiscordGuildDelete
func resourceDiscordGuildMemberDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	err := s.GuildMemberEdit(d.Get("guild_id").(string), d.Id(), []string{})
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to delete member: %s", err))
	}

	return nil
}
