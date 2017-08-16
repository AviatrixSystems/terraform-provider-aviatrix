# terraform-provider-aviatrix

A basic [Terraform](http://terraform.io) provider for Aviatrix.

## To observe the bug
```
# clone the repo to your $GOPATH:
git clone https://github.com/rakesh568/terraform-provider-aviatrix

# activate the provider by adding the following to `~/.terraformrc`
providers {
  "aviatrix" = "/YOUR_GOPATH/bin/terraform-provider-aviatrix"
}

# install the aviatrix provider
cd terraform-provider-aviatrix
make install

# install the aviatrix provider
cd terraform-provider-aviatrix
make install

# observe that the following works as expected:
terraform plan
```
