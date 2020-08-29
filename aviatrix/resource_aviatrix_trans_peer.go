package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransPeerCreate,
		Read:   resourceAviatrixTransPeerRead,
		Update: resourceAviatrixTransPeerUpdate,
		Delete: resourceAviatrixTransPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of source gateway.",
			},
			"nexthop": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of nexthop gateway.",
			},
			"reachable_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination CIDR.",
			},
			"source_original_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the source gateway when it was created.",
			},
			"nexthop_original_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the nexthop gateway when it was created.",
			},
		},
	}
}

func resourceAviatrixTransPeerCreate(d *schema.ResourceData, meta interface{}) error {
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
	return resourceAviatrixTransPeerRead(d, meta)
}

func resourceAviatrixTransPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	sourceGw := d.Get("source").(string)
	nextHopGw := d.Get("nexthop").(string)
	reachableCIDR := d.Get("reachable_cidr").(string)

	if sourceGw == "" || nextHopGw == "" || reachableCIDR == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names or reachable cidr received. "+
			"Import Id is %s", id)
		d.Set("source", strings.Split(id, "~")[0])
		d.Set("nexthop", strings.Split(id, "~")[1])
		d.Set("reachable_cidr", strings.Split(id, "~")[2])
		sourceGw = strings.Split(id, "~")[0]
		nextHopGw = strings.Split(id, "~")[1]
		reachableCIDR = strings.Split(id, "~")[2]
		d.SetId(id)
	}

	transPeer := &goaviatrix.TransPeer{
		Source:              sourceGw,
		Nexthop:             nextHopGw,
		ReachableCidr:       reachableCIDR,
		SourceOriginalName:  d.Get("source_original_name").(string),
		NexthopOriginalName: d.Get("nexthop_original_name").(string),
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
	d.Set("source_original_name", transPeer.SourceOriginalName)
	d.Set("nexthop_original_name", transPeer.NexthopOriginalName)

	return nil
}

func resourceAviatrixTransPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("source") {
		_, gwNameNew := d.GetChange("source")
		gateway := &goaviatrix.Gateway{
			GwName:         gwNameNew.(string),
			GwOriginalName: d.Get("source_original_name").(string),
		}
		err := client.IsGatewayNameUpdatable(gateway)
		if err != nil {
			return nil
		}
	}

	if d.HasChange("nexthop") {
		_, gwNameNew := d.GetChange("nexthop")
		gateway := &goaviatrix.Gateway{
			GwName:         gwNameNew.(string),
			GwOriginalName: d.Get("nexthop_original_name").(string),
		}
		err := client.IsGatewayNameUpdatable(gateway)
		if err != nil {
			return err
		}
	}

	transPeer := &goaviatrix.TransPeer{
		Source:        d.Get("source").(string),
		Nexthop:       d.Get("nexthop").(string),
		ReachableCidr: d.Get("reachable_cidr").(string),
	}

	d.SetId(transPeer.Source + "~" + transPeer.Nexthop + "~" + transPeer.ReachableCidr)
	return resourceAviatrixTransPeerRead(d, meta)
}

func resourceAviatrixTransPeerDelete(d *schema.ResourceData, meta interface{}) error {
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
