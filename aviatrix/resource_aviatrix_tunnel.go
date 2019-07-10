package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTunnelCreate,
		Read:   resourceAviatrixTunnelRead,
		Update: resourceAviatrixTunnelUpdate,
		Delete: resourceAviatrixTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"peering_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of the tunnel.",
			},
			"peering_hastatus": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of the HA tunnel.",
			},
			"peering_link": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the peering link.",
			},
			"enable_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether Peering HA is enabled. Valid inputs: true or false.",
			},
		},
	}
}

func resourceAviatrixTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tunnel := &goaviatrix.Tunnel{
		VpcName1:        d.Get("gw_name1").(string),
		VpcName2:        d.Get("gw_name2").(string),
		PeeringState:    d.Get("peering_state").(string),
		PeeringHaStatus: d.Get("peering_hastatus").(string),
		PeeringLink:     d.Get("peering_link").(string),
	}

	enableHA := d.Get("enable_ha").(bool)
	if enableHA {
		tunnel.EnableHA = "yes"
	} else {
		tunnel.EnableHA = "no"
	}

	log.Printf("[INFO] Creating Aviatrix tunnel: %#v", tunnel)

	err := client.CreateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Tunnel: %s", err)
	}

	d.SetId(tunnel.VpcName1 + "~" + tunnel.VpcName2)
	return resourceAviatrixTunnelRead(d, meta)
}

func resourceAviatrixTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcName1 := d.Get("gw_name1").(string)
	vpcName2 := d.Get("gw_name2").(string)

	if vpcName1 == "" || vpcName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		d.Set("gw_name1", strings.Split(id, "~")[0])
		d.Set("gw_name2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	tunnel := &goaviatrix.Tunnel{
		VpcName1: d.Get("gw_name1").(string),
		VpcName2: d.Get("gw_name2").(string),
	}
	tun, err := client.GetTunnel(tunnel)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Tunnel: %s", err)
	}
	log.Printf("[INFO] Found Aviatrix tunnel: %#v", tun)

	d.Set("peering_hastatus", tun.PeeringHaStatus)
	d.Set("peering_state", tun.PeeringState)
	d.Set("peering_link", tun.PeeringLink)

	if tun.PeeringHaStatus == "active" {
		d.Set("enable_ha", true)
	} else {
		d.Set("enable_ha", false)
	}

	d.SetId(tun.VpcName1 + "~" + tun.VpcName2)
	return nil
}

func resourceAviatrixTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tunnel := &goaviatrix.Tunnel{
		VpcName1:        d.Get("gw_name1").(string),
		VpcName2:        d.Get("gw_name2").(string),
		PeeringState:    d.Get("peering_state").(string),
		PeeringHaStatus: d.Get("peering_hastatus").(string),
		PeeringLink:     d.Get("peering_link").(string),
	}

	log.Printf("[INFO] Updating Aviatrix tunnel: %#v", tunnel)

	err := client.UpdateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Tunnel: %s", err)
	}

	d.SetId(tunnel.VpcName1 + "~" + tunnel.VpcName2)
	return resourceAviatrixTunnelRead(d, meta)
}

func resourceAviatrixTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tunnel := &goaviatrix.Tunnel{
		VpcName1: d.Get("gw_name1").(string),
		VpcName2: d.Get("gw_name2").(string),
		EnableHA: d.Get("enable_ha").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix tunnel: %#v", tunnel)
	if peeringHaStatus := d.Get("peering_hastatus").(string); peeringHaStatus == "active" {
		// parse the hagw name
		tunnel.VpcName1 += "-hagw"
		tunnel.VpcName2 += "-hagw"
		err := client.DeleteTunnel(tunnel)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix HA gateway: %s", err)
		}
	}

	tunnel.VpcName1 = d.Get("gw_name1").(string)
	tunnel.VpcName2 = d.Get("gw_name2").(string)

	err := client.DeleteTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Tunnel: %s", err)
	}

	return nil
}
