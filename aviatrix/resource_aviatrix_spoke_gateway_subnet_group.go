package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSpokeGatewaySubnetGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixSpokeGatewaySubnetGroupCreate,
		ReadWithoutTimeout:   resourceAviatrixSpokeGatewaySubnetGroupRead,
		UpdateWithoutTimeout: resourceAviatrixSpokeGatewaySubnetGroupUpdate,
		DeleteWithoutTimeout: resourceAviatrixSpokeGatewaySubnetGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet group name.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Spoke gateway name.",
			},
			"subnets": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "A set of subnets in the subnet group.",
			},
		},
	}
}

func marshalSpokeGatewaySubnetGroupInput(d *schema.ResourceData) *goaviatrix.SpokeGatewaySubnetGroup {
	var subnets []string
	for _, subnet := range d.Get("subnets").(*schema.Set).List() {
		subnets = append(subnets, subnet.(string))
	}

	spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
		SubnetGroupName: d.Get("name").(string),
		GatewayName:     d.Get("gw_name").(string),
		SubnetList:      subnets,
	}

	return spokeGatewaySubnetGroup
}

func resourceAviatrixSpokeGatewaySubnetGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	spokeGatewaySubnetGroup := marshalSpokeGatewaySubnetGroupInput(d)

	d.SetId(spokeGatewaySubnetGroup.GatewayName + "~" + spokeGatewaySubnetGroup.SubnetGroupName)
	flag := false
	defer resourceAviatrixAwsTgwSecurityDomainReadIfRequired(ctx, d, meta, &flag)

	if len(spokeGatewaySubnetGroup.SubnetList) == 0 {
		if err := client.AddSpokeGatewaySubnetGroup(ctx, spokeGatewaySubnetGroup); err != nil {
			return diag.Errorf("could not create an empty spoke gateway subnet group: %v", err)
		}
	} else {
		if err := client.UpdateSpokeGatewaySubnetGroup(ctx, spokeGatewaySubnetGroup); err != nil {
			return diag.Errorf("could not create spoke gateway subnet group: %v", err)
		}
	}

	return resourceAviatrixSpokeGatewaySubnetGroupReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixSpokeGatewaySubnetGroupReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeGatewaySubnetGroupRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixSpokeGatewaySubnetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)

	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("invalid ID, expected ID gw_name~name, instead got %s", d.Id())
		}
		d.Set("gw_name", parts[0])
		d.Set("name", parts[1])
		d.SetId(id)
	}

	name = d.Get("name").(string)
	spokeGatewayName := d.Get("gw_name").(string)

	spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
		GatewayName:     spokeGatewayName,
		SubnetGroupName: name,
	}

	err := client.GetSpokeGatewaySubnetGroup(ctx, spokeGatewaySubnetGroup)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get the spoke gateway subnet group %s due to %v", name, err)
	}

	if err := d.Set("subnets", spokeGatewaySubnetGroup.SubnetList); err != nil {
		return diag.Errorf("could not set subnets: %v", err)
	}

	d.SetId(spokeGatewayName + "~" + name)
	return nil
}

func resourceAviatrixSpokeGatewaySubnetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("subnets") {
		spokeGatewaySubnetGroup := marshalSpokeGatewaySubnetGroupInput(d)
		err := client.UpdateSpokeGatewaySubnetGroup(ctx, spokeGatewaySubnetGroup)
		if err != nil {
			return diag.Errorf("could not update spoke gateway subnet group: %v", err)
		}
	}

	return nil
}

func resourceAviatrixSpokeGatewaySubnetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	spokeGatewaySubnetGroup := marshalSpokeGatewaySubnetGroupInput(d)

	err := client.DeleteSpokeGatewaySubnetGroup(ctx, spokeGatewaySubnetGroup)
	if err != nil {
		return diag.Errorf("could not delete spoke gateway subnet group: %v", err)
	}

	return nil
}
