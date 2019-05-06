package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFQDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFQDNCreate,
		Read:   resourceAviatrixFQDNRead,
		Update: resourceAviatrixFQDNUpdate,
		Delete: resourceAviatrixFQDNDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixFQDNMigrateState,

		Schema: map[string]*schema.Schema{
			"fqdn_tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "FQDN Filter Tag Name.",
			},
			"fqdn_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "FQDN Filter Tag Status. Valid values: 'enabled', 'disabled'.",
			},
			"fqdn_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the tag color to be a white-list tag or black-list tag. 'white' or 'black'",
			},
			"gw_filter_tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of gateways to attach to the specific tag.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gw_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the gateway to attach to the specific tag.",
						},
						"source_ip_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "List of source IPs in the VPC qualified for a specific tag.",
						},
					},
				},
			},
			"domain_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of one or more domain names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fqdn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "FQDN.",
						},
						"proto": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port.",
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

	gwFilterTags := d.Get("gw_filter_tag_list").([]interface{})

	for _, gwFilterTag := range gwFilterTags {
		gFT := gwFilterTag.(map[string]interface{})
		gateway := &goaviatrix.Gateway{
			GwName: gFT["gw_name"].(string),
		}
		err := client.AttachTagToGw(fqdn, gateway)
		if err != nil {
			return fmt.Errorf("failed to add filter tag to gateway : %s", err)
		}
		sourceIPs := make([]string, 0)
		for _, sourceIP := range gFT["source_ip_list"].([]interface{}) {
			sourceIPs = append(sourceIPs, sourceIP.(string))
		}

		if len(sourceIPs) != 0 {
			err = client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
			if err != nil {
				return fmt.Errorf("failed to update source ips to gateway : %s", err)
			}
		}
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

	fqdnTag := d.Get("fqdn_tag").(string)
	if fqdnTag == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no fqdn tag received. Import Id is %s", id)
		d.Set("fqdn_tag", id)
		d.SetId(id)
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag: d.Get("fqdn_tag").(string),
	}

	fqdn, err := client.GetFQDNTag(fqdn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FQDN tag: %s", err)
	}
	d.Set("fqdn_status", fqdn.FQDNStatus)
	d.Set("fqdn_mode", fqdn.FQDNMode)

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
	log.Printf("[INOF] Enable FQDN tag status: %#v", newfqdn)

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

	newfqdn, err = client.GetGwFilterTagList(newfqdn)
	if err != nil {
		return fmt.Errorf("couldn't list FQDN Filter Tags: %s", err)
	}

	mGwFilterTags := make(map[string]map[string]interface{})
	for _, gwFilterTag := range newfqdn.GwFilterTagList {
		gFT := make(map[string]interface{})
		gFT["gw_name"] = gwFilterTag.Name

		var ipList []string
		for i := range gwFilterTag.SourceIPList {
			ipList = append(ipList, gwFilterTag.SourceIPList[i])
		}
		gFT["source_ip_list"] = ipList
		mGwFilterTags[gwFilterTag.Name] = gFT
	}

	var gwFilterTagList []map[string]interface{}
	gwFilterTags := d.Get("gw_filter_tag_list").([]interface{})

	mGwFilterTagsOld := make(map[string]bool)
	for _, gwFilterTag := range gwFilterTags {
		gFT := gwFilterTag.(map[string]interface{})

		mGwFilterTagsOld[gFT["gw_name"].(string)] = true

		if mGwFilterTags[gFT["gw_name"].(string)] != nil {
			mSourceIPsNew := make(map[string]bool)
			aSourceIPsNew := make([]string, 0)
			sourceIPs := mGwFilterTags[gFT["gw_name"].(string)]["source_ip_list"].([]string)

			for i := 0; i < len(sourceIPs); i++ {
				mSourceIPsNew[sourceIPs[i]] = true
			}

			sourceIPs1 := gFT["source_ip_list"].([]interface{})
			for i := 0; i < len(sourceIPs1); i++ {
				if mSourceIPsNew[sourceIPs1[i].(string)] {
					aSourceIPsNew = append(aSourceIPsNew, sourceIPs1[i].(string))
					mSourceIPsNew[sourceIPs1[i].(string)] = false
				}
			}

			for i := 0; i < len(sourceIPs); i++ {
				if mSourceIPsNew[sourceIPs[i]] {
					aSourceIPsNew = append(aSourceIPsNew, sourceIPs[i])
				}
			}

			gwFilterTagList = append(gwFilterTagList, mGwFilterTags[gFT["gw_name"].(string)])
		}
	}
	if err != nil {
		return fmt.Errorf("couldn't list attached gateways: %s", err)
	}
	for _, gwFilterTag := range newfqdn.GwFilterTagList {
		if !mGwFilterTagsOld[gwFilterTag.Name] {
			gwFilterTagList = append(gwFilterTagList, mGwFilterTags[gwFilterTag.Name])
		}
	}
	d.Set("gw_filter_tag_list", gwFilterTagList)

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
	if d.HasChange("fqdn_tag") {
		return fmt.Errorf("updating fqdn_tag is not allowed")
	}
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
	if d.HasChange("gw_filter_tag_list") {
		o, n := d.GetChange("gw_filter_tag_list")
		if o == nil {
			o = new([]interface{})
		}
		if n == nil {
			n = new([]interface{})
		}
		os := o.([]interface{})
		ns := n.([]interface{})

		oldGwList := make([]string, 0)
		for _, oldGwFilterTags := range os {
			oldGwFilterTag := oldGwFilterTags.(map[string]interface{})
			gwName := oldGwFilterTag["gw_name"].(string)
			oldGwList = append(oldGwList, gwName)
		}

		newGwList := make([]string, 0)
		for _, newGwFilterTags := range ns {
			newGwFilterTag := newGwFilterTags.(map[string]interface{})
			gwName := newGwFilterTag["gw_name"].(string)
			newGwList = append(newGwList, gwName)
		}

		gwToDelete := goaviatrix.Difference(oldGwList, newGwList)
		err := client.DetachGws(fqdn, gwToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete GWs for fqdn: %s", err)

		}

		gwToAdd := goaviatrix.Difference(newGwList, oldGwList)
		mGwToAdd := make(map[string]bool)
		for i := range gwToAdd {
			mGwToAdd[gwToAdd[i]] = true
		}

		for _, gwFilterTag := range ns {
			gFT := gwFilterTag.(map[string]interface{})
			gateway := &goaviatrix.Gateway{
				GwName: gFT["gw_name"].(string),
			}
			if mGwToAdd[gateway.GwName] {
				err := client.AttachTagToGw(fqdn, gateway)
				if err != nil {
					return fmt.Errorf("failed to add filter tag to gateway : %s", err)
				}
			}
			sourceIPs := make([]string, 0)
			for _, sourceIP := range gFT["source_ip_list"].([]interface{}) {
				sourceIPs = append(sourceIPs, sourceIP.(string))
			}

			if len(sourceIPs) != 0 {
				err = client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
				if err != nil {
					return fmt.Errorf("failed to update source ips to gateway : %s", err)
				}
			}
		}

		d.SetPartial("gw_filter_tag_list")
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

	gwList, err := client.ListGws(fqdn)
	if err != nil {
		return fmt.Errorf("failed to get GW list for fqdn: %s", err)
	}
	err = client.DetachGws(fqdn, gwList)
	if err != nil {
		return fmt.Errorf("failed to delete GWs for fqdn: %s", err)

	}

	err = client.DeleteFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FQDN: %s", err)
	}
	return nil
}
