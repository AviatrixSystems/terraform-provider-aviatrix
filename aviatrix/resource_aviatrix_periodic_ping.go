package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixPeriodicPing() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixPeriodicPingCreate,
		Read:   resourceAviatrixPeriodicPingRead,
		Delete: resourceAviatrixPeriodicPingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of gateway.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Interval between pings in seconds.",
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "IP Address to ping.",
			},
		},
	}
}

func marshalPeriodicPingInput(d *schema.ResourceData) *goaviatrix.PeriodicPing {
	return &goaviatrix.PeriodicPing{
		GwName:   d.Get("gw_name").(string),
		Interval: strconv.Itoa(d.Get("interval").(int)),
		IP:       d.Get("ip_address").(string),
	}
}

func resourceAviatrixPeriodicPingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	pp := marshalPeriodicPingInput(d)

	if err := client.CreatePeriodicPing(pp); err != nil {
		return err
	}

	d.SetId(pp.GwName)
	return nil
}

func resourceAviatrixPeriodicPingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no periodic_ping gw_name received. Import Id is %s", id)
		d.SetId(id)
		gwName = id
	}

	pp := &goaviatrix.PeriodicPing{
		GwName: gwName,
	}

	pp, err := client.GetPeriodicPing(pp)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find periodic_ping %s: %v", gwName, err)
	}

	d.Set("gw_name", gwName)
	d.Set("ip_address", pp.IP)
	d.Set("interval", pp.IntervalAsInt)

	d.SetId(pp.GwName)
	return nil
}

func resourceAviatrixPeriodicPingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	pp := marshalPeriodicPingInput(d)

	if err := client.DeletePeriodicPing(pp); err != nil {
		return err
	}

	return nil
}
