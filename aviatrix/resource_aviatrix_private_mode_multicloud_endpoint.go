package aviatrix

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixPrivateModeMulticloudEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixPrivateModeMulticloudEndpointCreate,
		ReadWithoutTimeout:   resourceAviatrixPrivateModeMulticloudEndpointRead,
		DeleteWithoutTimeout: resourceAviatrixPrivateModeMulticloudEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the access account.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the VPC region.",
			},
			"controller_lb_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC with the Controller load balancer.",
			},
			"dns_entry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS entry of this endpoint.",
			},
		},
	}
}

func resourceAviatrixPrivateModeMulticloudEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	privateModeMulticloudEndpoint := &goaviatrix.PrivateModeMulticloudEndpoint{
		AccountName:       getString(d, "account_name"),
		VpcId:             getString(d, "vpc_id"),
		Region:            getString(d, "region"),
		ControllerLbVpcId: getString(d, "controller_lb_vpc_id"),
	}

	flag := false
	defer resourceAviatrixPrivateModeMulticloudEndpointReadIfRequired(ctx, d, meta, &flag)

	err := client.CreatePrivateModeMulticloudEndpoint(ctx, privateModeMulticloudEndpoint)
	if err != nil {
		return diag.Errorf("failed to create Private Mode multicloud endpoint: %s", err)
	}

	d.SetId(privateModeMulticloudEndpoint.VpcId)

	return resourceAviatrixPrivateModeMulticloudEndpointReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixPrivateModeMulticloudEndpointReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixPrivateModeMulticloudEndpointRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixPrivateModeMulticloudEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if _, ok := d.GetOk("vpc_id"); !ok {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)
		mustSet(d, "vpc_id", id)
	}

	vpcId := getString(d, "vpc_id")
	privateModeMulticloudEndpoint, err := client.GetPrivateModeMulticloudEndpoint(ctx, vpcId)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get Private Mode multicloud endpoint: %s", err)
	}
	mustSet(d, "account_name", privateModeMulticloudEndpoint.AccountName)
	mustSet(d, "vpc_id", privateModeMulticloudEndpoint.VpcId)
	mustSet(d, "region", privateModeMulticloudEndpoint.Region)
	mustSet(d, "controller_lb_vpc_id", privateModeMulticloudEndpoint.ControllerLbVpcId)
	mustSet(d, "dns_entry", privateModeMulticloudEndpoint.DnsEntry)

	d.SetId(privateModeMulticloudEndpoint.VpcId)
	return nil
}

func resourceAviatrixPrivateModeMulticloudEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcId := getString(d, "vpc_id")
	err := client.DeletePrivateModeMulticloudEndpoint(ctx, vpcId)
	if err != nil {
		return diag.Errorf("failed to delete Private Mode multicloud endpoint: %s", err)
	}

	return nil
}
