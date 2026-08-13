package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransPeerCreate,
		Read:   resourceAviatrixTransPeerRead,
		Delete: resourceAviatrixTransPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of Source gateway.",
			},
			"nexthop": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of nexthop gateway.",
			},
			"reachable_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination CIDR.",
			},
		},
	}
}

func resourceAviatrixTransPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	transPeer := &goaviatrix.TransPeer{
		Source:        getString(d, "source"),
		Nexthop:       getString(d, "nexthop"),
		ReachableCidr: getString(d, "reachable_cidr"),
	}

	log.Printf("[INFO] Creating Aviatrix transitive peering: %#v", transPeer)

	d.SetId(transPeer.Source + "~" + transPeer.Nexthop + "~" + transPeer.ReachableCidr)
	flag := false
	defer func() { _ = resourceAviatrixTransPeerReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateTransPeer(transPeer)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transitive peering: %w", err)
	}

	return resourceAviatrixTransPeerReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransPeerReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransPeerRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	sourceGw := getString(d, "source")
	nestHopGw := getString(d, "nexthop")
	reachableCIDR := getString(d, "reachable_cidr")

	if sourceGw == "" || nestHopGw == "" || reachableCIDR == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names or reachable cidr received. "+
			"Import Id is %s", id)
		mustSet(d, "source", strings.Split(id, "~")[0])
		mustSet(d, "nexthop", strings.Split(id, "~")[1])
		mustSet(d, "reachable_cidr", strings.Split(id, "~")[2])
		d.SetId(id)
	}

	transPeer := &goaviatrix.TransPeer{
		Source:        getString(d, "source"),
		Nexthop:       getString(d, "nexthop"),
		ReachableCidr: getString(d, "reachable_cidr"),
	}
	transPeer, err := client.GetTransPeer(transPeer)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transitive peering: %w", err)
	}
	mustSet(d, "source", transPeer.Source)
	mustSet(d, "nexthop", transPeer.Nexthop)
	mustSet(d, "reachable_cidr", transPeer.ReachableCidr)

	return nil
}

func resourceAviatrixTransPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	transPeer := &goaviatrix.TransPeer{
		Source:        getString(d, "source"),
		Nexthop:       getString(d, "nexthop"),
		ReachableCidr: getString(d, "reachable_cidr"),
	}

	log.Printf("[INFO] Deleting Aviatrix transpeer: %#v", transPeer)

	err := client.DeleteTransPeer(transPeer)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transpeer: %w", err)
	}

	return nil
}
