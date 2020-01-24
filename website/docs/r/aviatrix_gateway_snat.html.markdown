---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway_snat"
description: |-
   Configure policies for NAT function mode of "customized_snat" for an Aviatrix gateway
---

# aviatrix_gateway_snat

The aviatrix_gateway_snat resource configures and manages policies for NAT function mode of "customized_snat" for Aviatrix gateways.

## Example Usage

```hcl
# Enable NAT function of mode "customized_snat" for an Aviatrix AWS Spoke Gateway
resource "aviatrix_gateway_snat" "test_snat" {
  gw_name   = "avtxgw1"
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

* `gw_name` - (Required) Aviatrix gateway unique name.
* `snat_mode` - (Optional) NAT mode. Valid values: "customized_snat". Default value: "customized_snat". 
* `dnat_policy` - (Required) Policy rule applied for enabling Destination NAT (DNAT), which allows you to change the destination to a virtual address range. Currently only supports AWS(1) and ARM(8).
  * `src_ip` - (Optional) A source IP address range where the policy rule applies.
  * `src_port` - (Optional) A source port that the policy rule applies.
  * `dst_ip` - (Optional) A destination IP address range where the policy rule applies.
  * `dst_port` - (Optional) A destination port where the policy rule applies.
  * `protocol` - (Optional) A destination port protocol where the policy rule applies.
  * `interface` - (Optional) An output interface where the policy rule applies.
  * `connection` - (Optional) Default value: "None".
  * `mark` - (Optional) A tag or mark of a TCP session where the policy rule applies.
  * `new_src_ip` - (Optional) The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.
  * `new_src_port` - (Optional) The translated destination port when all specified qualifier conditions meet. One of the rule field must be specified for this rule to take effect.
  * `exclude_rtb` - (Optional) This field specifies which VPC private route table will not be programmed with the default route entry.

## Import

Instance gateway can be imported using the gw_name, e.g.

```
$ terraform import aviatrix_gateway_snat.test gw_name
```