package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordChannelInvite() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordChannelInviteCreate,
		Read:   resourceDiscordChannelInviteRead,
		Update: resourceDiscordChannelInviteUpdate,
		Delete: resourceDiscordChannelInviteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"channel_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the channel of which the invite ought to be for",
				Required:    true,
			},
			"max_age": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Duration of invite in seconds before expiry, or 0 for never",
				Default:     86400,
				Optional:    true,
			},
			"max_users": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Max number of uses or 0 for unlimited",
				Optional:    true,
			},
			"temporary": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether this invite only grants temporary membership",
				Optional:    true,
			},
			"unique": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If true, don't try to reuse a similar invite (useful for creating many unique one time use invites)",
				Optional:    true,
			},

			// Computed
			"code": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Invite code",
				Computed:    true,
			},
			"uses": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"revoked": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDiscordChannelInviteCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	i := discordgo.Invite{
		MaxAge:    d.Get("max_age").(int),
		MaxUses:   d.Get("max_users").(int),
		Temporary: d.Get("temporary").(bool),
		Unique:    d.Get("unique").(bool),
	}

	invite, err := s.ChannelInviteCreate(
		d.Get("channel_id").(string),
		i)
	if err != nil {
		return err
	}

	d.SetId(invite.Code)
	d.Set("code", invite.Code)
	d.Set("uses", invite.Uses)
	d.Set("revoked", invite.Revoked)
	d.Set("created_at", invite.CreatedAt)

	return nil
}

func resourceDiscordChannelInviteRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	invites, err := s.ChannelInvites(d.Get("channel_id").(string))
	if err != nil {
		return err
	}

	var found bool
	for _, invite := range invites {
		if invite.Code == d.Id() {
			found = true
			d.Set("code", invite.Code)
			d.Set("uses", invite.Uses)
			d.Set("revoked", invite.Revoked)
			d.Set("created_at", invite.CreatedAt)
			break
		}
	}

	// Invite was used and thus deleted - mark resource for deletion
	if !found{
		d.SetId("")
	}

	return nil
}

func resourceDiscordChannelInviteUpdate(d *schema.ResourceData, meta interface{}) error {
	_, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	return nil
}

func resourceDiscordChannelInviteDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	_, err := s.InviteDelete(d.Id())
	return err
}
