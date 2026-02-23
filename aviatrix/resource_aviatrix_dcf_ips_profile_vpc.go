//revive:disable:var-naming
package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDCFIpsProfileVpc() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFIpsProfileVpcCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFIpsProfileVpcRead,
		UpdateWithoutTimeout: resourceAviatrixDCFIpsProfileVpcUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFIpsProfileVpcDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC ID to assign DCF IPS profiles to.",
			},
			"dcf_ips_profiles": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of DCF IPS profile UUIDs to assign to the VPC.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAviatrixDCFIpsProfileVpcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcId := getString(d, "vpc_id")
	profiles := expandStringList(getSet(d, "dcf_ips_profiles").List())

	_, err := client.SetIpsProfileVpc(ctx, vpcId, profiles)
	if err != nil {
		return diag.Errorf("failed to create DCF IPS profile VPC assignment: %v", err)
	}

	d.SetId(vpcId)
	return resourceAviatrixDCFIpsProfileVpcRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcId := d.Id()
	profileVpc, err := client.GetIpsProfileVpc(ctx, vpcId)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) || strings.Contains(err.Error(), "AVXERR-IPS-0003") {
			// If VPC is not found, clear the resource from state, so TF can destroy this IPS profile assignment
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF IPS profile VPC assignment: %v", err)
	}
	mustSet(d, "vpc_id", profileVpc.VpcId)
	mustSet(d, "dcf_ips_profiles", profileVpc.DcfIpsProfiles)

	return nil
}

func resourceAviatrixDCFIpsProfileVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcId := d.Id()
	profiles := expandStringList(getSet(d, "dcf_ips_profiles").List())

	_, err := client.SetIpsProfileVpc(ctx, vpcId, profiles)
	if err != nil {
		return diag.Errorf("failed to update DCF IPS profile VPC assignment: %v", err)
	}

	return resourceAviatrixDCFIpsProfileVpcRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcId := d.Id()

	// Clear all profiles by setting empty list
	_, err := client.SetIpsProfileVpc(ctx, vpcId, []string{})
	if err != nil {
		return diag.Errorf("failed to delete DCF IPS profile VPC assignment: %v", err)
	}

	d.SetId("")
	return nil
}
