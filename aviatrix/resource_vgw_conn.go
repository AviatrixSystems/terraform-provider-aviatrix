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
				Required: true,
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
	return nil
}

func resourceAviatrixVGWConnUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Aviatrix VGW Connection cannot be updated - delete and create new one")
}

func resourceAviatrixVGWConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgw_conn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
		VPCId:    d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix vgw_conn: %#v", vgw_conn)

	err := client.DeleteVGWConn(vgw_conn)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix VGWConn: %s", err)
	}
	return nil
}
