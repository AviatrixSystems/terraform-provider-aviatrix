package aviatrix

import (
	"context"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	client := meta.(*goaviatrix.Client)

	privateModeMulticloudEndpoint := &goaviatrix.PrivateModeMulticloudEndpoint{
		AccountName:       d.Get("account_name").(string),
		VpcId:             d.Get("vpc_id").(string),
		Region:            d.Get("region").(string),
		ControllerLbVpcId: d.Get("controller_lb_vpc_id").(string),
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
	client := meta.(*goaviatrix.Client)

	if _, ok := d.GetOk("vpc_id"); !ok {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc_id received. Import Id is %s", id)
		d.Set("vpc_id", id)
	}

	vpcId := d.Get("vpc_id").(string)
	privateModeMulticloudEndpoint, err := client.GetPrivateModeMulticloudEndpoint(ctx, vpcId)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get Private Mode multicloud endpoint: %s", err)
	}

	d.Set("account_name", privateModeMulticloudEndpoint.AccountName)
	d.Set("vpc_id", privateModeMulticloudEndpoint.VpcId)
	d.Set("region", privateModeMulticloudEndpoint.Region)
	d.Set("controller_lb_vpc_id", privateModeMulticloudEndpoint.ControllerLbVpcId)
	d.Set("dns_entry", privateModeMulticloudEndpoint.DnsEntry)

	return nil
}

func resourceAviatrixPrivateModeMulticloudEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Get("vpc_id").(string)
	err := client.DeletePrivateModeMulticloudEndpoint(ctx, vpcId)
	if err != nil {
		return diag.Errorf("failed to delete Private Mode multicloud endpoint: %s", err)
	}

	return nil
}
