package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSegmentationSecurityDomain() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Resource 'aviatrix_segmentation_security_domain' will be deprecated in future releases. Please use resource 'aviatrix_segmentation_network_domain' instead.",
		Create:             resourceAviatrixSegmentationSecurityDomainCreate,
		Read:               resourceAviatrixSegmentationSecurityDomainRead,
		Delete:             resourceAviatrixSegmentationSecurityDomainDelete,
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

	d.SetId(domain.DomainName)
	flag := false
	defer resourceAviatrixSegmentationSecurityDomainReadIfRequired(d, meta, &flag)

	if err := client.CreateSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not create security domain: %v", err)
	}

	return resourceAviatrixSegmentationSecurityDomainReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSegmentationSecurityDomainReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSegmentationSecurityDomainRead(d, meta)
	}
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
