# Maintenance status

!> Deprecated! Please use one of the more maintained forks, for example [DrFaust92/bitbucket](https://registry.terraform.io/providers/DrFaust92/bitbucket).

I just took the archived repository, and released it to the Terraform Registry. I'm not a Go developer, so I'm not going to maintain this repository anymore.

If you would like Atlassian to take ownership of the provider, you can voice your opinion at [BCLOUD-17209](https://jira.atlassian.com/browse/BCLOUD-17209).

## History

- Initally created during an Atlassian [24h hackathon](https://www.atlassian.com/company/shipit) by @cwood
- Maintained by HashiCorp and @cwood with contributions from the community
- [Archived](https://www.terraform.io/docs/internals/archiving.html) by HashiCorp due not being maintained (likely related to not being published to the registry)
- Forked and published to [Terraform Registry](https://registry.terraform.io) by @aeirola

Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-bitbucket`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-bitbucket
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-bitbucket
$ make build
```

Using the provider
----------------------

```hcl
# Configure the Bitbucket Provider
provider "bitbucket" {
  username = "GobBluthe"
  password = "idoillusions" # you can also use app passwords
}

# Manage your repository
resource "bitbucket_repository" "infrastructure" {
  owner = "myteam"
  name  = "terraform-code"
}

# Manage your project
resource "bitbucket_project" "infrastructure" {
  owner = "myteam" # must be a team
  name  = "terraform-project"
  key   = "TERRAFORMPROJ"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-bitbucket
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Terraform needs TF_ACC env variable set to run acceptance tests

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

About V1 APIs
------------------

This provider will not take any PRs about the v1 apis that dont have v2
equivalents. Please only focus on v2 apis when adding new featues to this
provider.
