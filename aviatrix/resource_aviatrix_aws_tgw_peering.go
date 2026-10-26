package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgwPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwPeeringCreate,
		Read:   resourceAviatrixAWSTgwPeeringRead,
		Delete: resourceAviatrixAWSTgwPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name1": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringTgwName1,
				Description:      "Name of the first AWS tgw to make a peer pair.",
			},
			"tgw_name2": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncAwsTgwPeeringTgwName2,
				Description:      "Name of the second AWS tgw to make a peer pair.",
			},
		},
	}
}

func resourceAviatrixAWSTgwPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: getString(d, "tgw_name1"),
		TgwName2: getString(d, "tgw_name2"),
	}

	log.Printf("[INFO] Creating Aviatrix AWS tgw peering: %#v", awsTgwPeering)

	d.SetId(awsTgwPeering.TgwName1 + "~" + awsTgwPeering.TgwName2)
	flag := false
	defer func() { _ = resourceAviatrixAWSTgwPeeringReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateAwsTgwPeering(awsTgwPeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS tgw peering: %w", err)
	}

	return resourceAviatrixAWSTgwPeeringReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAWSTgwPeeringReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSTgwPeeringRead(d, meta)
	}
	return nil
}

func resourceAviatrixAWSTgwPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tgwName1 := getString(d, "tgw_name1")
	tgwName2 := getString(d, "tgw_name2")

	if tgwName1 == "" || tgwName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		mustSet(d, "tgw_name1", strings.Split(id, "~")[0])
		mustSet(d, "tgw_name2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: getString(d, "tgw_name1"),
		TgwName2: getString(d, "tgw_name2"),
	}

	err := client.GetAwsTgwPeering(awsTgwPeering)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix AWS tgw peering: %w", err)
	}

	d.SetId(awsTgwPeering.TgwName1 + "~" + awsTgwPeering.TgwName2)
	return nil
}

func resourceAviatrixAWSTgwPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	awsTgwPeering := &goaviatrix.AwsTgwPeering{
		TgwName1: getString(d, "tgw_name1"),
		TgwName2: getString(d, "tgw_name2"),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS tgw peering: %#v", awsTgwPeering)

	err := client.DeleteAwsTgwPeering(awsTgwPeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS tgw peering: %w", err)
	}

	return nil
}
