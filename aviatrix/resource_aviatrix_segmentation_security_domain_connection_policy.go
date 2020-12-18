package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSegmentationSecurityDomainConnectionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSegmentationSecurityDomainConnectionPolicyCreate,
		Read:   resourceAviatrixSegmentationSecurityDomainConnectionPolicyRead,
		Delete: resourceAviatrixSegmentationSecurityDomainConnectionPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain_name_1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of security domain that will be connected to domain 2.",
			},
			"domain_name_2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of security domain that will be connected to domain 1.",
			},
		},
	}
}

func marshalSegmentationSecurityDomainConnectionPolicyInput(d *schema.ResourceData) *goaviatrix.SegmentationSecurityDomainConnectionPolicy {
	return &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
		Domain1: &goaviatrix.SegmentationSecurityDomain{
			DomainName: d.Get("domain_name_1").(string),
		},
		Domain2: &goaviatrix.SegmentationSecurityDomain{
			DomainName: d.Get("domain_name_2").(string),
		},
	}
}

func resourceAviatrixSegmentationSecurityDomainConnectionPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	policy := marshalSegmentationSecurityDomainConnectionPolicyInput(d)

	if err := client.CreateSegmentationSecurityDomainConnectionPolicy(policy); err != nil {
		return fmt.Errorf("could not create security domain connection policy: %v", err)
	}

	id := policy.Domain1.DomainName + "~" + policy.Domain2.DomainName
	d.SetId(id)
	return nil
}

func resourceAviatrixSegmentationSecurityDomainConnectionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainName1 := d.Get("domain_name_1").(string)
	domainName2 := d.Get("domain_name_2").(string)
	if domainName1 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no segmentation_security_domain_connection_policy domain_name received. Import Id is %s", id)
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
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find segmentation_security_domain_connection_policy %s: %v", domainName1+"~"+domainName2, err)
	}

	d.Set("domain_name_1", domainName1)
	d.Set("domain_name_2", domainName2)
	d.SetId(domainName1 + "~" + domainName2)
	return nil
}

func resourceAviatrixSegmentationSecurityDomainConnectionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	policy := marshalSegmentationSecurityDomainConnectionPolicyInput(d)

	if err := client.DeleteSegmentationSecurityDomainConnectionPolicy(policy); err != nil {
		return fmt.Errorf("could not delete segmentation_security_domain_connection_policy: %v", err)
	}

	return nil
}
