package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSegmentationNetworkDomainConnectionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationNetworkDomainConnectionPolicyCreate,
		Read:   resourceAviatrixSegmentationNetworkDomainConnectionPolicyRead,
		Delete: resourceAviatrixSegmentationNetworkDomainConnectionPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"domain_name_1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of network domain that will be connected to domain 2.",
			},
			"domain_name_2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of network domain that will be connected to domain 1.",
			},
		},
	}
}

func marshalSegmentationNetworkDomainConnectionPolicyInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomainConnectionPolicy {
	return &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
		Domain1: &goaviatrix.SegmentationSecurityDomain{
			DomainName: getString(d, "domain_name_1"),
		},
		Domain2: &goaviatrix.SegmentationSecurityDomain{
			DomainName: getString(d, "domain_name_2"),
		},
	}
}

func resourceAviatrixSegmentationNetworkDomainConnectionPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	policy := marshalSegmentationNetworkDomainConnectionPolicyInput(d)

	// validate domain names exist
	domainNames, err := client.ListSegmentationSecurityDomains()
	if err != nil {
		return fmt.Errorf("could not find segmentation_network_domains: %w", err)
	}
	for _, domainName := range []string{policy.Domain1.DomainName, policy.Domain2.DomainName} {
		if !slices.Contains(domainNames, domainName) {
			return fmt.Errorf("could not find segmentation_network_domain %s in %v", domainName, domainNames)
		}
	}

	d.SetId(policy.Domain1.DomainName + "~" + policy.Domain2.DomainName)
	flag := false
	defer func() { _ = resourceAviatrixSegmentationNetworkDomainConnectionPolicyReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if err := client.CreateSegmentationSecurityDomainConnectionPolicy(policy); err != nil {
		return fmt.Errorf("could not create network domain connection policy: %w", err)
	}

	return resourceAviatrixSegmentationNetworkDomainConnectionPolicyReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSegmentationNetworkDomainConnectionPolicyReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSegmentationNetworkDomainConnectionPolicyRead(d, meta)
	}
	return nil
}

func resourceAviatrixSegmentationNetworkDomainConnectionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	domainName1 := getString(d, "domain_name_1")
	domainName2 := getString(d, "domain_name_2")
	if domainName1 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no segmentation_network_domain_connection_policy domain_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		domainName1 = parts[0]
		domainName2 = parts[1]
	}

	policy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
		Domain1: &goaviatrix.SegmentationSecurityDomain{
			DomainName: domainName1,
		},
		Domain2: &goaviatrix.SegmentationSecurityDomain{
			DomainName: domainName2,
		},
	}

	_, err := client.GetSegmentationSecurityDomainConnectionPolicy(policy)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_network_domain_connection_policy %s: %w", domainName1+"~"+domainName2, err)
	}
	mustSet(d, "domain_name_1", domainName1)
	mustSet(d, "domain_name_2", domainName2)
	d.SetId(domainName1 + "~" + domainName2)
	return nil
}

func resourceAviatrixSegmentationNetworkDomainConnectionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	policy := marshalSegmentationNetworkDomainConnectionPolicyInput(d)

	if err := client.DeleteSegmentationSecurityDomainConnectionPolicy(policy); err != nil {
		return fmt.Errorf("could not delete segmentation_network_domain_connection_policy: %w", err)
	}

	return nil
}
