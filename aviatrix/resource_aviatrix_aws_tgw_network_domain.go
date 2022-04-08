package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixAwsTgwNetworkDomain() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixAwsTgwNetworkDomainCreate,
		ReadWithoutTimeout:   resourceAviatrixAwsTgwNetworkDomainRead,
		DeleteWithoutTimeout: resourceAviatrixAwsTgwNetworkDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Network domain name.",
				ValidateFunc: validation.StringDoesNotContainAny(":"),
			},
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW name.",
			},
			"aviatrix_firewall": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set to true if the network domain is an aviatrix firewall domain.",
			},
			"native_egress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set to true if the network domain is a native egress domain.",
			},
			"native_firewall": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set to true if the network domain is a native firewall domain.",
			},
		},
	}
}

func marshalNetworkDomainInput(d *schema.ResourceData) *goaviatrix.SecurityDomain {
	networkDomain := &goaviatrix.SecurityDomain{
		Name:                   d.Get("name").(string),
		AwsTgwName:             d.Get("tgw_name").(string),
		AviatrixFirewallDomain: d.Get("aviatrix_firewall").(bool),
		NativeEgressDomain:     d.Get("native_egress").(bool),
		NativeFirewallDomain:   d.Get("native_firewall").(bool),
	}

	return networkDomain
}

func resourceAviatrixAwsTgwNetworkDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	networkDomain := marshalNetworkDomainInput(d)

	num := 0
	if networkDomain.AviatrixFirewallDomain {
		num += 1
	}
	if networkDomain.NativeEgressDomain {
		num += 1
	}
	if networkDomain.NativeFirewallDomain {
		num += 1
	}
	if num > 1 {
		return diag.Errorf("only one or none of 'firewall_domain', 'native_egress' and 'native_firewall' could be set true")
	}

	d.SetId(networkDomain.AwsTgwName + "~" + networkDomain.Name)
	flag := false
	defer resourceAviatrixAwsTgwNetworkDomainReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateSecurityDomain(networkDomain); err != nil {
		return diag.Errorf("could not create security domain: %v", err)
	}

	return resourceAviatrixAwsTgwNetworkDomainReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixAwsTgwNetworkDomainReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAwsTgwNetworkDomainRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixAwsTgwNetworkDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)

	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("invalid ID, expected ID tgw_name~domain_name, instead got %s", d.Id())
		}
		d.Set("tgw_name", parts[0])
		d.Set("name", parts[1])
		d.SetId(id)
	}

	name = d.Get("name").(string)
	tgwName := d.Get("tgw_name").(string)

	networkDomain := &goaviatrix.SecurityDomain{
		Name:       name,
		AwsTgwName: tgwName,
	}

	networkDomainDetails, err := client.GetSecurityDomainDetails(ctx, networkDomain)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("couldn't get the details of the network domain %s due to %v", name, err)
	}

	d.Set("aviatrix_firewall", networkDomainDetails.AviatrixFirewallDomain)
	d.Set("native_egress", networkDomainDetails.NativeEgressDomain)
	d.Set("native_firewall", networkDomainDetails.NativeFirewallDomain)

	d.SetId(tgwName + "~" + name)
	return nil
}

func resourceAviatrixAwsTgwNetworkDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	networkDomain := &goaviatrix.SecurityDomain{
		Name:       d.Get("name").(string),
		AwsTgwName: d.Get("tgw_name").(string),
	}

	defaultDomains := []string{"Aviatrix_Edge_Domain", "Default_Domain", "Shared_Service_Domain"}

	for _, d := range defaultDomains {
		if networkDomain.Name == d {
			networkDomain.ForceDelete = true
		}
	}

	if err := client.DeleteSecurityDomain(networkDomain); err != nil {
		return diag.Errorf("could not delete network domain: %v", err)
	}

	return nil
}
