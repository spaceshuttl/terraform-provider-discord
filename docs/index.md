# <provider> Discord Provider

Interact with resources through Discord's REST api. Create or modify channels, roles, or emojis in a guild. Or even create new guilds!

## Example Usage

```hcl
# Configure the AWS Provider
provider "discord" {
  token = "341f...f6f9"
}

# Create a Guild
resource "discord_guild" "discord_as_code" {
  name = "Discord as Code"
  region = "us-central"
  icon = "${file("terraform.b64")}" # A locl base64 incoded image 
}

```

## Argument Reference

* `token` - (Required) Discord Bot token for authenticating with Discord's API
