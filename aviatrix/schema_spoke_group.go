package aviatrix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SpokeGroupRequiredSchema returns the required schema attributes for spoke group resources.
// These can be reused by other resources that share the same required fields.
func SpokeGroupRequiredSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the spoke gateway group.",
		},
		"cloud_type": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			Description:  "Type of cloud service provider.",
			ValidateFunc: validateCloudType,
		},
		"gw_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Gateway type for the spoke group.",
		},
		"group_instance_size": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Instance size for gateways in the group.",
		},
		"vpc_id": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      "VPC-ID/VNet-Name of cloud provider.",
			DiffSuppressFunc: DiffSuppressFuncGatewayVpcId,
		},
		"account_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the Cloud-Account in Aviatrix controller.",
		},
	}
}

// SpokeGroupComputedSchema returns the computed (read-only) schema attributes for spoke group resources.
func SpokeGroupComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gw_uuid_list": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of gateway UUIDs in the group.",
		},
		"vpc_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "VPC UUID.",
		},
		"vendor_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Vendor name.",
		},
		"software_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Software version.",
		},
		"image_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Image version.",
		},
	}
}

// SpokeGroupAzureComputedSchema returns Azure-specific computed schema attributes.
func SpokeGroupAzureComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_eip_name_resource_group": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Azure EIP name and resource group.",
		},
		"bgp_lan_ip_list": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of BGP LAN interface IPs. Only valid for Azure.",
		},
	}
}

// MergeSchemaMaps merges multiple schema maps into a single map.
// Later maps override earlier maps if there are key conflicts.
func MergeSchemaMaps(schemaMaps ...map[string]*schema.Schema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)
	for _, schemaMap := range schemaMaps {
		for k, v := range schemaMap {
			result[k] = v
		}
	}
	return result
}
