package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeGroupTransparentInspection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixSpokeGroupTransparentInspectionCreate,
		ReadContext:   resourceAviatrixSpokeGroupTransparentInspectionRead,
		UpdateContext: resourceAviatrixSpokeGroupTransparentInspectionUpdate,
		DeleteContext: resourceAviatrixSpokeGroupTransparentInspectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the spoke group. The group must have gw_type SPOKE and cloud_type AWS (1).",
			},
			"route_tables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Set of private route table IDs to use for transparent inspection. If not set, the controller " +
					"selects the VPC's private route tables.",
			},
			"egress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable egress through the spoke gateways for transparent inspection. Valid values: true, false. Default value: false.",
			},
		},
	}
}

// getSpokeGroupTransparentInspectionVpcID returns the VPC ID of the spoke group,
// validating that the group supports transparent inspection.
func getSpokeGroupTransparentInspectionVpcID(ctx context.Context, client *goaviatrix.Client, groupUUID string) (string, error) {
	group, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		return "", fmt.Errorf("failed to get spoke group %s: %w", groupUUID, err)
	}
	if !strings.EqualFold(strings.TrimPrefix(group.GwType, "GwGroupType."), "SPOKE") {
		return "", fmt.Errorf("transparent inspection is only supported for spoke groups with gw_type SPOKE")
	}
	if group.CloudType != goaviatrix.AWS {
		return "", fmt.Errorf("transparent inspection is only supported for AWS (1)")
	}
	if len(group.GwUUIDList) == 0 {
		return "", fmt.Errorf("spoke group %s has no gateways; transparent inspection can only be enabled after its spoke instances are launched", group.GroupName)
	}
	return strings.Split(group.VpcID, subnetSeparator)[0], nil
}

func buildSpokeGroupTransparentInspectionConfig(d *schema.ResourceData, vpcID string) *goaviatrix.GatewayInsertionVpcConfig {
	egressMode := goaviatrix.GatewayInsertionEgressModeNonEgress
	if getBool(d, "egress") {
		egressMode = goaviatrix.GatewayInsertionEgressModeLocalEgress
	}
	return &goaviatrix.GatewayInsertionVpcConfig{
		VpcID:       vpcID,
		RouteTables: getStringSet(d, "route_tables"),
		EgressMode:  egressMode,
	}
}

func resourceAviatrixSpokeGroupTransparentInspectionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	groupUUID := getString(d, "group_uuid")

	vpcID, err := getSpokeGroupTransparentInspectionVpcID(ctx, client, groupUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Enabling transparent inspection for spoke group %s (VPC %s)", groupUUID, vpcID)
	if err := client.EnableGatewayInsertion(ctx, buildSpokeGroupTransparentInspectionConfig(d, vpcID)); err != nil {
		return diag.Errorf("failed to enable transparent inspection for spoke group %s: %s", groupUUID, err)
	}

	d.SetId(groupUUID)
	return resourceAviatrixSpokeGroupTransparentInspectionRead(ctx, d, meta)
}

func resourceAviatrixSpokeGroupTransparentInspectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	groupUUID := d.Id()

	group, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get spoke group %s: %s", groupUUID, err)
	}

	vpcID := strings.Split(group.VpcID, subnetSeparator)[0]
	status, err := client.GetGatewayInsertionStatus(ctx, vpcID)
	if err != nil {
		return diag.Errorf("failed to get transparent inspection status for spoke group %s: %s", groupUUID, err)
	}
	if !status.Enabled {
		d.SetId("")
		return nil
	}

	mustSet(d, "group_uuid", groupUUID)
	if len(getStringSet(d, "route_tables")) != 0 {
		mustSet(d, "route_tables", status.RouteTableConfig)
	}
	mustSet(d, "egress",
		status.InsertionMode != "" && status.InsertionMode != goaviatrix.GatewayInsertionEgressModeNonEgress)
	return nil
}

func resourceAviatrixSpokeGroupTransparentInspectionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	groupUUID := d.Id()

	if d.HasChanges("route_tables", "egress") {
		vpcID, err := getSpokeGroupTransparentInspectionVpcID(ctx, client, groupUUID)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.EnableGatewayInsertion(ctx, buildSpokeGroupTransparentInspectionConfig(d, vpcID)); err != nil {
			return diag.Errorf("failed to update transparent inspection for spoke group %s: %s", groupUUID, err)
		}
	}

	return resourceAviatrixSpokeGroupTransparentInspectionRead(ctx, d, meta)
}

func resourceAviatrixSpokeGroupTransparentInspectionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	groupUUID := d.Id()

	group, err := client.GetGatewayGroup(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return nil
		}
		return diag.Errorf("failed to get spoke group %s: %s", groupUUID, err)
	}

	vpcID := strings.Split(group.VpcID, subnetSeparator)[0]
	enabled, err := client.IsGatewayInsertionEnabled(ctx, vpcID)
	if err != nil {
		return diag.Errorf("failed to get transparent inspection status for spoke group %s: %s", groupUUID, err)
	}
	if !enabled {
		return nil
	}

	log.Printf("[INFO] Disabling transparent inspection for spoke group %s (VPC %s)", groupUUID, vpcID)
	if err := client.DisableGatewayInsertion(ctx, vpcID); err != nil {
		return diag.Errorf("failed to disable transparent inspection for spoke group %s: %s", groupUUID, err)
	}
	return nil
}
