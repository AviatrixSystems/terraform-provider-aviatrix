---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_k8s_config"
description: |-
  Creates and manages an Aviatrix K8s Config
---

# aviatrix_k8s_config

This resource is deprecated. Use aviatrix_config_feature instead.

The **aviatrix_k8s_config** resource allows management of an Aviatrix K8s configuration. This resource is available as of provider version R3.0+.

## Example Usage

```hcl
# Create an Aviatrix K8s config
resource "aviatrix_k8s_config" "test" {
  enable_k8s = true
  enable_dcf_policies = true
}
```


## Argument Reference

The following arguments are supported:

### Optional
* `enable_k8s` - Whether to enable the K8s feature on an Aviatrix Controller. Valid values: true, false. Default value: false.
* `enable_dcf_policies` - Whether to enable DCF policies in K8s clusters. Can only be true if enable_k8s is also true. Valid values: true, false. Default value: false.

## Import

**aviatrix_k8s_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_k8s_config.test 10-11-12-13
```
