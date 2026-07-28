---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_smart_group"
description: |-
  Creates and manages an Aviatrix Smart Group
---

# aviatrix_smart_group

The **aviatrix_smart_group** resource handles the creation and management of Smart Groups. Available as of Provider R2.22.0+.

## Example Usage

```hcl
# Create an Aviatrix Smart Group
resource "aviatrix_smart_group" "test_smart_group_ip" {
  name = "smart-group"
  selector {
    match_expressions {
      type         = "vm"
      account_name = "devops"
      region       = "us-west-2"
      tags         = {
        k3 = "v3"
      }
    }

    match_expressions {
      cidr = "10.0.0.0/16"
    }

    match_expressions {
      fqdn = "www.aviatrix.com"
    }

    match_expressions {
      site = "site-test-0"
    }

    match_expressions {
      s2c = "remote-site-name"
    }

    # Below are Kubernetes type examples

    # Match all pods and services in a cluster
    match_expressions {
      type           = "k8s"
      k8s_cluster_id = resource.aviatrix_kubernetes_cluster.test_cluster.cluster_id
    }

    # Match all pods and services in a namespace across all clusters
    match_expressions {
      type           = "k8s"
      k8s_namespace  = "testnamespace"
    }

    # Match all pods and services in a cluster and namespace
    match_expressions {
      type           = "k8s"
      k8s_cluster_id = resource.aviatrix_kubernetes_cluster.test_cluster.cluster_id
      k8s_namespace  = "testnamespace"
    }

    # Match a service by name in a namespace of a cluster
    match_expressions {
      type           = "k8s"
      k8s_cluster_id = resource.aviatrix_kubernetes_cluster.test_cluster.cluster_id
      k8s_namespace  = "testnamespace"
      k8s_service    = "testservice"
    }

    # Match a pod by name in a namespace of a cluster
    match_expressions {
      type           = "k8s"
      k8s_cluster_id = resource.aviatrix_kubernetes_cluster.test_cluster.cluster_id
      k8s_namespace  = "testnamespace"
      k8s_pod        = "testpod"
    }

    # Match all nodes of a cluster
    match_expressions {
      type           = "k8s_node"
      k8s_cluster_id = resource.aviatrix_kubernetes_cluster.test_cluster.cluster_id
    }

    # Below are serverless examples

    match_expressions {
      type     = "serverless"
      name     = "my-function"
    }

    // Below are external group type examples

    // generic format
    match_expressions {
      external = "External_group_ID"
      ext_args = {
        external_group_ID_specific_field_1 = "value1"
        external_group_ID_specific_field_2 = "value2"
      }
    }

    match_expressions {
      external = "geo"
      ext_args = {
        country_iso_code = "US"
        is_in_eu = "1"
        continent_code = "NA"
      }
    }

    match_expressions {
      external = "threatiq"
      ext_args = {
        protocol = "tcp"
        type = "ciarmy"
        severity = "major"
      }
    }

    match_expressions {
      external = "azureips"
      ext_args = {
        service_name = "AzureCloud"
        region = "eastus"
      }
    }

    match_expressions {
      external = "githubips"
      ext_args = {
        service_name = "web"
      }
    }

    // End of external group type examples

  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Smart Group.
* `selector` - (Required) Block containing match expressions to filter the Smart Group.
  * `match_expressions` - (Required) List of match expressions. The Smart Group will be a union of all resources matched by each `match_expressions`.`match_expressions` blocks cannot be empty.
    * `cidr` - (Optional) - CIDR block or IP Address this expression matches. `cidr` cannot be used with any other filters in the same `match_expressions` block.
    * `fqdn` - (Optional) - FQDN address this expression matches. `fqdn` cannot be used with any other filters in the same `match_expressions` block.
    * `site` - (Optional) - Edge Site-ID this expression matches. `site` cannot be used with any other filters in the same `match_expressions` block.
    * `type` - (Optional) - Type of resource this expression matches. If not using the external selector it must be one of "vm", "vpc", "subnet" or "k8s". `type` is required when `cidr`, `fqdn` and `site` are all not used.
    * `res_id` - (Optional) - Resource ID this expression matches.
    * `account_id` - (Optional) - Account ID this expression matches.
    * `account_name` - (Optional) - Account name this expression matches.
    * `name` - (Optional) - Name this expression matches.
    * `region` - (Optional) - Region this expression matches.
    * `zone` - (Optional) - Zone this expression matches.
    * `k8s_cluster_id` - (Optional) - Resource ID of the Kubernetes cluster this expression matches. The resource ID can be found in the `cluster_id` attribute of the `aviatrix_kubernetes_cluster` resource.
      This property can only be used when `type` is set to `"k8s"`.
    * `k8s_namespace` - (Optional) - Kubernetes namespace this expression matches.
      This property can only be used when `type` is set to `"k8s"`.
    * `k8s_service` - (Optional) - Kubernetes service name this expression matches.
      This property can only be used when `type` is set to `"k8s"`.
      This property must not be used when `k8s_pod` is set.
    * `k8s_pod` - (Optional) - Kubernetes pod name this expression matches.
      This property can only be used when `type` is set to `"k8s"` and `k8s_cluster_id` and `k8s_namespace` are also set.
      This property must not be used when `k8s_service` is set.
    * `s2c` - (Optional) - Name of the remote site. Represents the CIDRs associated with the remote site.
    * `external` - (Optional) - Specifies an external feed, currently either "geo" or "threatiq".
    * `ext_args` - (Optional) - Map of the arguments associated with the external feed such as "country_iso_code" for the "geo" feed.
    * `tags` - (Optional) - Map of tags this expression matches.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the Smart Group.

## Import

**aviatrix_smart_group** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_smart_group.test 41984f8b-5a37-4272-89b3-57c79e9ff77c
```
