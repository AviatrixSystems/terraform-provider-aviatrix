---
layout: "aviatrix"
page_title: "Provider: Aviatrix"
description: |-
  The Aviatrix provider is used to interact with Aviatrix resources
---

# Aviatrix Provider

The Aviatrix provider is used to interact with the many resources supported by Aviatrix, which can be built upon various cloud infrastructure providers such as AWS, Azure, Google Cloud, and Oracle Cloud. It needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

-> **NOTE:** Although *version* is an optional argument, we highly recommend users to specify the proper Aviatrix provider release version corresponding to their Controller version in order to avoid potential compatibility issues. Please see the [compatibility chart](https://www.terraform.io/docs/providers/aviatrix/guides/release-compatibility.html) for full details. For more information on versioning, a native Terraform provider argument, see [here](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).


## Example Usage

```hcl
# Configure Aviatrix provider source and version
terraform {
  required_providers {
    aviatrix = {
      source = "AviatrixSystems/aviatrix"
      version = "2.20.1"
    }
  }
}

# Configure Aviatrix provider
provider "aviatrix" {
  controller_ip           = "1.2.3.4"
  username                = "admin"
  password                = "password"
  skip_version_validation = false
  verify_ssl_certificate  = true
  path_to_ca_certificate  = "/path/to/ca/cert.crt"
}

# Create an access account
resource "aviatrix_account" "my_acc" {
  # ...
}
```

## Authentication

The Aviatrix provider offers various means of providing credentials for authentication. The following methods are supported:

* Static credentials
* Environment variables

### Static credentials
!> **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

Static credentials can be provided by specifying the `controller_ip`, `username` and `password` arguments in-line in the Aviatrix provider block:

**Usage:**

```hcl
provider "aviatrix" {
  controller_ip = "1.2.3.4"
  username      = "admin"
  password      = "password"
}
```

### Environment variables
You can provide credentials via the `AVIATRIX_CONTROLLER_IP`, `AVIATRIX_USERNAME`, `AVIATRIX_PASSWORD` environment variables, representing your Aviatrix controller's public IP, username and password of your Aviatrix access account, respectively.

```hcl
provider "aviatrix" {}
```

**Usage:**

```sh
$ export AVIATRIX_CONTROLLER_IP = "1.2.3.4"
$ export AVIATRIX_USERNAME = "admin"
$ export AVIATRIX_PASSWORD = "password"
$ terraform plan
```

## Argument Reference

The following arguments are supported:

### Required

-> **NOTE:** It's recommended to verify the SSL certificate of the controller when `controller_ip` is a FQDN.

* `controller_ip` - (Required) Aviatrix controller's public IP, private IP or FQDN.
* `username` - (Required) Aviatrix account username which will be used to login to Aviatrix controller.
* `password` - (Required) Aviatrix account password corresponding to above username.

### Optional
* `skip_version_validation` - (Optional) Valid values: true, false. Default: false. If set to true, it skips checking whether current Terraform provider supports current Controller version.
* `version` - (Optional) Specify Aviatrix provider release version number. If not specified, Terraform will automatically pull and source the latest release. For Terraform version 0.13+, do not use this attribute. Instead, set provider version using a `required_providers` block like in the example above.
* `verify_ssl_certificate` - (Optional) Valid values: true, false. Default: false. If set to true, the SSL certificate of the controller will be verified.
* `path_to_ca_certificate` - (Optional) Specify the path to the root CA certificate. Valid only when `verify_ssl_certificate` is true. The CA certificate is required when the controller is using a self-signed certificate.
