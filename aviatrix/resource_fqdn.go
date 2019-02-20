package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	//"strings"
)

func resourceAviatrixFQDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFQDNCreate,
		Read:   resourceAviatrixFQDNRead,
		Update: resourceAviatrixFQDNUpdate,
		Delete: resourceAviatrixFQDNDelete,

		Schema: map[string]*schema.Schema{
			"fqdn_tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gw_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"domain_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"proto": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
		return fmt.Errorf("failed to create Aviatrix FQDN: %s", err)
	}
	d.SetId(fqdn.FQDNTag)

	if _, ok := d.GetOk("domain_names"); ok {
		names := d.Get("domain_names").([]interface{})
		for _, domain := range names {
			if domain != nil {
				dn := domain.(map[string]interface{})
				fqdnFilter := &goaviatrix.Filters{
					FQDN:     dn["fqdn"].(string),
					Protocol: dn["proto"].(string),
					Port:     dn["port"].(string),
				}

				fqdn.DomainList = append(fqdn.DomainList, fqdnFilter)
			}
		}

		err = client.UpdateDomains(fqdn)
		if err != nil {
			return fmt.Errorf("failed to add domain : %s", err)
		}
		d.Set("domain_names", fqdn.DomainList)
	}
	if _, ok := d.GetOk("gw_list"); ok {
		tagList := d.Get("gw_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
		fqdn.GwList = tagListStr
		err = client.AttachGws(fqdn)
		if err != nil {
			return fmt.Errorf("failed to attach GWs: %s", err)
		}
		d.Set("gw_list", fqdn.GwList)
	}
	if fqdnStatus := d.Get("fqdn_status").(string); fqdnStatus == "enabled" {
		log.Printf("[INOF] Enable FQDN tag status: %#v", fqdn)
		err := client.UpdateFQDNStatus(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN status : %s", err)
		}
	}
	// update fqdn_mode when set to non-default "blacklist" mode
	if fqdnMode := d.Get("fqdn_mode").(string); fqdnMode == "black" {
		log.Printf("[INFO] Enable FQDN Mode: %#v", fqdn)
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN mode : %s", err)
		}
	}
	return resourceAviatrixFQDNRead(d, meta)
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
		return fmt.Errorf("couldn't find FQDN tag: %s", err)
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
		return fmt.Errorf("couldn't list FQDN domains: %s", err)
	}
	log.Printf("[INOF] 2Enable FQDN tag status: %#v", newfqdn)

	if newfqdn != nil {
		// This is nothing IF ListDomains return empty
		var filter []map[string]interface{}
		for _, fqdnDomain := range newfqdn.DomainList {
			dn := make(map[string]interface{})
			dn["fqdn"] = fqdnDomain.FQDN
			dn["proto"] = fqdnDomain.Protocol
			dn["port"] = fqdnDomain.Port
			filter = append(filter, dn)
		}

		log.Printf("[INOF] 3Enable FQDN tag status: %#v", fqdn)

		d.Set("domain_names", filter)
	}
	tagList := d.Get("gw_list").([]interface{})
	tagListStr := goaviatrix.ExpandStringList(tagList)
	fqdn.GwList = tagListStr
	newfqdn, err = client.ListGws(fqdn)
	if err != nil {
		return fmt.Errorf("couldn't list attached gateways: %s", err)
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
			return fmt.Errorf("failed to update FQDN status : %s", err)
		}
		d.SetPartial("fqdn_status")
	}
	if d.HasChange("fqdn_mode") {
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN mode : %s", err)
		}
		d.SetPartial("fqdn_mode")
	}
	//Update Domain list
	if d.HasChange("domain_names") {
		if _, ok := d.GetOk("domain_names"); ok {
			names := d.Get("domain_names").([]interface{})
			for _, domain := range names {
				dn := domain.(map[string]interface{})
				fqdnDomain := &goaviatrix.Filters{
					FQDN:     dn["fqdn"].(string),
					Protocol: dn["proto"].(string),
					Port:     dn["port"].(string),
				}
				fqdn.DomainList = append(fqdn.DomainList, fqdnDomain)
			}
		}
		err := client.UpdateDomains(fqdn)
		if err != nil {
			return fmt.Errorf("failed to add domain : %s", err)
		}
		d.SetPartial("domain_names")
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
			return fmt.Errorf("failed to add GW : %s", err)
		}
		//Detach all the removed GWs
		toDelGws := goaviatrix.Difference(oldGwList, newGwList)
		log.Printf("[INFO] Gateways to be detached : %#v", toDelGws)
		fqdn.GwList = toDelGws
		err = client.DetachGws(fqdn)
		if err != nil {
			return fmt.Errorf("failed to add GW : %s", err)
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
			return fmt.Errorf("failed to detach GWs: %s", err)
		}
	}
	err := client.DeleteFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FQDN: %s", err)
	}
	return nil
}
