---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_transit_gateway"
description: |-
  Creates and manages the Aviatrix Transit Network gateways
---

# aviatrix_transit_gateway

The **aviatrix_transit_gateway** resource allows the creation and management of [Aviatrix Transit Network](https://docs.aviatrix.com/HowTos/transitvpc_faq.html#) gateways.

## Example Usage

```hcl
# Create an Aviatrix AWS Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
  cloud_type               = 1
  account_name             = "devops_aws"
  gw_name                  = "transit"
  vpc_id                   = "vpc-abcd1234"
  vpc_reg                  = "us-east-1"
  gw_size                  = "t2.micro"
  subnet                   = "10.1.0.0/24"
  ha_subnet                = "10.1.0.0/24"
  ha_gw_size               = "t2.micro"
  tag_list                 = [
    "name:value",
    "name1:value1",
    "name2:value2",
  ]
  enable_active_mesh       = true
  enable_hybrid_connection = true
  connected_transit        = true
}
```
```hcl
# Create an Aviatrix GCP Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_gcp" {
  cloud_type         = 4
  account_name       = "devops-gcp"
  gw_name            = "avtxgw-gcp"
  vpc_id             = "vpc-gcp-test"
  vpc_reg            = "us-west2-a"
  gw_size            = "n1-standard-1"
  subnet             = "10.8.0.0/16"
  ha_zone            = "us-west2-b"
  ha_subnet          = "10.8.0.0/16" // Optional
  ha_gw_size         = "n1-standard-1"
  enable_active_mesh = true
}
```
```hcl
# Create an Aviatrix Azure Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_azure" {
  cloud_type         = 8
  account_name       = "devops_azure"
  gw_name            = "transit"
  vpc_id             = "vnet1:hello"
  vpc_reg            = "West US"
  gw_size            = "Standard_B1s"
  subnet             = "10.30.0.0/24"
  ha_subnet          = "10.30.0.0/24"
  ha_gw_size         = "Standard_B1s"
  connected_transit  = true
  enable_active_mesh = true
}
```
```hcl
# Create an Aviatrix Oracle Transit Network Gateway
resource "aviatrix_transit_gateway" "test_transit_gateway_oracle" {
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
* `ha_subnet` - (Optional) HA Subnet CIDR. Required only if enabling HA for AWS/Azure gateway. Optional for GCP. Setting to empty/unsetting will disable HA. Setting to a valid subnet CIDR will create an HA gateway on the subnet. Example: "10.12.0.0/24".
* `ha_zone` - (Optional) HA Zone. Required only if enabling HA for GCP gateway. Setting to empty/unsetting will disable HA. Setting to a valid zone will create an HA gateway in the zone. Example: "us-west1-c".
* `ha_insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit HA Gateway. Required for AWS if `insane_mode` is enabled and `ha_subnet` is set. Example: AWS: "us-west-1a".
* `ha_eip` - (Optional) Public IP address that you want to assign to the HA peering instance. If no value is given, a new EIP will automatically be allocated. Only available for AWS.
* `ha_gw_size` - (Optional) HA Gateway Size. Mandatory if enabling HA. Example: "t2.micro".

### Insane Mode
* `insane_mode` - (Optional) Specify true for [Insane Mode](https://docs.aviatrix.com/HowTos/insane_mode.html) high performance gateway. Insane Mode gateway size must be at least c5 size (AWS) or Standard_D3_v2 (AZURE). If enabled, you must specify a valid /26 CIDR segment of the VPC to create a new subnet. Only available for AWS and Azure. Valid values: true, false.
* `insane_mode_az` - (Optional) AZ of subnet being created for Insane Mode Transit Gateway. Required for AWS if `insane_mode` is enabled. Example: AWS: "us-west-1a".

### SNAT
* `single_ip_snat` - (Optional) Enable "single_ip" mode Source NAT for this container. Valid values: true, false. **NOTE: Please see notes [here](#enable_snat-1) in regards to changes to this argument in R2.10.**

### Advanced Config
* `connected_transit` - (Optional) Specify Connected Transit status. If enabled, it allows spokes to run traffics to other spokes via transit gateway. Valid values: true, false. Default value: false.
* `enable_advertise_transit_cidr` - (Optional) Switch to enable/disable advertise transit VPC network CIDR for a vgw connection. Available as of R2.6. **NOTE: If previously enabled through vgw_conn resource prior to provider version R2.6, please see notes [here](#cidr-advertising).**
* `bgp_manual_spoke_advertise_cidrs` - (Optional) Intended CIDR list to advertise to VGW. Example: "10.2.0.0/16,10.4.0.0/16". Available as of R2.6. **NOTE: If previously enabled through vgw_conn resource prior to provider version R2.6, please see notes [here](#cidr-advertising).**
* `enable_hybrid_connection` - (Optional) Sign of readiness for TGW connection. Only supported for AWS. Example: false.
* `enable_firenet` - (Optional) Sign of readiness for FireNet connection. Valid values: true, false. Default value: false. **NOTE: If previously using an older provider version R2.5 where attribute name was `enable_firenet_interfaces`, please see notes [here](#enable_firenet-1).**

-> **NOTE:** Enabling FireNet will automatically enable hybrid connection. If `enable_firenet` is set to true, please set `enable_hybrid_connection` to true in the respective **aviatrix_transit_gateway** as well.

* `enable_transit_firenet` - (Optional) Sign of readiness for [Transit FireNet](https://docs.aviatrix.com/HowTos/transit_firenet_faq.html) connection. Valid values: true, false. Default value: false.

### Encryption
* `enable_encrypt_volume` - (Optional) Enable EBS volume encryption for Gateway. Only supports AWS. Valid values: true, false. Default value: false.
* `customer_managed_keys` - (Optional and Sensitive) Customer managed key ID.

### Route Customization
* `customized_spoke_vpc_routes` - (Optional) A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. It applies to all spoke gateways attached to this transit gateway. Example: "10.0.0.0/116,10.2.0.0/16".
* `filtered_spoke_vpc_routes` - (Optional) A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s routing table. It applies to all spoke gateways attached to this transit gateway. Example: "10.2.0.0/116,10.3.0.0/16".
* `excluded_advertised_spoke_routes` - (Optional) A list of comma separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'. When configured, it inspects all the advertised CIDRs from its spoke gateways and remove those included in the 'Excluded CIDR List'. Example: "10.4.0.0/116,10.5.0.0/16".

### Misc.
* `allocate_new_eip` - (Optional) When value is false, reuse an idle address in Elastic IP pool for this gateway. Otherwise, allocate a new Elastic IP and use it for this gateway. Available in Controller 4.7+. Valid values: true, false. Default: true. Option not available for GCP, Azure and OCI gateways, they will automatically allocate new EIPs.
* `eip` - (Optional) Required when `allocate_new_eip` is false. It uses the specified EIP for this gateway. Available in Controller version 4.7+. Only available for AWS.
* `tag_list` - (Optional) Instance tag of cloud provider. Only supported for AWS. Example: ["key1:value1","key2:value2"].
* `enable_active_mesh` - (Optional) Switch to enable/disable [Active Mesh Mode](https://docs.aviatrix.com/HowTos/activemesh_faq.html) for Transit Gateway. Valid values: true, false. Default value: false.
* `enable_vpc_dns_server` - (Optional) Enable VPC DNS Server for Gateway. Currently only supports AWS. Valid values: true, false. Default value: false.
* `enable_learned_cidrs_approval` - (Optional) Switch to enable/disable [encrypted transit approval](https://docs.aviatrix.com/HowTos/transit_approval.html) for transit Gateway. Valid values: true, false. Default value: false.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `eip` - Public IP address assigned to the gateway.
* `ha_eip` - Public IP address assigned to the HA gateway.
* `security_group_id` - Security group used for the transit gateway.
* `cloud_instance_id` - Cloud instance ID of the transit gateway.
* `private_ip` - Private IP address of the transit gateway created.
* `ha_cloud_instance_id` - Cloud instance ID of the HA transit gateway.
* `ha_gw_name` - Aviatrix transit gateway unique name of HA transit gateway.
* `ha_private_ip` - Private IP address of the HA transit gateway created.

The following arguments are deprecated:

* `enable_firenet_interfaces` - (Optional) Sign of readiness for FireNet connection. Valid values: true, false. Default value: false.
* `enable_snat` - (Optional) Enable Source NAT for this container. Valid values: true, false.

## Import

**transit_gateway** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_transit_gateway.test gw_name
```

## Notes
### CIDR advertising
`enable_advertise_transit_cidr` and `bgp_manual_spoke_advertise_cidrs` functionality has been migrated over to **aviatrix_transit_gateway** as of Aviatrix Terraform Provider R2.6. If you are using/upgraded to Aviatrix Terraform Provider R2.6+, and a **vgw_conn** resource was originally created with a provider version <R2.6, you must cut and paste these two arguments (and values) into the corresponding transit gateway resource referenced in the **vgw_conn**. A 'terraform refresh' will then successfully complete the migration and rectify the state file.

### enable_firenet
If you are using/upgraded to Aviatrix Terraform Provider R2.5+/UserConnect-5.0+ , and an AWS transit_gateway resource with `enable_firenet_interfaces` enabled was created with a provider version < R2.5/ UserConnect-5.0, you must replace `enable_firenet_interfaces` with `enable_firenet` in your configuration file, and do ‘terraform refresh’ to set its value to `enable_firenet` and apply it into the state file.

### insane_mode
If `insane_mode` is enabled, you must specify a valid /26 CIDR segment of the VPC specified for the `subnet`. This will then create a new subnet to be used for the corresponding gateway. You cannot specify an existing /26 subnet.

### enable_snat
If you are using/upgraded to Aviatrix Terraform Provider R2.10+, and a transit gateway with `enable_snat` set to true was originally created with a provider version <R2.10, you must do a ‘terraform refresh’ to update and apply the attribute’s value into the state. In addition, you must also change this attribute to `single_ip_snat` in your `.tf` file.

### ha_subnet
If you are using Aviatrix Terraform Provider R2.15+, and import a Google Cloud transit gateway with HA enabled then you must set a value for `ha_subnet` in your Terraform config.
