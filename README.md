# Terraform Provider
The [Terraform](https://terraform.io) provider for [Aviatrix](https://aviatrix.com/)

- To see official latest release, please check the [Terraform Registry](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs)
- To see official version compatibility between Terraform, Aviatrix provider and Aviatrix Controller, please see the [Compatibility Chart](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/release-compatibility)

## Requirements
-	Install [Terraform](https://www.terraform.io/downloads.html) 0.12.x/0.13.x/0.14.x/0.15.x (0.11.x or lower is incompatible)
-	Install [Go](https://golang.org/doc/install) 1.18+ (This will be used to build the provider plugin.)
  - For those using an M1 Macbook, please install the darwin-**arm64** package instead of the darwin-amd64 package
    - Use [this](https://go.dev/dl/go1.18.2.darwin-arm64.pkg) link for the darwin-arm package (1.18.2) (Latest as of 24 Feb 2023)
    - You may find the latest download options [here](https://go.dev/dl/) otherwise
-	Create a directory, `go/`, follow this [doc](https://github.com/golang/go/wiki/SettingGOPATH) to edit `~/.bash_profile` to setup the GOPATH environment variable)

## Instructions
### 1. Clone our Aviatrix Terraform repository
Clone repository to: `$GOPATH/src/github.com/AviatrixSystems/terraform-provider-aviatrix`

**MacOS / Linux**
```sh
$ mkdir -p $GOPATH/src/github.com/AviatrixSystems
$ cd $GOPATH/src/github.com/AviatrixSystems
$ git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix.git
```

**Windows**
```sh
mkdir %GOPATH%\src\github.com\AviatrixSystems
cd %GOPATH%\src\github.com\AviatrixSystems
git clone https://github.com/AviatrixSystems/terraform-provider-aviatrix.git
```

### 2. Enter the provider directory and build the provider
**MacOS / Linux**
```sh
$ cd $GOPATH/src/github.com/AviatrixSystems/terraform-provider-aviatrix
$ make fmt
$ make build13
```

**Windows**

Run the following commands for cmd:
```sh
cd %GOPATH%\src\github.com\AviatrixSystems\terraform-provider-aviatrix
go fmt
go install
xcopy "%GOPATH%\bin\terraform-provider-aviatrix.exe" "%APPDATA%\terraform.d\plugins\aviatrix.com\aviatrix\aviatrix\99.0.0\windows_amd64\" /Y
```
Run the following commands if using Powershell:
```sh
cd "$env:GOPATH\src\github.com\AviatrixSystems\terraform-provider-aviatrix"
go fmt
go install
xcopy "$env:GOPATH\bin\terraform-provider-aviatrix.exe" "$env:APPDATA\terraform.d\plugins\aviatrix.com\aviatrix\aviatrix\99.0.0\windows_amd64\" /Y
```

### 3. Using Aviatrix Provider (Terraform v0.13+)
For Terraform v0.13+, to use a locally built version of a provider you must add the following snippet to every module that you want to use the provider in.

```hcl
terraform {
  required_providers {
    aviatrix = {
      source  = "aviatrix.com/aviatrix/aviatrix"
      version = "99.0.0" # this MUST be set as seen here
    }
  }
}
```

### 4. OPTIONAL: Building out specific versions
In order to test out a specific version, you can build the corresponding branch. By default, step 2 will build whatever git branch you are currently on (typically default is master branch)

Check out the Git branch where the fix is, which corresponds to the Controller version, and repeat step 2 to build the provider

```sh
# Example - Want to locally build a version with a bug fix for 6.7-patch - corresponding to an unreleased R2.21.3
# assuming current working directory is the local git repo

$ git pull
$ git checkout UserConnect-6.7
$ make fmt
$ make build13
```

---
## Legacy Instructions (kept for reference)
**WARNING:** The following instructions are kept for legacy purposes. ONLY follow these instructions below IF you are trying to build an old version of the Aviatrix provider AND are using an old Terraform version (0.12)

### 1. Clone the repo - follow same steps as above

### 2. Building the provider (Terraform v0.12)
**MacOS / Linux**
```sh
$ cd $GOPATH/src/github.com/AviatrixSystems/terraform-provider-aviatrix
$ make fmt
$ make build13
```

**Windows**
```sh
cd %GOPATH%\src\github.com\AviatrixSystems\terraform-provider-aviatrix
go fmt
go install
```

### 3. Activating the Aviatrix provider (Terraform v0.12-)
**To use this provider after building it with Go, you must activate it.**

Activate the provider by adding the following to `~/.terraformrc` on Linux/Unix.

If the rc file is not present, it should be created

**MacOS / Linux**
```sh
providers {
  "aviatrix" = "$GOPATH/bin/terraform-provider-aviatrix"
}
```
**Windows**
```sh
providers {
  "aviatrix" = "$GOPATH/bin/terraform-provider-aviatrix.exe"
}
```
For Windows, the file should be at '%APPDATA%\terraform.rc'. Do not change $GOPATH to %GOPATH%.

- In Windows, for terraform 0.11.8 and lower use the above text.
- In Windows, for terraform 0.11.9 and higher use the following at '%APPDATA%\terraform.rc'

### 4. Using Aviatrix Provider (Terraform v0.12-)
In versions prior to 0.13, the `terraform {}` was not required


---
## Examples

- Check resource examples and documentation [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs)

- Visit [here](https://github.com/AviatrixSystems/terraform-provider-aviatrix/tree/master/docs) for the complete documentation for all resources on Github


## Controller version
Due to some non-backward compatible changes in REST API, not all controller versions are supported.

If you find a branch with the controller version please use that branch
Controller versions older than 3.3 are not supported
For example:
 * UserConnect-3.3 for 3.3.x controller version
 * UserConnect-3.4 for 3.4.x controller version

`master` branch is a development branch, thus it is unverified for all use cases and supports features that may not be available in the latest released version of the controller. For production use cases, please only use a released version from the Terraform Registry.

We also recommend you to update to the latest controller version to stay on top of fixes/features.
