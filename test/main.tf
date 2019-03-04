resource "discord_guild" "test" {
  name = "Terraform Test"
}

resource "discord_guild_channel" "test_channel" {
  guild_id = "${discord_guild.test.id}"
  name = "Welcome"
}

resource "discord_channel_invite" "test_invite" {
  channel_id = "${discord_guild_channel.test_channel.id}"
  unique = true
  max_users = 1
}

output "invite" {
  value = "${discord_channel_invite.test_invite.code}"
}