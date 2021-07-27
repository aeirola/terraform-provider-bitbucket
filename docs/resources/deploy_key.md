--
layout: "bitbucket"
page_title: "Bitbucket: deploy_key"
sidebar_current" "docs-bitbucket-resource-deploy-key"
description: |-
  Add a Deploy Key to a Repository
---

# bitbucket\_deploy\_key

Add a [Deploy Key](https://support.atlassian.com/bitbucket-cloud/docs/add-access-keys/)
to a Repository.

This resource allows you to manage public keys given read-only access to the repository,
such as servers or pipelines.

## Example Usage
```hcl
# Add a repository and public key
resource "bitbucket_repository" "demo" {
  owner = "myteam"
  name = "example-project"
}

resource "tls_private_key" "my_key" {
  algorithm = "RSA"
  rsa_bits = 4096
}

reosurce "bitbucket_deploy_key" {
  name = "My deploy key"
  key = tls_private_key.my_key.public_key_openssh
}
```

Note that keys cannot be updated via the Bitbucket API, so changing any property
will cause the key to be re-created.

## Argument Reference

* `name` - (Required) The name given to this Access Key
* `key` - (Required) The public key, in OpenSSH format.
* `owner` - (Optional) An additional text comment to assign to the key. Only visible via the Bitbucket API.

## Import

Deploy keys have a numeric ID that can be found using the Bitbucket API. Given this ID,
a key can be imported as
```
$ terraform import bitbucket_deploy_key.my_key my_account/my_repo:<KEY_ID>
```
