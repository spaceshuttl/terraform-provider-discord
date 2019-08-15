package discord

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscordMessage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscordMessageCreate,
		Read:   resourceDiscordMessageRead,
		Update: resourceDiscordMessageUpdate,
		Delete: resourceDiscordMessageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"channel_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// =====================================================================================================================
// Discord Guild CRUD Operations
// =====================================================================================================================

// resourceDiscordGuildCreate
func resourceDiscordMessageCreate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	m, err := s.ChannelMessageSend(d.Get("channel_id").(string), d.Get("content").(string))
	if err != nil {
		return err
	}

	d.SetId(m.ID)

	return resourceDiscordMessageRead(d, meta)
}

// resourceDiscordGuildRead
func resourceDiscordMessageRead(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	m, err := s.ChannelMessage(d.Get("channel_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.Set("content", m.Content)

	return nil
}

// resourceDiscordGuildUpdate
func resourceDiscordMessageUpdate(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	_, err := s.ChannelMessageEdit(d.Get("channel_id").(string), d.Id(), d.Get("content").(string))
	if err != nil {
		return err
	}

	return resourceDiscordMessageRead(d, meta)
}

// resourceDiscordGuildDelete
func resourceDiscordMessageDelete(d *schema.ResourceData, meta interface{}) error {
	s, ok := meta.(*discordgo.Session)
	if !ok {
		return ErrClientNotConfigured
	}

	err := s.ChannelMessageDelete(d.Get("channel_id").(string), d.Id())
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to delete member: %s", err))
	}

	return nil
}
