package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSecurityDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityDomainCreate,
		Read:   resourceSecurityDomainRead,
		Update: resourceSecurityDomainUpdate,
		Delete: resourceSecurityDomainDelete,

		Schema: map[string]*schema.Schema{
			"route_domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_tgw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	securityDomain := &goaviatrix.SecurityDomain{
		Name:        d.Get("route_domain_name").(string),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		AwsTgwName:  d.Get("aws_tgw_name").(string),
	}

	log.Printf("[INFO] Creating Security Domain: %#v", securityDomain)

	err := client.CreateSecurityDomain(securityDomain)
	if err != nil {
		return fmt.Errorf("failed to create Security Domain: %s", err)
	}

	d.SetId(securityDomain.Name)

	return resourceSecurityDomainRead(d, meta)
}

func resourceSecurityDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	securityDomain := &goaviatrix.SecurityDomain{
		Name:        d.Get("route_domain_name").(string),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		AwsTgwName:  d.Get("aws_tgw_name").(string),
	}

	securityDomainResp, err := client.GetSecurityDomain(securityDomain)
	if err != nil {
		return fmt.Errorf("couldn't find Security Domain: %s", err)
	}
	log.Printf("[TRACE] reading Security Domain %s: %#v", d.Get("route_domain_name").(string),
		securityDomainResp)

	d.Set("account_name", securityDomain.AccountName)
	d.Set("aws_tgw_name", securityDomain.AwsTgwName)
	d.Set("region", securityDomain.Region)
	d.Set("route_domain_name", securityDomain.Name)

	return nil
}

func resourceSecurityDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceSecurityDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	securityDomain := &goaviatrix.SecurityDomain{
		Name:        d.Get("route_domain_name").(string),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		AwsTgwName:  d.Get("aws_tgw_name").(string),
	}

	err := client.DeleteSecurityDomain(securityDomain)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Security Domain: %s", err)
	}

	return nil
}
