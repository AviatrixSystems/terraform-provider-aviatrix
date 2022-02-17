---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_instance_images"
description: |-
  Gets the Aviatrix firewall instance images information.
---

# aviatrix_firewall_instance_images

Use this data source to get the list of firewall instance images for use in other resources.

**NOTE:** A firenet enabled gateway in a security VPC is required for this data source. 

## Example Usage

```hcl
# Aviatrix Firewall Instance Images Data Source
data "aviatrix_firewall_instance_images" "foo" {
  vpc_id = "vpc-1234567"
}
```

## Argument Reference

The following argument is supported:

* `vpc_id` - (Required) VPC ID. Example: AWS: "vpc-abcd1234", GCP: "vpc-gcp-test~-~project_id", Azure: "vnet_name:rg_name:resource_guid", OCI: "vpc-oracle-test1".

## Attribute Reference

In addition to the argument above, the following attributes are exported:

* `firewall_images` - List of firewall images.
    * `firewall_image` - Name of the firewall image.
    * `firewall_image_version` - List of firewall image versions.
    * `firewall_size` - List of firewall instance sizes.
