Terraform Provider
==================
A basic [Terraform](http://terraform.io) provider for Aviatrix.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (This will be used to build the provider plugin. Check this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to setup GOPATH)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-$PROVIDER_NAME`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-$PROVIDER_NAME
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-$PROVIDER_NAME
$ make build
```

Using Aviatrix Provider
-----------------------

Activate the provider by adding the following to `~/.terraformrc`
```sh
providers {
  "aviatrix" = "/YOUR_GOPATH/bin/terraform-provider-aviatrix"
}
```
Examples
--------
Check examples [here](http://docs.aviatrix.com/HowTos/aviatrix_terraform.html).