package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixAwsTgwSecurityDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsTgwSecurityDomainCreate,
		Read:   resourceAviatrixAwsTgwSecurityDomainRead,
		Delete: resourceAviatrixAwsTgwSecurityDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceAviatrixAwsTgwSecurityDomainCreate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("only one or none of 'firewall_domain', 'native_egress' and 'native_firewall' could be set true")
	}

	if err := client.CreateSecurityDomain(securityDomain); err != nil {
		return fmt.Errorf("could not create security domain: %v", err)
	}

	d.SetId(securityDomain.AwsTgwName + "~" + securityDomain.Name)
	return resourceAviatrixAwsTgwSecurityDomainRead(d, meta)
}

func resourceAviatrixAwsTgwSecurityDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)

	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	name = d.Get("name").(string)
	tgwName := d.Get("tgw_name").(string)

	awsTgw := &goaviatrix.AWSTgw{
		Name: tgwName,
	}

	awsTgw, err := client.GetAWSTgw(awsTgw)
	if err != nil {
		return fmt.Errorf("couldn't find AWS TGW %s: %v", tgwName, err)
	}

	notFound := true
	for _, sd := range awsTgw.SecurityDomains {
		if sd.Name == name {
			d.Set("aviatrix_firewall", sd.AviatrixFirewallDomain)
			d.Set("native_egress", sd.NativeEgressDomain)
			d.Set("native_firewall", sd.NativeFirewallDomain)
			notFound = false
			break
		}
	}

	if notFound {
		d.SetId("")
		return nil
	}

	d.SetId(tgwName + "~" + name)
	return nil
}

func resourceAviatrixAwsTgwSecurityDomainDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("could not delete security domain: %v", err)
	}

	return nil
}
