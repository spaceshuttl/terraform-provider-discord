package discord

import "github.com/bwmarrin/discordgo"

func flattenRoles(roles []*discordgo.Role) []interface{} {
	var roleList []interface{}

	for _, role := range roles {
		values := map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"managed":     role.Managed,
			"mentionable": role.Mentionable,
			"hoist":       role.Hoist,
			"color":       role.Color,
			"position":    role.Position,
			"permissions": role.Permissions,
		}

		roleList = append(roleList, values)
	}

	return roleList
}

func flattenEmojis(emojis []*discordgo.Emoji) []interface{} {
	var emojiList []interface{}

	for _, emoji := range emojis {
		values := map[string]interface{}{
			"id":              emoji.ID,
			"managed":         emoji.Managed,
			"name":            emoji.Name,
			"animated":        emoji.Animated,
			"required_colons": emoji.RequireColons,
		}

		emojiList = append(emojiList, values)
	}

	return emojiList
}

func flattenChannels(channels []*discordgo.Channel) []interface{} {
	var channelList []interface{}

	for _, channel := range channels {
		values := map[string]interface{}{
			"id":                    channel.ID,
			"guild_id":              channel.GuildID,
			"name":                  channel.Name,
			"topic":                 channel.Topic,
			"type":                  channel.Type,
			"last_message_id":       channel.LastMessageID,
			"nsfw":                  channel.NSFW,
			"icon":                  channel.Icon,
			"position":              channel.Position,
			"bitrate":               channel.Bitrate,
			"recipients":            channel.Recipients,
			"permission_overwrites": channel.PermissionOverwrites,
			"user_limit":            channel.UserLimit,
			"parent_id":             channel.ParentID,
		}

		channelList = append(channelList, values)
	}

	return channelList
}

func flattenMembers(members []*discordgo.Member) []interface{} {
	var memberList []interface{}

	for _, member := range members {
		values := map[string]interface{}{
			"guild_id":  member.GuildID,
			"joined_at": member.JoinedAt,
			"nick":      member.Nick,
			"deaf":      member.Deaf,
			"mute":      member.Mute,
			"roles":     member.Roles,
			"user":      flattenUser(member.User),
		}

		memberList = append(memberList, values)
	}

	return memberList
}

func flattenUser(user *discordgo.User) map[string]interface{} {
	return map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"username":      user.Username,
		"avatar":        user.Avatar,
		"locale":        user.Locale,
		"discriminator": user.Discriminator,
		"token":         user.Token,
		"verified":      user.Verified,
		"mfa_enabled":   user.MFAEnabled,
		"bot":           user.Bot,
	}
}
