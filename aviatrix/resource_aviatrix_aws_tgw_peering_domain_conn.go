package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgwPeeringDomainConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwPeeringDomainConnCreate,
		Read:   resourceAviatrixAWSTgwPeeringDomainConnRead,
		Delete: resourceAviatrixAWSTgwPeeringDomainConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AWS tgw name of the source domain to make a connection.",
			},
			"domain_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the source domain to make a connection.",
			},
			"tgw_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AWS tgw name of the destination domain to make a connection.",
			},
			"domain_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the destination domain to make a connection.",
			},
		},
	}
}

func resourceAviatrixAWSTgwPeeringDomainConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainConn := &goaviatrix.DomainConn{
		TgwName1:    d.Get("tgw_name1").(string),
		DomainName1: d.Get("domain_name1").(string),
		TgwName2:    d.Get("tgw_name2").(string),
		DomainName2: d.Get("domain_name2").(string),
	}

	log.Printf("[INFO] Creating Aviatrix domain connection between tgw: %s and %s", domainConn.TgwName1, domainConn.TgwName2)

	err := client.CreateDomainConn(domainConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix domain connection between two tgws: %s", err)
	}

	d.SetId(domainConn.TgwName1 + ":" + domainConn.DomainName1 + "~" + domainConn.TgwName2 + ":" + domainConn.DomainName2)
	return resourceAviatrixAWSTgwPeeringDomainConnRead(d, meta)
}

func resourceAviatrixAWSTgwPeeringDomainConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName1 := d.Get("tgw_name1").(string)
	domainName1 := d.Get("domain_name1").(string)
	tgwName2 := d.Get("tgw_name2").(string)
	domainName2 := d.Get("domain_name2").(string)

	if tgwName1 == "" || domainName1 == "" || tgwName2 == "" || domainName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		tgwDomain1 := strings.Split(id, "~")[0]
		tgwDomain2 := strings.Split(id, "~")[1]
		d.Set("tgw_name1", strings.Split(tgwDomain1, ":")[0])
		d.Set("domain_name1", strings.Split(tgwDomain1, ":")[1])
		d.Set("tgw_name2", strings.Split(tgwDomain2, ":")[0])
		d.Set("domain_name2", strings.Split(tgwDomain2, ":")[1])
		d.SetId(id)
	}

	domainConn := &goaviatrix.DomainConn{
		TgwName1:    d.Get("tgw_name1").(string),
		DomainName1: d.Get("domain_name1").(string),
		TgwName2:    d.Get("tgw_name2").(string),
		DomainName2: d.Get("domain_name2").(string),
	}

	err := client.GetDomainConn(domainConn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix domain connection: %s", err)
	}

	d.SetId(domainConn.TgwName1 + ":" + domainConn.DomainName1 + "~" + domainConn.TgwName2 + ":" + domainConn.DomainName2)
	return nil
}

func resourceAviatrixAWSTgwPeeringDomainConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainConn := &goaviatrix.DomainConn{
		TgwName1:    d.Get("tgw_name1").(string),
		DomainName1: d.Get("domain_name1").(string),
		TgwName2:    d.Get("tgw_name2").(string),
		DomainName2: d.Get("domain_name2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix domain connection: %#v", domainConn)

	err := client.DeleteDomainConn(domainConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix domain connection: %s", err)
	}

	return nil
}
