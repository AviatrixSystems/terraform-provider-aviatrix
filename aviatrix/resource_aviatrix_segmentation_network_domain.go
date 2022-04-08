package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSegmentationNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationNetworkDomainCreate,
		Read:   resourceAviatrixSegmentationNetworkDomainRead,
		Delete: resourceAviatrixSegmentationNetworkDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network domain name.",
			},
		},
	}
}

func marshalSegmentationNetworkDomainInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomain {
	return &goaviatrix.SegmentationSecurityDomain{
		DomainName: d.Get("domain_name").(string),
	}
}

func resourceAviatrixSegmentationNetworkDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domain := marshalSegmentationNetworkDomainInput(d)

	d.SetId(domain.DomainName)
	flag := false
	defer resourceAviatrixSegmentationNetworkDomainReadIfRequired(d, meta, &flag)

	if err := client.CreateSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not create network domain: %v", err)
	}

	return resourceAviatrixSegmentationNetworkDomainReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSegmentationNetworkDomainReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSegmentationNetworkDomainRead(d, meta)
	}
	return nil
}

func resourceAviatrixSegmentationNetworkDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainName := d.Get("domain_name").(string)
	if domainName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no segmentation_network_domain domain_name received. Import Id is %s", id)
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
		return fmt.Errorf("could not find segmentation_network_domain %s: %v", domainName, err)
	}

	d.Set("domain_name", domain.DomainName)
	d.SetId(domain.DomainName)
	return nil
}

func resourceAviatrixSegmentationNetworkDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domain := marshalSegmentationNetworkDomainInput(d)

	if err := client.DeleteSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not delete segmentation_network_domain: %v", err)
	}

	return nil
}
