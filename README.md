# terraform-provider-example

A basic [Terraform](http://terraform.io) provider example.

This repo also seeks to illustrate a feature opportunity within Terraform
wherein a provider resource with a `connection` key of type `TypeList` seemingly yields
in misleading error message, as [`connection` is a reserved word](https://www.terraform.io/docs/provisioners/connection.html);
the key must be named something other than `connection`. In such a scenario,
error output other than `"connection": required field is not set` could be helpful.

## To observe the bug
```
# clone the repo to your $GOPATH:
git clone https://github.com/mdb/terraform-provider-example

# activate the provider by adding the following to `~/.terraformrc`
providers {
  "example" = "/YOUR_GOPATH/bin/terraform-provider-example"
}

# install the example provider
cd terraform-provider-example
make install

# observe that the following outputs an error:
# * example_thing.test_thing: "connection": required field is not set
terraform plan
```

To observe that changing the key name to something other than `connection` --
in this case `thing_connection` -- works as expected:

```
# check out the `works` branch
git checkout works

# install the example provider
cd terraform-provider-example
make install

# observe that the following works as expected:
terraform plan
```
