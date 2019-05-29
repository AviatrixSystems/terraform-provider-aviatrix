package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceTransPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransPeerCreate,
		Read:   resourceTransPeerRead,
		Update: resourceTransPeerUpdate,
		Delete: resourceTransPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Source gateway.",
			},
			"nexthop": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of nexthop gateway.",
			},
			"reachable_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination CIDR.",
			},
		},
	}
}

func resourceTransPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	log.Printf("[INFO] Creating Aviatrix transitive peering: %#v", transPeer)

	err := client.CreateTransPeer(transPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transitive peering: %s", err)
	}
	d.SetId(transPeer.Source + "~" + transPeer.Nexthop + "~" + transPeer.ReachableCidr)
	//return nil
	return resourceTransPeerRead(d, meta)
}

func resourceTransPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	sourceGw := d.Get("source").(string)
	nestHopGw := d.Get("nexthop").(string)
	reachableCIDR := d.Get("reachable_cidr").(string)

	if sourceGw == "" || nestHopGw == "" || reachableCIDR == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names or reachable cidr received. "+
			"Import Id is %s", id)
		d.Set("source", strings.Split(id, "~")[0])
		d.Set("nexthop", strings.Split(id, "~")[1])
		d.Set("reachable_cidr", strings.Split(id, "~")[2])
		d.SetId(id)
	}

	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}
	transPeer, err := client.GetTransPeer(transPeer)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transitive peering: %s", err)
	}

	d.Set("source", transPeer.Source)
	d.Set("nexthop", transPeer.Nexthop)
	d.Set("reachable_cidr", transPeer.ReachableCidr)
	return nil
}

func resourceTransPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("aviatrix transitive peering cannot be updated - delete and create new one")
}

func resourceTransPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix transpeer: %#v", transPeer)

	err := client.DeleteTransPeer(transPeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transpeer: %s", err)
	}
	return nil
}
