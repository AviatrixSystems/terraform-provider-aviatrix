package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSegmentationSecurityDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationSecurityDomainCreate,
		Read:   resourceAviatrixSegmentationSecurityDomainRead,
		Delete: resourceAviatrixSegmentationSecurityDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security domain name.",
			},
		},
	}
}

func marshalSegmentationSecurityDomainInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomain {
	return &goaviatrix.SegmentationSecurityDomain{
		DomainName: d.Get("domain_name").(string),
	}
}

func resourceAviatrixSegmentationSecurityDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domain := marshalSegmentationSecurityDomainInput(d)

	if err := client.CreateSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not create security domain: %v", err)
	}

	d.SetId(domain.DomainName)
	return nil
}

func resourceAviatrixSegmentationSecurityDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainName := d.Get("domain_name").(string)
	if domainName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no segmentation_security_domain domain_name received. Import Id is %s", id)
		d.SetId(id)
		domainName = id
	}

	domain := &goaviatrix.SegmentationSecurityDomain{
		DomainName: domainName,
	}

	domain, err := client.GetSegmentationSecurityDomain(domain)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_security_domain %s: %v", domainName, err)
	}

	d.Set("domain_name", domain.DomainName)
	d.SetId(domain.DomainName)
	return nil
}

func resourceAviatrixSegmentationSecurityDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domain := marshalSegmentationSecurityDomainInput(d)

	if err := client.DeleteSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not delete segmentation_security_domain: %v", err)
	}

	return nil
}
