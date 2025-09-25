//revive:disable:var-naming
package aviatrix

import (
	"context"
	"errors"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of DCF IPS profile UUIDs to assign to the VPC.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAviatrixDCFIpsProfileVpcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Get("vpc_id").(string)
	profiles := expandStringList(d.Get("dcf_ips_profiles").([]interface{}))

	_, err := client.SetIpsProfileVpc(ctx, vpcId, profiles)
	if err != nil {
		return diag.Errorf("failed to create DCF IPS profile VPC assignment: %v", err)
	}

	d.SetId(vpcId)
	return resourceAviatrixDCFIpsProfileVpcRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Id()
	profileVpc, err := client.GetIpsProfileVpc(ctx, vpcId)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF IPS profile VPC assignment: %v", err)
	}

	d.Set("vpc_id", profileVpc.VpcId)
	d.Set("dcf_ips_profiles", profileVpc.DcfIpsProfiles)

	return nil
}

func resourceAviatrixDCFIpsProfileVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Id()
	profiles := expandStringList(d.Get("dcf_ips_profiles").([]interface{}))

	_, err := client.SetIpsProfileVpc(ctx, vpcId, profiles)
	if err != nil {
		return diag.Errorf("failed to update DCF IPS profile VPC assignment: %v", err)
	}

	return resourceAviatrixDCFIpsProfileVpcRead(ctx, d, meta)
}

func resourceAviatrixDCFIpsProfileVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Id()

	// Clear all profiles by setting empty list
	_, err := client.SetIpsProfileVpc(ctx, vpcId, []string{})
	if err != nil {
		return diag.Errorf("failed to delete DCF IPS profile VPC assignment: %v", err)
	}

	d.SetId("")
	return nil
}
