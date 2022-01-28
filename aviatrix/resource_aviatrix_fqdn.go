package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

		SchemaVersion: 2,
		MigrateState:  resourceAviatrixFQDNMigrateState,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixFQDNResourceV1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixFQDNStateUpgradeV1,
				Version: 1,
			},
		},

		Schema: map[string]*schema.Schema{
			"fqdn_tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "FQDN Filter Tag Name.",
			},
			"fqdn_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "FQDN Filter Tag Status. Valid values: true or false.",
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
							Required:    true,
							Description: "Name of the gateway to attach to the specific tag.",
						},
						"source_ip_list": {
							Type:        schema.TypeSet,
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
				Deprecated:  "Please set `manage_domain_names` to false, and use the standalone aviatrix_fqdn_tag_rule resource instead.",
				Description: "A list of one or more domain names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fqdn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "FQDN.",
						},
						"proto": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol.",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Port.",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Base Policy",
							ValidateFunc: validation.StringInSlice([]string{"Base Policy", "Allow", "Deny"}, false),
							Description: "What action should happen to matching requests. " +
								"Possible values are: 'Base Policy', 'Allow' or 'Deny'. Defaults to 'Base Policy' if no value is provided.",
						},
					},
				},
			},
			"manage_domain_names": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Enable to manage domain name rules in-line. If false, domain name rules must be managed " +
					"using `aviatrix_fqdn_tag_rule` resources.",
			},
		},
	}
}

func resourceAviatrixFQDNCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	_, hasSetDomainNames := d.GetOk("domain_names")
	enabledInlineDomainNames := d.Get("manage_domain_names").(bool)
	if hasSetDomainNames && !enabledInlineDomainNames {
		return fmt.Errorf("manage_domain_names must be set to true to set in-line domain names")
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag:  d.Get("fqdn_tag").(string),
		FQDNMode: d.Get("fqdn_mode").(string),
	}

	fqdnStatus := d.Get("fqdn_enabled").(bool)
	if fqdnStatus {
		fqdn.FQDNStatus = "enabled"
	} else {
		fqdn.FQDNStatus = "disabled"
	}

	log.Printf("[INFO] Creating Aviatrix FQDN: %#v", fqdn)

	d.SetId(fqdn.FQDNTag)
	flag := false
	defer resourceAviatrixFQDNReadIfRequired(d, meta, &flag)

	err := client.CreateFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix FQDN: %s", err)
	}

	if hasSetDomainNames && enabledInlineDomainNames {
		names := d.Get("domain_names").([]interface{})
		for _, domain := range names {
			if domain != nil {
				dn := domain.(map[string]interface{})
				fqdnFilter := &goaviatrix.Filters{
					FQDN:     dn["fqdn"].(string),
					Protocol: dn["proto"].(string),
					Port:     dn["port"].(string),
					Verdict:  dn["action"].(string),
				}

				fqdn.DomainList = append(fqdn.DomainList, fqdnFilter)
			}
		}

		err = client.UpdateDomains(fqdn)
		if err != nil {
			return fmt.Errorf("failed to add domain : %s", err)
		}
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
		for _, sourceIP := range gFT["source_ip_list"].(*schema.Set).List() {
			sourceIPs = append(sourceIPs, sourceIP.(string))
		}

		if len(sourceIPs) != 0 {
			err = client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
			if err != nil {
				return fmt.Errorf("failed to update source ips to gateway : %s", err)
			}
		}
	}

	if fqdnStatus := d.Get("fqdn_enabled").(bool); fqdnStatus {
		fqdn.FQDNStatus = "enabled"
		log.Printf("[INFO] Enable FQDN tag status: %#v", fqdn)

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

	return resourceAviatrixFQDNReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFQDNReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFQDNRead(d, meta)
	}
	return nil
}

func resourceAviatrixFQDNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fqdnTag := d.Get("fqdn_tag").(string)
	if fqdnTag == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no fqdn tag received. Import Id is %s", id)
		d.Set("fqdn_tag", id)
		d.Set("manage_domain_names", true)
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

	if fqdn.FQDNStatus == "enabled" {
		d.Set("fqdn_enabled", true)
	} else {
		d.Set("fqdn_enabled", false)
	}

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
		if _, ok := d.GetOk("fqdn_enabled"); ok {
			if fqdn.FQDNStatus == "enabled" {
				d.Set("fqdn_enabled", true)
			} else {
				d.Set("fqdn_enabled", false)
			}
		}
		if _, ok := d.GetOk("fqdn_mode"); ok {
			d.Set("fqdn_mode", newfqdn.FQDNMode)
		}
	}
	newfqdn, err = client.ListDomains(fqdn)
	if err != nil {
		return fmt.Errorf("couldn't list FQDN domains: %s", err)
	}
	log.Printf("[INFO] Enable FQDN tag status: %#v", newfqdn)

	if newfqdn != nil {
		// This is nothing IF ListDomains return empty
		var filter []map[string]interface{}
		for _, fqdnDomain := range newfqdn.DomainList {
			dn := make(map[string]interface{})
			dn["fqdn"] = fqdnDomain.FQDN
			dn["proto"] = fqdnDomain.Protocol
			dn["port"] = fqdnDomain.Port
			dn["action"] = fqdnDomain.Verdict
			filter = append(filter, dn)
		}

		log.Printf("[INFO] Enable FQDN tag status: %#v", fqdn)

		// Only write domain names to state if the user has enabled in-line domain names.
		if d.Get("manage_domain_names").(bool) {
			if err = d.Set("domain_names", filter); err != nil {
				log.Printf("[WARN] Error setting domain_names for (%s): %s", d.Id(), err)
			}
		}
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

			sourceIPs1 := gFT["source_ip_list"].(*schema.Set).List()
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

	for _, gwFilterTag := range newfqdn.GwFilterTagList {
		if !mGwFilterTagsOld[gwFilterTag.Name] {
			gwFilterTagList = append(gwFilterTagList, mGwFilterTags[gwFilterTag.Name])
		}
	}

	if err := d.Set("gw_filter_tag_list", gwFilterTagList); err != nil {
		log.Printf("[WARN] Error setting gw_filter_tag_list for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAviatrixFQDNUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	_, hasSetDomainNames := d.GetOk("domain_names")
	enabledInlineDomainNames := d.Get("manage_domain_names").(bool)
	if hasSetDomainNames && !enabledInlineDomainNames {
		return fmt.Errorf("manage_domain_names must be set to true to set in-line domain names")
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag:  d.Get("fqdn_tag").(string),
		FQDNMode: d.Get("fqdn_mode").(string),
	}

	fqdnStatus := d.Get("fqdn_enabled").(bool)
	if fqdnStatus {
		fqdn.FQDNStatus = "enabled"
	} else {
		fqdn.FQDNStatus = "disabled"
	}

	d.Partial(true)
	if d.HasChange("fqdn_tag") {
		return fmt.Errorf("updating fqdn_tag is not allowed")
	}
	if d.HasChange("fqdn_enabled") {
		err := client.UpdateFQDNStatus(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN status : %s", err)
		}
	}
	if d.HasChange("fqdn_mode") {
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN mode : %s", err)
		}
	}
	// Update Domain list
	if d.HasChange("domain_names") && enabledInlineDomainNames {
		if hasSetDomainNames {
			names := d.Get("domain_names").([]interface{})
			for _, domain := range names {
				dn := domain.(map[string]interface{})
				fqdnDomain := &goaviatrix.Filters{
					FQDN:     dn["fqdn"].(string),
					Protocol: dn["proto"].(string),
					Port:     dn["port"].(string),
					Verdict:  dn["action"].(string),
				}
				fqdn.DomainList = append(fqdn.DomainList, fqdnDomain)
			}
		}
		err := client.UpdateDomains(fqdn)
		if err != nil {
			return fmt.Errorf("failed to add domain : %s", err)
		}
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

		mapOldTagList := make(map[string][]string)
		for _, oldGwFilterTags := range os {
			oldGwFilterTag := oldGwFilterTags.(map[string]interface{})
			sourceIPs := make([]string, 0)
			for _, sourceIP := range oldGwFilterTag["source_ip_list"].(*schema.Set).List() {
				sourceIPs = append(sourceIPs, sourceIP.(string))
			}
			mapOldTagList[oldGwFilterTag["gw_name"].(string)] = sourceIPs
		}

		for _, newGwFilterTags := range ns {
			newGwFilterTag := newGwFilterTags.(map[string]interface{})
			gateway := &goaviatrix.Gateway{
				GwName: newGwFilterTag["gw_name"].(string),
			}

			sourceIPs := make([]string, 0)
			for _, sourceIP := range newGwFilterTag["source_ip_list"].(*schema.Set).List() {
				sourceIPs = append(sourceIPs, sourceIP.(string))
			}

			val, ok := mapOldTagList[gateway.GwName]
			if !ok {
				err := client.AttachTagToGw(fqdn, gateway)
				if err != nil {
					return fmt.Errorf("failed to add filter tag to gateway : %s", err)
				}
				continue
			}

			if !goaviatrix.Equivalent(val, sourceIPs) {
				err := client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
				if err != nil {
					return fmt.Errorf("failed to update source ips to gateway : %s", err)
				}
			}
			delete(mapOldTagList, gateway.GwName)
		}

		keys := make([]string, 0)
		for key := range mapOldTagList {
			keys = append(keys, key)
		}
		if len(keys) != 0 {
			err := client.DetachGws(fqdn, keys)
			if err != nil {
				return fmt.Errorf("failed to delete GWs for fqdn in update: %s", err)
			}
		}
	}

	d.Partial(false)
	d.SetId(fqdn.FQDNTag)
	return resourceAviatrixFQDNRead(d, meta)
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
