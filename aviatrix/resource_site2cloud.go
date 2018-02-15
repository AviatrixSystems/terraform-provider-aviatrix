package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixSite2Cloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSite2CloudCreate,
		Read:   resourceAviatrixSite2CloudRead,
		Update: resourceAviatrixSite2CloudUpdate,
		Delete: resourceAviatrixSite2CloudDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"conn_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_gw_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gw_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pre_shared_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_gw_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_subnet": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixSite2CloudCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	s2c := &goaviatrix.Site2Cloud{
		GwName:       d.Get("gw_name").(string),
		VpcID:        d.Get("vpc_id").(string),
		ConnName:     d.Get("conn_name").(string),
		RemoteGwType: d.Get("remote_gw_type").(string),
		TunnelType:   d.Get("tunnel_type").(string),
		RemoteGwIP:   d.Get("remote_gw_ip").(string),
		PreSharedKey: d.Get("pre_shared_key").(string),
		RemoteSubnet: d.Get("remote_subnet").(string),
		LocalSubnet:  d.Get("local_subnet").(string),
	}

	log.Printf("[INFO] Creating Aviatrix Site2Cloud: %#v", s2c)

	err := client.CreateSite2Cloud(s2c)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Site2Cloud: %s", err)
	}
	d.SetId(s2c.ConnName)
	return nil
}

func resourceAviatrixSite2CloudRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	site2cloud := &goaviatrix.Site2Cloud{
		ConnName: d.Get("conn_name").(string),
	}
	s2c, err := client.GetSite2Cloud(site2cloud)
	if err != nil {
		return fmt.Errorf("Couldn't find Aviatrix Site2Cloud: %s", err)
	}
	log.Printf("[TRACE] Reading Aviatrix Site2Cloud %s: %#v",
		d.Get("gw_name").(string), s2c)
	if s2c != nil {
		d.Set("vpc_id", s2c.VpcID)
		d.Set("remote_gw_type", s2c.RemoteGwType)
		d.Set("tunnel_type", s2c.TunnelType)
		d.Set("remote_gw_ip", s2c.RemoteGwIP)
		d.Set("remote_subnet", s2c.RemoteSubnet)
		d.Set("local_subnet", s2c.LocalSubnet)
	}
	return nil
}

func resourceAviatrixSite2CloudUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	site2cloud := &goaviatrix.Site2Cloud{
		GwName:   d.Get("gw_name").(string),
		VpcID:    d.Get("vpc_id").(string),
		ConnName: d.Get("conn_name").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", site2cloud)
	if ok := d.HasChange("remote_subnet"); ok {
		site2cloud.RemoteSubnet = d.Get("remote_subnet").(string)
		err := client.UpdateSite2Cloud(site2cloud)
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix Site2Cloud: %s", err)
		}
		d.SetPartial("remote_subnet")
	}
	if ok := d.HasChange("local_subnet"); ok {
		site2cloud.LocalSubnet = d.Get("local_subnet").(string)
		err := client.UpdateSite2Cloud(site2cloud)
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix Site2Cloud: %s", err)
		}
		d.SetPartial("local_subnet")
	}
	d.Partial(false)
	return nil
}

func resourceAviatrixSite2CloudDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	site2cloud := &goaviatrix.Site2Cloud{
		ConnName: d.Get("conn_name").(string),
		VpcID:    d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix s2c: %#v", site2cloud)

	err := client.DeleteSite2Cloud(site2cloud)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Site2Cloud: %s", err)
	}
	return nil
}
