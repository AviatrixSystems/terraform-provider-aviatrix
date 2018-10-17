---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dc_extn"
sidebar_current: "docs-aviatrix-resource-dc-extn"
description: |-
  Creates and manages Aviatrix Datacenter Extension. Only available in CloudN
---

# aviatrix_dc_extn

The DCExtn resource allows the creation and management of Aviatrix Datacenter Extension. Only available in CloudN

## Example Usage

```hcl
# Create Aviatrix Datacenter Extension
resource "aviatrix_dc_extn" "test_dc_extn" {
  cloud_type = "1"
  account_name = "test-account"
  gw_name = "gateway-1"
  vpc_reg = "us-east-1"
  gw_size = "t2.micro"
  subnet_cidr = "172.16.32.0/20"
  internet_access = "no"
  public_subnet = "no"
  tunnel_type = "tcp"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. 1: AWS.
* `account_name` - (Required) Account Name to use for launching this container.
* `gw_name` - (Required) A unique name of the container.
* `vpc_reg` - (Required) A region where this container will be launched.
* `gw_size` - (Required) Size of Gateway Instance. e.g.: "t2.micro".
* `subnet_cidr` - (Required) A VPC Network address range selected from one of the available network ranges. e.g.: "172.16.32.0/20".
* `internet_access` - (Optional) Indicates whether internet access would be enabled or not. Valid Values: "yes", "no".
* `public_subnet` - (Optional) Indicates whether a public subnet CIDR would be assigned or not. Valid Values: "yes", "no".
* `tunnel_type` - (Optional) Indicates whether tunnel is TCP-based or UDP-based. Valid Values: "tcp", "udp".
