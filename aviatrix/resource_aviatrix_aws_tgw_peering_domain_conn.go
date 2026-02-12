package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgwPeeringDomainConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwPeeringDomainConnCreate,
		Read:   resourceAviatrixAWSTgwPeeringDomainConnRead,
		Delete: resourceAviatrixAWSTgwPeeringDomainConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name1": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringDomainConnTgwName1,
				Description:      "The AWS tgw name of the source domain to make a connection.",
			},
			"domain_name1": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringDomainConnDomainName1,
				Description:      "The name of the source domain to make a connection.",
			},
			"tgw_name2": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringDomainConnTgwName2,
				Description:      "The AWS tgw name of the destination domain to make a connection.",
			},
			"domain_name2": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringDomainConnDomainName2,
				Description:      "The name of the destination domain to make a connection.",
			},
		},
	}
}

func resourceAviatrixAWSTgwPeeringDomainConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	domainConn := &goaviatrix.DomainConn{
		TgwName1:    getString(d, "tgw_name1"),
		DomainName1: getString(d, "domain_name1"),
		TgwName2:    getString(d, "tgw_name2"),
		DomainName2: getString(d, "domain_name2"),
	}

	log.Printf("[INFO] Creating Aviatrix domain connection between tgw: %s and %s", domainConn.TgwName1, domainConn.TgwName2)

	d.SetId(domainConn.TgwName1 + ":" + domainConn.DomainName1 + "~" + domainConn.TgwName2 + ":" + domainConn.DomainName2)
	flag := false
	defer func() { _ = resourceAviatrixAWSTgwPeeringDomainConnReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateDomainConn(domainConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix domain connection between two tgws: %w", err)
	}

	return resourceAviatrixAWSTgwPeeringDomainConnReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAWSTgwPeeringDomainConnReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSTgwPeeringDomainConnRead(d, meta)
	}
	return nil
}

func resourceAviatrixAWSTgwPeeringDomainConnRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tgwName1 := getString(d, "tgw_name1")
	domainName1 := getString(d, "domain_name1")
	tgwName2 := getString(d, "tgw_name2")
	domainName2 := getString(d, "domain_name2")

	if tgwName1 == "" || domainName1 == "" || tgwName2 == "" || domainName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		tgwDomain1 := strings.Split(id, "~")[0]
		tgwDomain2 := strings.Split(id, "~")[1]
		mustSet(d, "tgw_name1", strings.Split(tgwDomain1, ":")[0])
		mustSet(d, "domain_name1", strings.Split(tgwDomain1, ":")[1])
		mustSet(d, "tgw_name2", strings.Split(tgwDomain2, ":")[0])
		mustSet(d, "domain_name2", strings.Split(tgwDomain2, ":")[1])
		d.SetId(id)
	}

	domainConn := &goaviatrix.DomainConn{
		TgwName1:    getString(d, "tgw_name1"),
		DomainName1: getString(d, "domain_name1"),
		TgwName2:    getString(d, "tgw_name2"),
		DomainName2: getString(d, "domain_name2"),
	}

	err := client.GetDomainConn(domainConn)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix domain connection: %w", err)
	}

	d.SetId(domainConn.TgwName1 + ":" + domainConn.DomainName1 + "~" + domainConn.TgwName2 + ":" + domainConn.DomainName2)
	return nil
}

func resourceAviatrixAWSTgwPeeringDomainConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	domainConn := &goaviatrix.DomainConn{
		TgwName1:    getString(d, "tgw_name1"),
		DomainName1: getString(d, "domain_name1"),
		TgwName2:    getString(d, "tgw_name2"),
		DomainName2: getString(d, "domain_name2"),
	}

	log.Printf("[INFO] Deleting Aviatrix domain connection: %#v", domainConn)

	err := client.DeleteDomainConn(domainConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix domain connection: %w", err)
	}

	return nil
}
