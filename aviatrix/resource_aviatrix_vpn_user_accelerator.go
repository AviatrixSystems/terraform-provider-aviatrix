package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVPNUserAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVPNUserAcceleratorCreate,
		Read:   resourceAviatrixVPNUserAcceleratorRead,
		Delete: resourceAviatrixVPNUserAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"elb_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ELB to include into the VPN User Accelerator.",
			},
		},
	}
}

func resourceAviatrixVPNUserAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	elb := getString(d, "elb_name")
	// compare if elb is in elb list for current elbs
	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		return fmt.Errorf("unable to read endpoint list for User Accelerator due to %w", err)
	}
	if goaviatrix.Contains(elbList, elb) {
		return fmt.Errorf("elb is already included in the VPN User Accelerator. Import into terraform rather than create a new resource")
	} else {
		elbList = append(elbList, elb)
		elbListStrings := strings.Join(elbList, "\",\"")
		str := "[\"" + elbListStrings + "\"]"
		xlr := &goaviatrix.VpnUserXlr{
			Endpoints: str,
		}
		log.Printf("[DEBUG] Endpoint list: %s", xlr.Endpoints)

		log.Printf("[INFO] Creating User Accelerator.")
		var err error
		for i := 0; ; i++ {
			err = client.UpdateVpnUserAccelerator(xlr)
			if err == nil {
				break
			}
			if i <= 10 && (strings.Contains(err.Error(), "Endpoint not found") || strings.Contains(err.Error(), "not active")) {
				time.Sleep(60 * time.Second)
			} else {
				return fmt.Errorf("failed to create Vpn User Accelerator: %w", err)
			}
		}
	}

	d.SetId(elb)
	return resourceAviatrixVPNUserAcceleratorRead(d, meta)
}

func resourceAviatrixVPNUserAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	elbName := getString(d, "elb_name")
	if elbName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no elb name received. Import id is %s", id)
		mustSet(d, "elb_name", id)
		elbName = id
		d.SetId(id)
	}

	log.Printf("[DEBUG] elbName: %s", elbName)

	log.Printf("[INFO] Reading User Accelerator LB List")

	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read endpoint list for User Accelerator due to %w", err)
	}

	log.Printf("[DEBUG] elbList: %s", elbList)

	if elbList != nil {
		if goaviatrix.Contains(elbList, elbName) {
			mustSet(d, "elb_name", elbName)
		} else {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceAviatrixVPNUserAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	elbName := getString(d, "elb_name")
	toDelete := []string{elbName}

	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		return fmt.Errorf("unable to read endpoint list for User Accelerator due to %w", err)
	}
	if elbList != nil {
		xlr := &goaviatrix.VpnUserXlr{}

		// Removes the specified elb from the Vpn User Accelerator List
		newList := goaviatrix.Difference(elbList, toDelete)
		if len(newList) > 0 {
			elbListStrings := strings.Join(newList, "\",\"")
			str := "[\"" + elbListStrings + "\"]"
			xlr.Endpoints = str
		} else {
			xlr.Endpoints = "[]"
		}

		err := client.UpdateVpnUserAccelerator(xlr)
		if err != nil {
			return fmt.Errorf("unable to remove elb in Vpn User Accelerator due to %w", err)
		}
	}

	return nil
}
