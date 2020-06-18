---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway"
description: |-
  Creates and manages Aviatrix spoke gateways
---

# aviatrix_spoke_gateway

The **aviatrix_spoke_gateway** resource allows the creation and management of Aviatrix spoke gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_aws" {
  cloud_type         = 1
  account_name       = "my-aws"
  gw_name            = "spoke-gw-aws"
  vpc_id             = "vpc-abcd123"
  vpc_reg            = "us-west-1"
  gw_size            = "t2.micro"
  subnet             = "10.11.0.0/24"
  enable_snat        = false
  enable_active_mesh = true
  tag_list           = [
    "k1:v1",
    "k2:v2",
  ]
}
```
```hcl
# Create an Aviatrix GCP Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_gcp" {
  cloud_type         = 4
  account_name       = "my-gcp"
  gw_name            = "spoke-gw-gcp"
  vpc_id             = "gcp-spoke-vpc"
  vpc_reg            = "us-west1-b"
  gw_size            = "n1-standard-1"
  subnet             = "10.12.0.0/24"
  enable_snat        = false
  enable_active_mesh = true
}
```
```hcl
# Create an Aviatrix Azure Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_azure" {
  cloud_type         = 8
  account_name       = "my-azure"
  gw_name            = "spoke-gw-01"
  vpc_id             = "spoke:test-spoke-gw-123"
  vpc_reg            = "West US"
  gw_size            = "Standard_B1s"
  subnet             = "10.13.0.0/24"
  enable_snat        = false
  enable_active_mesh = true
}
```
```hcl
# Create an Aviatrix Oracle Spoke Gateway
resource "aviatrix_spoke_gateway" "test_spoke_gateway_oracle" {
  cloud_type         = 16
  account_name       = "devops-oracle"
  gw_name            = "avtxgw-oracle"
  vpc_id             = "vpc-oracle-test"
  vpc_reg            = "us-ashburn-1"
  gw_size            = "VM.Standard2.2"
  subnet             = "10.7.0.0/16"
  enable_active_mesh = true
}
```

## Argument Reference

The following arguments are supported:

### Required
* `cloud_type` - (Required) Type of cloud service provider, requires an integer value. Currently only AWS(1), GCP(4), AZURE(8), and OCI(16) are supported.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `gw_name` - (Required) Name of the gateway which is going to be created.
* `vpc_id` - (Required) VPC-ID/VNet-Name of cloud provider. Example: AWS: "vpc-abcd1234", GCP: "vpc-gcp-test", AZURE: "vnet1:hello", OCI: "vpc-oracle-test1".
* `vpc_reg` - (Required) Region of cloud provider. Example: AWS: "us-east-1", GCP: "us-west2-a", AZURE: "East US 2", Oracle: "us-ashburn-1".
* `gw_size` - (Required) Size of the gateway instance. Example: AWS: "t2.large", AZURE: "Standard_B1s", Oracle: "VM.Standard2.2", GCP: "n1-standard-1".
* `subnet` - (Required) A VPC Network address range selected from one of the available network ranges. Example: "172.31.0.0/20". **NOTE: If using `insane_mode`, please see notes [here](#insane_mode-1).**

### HA
* `single_az_ha` (Optional) Set to true if this [feature](https://docs.aviatrix.com/Solutions/gateway_ha.html#single-az-gateway) is desired. Valid values: true, false.
* `ha_subnet` - (Optional) HA Subnet. Required if enabling HA for AWS/Azure gateway. Optional for GCP. Setting to empty/unsetting will disable HA. Setting to a valid subnet CIDR will create an HA gateway on the subnet. Example: "10.12.0.0/24"
* `ha_zone` - (Optional) HA Zone. Required only if enabling HA for GCP gateway. Setting to empty/unsetting will disable HA. Setting to a valid zone will create an HA gateway in the zone. Example: "us-west1-c".
* `ha_insane_mode_az` (Optional) AZ of subnet being created for Insane Mode Spoke HA Gateway. Required for AWS if `insane_mode` is enabled and `ha_subnet` is set. Example: AWS: "us-west-1a".
* `ha_eip` - (Optional) Public IP address that you want to assign to the HA peering instance. If no value is given, a new EIP will automatically be allocated. Only available for AWS.
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if enabling HA. Example: "t2.micro".

### Insane Mode
* `insane_mode` - (Optional) Enable [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) for Spoke Gateway. Insane Mode gateway size has to be at least c5 (AWS) or Standard_D3_v2 (AZURE). If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet. Only supported for AWS and Azure. Valid values: true, false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Spoke Gateway. Required for AWS if `insane_mode` is enabled. Example: AWS: "us-west-1a".

### SNAT/DNAT
* `single_ip_snat` - (Optional) Specify whether to enable Source NAT feature in "single_ip" mode on the gateway or not. Please disable AWS NAT instance before enabling this feature. Currently only supports AWS(1) and AZURE(8). Valid values: true, false.

-> **NOTE:** `enable_snat` has been renamed to `single_ip_snat` in provider version R2.10. Please see notes [here](#enable_snat-1) for more information.

~> **NOTE:** Custom SNAT and DNAT support have been deprecated and functionality has been moved to **aviatrix_gateway_snat** and **aviatrix_gateway_dnat** respectively, in provider version R2.10. Please see notes for `snat_mode`, `snat_policy` and `dnat_policy` in the Notes section below.

### Encryption
* `enable_encrypt_volume` - (Optional) Enable EBS volume encryption for Gateway. Only supports AWS. Valid values: true, false. Default value: false.
* `customer_managed_keys` - (Optional and Sensitive) Customer managed key ID.

### Route Customization
* `customized_spoke_vpc_routes` - (Optional) A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. It applies to this spoke gateway only​. Example: "10.0.0.0/116,10.2.0.0/16".
* `filtered_spoke_vpc_routes` - (Optional) A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s routing table. It applies to this spoke gateway only. Example: "10.2.0.0/116,10.3.0.0/16".
* `included_advertised_spoke_routes` - (Optional) A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC. Example: "10.4.0.0/116,10.5.0.0/16".

### Misc.
* `transit_gw` - (Optional) Specify the Aviatrix transit gateway to attach this spoke gateway to.
* `allocate_new_eip` - (Optional) When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in Controller 4.7+. Valid values: true, false. Default: true. Option not available for GCP, AZURE and OCI gateways, they will automatically allocate new EIPs.
* `eip` - (Optional) Required when `allocate_new_eip` is false. It uses the specified EIP for this gateway. Available in Controller 4.7+. Only available for AWS.
* `tag_list` - (Optional) Instance tag of cloud provider. Only AWS, cloud_type is "1", is supported. Example: ["key1:value1", "key2:value2"].
* `enable_active_mesh` - (Optional) Switch to enable/disable [Active Mesh Mode](https://docs.aviatrix.com/HowTos/activemesh_faq.html) for Spoke Gateway. Valid values: true, false. Default value: false.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server for Gateway. Currently only supports AWS. Valid values: true, false. Default value: false.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `eip` - Public IP address assigned to the gateway.
* `ha_eip` - Public IP address assigned to the HA gateway.
* `security_group_id` - Security group used for the spoke gateway.
* `cloud_instance_id` - Cloud instance ID of the spoke gateway.
* `private_ip` - Private IP address of the spoke gateway created.
* `ha_cloud_instance_id` - Cloud instance ID of the HA spoke gateway.
* `ha_gw_name` - Aviatrix spoke gateway unique name of HA spoke gateway.
* `ha_private_ip` - Private IP address of HA spoke gateway.

The following arguments are deprecated:

* `enable_snat` - (Optional) Specify whether enabling Source NAT feature on the gateway or not. Please disable AWS NAT instance before enabling this feature. Currently only supports AWS(1) and AZURE(8). Valid values: true, false.
* `snat_mode` - (Optional) Valid values: "primary", "secondary" and "custom". Default value: "primary".
* `snat_policy` - (Optional) Policy rule applied for "snat_mode" of "custom".
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
* `dnat_policy` - (Optional) Policy rule applied for enabling Destination NAT (DNAT), which allows you to change the destination to a virtual address range. Currently only supports AWS(1) and AZURE(8).
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

**spoke_gateway** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_spoke_gateway.test gw_name
```

## Notes
### insane_mode
If `insane_mode` is enabled, you must specify a valid /26 CIDR segment of the VPC specified for the `subnet`. This will then create a new subnet to be used for the corresponding gateway. You cannot specify an existing /26 subnet.

### enable_snat
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a spoke gateway with `enable_snat` set to true was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also change this attribute to `single_ip_snat` in your `.tf` file.

### snat_mode & snat_policy
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a spoke gateway with `snat_mode` and `snat_policy` was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to remove attribute’s value from the state. In addition, you must transfer its corresponding values to the **aviatrix_gateway_snat** resource in your `.tf` file and perform a 'terraform import' to rectify the state file.

### dnat_policy
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a spoke gateway with `dnat_policy` was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to remove attribute’s value from the state. In addition, you must its value to its corresponding **aviatrix_gateway_dnat** resource in your `.tf` file and perform a 'terraform import' to rectify the state file.

### ha_subnet
If you are using Aviatrix Terraform Provider R2.15+, and import a Google Cloud spoke gateway with HA enabled then you must set a value for `ha_subnet` in your Terraform config.
