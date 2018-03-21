package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTunnelCreate,
		Read:   resourceTunnelRead,
		Update: resourceTunnelUpdate,
		Delete: resourceTunnelDelete,

		Schema: map[string]*schema.Schema{
			"vpc_name1": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_name2": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"over_aws_peering": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peering_hastatus": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_link": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_ha": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
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
	}

	log.Printf("[INFO] Creating Aviatrix tunnel: %#v", tunnel)

	err := client.CreateTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Tunnel: %s", err)
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
		return fmt.Errorf("Couldn't find Aviatrix Tunnel: %s", err)
	}
	log.Printf("[INFO] Found Aviatrix tunnel: %#v", tun)
	d.Set("over_aws_peering", tun.OverAwsPeering)
	d.Set("peering_state", tun.PeeringState)
	d.Set("peering_hastatus", tun.PeeringHaStatus)
	d.Set("cluster", tun.Cluster)
	d.Set("peering_link", tun.PeeringLink)
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
		return fmt.Errorf("Failed to update Aviatrix Tunnel: %s", err)
	}
	d.SetId(tunnel.VpcName1 + "<->" + tunnel.VpcName2)
	return resourceTunnelRead(d, meta)
}

func resourceTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	tunnel := &goaviatrix.Tunnel{
		VpcName1: d.Get("vpc_name1").(string),
		VpcName2: d.Get("vpc_name2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix tunnel: %#v", tunnel)

	err := client.DeleteTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Tunnel: %s", err)
	}
	return nil
}
