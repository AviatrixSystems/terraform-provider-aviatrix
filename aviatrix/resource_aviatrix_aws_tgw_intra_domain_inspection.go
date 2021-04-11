package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAwsTgwIntraDomainInspection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixAwsTgwIntraDomainInspectionCreate,
		ReadContext:   resourceAviatrixAwsTgwIntraDomainInspectionRead,
		DeleteContext: resourceAviatrixAwsTgwIntraDomainInspectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW name.",
			},
			"route_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Route domain name.",
			},
			"firewall_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall domain name.",
			},
		},
	}
}

func marshalAwsTgwIntraDomainInspectionInput(d *schema.ResourceData) *goaviatrix.IntraDomainInspection {
	return &goaviatrix.IntraDomainInspection{
		TgwName:            d.Get("tgw_name").(string),
		RouteDomainName:    d.Get("route_domain_name").(string),
		FirewallDomainName: d.Get("firewall_domain_name").(string),
	}
}

func resourceAviatrixAwsTgwIntraDomainInspectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	intraDomainInspection := marshalAwsTgwIntraDomainInspectionInput(d)

	err := client.EnableIntraDomainInspection(ctx, intraDomainInspection)
	if err != nil {
		return diag.Errorf("could not enable intra domain inspection: %v", err)
	}

	d.SetId(intraDomainInspection.TgwName + "~" + intraDomainInspection.RouteDomainName)
	return resourceAviatrixAwsTgwIntraDomainInspectionRead(ctx, d, meta)
}

func resourceAviatrixAwsTgwIntraDomainInspectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("tgw_name") == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("invalid ID format, expected ID in format tgw_name~route_domain_name, instead got %s", d.Id())
		}

		tgwName := parts[0]
		routeDomainName := parts[1]

		if tgwName == "" || routeDomainName == "" {
			return diag.Errorf("tgw_name or route_domain_name cannot be empty")
		}

		d.Set("tgw_name", tgwName)
		d.Set("route_domain_name", routeDomainName)

		d.SetId(tgwName + "~" + routeDomainName)
	}

	intraDomainInspection := marshalAwsTgwIntraDomainInspectionInput(d)

	err := client.GetIntraDomainInspectionStatus(ctx, intraDomainInspection)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not get intra domain inspection status: %v", err)
	}

	d.Set("firewall_domain_name", intraDomainInspection.FirewallDomainName)

	d.SetId(intraDomainInspection.TgwName + "~" + intraDomainInspection.RouteDomainName)
	return nil
}

func resourceAviatrixAwsTgwIntraDomainInspectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	intraDomainInspection := marshalAwsTgwIntraDomainInspectionInput(d)

	err := client.DisableIntraDomainInspection(ctx, intraDomainInspection)
	if err != nil {
		return diag.Errorf("failed to disable intra domain inspection: %v", err)
	}

	return nil
}
