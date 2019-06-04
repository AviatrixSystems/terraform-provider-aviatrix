package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFirewallTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallTagCreate,
		Read:   resourceAviatrixFirewallTagRead,
		Update: resourceAviatrixFirewallTagUpdate,
		Delete: resourceAviatrixFirewallTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
							Optional:    true,
							Description: "The name attribute of a policy.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CIDR attribute of a policy.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixFirewallTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	if firewallTag.Name == "" {
		return fmt.Errorf("invalid choice: firewall tag can't be empty")
	}
	err := client.CreateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to create firewall tag: %s", err)
	}
	d.SetId(firewallTag.Name)
	//If cidr list is present, update cidr list
	if _, ok := d.GetOk("cidr_list"); ok {
		cidrList := d.Get("cidr_list").([]interface{})
		for _, currCIDR := range cidrList {
			cm := currCIDR.(map[string]interface{})
			cidrMember := goaviatrix.CIDRMember{
				CIDRTag: cm["cidr_tag_name"].(string),
				CIDR:    cm["cidr"].(string),
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
			return fmt.Errorf("failed to update Aviatrix FirewallTag: %s", err)
		}
	}
	return resourceAviatrixFirewallTagRead(d, meta)
}

func resourceAviatrixFirewallTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fTag := d.Get("firewall_tag").(string)
	if fTag == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no firewall tag name received. Import Id is %s", id)
		d.Set("firewall_tag", id)
		d.SetId(id)
	}

	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	fwt, err := client.GetFirewallTag(firewallTag)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error fetching firewall tag %s: %s", firewallTag.Name, err)
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
		d.Set("cidr_list", cidrList)
	}

	return nil
}

func resourceAviatrixFirewallTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewallTag)
	//Update cidr list
	cidrList := d.Get("cidr_list").([]interface{})
	for _, currCIDR := range cidrList {
		cm := currCIDR.(map[string]interface{})
		cidrMember := goaviatrix.CIDRMember{
			CIDRTag: cm["cidr_tag_name"].(string),
			CIDR:    cm["cidr"].(string),
		}
		firewallTag.CIDRList = append(firewallTag.CIDRList, cidrMember)
	}
	err := client.UpdateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix FirewallTag: %s", err)
	}
	if _, ok := d.GetOk("cidr_list"); ok {
		d.SetPartial("cidr_list")
	}
	d.Partial(false)
	return nil
}

func resourceAviatrixFirewallTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewallTag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	err := client.UpdateFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FirewallTag policy list: %s", err)
	}
	err = client.DeleteFirewallTag(firewallTag)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix FirewallTag policy list: %s", err)
	}
	return nil
}
