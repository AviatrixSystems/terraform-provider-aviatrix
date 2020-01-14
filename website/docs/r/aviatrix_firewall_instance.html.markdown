---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_instance"
description: |-
  Creates and deletes Aviatrix Firewall Instances
---

# aviatrix_firewall_instance

The aviatrix_firewall_instance resource allows the creation and deletion of Aviatrix Firewall Instances.

## Example Usage

```hcl
# Create an Aviatrix Firewall Instance
resource "aviatrix_firewall_instance" "test_firewall_instance" {
  vpc_id            = "vpc-032005cc371"
  firenet_gw_name   = "avx_firenet_gw"
  firewall_name     = "avx_firewall_instance"
  firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
  firewall_size     = "m5.xlarge"
  management_subnet = "10.4.0.16/28"
  egress_subnet     = "10.4.0.32/28"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the Security VPC.
* `firenet_gw_name` - (Required) Name of the primary FireNet gateway.
* `firewall_name` - (Required) Name of the firewall instance to be created.
* `firewall_image` - (Required) One of the AWS AMIs from Palo Alto Networks.
* `firewall_size` - (Required) Instance size of the firewall. Example: "m5.xlarge".  
* `management_subnet` - (Required) Management Interface Subnet. Select the subnet whose name contains “gateway and firewall management”.
* `egress_subnet` - (Required) Egress Interface Subnet. Select the subnet whose name contains “FW-ingress-egress”.
* `firewall_image_version` - (Optional) Version of firewall image.
* `key_name`- (Optional) The .pem file name for SSH access to the firewall instance.
* `iam_role` - (Optional) In advanced mode, create an IAM Role on the AWS account that launched the FireNet gateway. Create a policy to attach to the role. The policy is to allow access to “Bootstrap Bucket”.
* `bootstrap_bucket_name`- (Optional) In advanced mode, specify a bootstrap bucket name where the initial configuration and policy file is stored.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `instance_id`- ID of the firewall instance created.
* `lan_interface`- ID of Lan Interface created.
* `management_interface`- ID of Management Interface created.
* `egress_interface`- ID of Egress Interface created.
* `public_ip`- Management Public IP.
                                                                      
## Import

Instance firewall_instance can be imported using the instance_id, e.g.

```
$ terraform import aviatrix_firewall_instance.test instance_id
```
