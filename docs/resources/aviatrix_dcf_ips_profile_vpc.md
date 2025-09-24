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
    aviatrix_dcf_ips_profile.profile1.uuid,
    aviatrix_dcf_ips_profile.profile2.uuid
  ]
}

# Example with multiple profiles
resource "aviatrix_dcf_ips_profile_vpc" "multi_profile" {
  vpc_id = "vpc-0123456789abcdef0"
  dcf_ips_profiles = [
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440002"
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
- `vpc_id` - (Required) The VPC ID to assign DCF IPS profiles to. Type: String.
- `dcf_ips_profiles` - (Required) List of DCF IPS profile UUIDs to assign to the VPC. Use an empty list to clear all profiles. Type: List(String).

## Import

**aviatrix_dcf_ips_profile_vpc** can be imported using the VPC ID:

```
$ terraform import aviatrix_dcf_ips_profile_vpc.example vpc-0a1b2c3d4e5f67890
```
