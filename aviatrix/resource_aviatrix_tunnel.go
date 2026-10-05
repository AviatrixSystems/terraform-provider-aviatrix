package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTunnelCreate,
		Read:   resourceAviatrixTunnelRead,
		Update: resourceAviatrixTunnelUpdate,
		Delete: resourceAviatrixTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"gw_name1": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The first VPC Container name to make a peer pair.",
			},
			"gw_name2": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The second VPC Container name to make a peer pair.",
			},
			"enable_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether Peering HA is enabled. Valid inputs: true or false.",
			},
			"peering_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the tunnel.",
			},
			"peering_hastatus": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the HA tunnel.",
			},
			"peering_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the peering link.",
			},
		},
	}
}

func resourceAviatrixTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tunnel := &goaviatrix.Tunnel{
		VpcName1:        getString(d, "gw_name1"),
		VpcName2:        getString(d, "gw_name2"),
		PeeringState:    getString(d, "peering_state"),
		PeeringHaStatus: getString(d, "peering_hastatus"),
		PeeringLink:     getString(d, "peering_link"),
	}

	enableHA := getBool(d, "enable_ha")
	if enableHA {
		tunnel.EnableHA = "yes"
	} else {
		tunnel.EnableHA = "no"
	}

	log.Printf("[INFO] Creating Aviatrix tunnel: %#v", tunnel)

	d.SetId(tunnel.VpcName1 + "~" + tunnel.VpcName2)
	flag := false
	defer func() { _ = resourceAviatrixTunnelReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Tunnel: %w", err)
	}

	return resourceAviatrixTunnelReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTunnelReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTunnelRead(d, meta)
	}
	return nil
}

func resourceAviatrixTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	vpcName1 := getString(d, "gw_name1")
	vpcName2 := getString(d, "gw_name2")

	if vpcName1 == "" || vpcName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		mustSet(d, "gw_name1", strings.Split(id, "~")[0])
		mustSet(d, "gw_name2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	tunnel := &goaviatrix.Tunnel{
		VpcName1: getString(d, "gw_name1"),
		VpcName2: getString(d, "gw_name2"),
	}
	tun, err := client.GetTunnel(tunnel)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Tunnel: %w", err)
	}
	log.Printf("[INFO] Found Aviatrix tunnel: %#v", tun)
	mustSet(d, "peering_hastatus", tun.PeeringHaStatus)
	mustSet(d, "peering_state", tun.PeeringState)
	mustSet(d, "peering_link", tun.PeeringLink)

	if tun.PeeringHaStatus == "active" {
		mustSet(d, "enable_ha", true)
	} else {
		mustSet(d, "enable_ha", false)
	}

	d.SetId(tun.VpcName1 + "~" + tun.VpcName2)
	return nil
}

func resourceAviatrixTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tunnel := &goaviatrix.Tunnel{
		VpcName1:        getString(d, "gw_name1"),
		VpcName2:        getString(d, "gw_name2"),
		PeeringState:    getString(d, "peering_state"),
		PeeringHaStatus: getString(d, "peering_hastatus"),
		PeeringLink:     getString(d, "peering_link"),
	}

	log.Printf("[INFO] Updating Aviatrix tunnel: %#v", tunnel)

	err := client.UpdateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Tunnel: %w", err)
	}

	d.SetId(tunnel.VpcName1 + "~" + tunnel.VpcName2)
	return resourceAviatrixTunnelRead(d, meta)
}

func resourceAviatrixTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tunnel := &goaviatrix.Tunnel{
		VpcName1: getString(d, "gw_name1"),
		VpcName2: getString(d, "gw_name2"),
	}

	log.Printf("[INFO] Deleting Aviatrix tunnel: %#v", tunnel)

	if peeringHaStatus := getString(d, "peering_hastatus"); peeringHaStatus == "active" {
		// parse the hagw name
		tunnel.VpcName1 += "-hagw"
		tunnel.VpcName2 += "-hagw"
		err := client.DeleteTunnel(tunnel)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix HA gateway: %w", err)
		}
	}

	tunnel.VpcName1 = getString(d, "gw_name1")
	tunnel.VpcName2 = getString(d, "gw_name2")

	err := client.DeleteTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Tunnel: %w", err)
	}

	return nil
}
