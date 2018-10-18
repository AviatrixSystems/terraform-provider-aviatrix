Terraform Provider
==================
A basic [Terraform](http://terraform.io) provider for Aviatrix. Read this [tutorial](https://docs.aviatrix.com/HowTos/tf_aviatrix_howto.html) as an alternative to the README, only if the instructions are unclear.

Requirements
------------

-	Install [Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	Install [Go](https://golang.org/doc/install) 1.8 (This will be used to build the provider plugin.) 
-	Create a directory, go, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit ~/.bash_profile to setup the GOPATH environment variable)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-aviatrix`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix
```

To clone on windows
```sh
mkdir %GOPATH%\src\github.com\terraform-providers\terraform-provider-aviatrix
cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-aviatrix
git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-aviatrix
$ make fmt
$ make build
```

To build on Windows
```sh
cd %GOPATH%\src\github.com\terraform-providers\terraform-provider-aviatrix
go fmt
go install
```

Using Aviatrix Provider
-----------------------

Activate the provider by adding the following to `~/.terraformrc` on Linux/Unix.
```sh
providers {
  "aviatrix" = "$GOPATH/bin/terraform-provider-aviatrix"
}
```
For Windows, the file should be at '%APPDATA%\terraform.rc'. Do not change $GOPATH to %GOPATH%

If the file is not present, it should be created

Examples
--------

Check examples [here](http://docs.aviatrix.com/HowTos/aviatrix_terraform.html).
