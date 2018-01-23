package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixVGWConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVGWConnCreate,
		Read:   resourceAviatrixVGWConnRead,
		Update: resourceAviatrixVGWConnUpdate,
		Delete: resourceAviatrixVGWConnDelete,

		Schema: map[string]*schema.Schema{
			"conn_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bgp_vgw_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bgp_local_as_num": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ha": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixVGWConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgw_conn := &goaviatrix.VGWConn{
		ConnName:      d.Get("conn_name").(string),
		GwName:        d.Get("gw_name").(string),
		VPCId:         d.Get("vpc_id").(string),
		BgpVGWId:      d.Get("bgp_vgw_id").(string),
		BgpLocalAsNum: d.Get("bgp_local_as_num").(string),
		EnableHa:      d.Get("enable_ha").(string),
	}

	log.Printf("[INFO] Creating Aviatrix VGW Connection: %#v", vgw_conn)

	err := client.CreateVGWConn(vgw_conn)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix VGWConn: %s", err)
	}
	d.SetId(vgw_conn.ConnName)
	return nil
	//return resourceAviatrixVGWConnRead(d, meta)
}

func resourceAviatrixVGWConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgw_conn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
	}
	conn, err := client.GetVGWConn(vgw_conn)
	if err != nil {
		return fmt.Errorf("Couldn't find Aviatrix VGWConn: %s", err)
	}
	log.Printf("[TRACE] reading vgw_conn %s: %#v",
		d.Get("conn_name").(string), conn)
	if conn != nil {
		d.Set("conn_name", conn.ConnName)
		//d.Set("gw_name", conn.GwName)
		d.Set("vpc_id", conn.VPCId)
		//d.Set("bgp_vgw_id", conn.BgpVGWId)
		//d.Set("bgp_local_as_num", conn.BgpLocalAsNum)
		//d.Set("enable_ha", conn.EnableHa)
	}
	return nil
}

func resourceAviatrixVGWConnUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Aviatrix VGW Connection cannot be updated - delete and create new one")
}

func resourceAviatrixVGWConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgw_conn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
		GwName:   d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix vgw_conn: %#v", vgw_conn)

	err := client.DeleteVGWConn(vgw_conn)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix VGWConn: %s", err)
	}
	return nil
}
