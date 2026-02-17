package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixPeriodicPing() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixPeriodicPingCreate,
		Read:   resourceAviatrixPeriodicPingRead,
		Delete: resourceAviatrixPeriodicPingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
		GwName:   getString(d, "gw_name"),
		Interval: strconv.Itoa(getInt(d, "interval")),
		IP:       getString(d, "ip_address"),
	}
}

func resourceAviatrixPeriodicPingCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	pp := marshalPeriodicPingInput(d)

	d.SetId(pp.GwName)
	flag := false
	defer func() { _ = resourceAviatrixPeriodicPingReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if err := client.CreatePeriodicPing(pp); err != nil {
		return err
	}

	return resourceAviatrixPeriodicPingReadIfRequired(d, meta, &flag)
}

func resourceAviatrixPeriodicPingReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixPeriodicPingRead(d, meta)
	}
	return nil
}

func resourceAviatrixPeriodicPingRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gwName := getString(d, "gw_name")
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
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find periodic_ping %s: %w", gwName, err)
	}
	mustSet(d, "gw_name", gwName)
	mustSet(d, "ip_address", pp.IP)
	mustSet(d, "interval", pp.IntervalAsInt)

	d.SetId(pp.GwName)
	return nil
}

func resourceAviatrixPeriodicPingDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	pp := marshalPeriodicPingInput(d)

	if err := client.DeletePeriodicPing(pp); err != nil {
		return err
	}

	return nil
}
