---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway_dnat"
description: |-
   Configure policies for destination NAT for an Aviatrix gateway
---

# aviatrix_gateway_dnat

The **aviatrix_gateway_dnat** resource configures and manages policies for destination NAT function for Aviatrix gateways.

## Example Usage

```hcl
# Add policy for destination NAT function for an Aviatrix AWS Spoke Gateway
resource "aviatrix_gateway_dnat" "test_dnat" {
  gw_name   = "avtx-gw-1"
  dnat_policy {
    src_cidr    = "13.0.0.0/16"
    src_port    = "22"
    dst_cidr    = "14.0.0.0/16"
    dst_port    = "222"
    protocol    = "tcp"
    interface   = "eth0"
    connection  = "None"
    mark        = "22"
    dnat_ips    = "175.32.12.12"
    dnat_port   = "12"
    exclude_rtb = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Name of the Aviatrix gateway the custom DNAT will be configured for.
* `dnat_policy` - (Required) Policy rule applied for enabling Destination NAT (DNAT), which allows you to change the destination to a virtual address range. Currently only supports AWS(1) and ARM(8).
  * `src_cidr` - (Optional) This is a qualifier condition that specifies a source IP address range where the rule applies. When left blank, this field is not used.
  * `src_port` - (Optional) This is a qualifier condition that specifies a source port that the rule applies. When left blank, this field is not used.
  * `dst_cidr` - (Optional) This is a qualifier condition that specifies a destination IP address range where the rule applies. When left blank, this field is not used.
  * `dst_port` - (Optional) This is a qualifier condition that specifies a destination port where the rule applies. When left blank, this field is not used.
  * `protocol` - (Optional) This is a qualifier condition that specifies a destination port protocol where the rule applies. When left blank, this field is not used.
  * `interface` - (Optional) This is a qualifier condition that specifies output interface where the rule applies. When left blank, this field is not used.
  * `connection` - (Optional) Default value: "None".
  * `mark` - (Optional) This is a rule field that specifies a tag or mark of a TCP session when all qualifier conditions meet. When left blank, this field is not used.
  * `dnat_ips` - (Optional) This is a rule field that specifies the translated destination IP address when all specified qualifier conditions meet. When left blank, this field is not used. One of the rule field must be specified for this rule to take effect.
  * `dnat_port` - (Optional) This is a rule field that specifies the translated destination port when all specified qualifier conditions meet. When left blank, this field is not used. One of the rule field must be specified for this rule to take effect.
  * `exclude_rtb` - (Optional) This field specifies which VPC private route table will not be programmed with the default route entry.

## Import

**gateway_dnat** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_gateway_dnat.test gw_name
```
