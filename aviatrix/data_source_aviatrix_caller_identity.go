package aviatrix

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	client := mustClient(meta)

	log.Printf("[DEBUG] CID is '%s'", client.CID)

	d.SetId(time.Now().UTC().String())
	mustSet(d, "cid", client.CID)
	return nil
}
