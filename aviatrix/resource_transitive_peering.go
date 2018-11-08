package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTranspeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTranspeerCreate,
		Read:   resourceTranspeerRead,
		Update: resourceTranspeerUpdate,
		Delete: resourceTranspeerDelete,

		Schema: map[string]*schema.Schema{
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"nexthop": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"reachable_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTranspeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transpeer := &goaviatrix.Transpeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	log.Printf("[INFO] Creating Aviatrix transitive peering: %#v", transpeer)

	err := client.CreateTranspeer(transpeer)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Transitive peering: %s", err)
	}
	d.SetId(transpeer.Source + transpeer.Nexthop + transpeer.ReachableCidr)
	//return nil
	return resourceTranspeerRead(d, meta)
}

func resourceTranspeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transpeer := &goaviatrix.Transpeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}
	transpeer, err := client.GetTranspeer(transpeer)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find Aviatrix Transitive peering: %s", err)
	}

	d.Set("source", transpeer.Source)
	d.Set("nexthop", transpeer.Nexthop)
	d.Set("reachable_cidr", transpeer.ReachableCidr)
	return nil
}

func resourceTranspeerUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Aviatrix transitive peering cannot be updated - delete and create new one")
}

func resourceTranspeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transpeer := &goaviatrix.Transpeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix transpeer: %#v", transpeer)

	err := client.DeleteTranspeer(transpeer)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Transpeer: %s", err)
	}
	return nil
}
