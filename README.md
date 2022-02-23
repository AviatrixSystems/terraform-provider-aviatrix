Terraform Provider
==================
A basic [Terraform](http://terraform.io) provider for Aviatrix. Read this [tutorial](https://docs.aviatrix.com/HowTos/tf_aviatrix_howto.html) as an alternative to the README, only if the instructions are unclear.

-> **NOTE:** This release has a big structure change from release v1.*, please read this [changelist-v2](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/feature-changelist-v2) first, and change your cloud infrastructures accordingly.

Requirements
------------

-	Install [Terraform](https://www.terraform.io/downloads.html) 0.12.x/0.13.x/0.14.x/0.15.x (0.11.x or lower is incompatible)
-	Install [Go](https://golang.org/doc/install) 1.16+ (This will be used to build the provider plugin.)
-	Create a directory, go, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit ~/.bash_profile to setup the GOPATH environment variable)

Building The Provider (Terraform v0.12+)
---------------------

Clone repository to: `$GOPATH/src/github.com/AviatrixSystems/terraform-provider-aviatrix`

```sh
$ mkdir -p $GOPATH/src/github.com/AviatrixSystems
$ cd $GOPATH/src/github.com/AviatrixSystems
$ git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix.git
```

To clone on windows
```sh
mkdir %GOPATH%\src\github.com\AviatrixSystems
cd %GOPATH%\src\github.com\AviatrixSystems
git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/AviatrixSystems/terraform-provider-aviatrix
$ make fmt
$ make build
```

To build on Windows
```sh
cd %GOPATH%\src\github.com\AviatrixSystems\terraform-provider-aviatrix
go fmt
go install
```

Building The Provider (Terraform v0.13+)
-----------------------

### MacOS / Linux
Run the following command:
```sh
$ make build13
```

### Windows
Run the following commands for cmd:
```sh
cd %GOPATH%\src\github.com\AviatrixSystems\terraform-provider-aviatrix
go fmt
go install
xcopy "%GOPATH%\bin\terraform-provider-aviatrix.exe" "%APPDATA%\terraform.d\plugins\aviatrix.com\aviatrix\aviatrix\99.0.0\windows_amd64\" /Y
```
Run the following commands if using powershell:
```sh
cd "$env:GOPATH\src\github.com\AviatrixSystems\terraform-provider-aviatrix"
go fmt
go install
xcopy "$env:GOPATH\bin\terraform-provider-aviatrix.exe" "$env:APPDATA\terraform.d\plugins\aviatrix.com\aviatrix\aviatrix\99.0.0\windows_amd64\" /Y
```
Using Aviatrix Provider (Terraform v0.12+)
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

Using Aviatrix Provider (Terraform v0.13+)
-----------------------

For Terraform v0.13+, to use a locally built version of a provider you must add the following snippet to every module
that you want to use the provider in.

```hcl
terraform {
  required_providers {
    aviatrix = {
      source  = "aviatrix.com/aviatrix/aviatrix"
      version = "99.0.0"
    }
  }
}
```

Examples
--------

Check examples and documentation [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs)

Visit [here](https://github.com/AviatrixSystems/terraform-provider-aviatrix/tree/master/docs) for the complete documentation for all resources on github


Controller version
------------------
Due to some non-backward compatible changes in REST API not all controller versions are supported.
If you find a branch with the controller version please use that branch
Controller versions older than 3.3 are not supported
For example:
 * UserConnect-3.3 for 3.3.x controller version
 * UserConnect-3.4 for 3.4.x controller version

`master` branch is a development branch, thus it is unverified for all use cases and supports features that may not be available in the latest released version of the controller. For production use cases, please only use a released version from the Terraform Registry.

We also recommend you to update to the latest controller version to stay on top of fixes/features.
