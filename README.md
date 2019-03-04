Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

`terraform-provider-discord` uses Go moduels, so feel you have [enabled modules](https://github.com/golang/go/wiki/Modules#how-to-use-modules).

```sh
$ git clone git@github.com:tobbbles/terraform-provider-discord
```

Enter the provider directory and build the provider

```sh
$ cd ./terraform-provider-discord
$ make build
```

Using the provider
----------------------

The provider must authenticate with a Discord Application Bot, that can be created from the [Discord Developer Portal](https://discordapp.com/developers/applications/). From here please Create a 'Bot' and utilise the token.

`DISCORD_TOKEN`


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-discord
...
```

### Running tests
TODO