package aviatrix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixAwsTgwSecurityDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixAwsTgwSecurityDomainCreate,
		ReadContext:   resourceAviatrixAwsTgwSecurityDomainRead,
		DeleteContext: resourceAviatrixAwsTgwSecurityDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Security domain name.",
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
				Description: "Set to true if the security domain is an aviatrix firewall domain.",
			},
			"native_egress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set to true if the security domain is a native egress domain.",
			},
			"native_firewall": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set to true if the security domain is a native firewall domain.",
			},
		},
	}
}

func marshalSecurityDomainInput(d *schema.ResourceData) *goaviatrix.SecurityDomain {
	securityDomain := &goaviatrix.SecurityDomain{
		Name:                   d.Get("name").(string),
		AwsTgwName:             d.Get("tgw_name").(string),
		AviatrixFirewallDomain: d.Get("aviatrix_firewall").(bool),
		NativeEgressDomain:     d.Get("native_egress").(bool),
		NativeFirewallDomain:   d.Get("native_firewall").(bool),
	}

	return securityDomain
}

func resourceAviatrixAwsTgwSecurityDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	securityDomain := marshalSecurityDomainInput(d)

	num := 0
	if securityDomain.AviatrixFirewallDomain {
		num += 1
	}
	if securityDomain.NativeEgressDomain {
		num += 1
	}
	if securityDomain.NativeFirewallDomain {
		num += 1
	}
	if num > 1 {
		return diag.Errorf("only one or none of 'firewall_domain', 'native_egress' and 'native_firewall' could be set true")
	}

	if err := client.CreateSecurityDomain(securityDomain); err != nil {
		return diag.Errorf("could not create security domain: %v", err)
	}

	d.SetId(securityDomain.AwsTgwName + "~" + securityDomain.Name)
	return resourceAviatrixAwsTgwSecurityDomainRead(ctx, d, meta)
}

func resourceAviatrixAwsTgwSecurityDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)

	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("invalid ID, expected ID tgw_name~domain_name, instead got %s", d.Id())
		}
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	name = d.Get("name").(string)
	tgwName := d.Get("tgw_name").(string)

	securityDomain := &goaviatrix.SecurityDomain{
		Name:       name,
		AwsTgwName: tgwName,
	}

	securityDomainRule, err := client.GetSecurityDomainDetails(securityDomain)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("couldn't get the details of the security domain %s due to %v", name, err)
	}

	d.Set("aviatrix_firewall", securityDomainRule.AviatrixFirewallDomain)
	d.Set("native_egress", securityDomainRule.NativeEgressDomain)
	d.Set("native_firewall", securityDomainRule.NativeFirewallDomain)

	d.SetId(tgwName + "~" + name)
	return nil
}

func resourceAviatrixAwsTgwSecurityDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	securityDomain := &goaviatrix.SecurityDomain{
		Name:       d.Get("name").(string),
		AwsTgwName: d.Get("tgw_name").(string),
	}

	defaultDomains := []string{"Aviatrix_Edge_Domain", "Default_Domain", "Shared_Service_Domain"}

	for _, d := range defaultDomains {
		if securityDomain.Name == d {
			securityDomain.ForceDelete = true
		}
	}

	if err := client.DeleteSecurityDomain(securityDomain); err != nil {
		return diag.Errorf("could not delete security domain: %v", err)
	}

	return nil
}
