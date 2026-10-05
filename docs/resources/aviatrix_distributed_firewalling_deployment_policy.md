---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_distributed_firewalling_deployment_policy"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling deployment policy
---

# aviatrix_distributed_firewalling_deployment_policy

The **aviatrix_distributed_firewalling_deployment_policy** resource handles the creation and management of Distributed-firewalling deployment policy. Available as of Provider 3.2.2+.

Default CSP providers: `["GCP", "AWS", "AWS-GOV", "AZURE-GOV", "AZURE", "AVX-TEST"]`
Once the Distributed Cloud Firewall (DCF) is enabled a default deployment policy is created with default CSP providers. We can override this by adding a deployment policy to limit the use of DCF for CSPs. The policy has creation has two fields:
-providers: This is a list of CSPs where DCF should be enabled.
-set_defaults: If this is set to true, the deployment policy will be enabled in default CSPs. The providers set will be ignoerd. If this is set to false, then it will enable DCF only in the CSPs mentioned in the providers list.

Once this policy is destroyed the providers will be reset back to the default set of providers.

Deployment policy list providers cannot be an empty list, if DCF is not needed in any CSPs then it is best to disable DCF entirely.

## Example Usage

```hcl
# Create an Aviatrix Distributed Firewalling Deployment Policy
resource "aviatrix_distributed_firewalling_deployment_policy" "test" {
  providers = ["AWS", "GCP"]
  set_defaults = false
}
```

## Argument Reference

The following arguments are supported:

### Required
    * `providers` - (Optional) This is the list of CSPs where we want DCF to be enabled. This will be ignored when set_defaults is set to True Type: List of Strings
    * `set_defaults` - (Optional) Set to False if not provided. If we want to override the providers argument and set to default providers set this argument to True. Type: Boolean.

## Import

**aviatrix_distributed_firewalling_deployment_policy** can be imported using the controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_distributed_firewalling_deployment_policy.test 10-11-12-13
```
