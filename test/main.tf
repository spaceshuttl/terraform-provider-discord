resource "discord_guild" "dsac" {
  name = "Discord as Code"
  region = "us-central"
  icon = "${file("terraform.b64")}"
}

resource "discord_guild_role" "admin" {
  guild_id = "${discord_guild.dsac.id}"
  name = "Council"
  hoist = true
  color = 0x4286f4
}

resource "discord_guild_role" "moderator" {
  guild_id = "${discord_guild.dsac.id}"
  name = "Moderators"
  hoist = true
  color = 0x42f442
  permissions = 230695046 // 0x2 | 0x4 | 0x4000000 | 0x8000000 | 0x80 | 0x2000 | 0x400000 | 0x800000 | 0x1000000
}

resource "discord_guild_channel" "announcements" {
  guild_id = "${discord_guild.dsac.id}"
  name = "Announcements & Info"
  permission_overwrites = [
    {
      id = "${discord_guild_role.admin.id}"
      type = "role"
      allow = 2048 // Send message perms
    },
    {
      id = "${discord_guild.dsac.id}" // Should be the @everyone role
      type = "role"
      deny = 2048
    }
  ]
  type = 4 // Category
  position = 0
}

resource "discord_guild_channel" "about" {
  guild_id = "${discord_guild.dsac.id}"
  name = "about"
  parent_id = "${discord_guild_channel.announcements.id}"
  position = 0
}

resource "discord_guild_channel" "general_category" {
  guild_id = "${discord_guild.dsac.id}"
  name = "General"
  type = 4 // Category
  position = 1
}

resource "discord_guild_channel" "general" {
  guild_id = "${discord_guild.dsac.id}"
  name = "general"
  parent_id = "${discord_guild_channel.general_category.id}"
  position = 1
}

resource "discord_channel_invite" "test_invite" {
  channel_id = "${discord_guild_channel.about.id}"
  unique = true
  max_users = 0
  max_age = 0
}

//resource "discord_guild_member" "kelwing" {
//  guild_id = "${discord_guild.dsac.id}"
//  member_id = "109710323094683648"
//  roles = ["${discord_guild_role.admin.id}"]
//}

output "invite" {
  value = "${discord_channel_invite.test_invite.code}"
}