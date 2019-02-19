package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTunnelCreate,
		Read:   resourceTunnelRead,
		Update: resourceTunnelUpdate,
		Delete: resourceTunnelDelete,

		Schema: map[string]*schema.Schema{
			"vpc_name1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_name2": {
				Type:     schema.TypeString,
				Required: true,
			},
			"over_aws_peering": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peering_hastatus": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_link": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_ha": {
				Type:     schema.TypeString,
				Optional: true,
			},
			//FIXME : Some of the above are computed. Set them correctly. Boolean valus should not be Optional to
			// prevent tf state corruption
		},
	}
}

func resourceTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	tunnel := &goaviatrix.Tunnel{
		VpcName1:        d.Get("vpc_name1").(string),
		VpcName2:        d.Get("vpc_name2").(string),
		OverAwsPeering:  d.Get("over_aws_peering").(string),
		PeeringState:    d.Get("peering_state").(string),
		PeeringHaStatus: d.Get("peering_hastatus").(string),
		Cluster:         d.Get("cluster").(string),
		PeeringLink:     d.Get("peering_link").(string),
		EnableHA:        d.Get("enable_ha").(string),
	}

	log.Printf("[INFO] Creating Aviatrix tunnel: %#v", tunnel)

	err := client.CreateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Tunnel: %s", err)
	}
	d.SetId(tunnel.VpcName1 + "<->" + tunnel.VpcName2)
	return resourceTunnelRead(d, meta)
}

func resourceTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	tunnel := &goaviatrix.Tunnel{
		VpcName1: d.Get("vpc_name1").(string),
		VpcName2: d.Get("vpc_name2").(string),
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

	log.Printf("zjin00: tun.EnableHA is %v", tun.EnableHA)

	d.Set("cluster", tun.Cluster)
	d.Set("over_aws_peering", tun.OverAwsPeering)
	d.Set("peering_hastatus", tun.PeeringHaStatus)
	d.Set("peering_state", tun.PeeringState)
	d.Set("peering_link", tun.PeeringLink)
	d.Set("enable_ha", tun.EnableHA)
	d.SetId(tun.VpcName1 + "<->" + tun.VpcName2)
	log.Printf("[INFO] Found tunnel: %#v", d)
	return nil
}

func resourceTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	tunnel := &goaviatrix.Tunnel{
		VpcName1:        d.Get("vpc_name1").(string),
		VpcName2:        d.Get("vpc_name2").(string),
		OverAwsPeering:  d.Get("over_aws_peering").(string),
		PeeringState:    d.Get("peering_state").(string),
		PeeringHaStatus: d.Get("peering_hastatus").(string),
		Cluster:         d.Get("cluster").(string),
		PeeringLink:     d.Get("peering_link").(string),
	}

	log.Printf("[INFO] Updating Aviatrix tunnel: %#v", tunnel)

	err := client.UpdateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Tunnel: %s", err)
	}
	d.SetId(tunnel.VpcName1 + "<->" + tunnel.VpcName2)
	return resourceTunnelRead(d, meta)
}

func resourceTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	tunnel := &goaviatrix.Tunnel{
		VpcName1: d.Get("vpc_name1").(string),
		VpcName2: d.Get("vpc_name2").(string),
		EnableHA: d.Get("enable_ha").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix tunnel: %#v", tunnel)
	//if enableHA := d.Get("enable_ha").(string); enableHA == "yes" {
	//	// parse the hagw name
	//	tunnel.VpcName1 += "-hagw"
	//	tunnel.VpcName2 += "-hagw"
	//	err := client.DeleteTunnel(tunnel)
	//	if err != nil {
	//		return fmt.Errorf("failed to delete Aviatrix HA gateway: %s", err)
	//	}
	//}

	tunnelHa := &goaviatrix.Tunnel{
		VpcName1: tunnel.VpcName1 + "-hagw",
		VpcName2: tunnel.VpcName2 + "-hagw",
	}
	_, err := client.GetTunnel(tunnelHa)
	if err == nil {
		err1 := client.DeleteTunnel(tunnelHa)
		if err1 != nil {
			return fmt.Errorf("failed to delete Aviatrix HA gateway: %s", err1)
		}
	} else if err != goaviatrix.ErrNotFound {
		return fmt.Errorf("failed to delete Aviatrix HA gateway: %s", err)
	}

	err = client.DeleteTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Tunnel: %s", err)
	}
	return nil
}
