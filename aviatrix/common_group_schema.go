package aviatrix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GroupRequiredSchema returns the required schema attributes for group resources.
// These can be reused by other resources that share the same required fields.
func GroupRequiredSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the gateway group.",
		},
		"cloud_type": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			Description:  "Type of cloud service provider.",
			ValidateFunc: validateCloudType,
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
			Description:      "VPC-ID/VNet-Name of cloud provider or Site-ID for edge providers.",
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

// GroupOptionalSchema returns the optional schema attributes for group resources.
// These can be reused by other resources that share the same optional fields.
func GroupOptionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customized_spoke_vpc_routes": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateCIDR,
			},
			Description: "A list of comma-separated CIDRs to be customized for the spoke VPC routes.",
		},
		"vpc_region": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Region of cloud provider. Required for CSP.",
		},
		"domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Network domain for the spoke group.",
		},
		"private_route_table_config": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Computed:    true,
			Description: "Set of Azure route table selectors to treat as private route tables for the group VNet. Each entry is in the format \"<route_table_name>:<resource_group_name>\". Only applicable for Azure (8), AzureGov (32) and AzureChina (2048).",
		},
	}
}

// GroupComputedSchema returns the computed (read-only) schema attributes for group resources.
func GroupComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Gateway group UUID.",
		},
		"gw_uuid_list": {
			Type:        schema.TypeSet,
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
		"explicitly_created": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates if the group was explicitly created.",
		},
	}
}

// MergeSchemaMaps merges multiple schema maps into a single map.
// Panics if duplicate keys are found across maps, as this indicates a programming error.
func MergeSchemaMaps(schemaMaps ...map[string]*schema.Schema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)
	for _, schemaMap := range schemaMaps {
		for k, v := range schemaMap {
			if _, exists := result[k]; exists {
				panic("duplicate schema key found: " + k)
			}
			result[k] = v
		}
	}
	return result
}
