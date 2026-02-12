---
subcategory: "Useful Tools"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_config_feature"
description: |-
  Enable/Disable aviatrix features
---

# aviatrix_config_feature


The **aviatrix_config_feature** resource allows management of an Aviatrix Features. This resource is available as of provider version R9.0.0+.

## Example Usage

```hcl
# Create a Microseg feature config
resource "aviatrix_config_feature" "test" {
  feature_name = "microseg"
  is_enabled = true
}
```


## Argument Reference

The following arguments are supported:

### Required
* `feature_name` - (Required) Which feature to enable, this should be one of the following:
"microseg", "cost_iq", "cai", "ipv6", "nfq_enforce_tls", "dcf_on_s2c", "dcf_on_psf", "dcf_stats_obs_sink", "dcf_logs_obs_sink", "k8s", "sre_metrics_export", "k8s_dcf_policies", "dcf_on_firenet", "primary_gateway_deletion", "enable_k8s", "enable_dcf_policies",.
* `is_enabled` - (Required) If set to true, the feature is enabled, set to false the feature is disabled.

## Import

**aviatrix_config_feature** can be imported using feature name, e.g. feature_name is : microseg

```
$ terraform import aviatrix_config_feature.test microseg
```
