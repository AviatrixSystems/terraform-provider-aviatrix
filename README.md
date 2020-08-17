Terraform Provider
==================
A basic [Terraform](http://terraform.io) provider for Aviatrix. Read this [tutorial](https://docs.aviatrix.com/HowTos/tf_aviatrix_howto.html) as an alternative to the README, only if the instructions are unclear.

-> **NOTE:** This release has a big structure change from release v1.*, please read this [changelist-v2](https://www.terraform.io/docs/providers/aviatrix/guides/feature-changelist-v2.html) first, and change your cloud infrastructures accordingly.

Requirements
------------

-	Install [Terraform](https://www.terraform.io/downloads.html) 0.12.x (0.11.x or lower is incompatible)
-	Install [Go](https://golang.org/doc/install) 1.12+ (This will be used to build the provider plugin.)
-	Create a directory, go, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit ~/.bash_profile to setup the GOPATH environment variable)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-aviatrix`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix.git
```

To clone on windows
```sh
mkdir %GOPATH%\src\github.com\terraform-providers
cd %GOPATH%\src\github.com\terraform-providers
git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix.git
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
For Windows, the file should be at '%APPDATA%\terraform.rc'. Do not change $GOPATH to %GOPATH%.

In Windows, for terraform 0.11.8 and lower use the above text.

In Windows, for terraform 0.11.9 and higher use the following at '%APPDATA%\terraform.rc'
```sh
providers {
  "aviatrix" = "$GOPATH/bin/terraform-provider-aviatrix.exe"
}
```

If the rc file is not present, it should be created

Examples
--------

Check examples and documentation [here](https://www.terraform.io/docs/providers/aviatrix/)

Visit [here](https://github.com/terraform-providers/terraform-provider-aviatrix/tree/master/website/docs/) for the complete documentation for all resources on github


Controller version
------------------
Due to some non-backward compatible changes in REST API not all controller versions are supported.
If you find a branch with the controller version please use that branch
Controller versions older than 3.3 are not supported
For example:
 * UserConnect-3.3 for 3.3.x controller version
 * UserConnect-3.4 for 3.4.x controller version

master branch supports latest controller version but please use the branch specific to your controller version. This is so that you do not go out of sync with the controller. (For instance, the master branch code may get updated to be only 4.2 compatible but your controller may still be running 4.1)

We also recommend you to update to the latest controller version to stay on top of fixes/features.
