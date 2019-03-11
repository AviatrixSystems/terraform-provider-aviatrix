package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixVGWConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVGWConnCreate,
		Read:   resourceAviatrixVGWConnRead,
		Update: resourceAviatrixVGWConnUpdate,
		Delete: resourceAviatrixVGWConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"conn_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bgp_vgw_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bgp_local_as_num": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAviatrixVGWConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgwConn := &goaviatrix.VGWConn{
		ConnName:      d.Get("conn_name").(string),
		GwName:        d.Get("gw_name").(string),
		VPCId:         d.Get("vpc_id").(string),
		BgpVGWId:      d.Get("bgp_vgw_id").(string),
		BgpLocalAsNum: d.Get("bgp_local_as_num").(string),
	}

	log.Printf("[INFO] Creating Aviatrix VGW Connection: %#v", vgwConn)

	err := client.CreateVGWConn(vgwConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix VGWConn: %s", err)
	}
	d.SetId(vgwConn.ConnName)

	return nil
}

func resourceAviatrixVGWConnRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAviatrixVGWConnUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("aviatrix VGW Connection cannot be updated - delete and create new one")
}

func resourceAviatrixVGWConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vgwConn := &goaviatrix.VGWConn{
		ConnName: d.Get("conn_name").(string),
		VPCId:    d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix vgw_conn: %#v", vgwConn)

	err := client.DeleteVGWConn(vgwConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix VGWConn: %s", err)
	}
	return nil
}
