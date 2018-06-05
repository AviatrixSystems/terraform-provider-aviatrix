package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixFQDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFQDNCreate,
		Read:   resourceAviatrixFQDNRead,
		Update: resourceAviatrixFQDNUpdate,
		Delete: resourceAviatrixFQDNDelete,

		Schema: map[string]*schema.Schema{
			"fqdn_tag": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gw_list": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"domain_list": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAviatrixFQDNCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	fqdn := &goaviatrix.FQDN{
		FQDNTag:    d.Get("fqdn_tag").(string),
		FQDNStatus: d.Get("fqdn_status").(string),
		FQDNMode:   d.Get("fqdn_mode").(string),
	}
	log.Printf("[INFO] Creating Aviatrix FQDN: %#v", fqdn)
	err := client.CreateFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix FQDN: %s", err)
	}
	if _, ok := d.GetOk("domain_list"); ok {
		fqdn.DomainList = goaviatrix.ExpandStringList(d.Get("domain_list").([]interface{}))
		err = client.UpdateDomains(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to add domain : %s", err)
		}
	}
	if _, ok := d.GetOk("gw_list"); ok {
		fqdn.GwList = goaviatrix.ExpandStringList(d.Get("gw_list").([]interface{}))
		err = client.AttachGws(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to attach GWs: %s", err)
		}
	}
	if fqdn_status := d.Get("fqdn_status").(string); fqdn_status == "enabled" {
		log.Printf("[INFO] Enable FQDN tag status: %#v", fqdn)
		err := client.UpdateFQDNStatus(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to update FQDN status : %s", err)
		}
	}
	// update fqdn_mode when set to non-default "blacklist" mode
	if fqdn_mode := d.Get("fqdn_mode").(string); fqdn_mode == "black" {
		log.Printf("[INFO] Enable FQDN Mode: %#v", fqdn)
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to update FQDN mode : %s", err)
		}
	}
	d.SetId(fqdn.FQDNTag)
	return nil
}

func resourceAviatrixFQDNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	fqdn := &goaviatrix.FQDN{
		FQDNTag:    d.Get("fqdn_tag").(string),
		FQDNStatus: d.Get("fqdn_status").(string),
		FQDNMode:   d.Get("fqdn_mode").(string),
	}
	log.Printf("[INFO] Reading Aviatrix FQDN: %#v", fqdn)
	newfqdn, err := client.GetFQDNTag(fqdn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find FQDN tag: %s", err)
	}
	if newfqdn != nil {
		if _, ok := d.GetOk("fqdn_status"); ok {
			d.Set("fqdn_status", newfqdn.FQDNStatus)
		}
		if _, ok := d.GetOk("fqdn_mode"); ok {
			d.Set("fqdn_mode", newfqdn.FQDNMode)
		}
	}

	newfqdn, err = client.ListDomains(fqdn)
	if err != nil {
		return fmt.Errorf("Couldn't list FQDN domains: %s", err)
	}
	if newfqdn != nil {
		d.Set("domain_list", newfqdn.DomainList)
	}
	newfqdn, err = client.ListGws(fqdn)
	if err != nil {
		return fmt.Errorf("Couldn't list attached gateways: %s", err)
	}
	if newfqdn != nil {
		d.Set("gw_list", newfqdn.GwList)
	}
	return nil
}

func resourceAviatrixFQDNUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	fqdn := &goaviatrix.FQDN{
		FQDNTag:    d.Get("fqdn_tag").(string),
		FQDNStatus: d.Get("fqdn_status").(string),
		FQDNMode:   d.Get("fqdn_mode").(string),
	}
	d.Partial(true)
	if d.HasChange("fqdn_status") {
		err := client.UpdateFQDNStatus(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to update FQDN status : %s", err)
		}
		d.SetPartial("fqdn_status")
	}
	if d.HasChange("fqdn_mode") {
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to update FQDN mode : %s", err)
		}
		d.SetPartial("fqdn_mode")
	}
	//Update Domain list
	if d.HasChange("domain_list") {
		if _, ok := d.GetOk("domain_list"); ok {
			fqdn.DomainList = goaviatrix.ExpandStringList(d.Get("domain_list").([]interface{}))
		}
		err := client.UpdateDomains(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to add domain : %s", err)
		}
		d.SetPartial("domain_list")
	}
	//Update attached GW list
	if d.HasChange("gw_list") {
		o, n := d.GetChange("gw_list")
		if o == nil {
			o = new([]interface{})
		}
		if n == nil {
			n = new([]interface{})
		}
		os := o.([]interface{})
		ns := n.([]interface{})
		oldGwList := goaviatrix.ExpandStringList(os)
		newGwList := goaviatrix.ExpandStringList(ns)
		//Attach all the newly added GWs
		toAddGws := goaviatrix.Difference(newGwList, oldGwList)
		log.Printf("[INFO] Gateways to be attached : %#v", toAddGws)
		fqdn.GwList = toAddGws
		err := client.AttachGws(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to add GW : %s", err)
		}
		//Detach all the removed GWs
		toDelGws := goaviatrix.Difference(oldGwList, newGwList)
		log.Printf("[INFO] Gateways to be detached : %#v", toDelGws)
		fqdn.GwList = toDelGws
		err = client.DetachGws(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to add GW : %s", err)
		}
		d.SetPartial("gw_list")
	}
	d.Partial(false)
	return nil
}

func resourceAviatrixFQDNDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	fqdn := &goaviatrix.FQDN{
		FQDNTag: d.Get("fqdn_tag").(string),
	}
	log.Printf("[INFO] Deleting Aviatrix FQDN: %#v", fqdn)
	if _, ok := d.GetOk("gw_list"); ok {
		log.Printf("[INFO] Found GWs: %#v", fqdn)
		fqdn.GwList = goaviatrix.ExpandStringList(d.Get("gw_list").([]interface{}))
		err := client.DetachGws(fqdn)
		if err != nil {
			return fmt.Errorf("Failed to detach GWs: %s", err)
		}
	}
	err := client.DeleteFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix FQDN: %s", err)
	}

	return nil
}
