---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_ips_profile_vpc"
description: |-
  Manages DCF IPS profiles for a VPC
---

# aviatrix_dcf_ips_profile_vpc

The **aviatrix_dcf_ips_profile_vpc** resource allows you to manage the list of DCF IPS profiles assigned to a specific VPC for Distributed Cloud Firewall (DCF) Intrusion Prevention System. Available as of Provider 3.2.2+.

## Example Usage

```hcl
# Assign IPS profiles to a VPC
resource "aviatrix_dcf_ips_profile_vpc" "example" {
  vpc_id = "vpc-0a1b2c3d4e5f67890"
  dcf_ips_profiles = [
    aviatrix_dcf_ips_profile.profile1.uuid
  ]
}


# Clear all profiles from a VPC
resource "aviatrix_dcf_ips_profile_vpc" "clear_profiles" {
  vpc_id           = "vpc-0a1b2c3d4e5f67890"
  dcf_ips_profiles = []
}
```

## Argument Reference

The following arguments are supported:

### Required
- `vpc_id` - (Required) The VPC ID to which DCF IPS Profiles will be assigned.
The VPC must have a DCF-applied gateway and must be created before the `aviatrix_dcf_ips_profile_vpc` resource is defined. Type: String

- `dcf_ips_profiles` – (Required) List of DCF IPS profile UUIDs to assign to the VPC. Set to an empty list (`[]`) to remove all profiles. Only one IPS profile can be assigned per VPC on Controller version 8.2. Type: `list(string)`.

## Import

**aviatrix_dcf_ips_profile_vpc** can be imported using the VPC ID:

```
$ terraform import aviatrix_dcf_ips_profile_vpc.example vpc-0a1b2c3d4e5f67890
```
