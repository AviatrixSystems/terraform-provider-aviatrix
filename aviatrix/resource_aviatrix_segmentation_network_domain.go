package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSegmentationNetworkDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationNetworkDomainCreate,
		Read:   resourceAviatrixSegmentationNetworkDomainRead,
		Delete: resourceAviatrixSegmentationNetworkDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
		DomainName: getString(d, "domain_name"),
	}
}

func resourceAviatrixSegmentationNetworkDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	domain := marshalSegmentationNetworkDomainInput(d)

	d.SetId(domain.DomainName)
	flag := false
	defer func() { _ = resourceAviatrixSegmentationNetworkDomainReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if err := client.CreateSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not create network domain: %w", err)
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
	client := mustClient(meta)

	domainName := getString(d, "domain_name")
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
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_network_domain %s: %w", domainName, err)
	}
	mustSet(d, "domain_name", domain.DomainName)
	d.SetId(domain.DomainName)
	return nil
}

func resourceAviatrixSegmentationNetworkDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	domain := marshalSegmentationNetworkDomainInput(d)

	if err := client.DeleteSegmentationSecurityDomain(domain); err != nil {
		return fmt.Errorf("could not delete segmentation_network_domain: %w", err)
	}

	return nil
}
