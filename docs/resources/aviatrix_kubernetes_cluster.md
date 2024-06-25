---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_kubernetes_cluster"
description: |-
  Configures a Kubernetes cluster for building Aviatrix Smart Groups from applications running in Kubernetes.
---

# aviatrix_kubernetes_cluster

The **aviatrix_kubernetes_cluster** resource allows to configure a Kubernetes cluster.
Aviatrix Smart Groups can be built based on applications running in these clusters .
<!-- TODO: Add version, 7.2? -->
This resource is available as of provider version R3.0+.

## Example Usage


```hcl
# Create a Aviatrix Kubernetes Cluster so that the controller allows building Aviatrix Smart Groups from an AWS EKS cluster
resource "aviatrix_kubernetes_cluster" "rptest" {
  arn                 = data.aws_eks_cluster.eks_cluster.arn
  use_csp_credentials = true
}

data "aws_eks_cluster" "eks_cluster" {
  name = "mycluster"
}
```

```hcl
# Create a Aviatrix Kubernetes Cluster so that the controller allows building Aviatrix Smart Groups from an Azure AKS cluster
resource "aviatrix_kubernetes_cluster" "rptest" {
  cluster_id          = data.azurerm_kubernetes_cluster.mycluster.id
  use_csp_credentials = true
}

data "azurerm_kubernetes_cluster" "mycluster" {
  name                = "myakscluster"
  resource_group_name = "my-example-resource-group"
}
```

```hcl
# Create a Aviatrix Kubernetes Cluster for a custom built cluster
resource "aviatrix_kubernetes_cluster" "my_cluster" {
  cluster_id = "my-cluster-id"

  kube_config = var.kubeconfig

  account_name = data.aviatrix_account.my_account.account_name
  account_id = data.aviatrix_account.my_account.aws_account_number
  name        = "my_cluster"
  region       = "us-east-2"
  vpc_id       = data.aws_vpc.vpc.id
  public       = true
  platform     = "kops"
  version      = "1.30"
  network_mode = "FLAT"
  tags = {
    "type" = "prod"
  }
}
```


## Argument Reference

The following arguments are supported:

### Required

* Exactly one of cluster_id or arn must be provided.
  * `cluster_id` - (Optional) The ID of the Kubernetes cluster. If the cluster to be configured is an AKS cluster this should be the full resource ID of the AKS cluster. If the cluster is a custom built cluster this can be any unique identifier.
  * `arn` - (Optional) The ARN of the Kubernetes cluster if the cluster to be configured is an AWS EKS cluster.

### Optional

* `kube_config` - (Optional) The kubeconfig file for the Kubernetes cluster. This is a sensitive value.
* `use_csp_credentials` - (Optional) Whether to use the CSP credentials for the Kubernetes cluster. Valid values: true, false. Default value: false.

* `account_name` - (Optional) The name of the Aviatrix account to use for the Kubernetes cluster.
* `account_id` - (Optional) The ID of the Aviatrix account to use for the Kubernetes cluster.
* `name` - (Optional) The name of the Kubernetes cluster.
* `region` - (Optional) The region of the Kubernetes cluster.
* `vpc_id` - (Optional) The ID of the VPC where the Kubernetes cluster is running.
* `public` - (Optional) Whether the API server of Kubernetes cluster is publicly accessible over the internet. Valid values: true, false. Default value: false.
* `platform` - (Optional) The platform of the Kubernetes cluster. Any string is allowed.
* `version` - (Optional) The Kubernetes version of the cluster.
* `network_mode` - (Optional) The network mode of the Kubernetes cluster. Valid values: "FLAT", "OVERLAY".
* `tags` - (Optional) A map of tags to assign to the Kubernetes cluster.

