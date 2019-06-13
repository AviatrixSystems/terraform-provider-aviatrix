package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVPNUserAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVPNUserAcceleratorCreate,
		Read:   resourceAviatrixVPNUserAcceleratorRead,
		Delete: resourceAviatrixVPNUserAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"elb_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ELB to include into the VPN User Acclerator.",
			},
		},
	}
}

func resourceAviatrixVPNUserAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	elb := d.Get("elb_name").(string)
	// compare if elb is in elb list for current elbs
	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		return fmt.Errorf("unable to read endpoint list for User Accelerator due to %v", err)
	}
	if goaviatrix.Contains(elbList, elb) {
		return fmt.Errorf("ELB is already included in the VPN User Accelerator. Import into terraform rather than create a new resource.")
	} else {
		elbList = append(elbList, elb)
		elbListStrings := strings.Join(elbList, "\",\"")
		str := "[\"" + elbListStrings + "\"]"
		xlr := &goaviatrix.VpnUserXlr{
			Endpoints: str,
		}
		log.Printf("[DEBUG] Endpoint list: %s", xlr.Endpoints)

		log.Printf("[INFO] Creating User Accelerator.")
		err := client.UpdateVpnUserAccelerator(xlr)
		if err != nil {
			// retry in case of not found elb after waiting
			if strings.Contains(err.Error(), "Endpoint not found") {
				time.Sleep(180 * time.Second)
				err := client.UpdateVpnUserAccelerator(xlr)
				if err != nil {
					return fmt.Errorf("failed to create Vpn User Accelerator: %s", err)
				}
			} else {
				return fmt.Errorf("failed to create Vpn User Accelerator: %s", err)
			}
		}
	}
	d.SetId(elb)
	return resourceAviatrixVPNUserAcceleratorRead(d, meta)
}

func resourceAviatrixVPNUserAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	elbName := d.Get("elb_name").(string)
	if elbName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no elb name received. Import id is %s", id)
		d.Set("elb_name", id)
		elbName = id
		d.SetId(id)
	}
	log.Printf("[DEBUG] elbName: %s", elbName)

	log.Printf("[INFO] Reading User Accelerator LB List")
	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read endpoint list for User Accelerator due to %v", err)
	}
	log.Printf("[DEBUG] elbList: %s", elbList)
	if elbList != nil {
		if goaviatrix.Contains(elbList, elbName) {
			d.Set("elb_name", elbName)
		} else {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceAviatrixVPNUserAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	elbName := d.Get("elb_name").(string)
	toDelete := []string{elbName}

	elbList, err := client.GetVpnUserAccelerator()
	if err != nil {
		return fmt.Errorf("unable to read endpoint list for User Accelerator due to %v", err)
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
			return fmt.Errorf("unable to remove elb in Vpn User Accelerator due to %v", err)
		}
	}
	return nil
}
