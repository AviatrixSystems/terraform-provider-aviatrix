package aviatrix

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixCallerIdentityRead,

		Schema: map[string]*schema.Schema{
			"cid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix caller identity.",
			},
		},
	}
}

func dataSourceAviatrixCallerIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[DEBUG] CID is '%s'", client.CID)
	d.SetId(time.Now().UTC().String())
	d.Set("cid", client.CID)
	return nil
}
