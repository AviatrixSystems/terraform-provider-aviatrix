---
subcategory: "Gateway"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway_snat"
description: |-
   Configure customized SNAT policies for an Aviatrix gateway
---

# aviatrix_gateway_snat

The **aviatrix_gateway_snat** resource configures and manages policies for customized source NAT for Aviatrix gateways.

## Example Usage

```hcl
# Enable NAT function of mode "customized_snat" for an Aviatrix AWS Spoke Gateway
resource "aviatrix_gateway_snat" "test_snat" {
  gw_name   = "avtx-gw-1"
  snat_mode = "customized_snat"
  snat_policy {
    src_cidr    = "13.0.0.0/16"
    src_port    = "22"
    dst_cidr    = "14.0.0.0/16"
    dst_port    = "222"
    protocol    = "tcp"
    interface   = "eth0"
    connection  = "None"
    mark        = "22"
    snat_ips    = "175.32.12.12"
    snat_port   = "12"
    exclude_rtb = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Name of the Aviatrix gateway the custom SNAT will be configured for.
* `snat_mode` - (Optional) NAT mode. Valid values: "customized_snat". Default value: "customized_snat".
* `sync_to_ha` - (Optional) Sync the policies to the HA gateway. Valid values: true, false. Default: false.
* `snat_policy` - (Required) Policy rule applied for enabling source NAT (mode: "customized_snat"). Currently only supports AWS(1) and Azure(8).
  * `src_cidr` - (Optional) This is a qualifier condition that specifies a source IP address range where the rule applies. When not specified, this field is not used.
  * `src_port` - (Optional) This is a qualifier condition that specifies a source port that the rule applies. When not specified, this field is not used.
  * `dst_cidr` - (Optional) This is a qualifier condition that specifies a destination IP address range where the rule applies. When not specified, this field is not used.
  * `dst_port` - (Optional) This is a qualifier condition that specifies a destination port where the rule applies. When not specified, this field is not used.
  * `protocol` - (Optional) This is a qualifier condition that specifies a destination port protocol where the rule applies. Valid values: 'all', 'tcp', 'udp', 'icmp'. 'Default: 'all'.
  * `interface` - (Optional) This is a qualifier condition that specifies output interface where the rule applies. When not specified, this field is not used. Must be empty when `connection` is set.
  * `connection` - (Optional) This is a qualifier condition that specifies output connection where the rule applies. Default value: "None".
  * `mark` - (Optional) This is a qualifier condition that specifies a tag or mark of a TCP session where the rule applies. When not specified, this field is not used.
  * `snat_ips` - (Optional) This is a rule field that specifies the changed source IP address when all specified qualifier conditions meet. When not specified, this field is not used. One of the rule fields must be specified for this rule to take effect.
  * `snat_port` - (Optional) This is a rule field that specifies the changed source port when all specified qualifier conditions meet. When not specified, this field is not used. One of the rule fields must be specified for this rule to take effect.
  * `exclude_rtb` - (Optional) This field specifies which VPC private route table will not be programmed with the default route entry.

!> **WARNING:** Apply Route Entry was added in Provider version R2.21.0 with a default value of True. Any existing source NAT policy resources created in Terraform before R2.21.0 will see a diff. Please add `apply_route_entry = false` to any of these Terraform configurations to prevent any changes.

  * `apply_route_entry` - (Optional) This is an option to program the route entry 'DST CIDR pointing to Aviatrix Gateway' into Cloud platform routing table. Type: Boolean. Default: True. Available as of provider version R2.21.0+.

## Import

**gateway_snat** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_gateway_snat.test gw_name
```

## Notes
### snat_policy
When an attribute is referred to as 'left blank', or if an attribute is intended to not be specified in the configuration, it should not be written in the .tf configuration. For example, if `interface` was intended to not be specified, the attribute should not be present in the .tf config. Setting `interface = ""` (an empty string), is not the same as not specifying the interface value, and will result in deltas in the terraform state.
