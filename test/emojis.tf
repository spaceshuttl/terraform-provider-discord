resource "discord_guild_emoji" "test_emoji" {
  guild_id = "${discord_guild.dsac.id}"
  name = "doggo"
  image = "${file("doggo.b64")}"
}

resource "discord_guild_emoji" "terraform" {
  guild_id = "${discord_guild.dsac.id}"
  name = "terraform"
  image = "${file("terraform.b64")}"
}