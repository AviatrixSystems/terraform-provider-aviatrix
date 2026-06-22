package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFirewallTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallTagCreate,
		Read:   resourceAviatrixFirewallTagRead,
		Update: resourceAviatrixFirewallTagUpdate,
		Delete: resourceAviatrixFirewallTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"firewall_tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"cidr_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A JSON file with information of 'cidr_tag_name' and 'cidr'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_tag_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name attribute of a policy.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CIDR attribute of a policy.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixFirewallTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	firewallTag := &goaviatrix.FirewallTag{
		Name: getString(d, "firewall_tag"),
	}

	if firewallTag.Name == "" {
		return fmt.Errorf("invalid choice: firewall tag can't be empty")
	}

	d.SetId(firewallTag.Name)
	flag := false
	defer func() { _ = resourceAviatrixFirewallTagReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to create firewall tag: %w", err)
	}

	// If cidr list is present, update cidr list
	if _, ok := d.GetOk("cidr_list"); ok {
		cidrList := getList(d, "cidr_list")
		for _, currCIDR := range cidrList {
			cm := mustMap(currCIDR)
			cidrMember := goaviatrix.CIDRMember{
				CIDRTag: mustString(cm["cidr_tag_name"]),
				CIDR:    mustString(cm["cidr"]),
			}
			if cidrMember.CIDRTag == "" {
				return fmt.Errorf("invalid choice: cidr_tag_name can't be empty")
			}
			if cidrMember.CIDR == "" {
				return fmt.Errorf("invalid choice: cidr can't be empty")
			}
			firewallTag.CIDRList = append(firewallTag.CIDRList, cidrMember)
		}

		err := client.UpdateFirewallTag(firewallTag)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix FirewallTag: %w", err)
		}
	}

	return resourceAviatrixFirewallTagReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFirewallTagReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFirewallTagRead(d, meta)
	}
	return nil
}

func resourceAviatrixFirewallTagRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	fTag := getString(d, "firewall_tag")
	if fTag == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no firewall tag name received. Import Id is %s", id)
		mustSet(d, "firewall_tag", id)
		d.SetId(id)
	}

	firewallTag := &goaviatrix.FirewallTag{
		Name: getString(d, "firewall_tag"),
	}
	fwt, err := client.GetFirewallTag(firewallTag)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error fetching firewall tag %s: %w", firewallTag.Name, err)
	}

	log.Printf("[TRACE] Reading cidr list for tag %s: %#v", firewallTag.Name, fwt)

	if fwt != nil {
		var cidrList []map[string]interface{}
		for _, cidrMember := range fwt.CIDRList {
			cm := make(map[string]interface{})
			cm["cidr_tag_name"] = cidrMember.CIDRTag
			cm["cidr"] = cidrMember.CIDR

			cidrList = append(cidrList, cm)
		}

		if err := d.Set("cidr_list", cidrList); err != nil {
			log.Printf("[WARN] Error setting cidr_list for (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceAviatrixFirewallTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	firewallTag := &goaviatrix.FirewallTag{
		Name: getString(d, "firewall_tag"),
	}

	d.Partial(true)

	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewallTag)

	// Update cidr list
	cidrList := getList(d, "cidr_list")
	for _, currCIDR := range cidrList {
		cm := mustMap(currCIDR)
		cidrMember := goaviatrix.CIDRMember{
			CIDRTag: mustString(cm["cidr_tag_name"]),
			CIDR:    mustString(cm["cidr"]),
		}
		firewallTag.CIDRList = append(firewallTag.CIDRList, cidrMember)
	}

	err := client.UpdateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix FirewallTag: %w", err)
	}

	d.Partial(false)
	return resourceAviatrixFirewallTagRead(d, meta)
}

func resourceAviatrixFirewallTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	firewallTag := &goaviatrix.FirewallTag{
		Name: getString(d, "firewall_tag"),
	}

	err := client.DeleteFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FirewallTag policy list: %w", err)
	}

	return nil
}
