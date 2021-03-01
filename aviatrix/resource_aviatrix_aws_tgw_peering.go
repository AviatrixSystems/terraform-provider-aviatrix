package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAWSTgwPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwPeeringCreate,
		Read:   resourceAviatrixAWSTgwPeeringRead,
		Delete: resourceAviatrixAWSTgwPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the first AWS tgw to make a peer pair.",
			},
			"tgw_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the second AWS tgw to make a peer pair.",
			},
		},
	}
}

func resourceAviatrixAWSTgwPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: d.Get("tgw_name1").(string),
		TgwName2: d.Get("tgw_name2").(string),
	}

	log.Printf("[INFO] Creating Aviatrix AWS tgw peering: %#v", awsTgwPeering)

	err := client.CreateAwsTgwPeering(awsTgwPeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS tgw peering: %s", err)
	}

	d.SetId(awsTgwPeering.TgwName1 + "~" + awsTgwPeering.TgwName2)
	return resourceAviatrixAWSTgwPeeringRead(d, meta)
}

func resourceAviatrixAWSTgwPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName1 := d.Get("tgw_name1").(string)
	tgwName2 := d.Get("tgw_name2").(string)

	if tgwName1 == "" || tgwName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		d.Set("tgw_name1", strings.Split(id, "~")[0])
		d.Set("tgw_name2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: d.Get("tgw_name1").(string),
		TgwName2: d.Get("tgw_name2").(string),
	}

	err := client.GetAwsTgwPeering(awsTgwPeering)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix AWS tgw peering: %s", err)
	}

	d.SetId(awsTgwPeering.TgwName1 + "~" + awsTgwPeering.TgwName2)
	return nil
}

func resourceAviatrixAWSTgwPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: d.Get("tgw_name1").(string),
		TgwName2: d.Get("tgw_name2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS tgw peering: %#v", awsTgwPeering)

	err := client.DeleteAwsTgwPeering(awsTgwPeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS tgw peering: %s", err)
	}

	return nil
}
