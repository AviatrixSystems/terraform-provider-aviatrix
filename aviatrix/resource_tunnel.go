package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTunnelCreate,
		Read:   resourceTunnelRead,
		Update: resourceTunnelUpdate,
		Delete: resourceTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_name1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_name2": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peering_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peering_hastatus": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	d.SetId(tunnel.VpcName1 + "~" + tunnel.VpcName2)
	return resourceTunnelRead(d, meta)
}

func resourceTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcName1 := d.Get("vpc_name1").(string)
	vpcName2 := d.Get("vpc_name2").(string)

	if vpcName1 == "" || vpcName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		d.Set("vpc_name1", strings.Split(id, "~")[0])
		d.Set("vpc_name2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

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

	d.Set("cluster", tun.Cluster)
	d.Set("peering_hastatus", tun.PeeringHaStatus)
	d.Set("peering_state", tun.PeeringState)
	d.Set("peering_link", tun.PeeringLink)
	d.Set("enable_ha", tun.EnableHA)
	d.SetId(tun.VpcName1 + "~" + tun.VpcName2)
	log.Printf("[INFO] Found tunnel: %#v", d)
	return nil
}

func resourceTunnelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	tunnel := &goaviatrix.Tunnel{
		VpcName1:        d.Get("vpc_name1").(string),
		VpcName2:        d.Get("vpc_name2").(string),
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
	d.SetId(tunnel.VpcName1 + "~" + tunnel.VpcName2)
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
	if peeringHaStatus := d.Get("peering_hastatus").(string); peeringHaStatus == "active" {
		// parse the hagw name
		tunnel.VpcName1 += "-hagw"
		tunnel.VpcName2 += "-hagw"
		err := client.DeleteTunnel(tunnel)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix HA gateway: %s", err)
		}
	}

	tunnel.VpcName1 = d.Get("vpc_name1").(string)
	tunnel.VpcName2 = d.Get("vpc_name2").(string)
	err := client.DeleteTunnel(tunnel)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Tunnel: %s", err)
	}
	return nil
}
