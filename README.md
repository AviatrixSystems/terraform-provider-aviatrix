Terraform Provider
==================
A basic [Terraform](http://terraform.io) provider for Aviatrix.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Using Aviatrix Provider
---------------------

Clone repository to your $GOPATH: `git clone https://github.com/rakesh568/terraform-provider-aviatrix`

Activate the provider by adding the following to `~/.terraformrc`
```sh
providers {
  "aviatrix" = "/YOUR_GOPATH/bin/terraform-provider-aviatrix"
}
```
Install the aviatrix provider
------------------------------
```sh
cd terraform-provider-aviatrix
make install
```