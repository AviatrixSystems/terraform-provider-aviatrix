package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFQDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFQDNCreate,
		Read:   resourceAviatrixFQDNRead,
		Update: resourceAviatrixFQDNUpdate,
		Delete: resourceAviatrixFQDNDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
				Description: "A list of one or more domain names/tag rules.",
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
	client := mustClient(meta)

	_, hasSetDomainNames := d.GetOk("domain_names")
	enabledInlineDomainNames := getBool(d, "manage_domain_names")
	if hasSetDomainNames && !enabledInlineDomainNames {
		return fmt.Errorf("manage_domain_names must be set to true to set in-line domain names")
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag:  getString(d, "fqdn_tag"),
		FQDNMode: getString(d, "fqdn_mode"),
	}

	fqdnStatus := getBool(d, "fqdn_enabled")
	if fqdnStatus {
		fqdn.FQDNStatus = "enabled"
	} else {
		fqdn.FQDNStatus = "disabled"
	}

	log.Printf("[INFO] Creating Aviatrix FQDN: %#v", fqdn)

	d.SetId(fqdn.FQDNTag)
	flag := false
	defer func() { _ = resourceAviatrixFQDNReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix FQDN: %w", err)
	}

	if hasSetDomainNames && enabledInlineDomainNames {
		names := getList(d, "domain_names")
		mapDomains := make(map[string]bool)
		for _, domain := range names {
			if domain != nil {
				dn := mustMap(domain)
				fqdnFilter := &goaviatrix.Filters{
					FQDN:     mustString(dn["fqdn"]),
					Protocol: mustString(dn["proto"]),
					Port:     mustString(dn["port"]),
					Verdict:  mustString(dn["action"]),
				}
				str := fqdnFilter.FQDN + fqdnFilter.Protocol + fqdnFilter.Port + fqdnFilter.Verdict
				if mapDomains[str] {
					return fmt.Errorf("validation on domain_names failed: duplicate rules are not allowed")
				}
				mapDomains[str] = true
				fqdn.DomainList = append(fqdn.DomainList, fqdnFilter)
			}
		}
		if err := client.UpdateDomains(fqdn); err != nil {
			return fmt.Errorf("failed to set domain names: %w", err)
		}
	}

	gwFilterTags := getList(d, "gw_filter_tag_list")
	for _, gwFilterTag := range gwFilterTags {
		gFT := mustMap(gwFilterTag)
		gateway := &goaviatrix.Gateway{
			GwName: mustString(gFT["gw_name"]),
		}
		err := client.AttachTagToGw(fqdn, gateway)
		if err != nil {
			return fmt.Errorf("failed to add filter tag to gateway : %w", err)
		}
		sourceIPs := make([]string, 0)
		for _, sourceIP := range mustSchemaSet(gFT["source_ip_list"]).List() {
			sourceIPs = append(sourceIPs, mustString(sourceIP))
		}

		if len(sourceIPs) != 0 {
			err = client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
			if err != nil {
				return fmt.Errorf("failed to update source ips to gateway : %w", err)
			}
		}
	}

	if fqdnStatus := getBool(d, "fqdn_enabled"); fqdnStatus {
		fqdn.FQDNStatus = "enabled"
		log.Printf("[INFO] Enable FQDN tag status: %#v", fqdn)

		err := client.UpdateFQDNStatus(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN status : %w", err)
		}
	}

	// update fqdn_mode when set to non-default "blacklist" mode
	if fqdnMode := getString(d, "fqdn_mode"); fqdnMode == "black" {
		log.Printf("[INFO] Enable FQDN Mode: %#v", fqdn)
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN mode : %w", err)
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
	client := mustClient(meta)

	fqdnTag := getString(d, "fqdn_tag")
	if fqdnTag == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no fqdn tag received. Import Id is %s", id)
		mustSet(d, "fqdn_tag", id)
		mustSet(d, "manage_domain_names", true)
		d.SetId(id)
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag: getString(d, "fqdn_tag"),
	}

	fqdn, err := client.GetFQDNTag(fqdn)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FQDN tag: %w", err)
	}

	if fqdn.FQDNStatus == "enabled" {
		mustSet(d, "fqdn_enabled", true)
	} else {
		mustSet(d, "fqdn_enabled", false)
	}
	mustSet(d, "fqdn_mode", fqdn.FQDNMode)

	log.Printf("[INFO] Reading Aviatrix FQDN: %#v", fqdn)
	newfqdn, err := client.GetFQDNTag(fqdn)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find FQDN tag: %w", err)
	}

	if newfqdn != nil {
		if _, ok := d.GetOk("fqdn_enabled"); ok {
			if fqdn.FQDNStatus == "enabled" {
				mustSet(d, "fqdn_enabled", true)
			} else {
				mustSet(d, "fqdn_enabled", false)
			}
		}
		if _, ok := d.GetOk("fqdn_mode"); ok {
			mustSet(d, "fqdn_mode", newfqdn.FQDNMode)
		}
	}
	newfqdn, err = client.ListDomains(fqdn)
	if err != nil {
		return fmt.Errorf("couldn't list FQDN domains: %w", err)
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
		if getBool(d, "manage_domain_names") {
			if err = d.Set("domain_names", filter); err != nil {
				log.Printf("[WARN] Error setting domain_names for (%s): %s", d.Id(), err)
			}
		}
	}

	newfqdn, err = client.GetGwFilterTagList(newfqdn)
	if err != nil {
		return fmt.Errorf("couldn't list FQDN Filter Tags: %w", err)
	}

	mGwFilterTags := make(map[string]map[string]interface{})
	for _, gwFilterTag := range newfqdn.GwFilterTagList {
		gFT := make(map[string]interface{})
		gFT["gw_name"] = gwFilterTag.Name

		var ipList []string
		ipList = append(ipList, gwFilterTag.SourceIPList...)
		gFT["source_ip_list"] = ipList
		mGwFilterTags[gwFilterTag.Name] = gFT
	}

	var gwFilterTagList []map[string]interface{}
	gwFilterTags := getList(d, "gw_filter_tag_list")

	mGwFilterTagsOld := make(map[string]bool)
	for _, gwFilterTag := range gwFilterTags {
		gFT := mustMap(gwFilterTag)

		gwName := mustString(gFT["gw_name"])
		mGwFilterTagsOld[gwName] = true

		if mGwFilterTags[gwName] != nil {
			gwFilterTagList = append(gwFilterTagList, mGwFilterTags[gwName])
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
	client := mustClient(meta)

	_, hasSetDomainNames := d.GetOk("domain_names")
	enabledInlineDomainNames := getBool(d, "manage_domain_names")
	if hasSetDomainNames && !enabledInlineDomainNames {
		return fmt.Errorf("manage_domain_names must be set to true to set in-line domain names")
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag:  getString(d, "fqdn_tag"),
		FQDNMode: getString(d, "fqdn_mode"),
	}

	fqdnStatus := getBool(d, "fqdn_enabled")
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
			return fmt.Errorf("failed to update FQDN status : %w", err)
		}
	}
	if d.HasChange("fqdn_mode") {
		err := client.UpdateFQDNMode(fqdn)
		if err != nil {
			return fmt.Errorf("failed to update FQDN mode : %w", err)
		}
	}
	// Update Domain list
	if d.HasChange("domain_names") && enabledInlineDomainNames {
		if hasSetDomainNames {
			names := getList(d, "domain_names")
			mapDomains := make(map[string]bool)
			for _, domain := range names {
				dn := mustMap(domain)
				fqdnDomain := &goaviatrix.Filters{
					FQDN:     mustString(dn["fqdn"]),
					Protocol: mustString(dn["proto"]),
					Port:     mustString(dn["port"]),
					Verdict:  mustString(dn["action"]),
				}
				str := fqdnDomain.FQDN + fqdnDomain.Protocol + fqdnDomain.Port + fqdnDomain.Verdict
				if mapDomains[str] {
					return fmt.Errorf("validation on domain_names failed in update: duplicate rules are not allowed")
				}
				mapDomains[str] = true
				fqdn.DomainList = append(fqdn.DomainList, fqdnDomain)
			}
		}
		if err := client.UpdateDomains(fqdn); err != nil {
			return fmt.Errorf("failed to set domain names in update : %w", err)
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
		os := mustSlice(o)
		ns := mustSlice(n)

		mapOldTagList := make(map[string][]string)
		for _, oldGwFilterTags := range os {
			oldGwFilterTag := mustMap(oldGwFilterTags)
			sourceIPs := make([]string, 0)
			for _, sourceIP := range mustSchemaSet(oldGwFilterTag["source_ip_list"]).List() {
				sourceIPs = append(sourceIPs, mustString(sourceIP))
			}
			mapOldTagList[mustString(oldGwFilterTag["gw_name"])] = sourceIPs
		}

		for _, newGwFilterTags := range ns {
			newGwFilterTag := mustMap(newGwFilterTags)
			gateway := &goaviatrix.Gateway{
				GwName: mustString(newGwFilterTag["gw_name"]),
			}

			sourceIPs := make([]string, 0)
			for _, sourceIP := range mustSchemaSet(newGwFilterTag["source_ip_list"]).List() {
				sourceIPs = append(sourceIPs, mustString(sourceIP))
			}

			val, ok := mapOldTagList[gateway.GwName]
			if !ok {
				err := client.AttachTagToGw(fqdn, gateway)
				if err != nil {
					return fmt.Errorf("failed to add filter tag to gateway : %w", err)
				}
				err = client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
				if err != nil {
					return fmt.Errorf("failed to update source ips to gateway : %w", err)
				}
				continue
			}

			if !goaviatrix.Equivalent(val, sourceIPs) {
				err := client.UpdateSourceIPFilters(fqdn, gateway, sourceIPs)
				if err != nil {
					return fmt.Errorf("failed to update source ips to gateway : %w", err)
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
				return fmt.Errorf("failed to delete GWs for fqdn in update: %w", err)
			}
		}
	}

	d.Partial(false)
	d.SetId(fqdn.FQDNTag)
	return resourceAviatrixFQDNRead(d, meta)
}

func resourceAviatrixFQDNDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	fqdn := &goaviatrix.FQDN{
		FQDNTag: getString(d, "fqdn_tag"),
	}

	log.Printf("[INFO] Deleting Aviatrix FQDN: %#v", fqdn)

	gwList, err := client.ListGws(fqdn)
	if err != nil {
		return fmt.Errorf("failed to get GW list for fqdn: %w", err)
	}
	err = client.DetachGws(fqdn, gwList)
	if err != nil {
		return fmt.Errorf("failed to delete GWs for fqdn: %w", err)
	}

	err = client.DeleteFQDN(fqdn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FQDN: %w", err)
	}

	return nil
}
