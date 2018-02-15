package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixFirewallTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallTagCreate,
		Read:   resourceAviatrixFirewallTagRead,
		Update: resourceAviatrixFirewallTagUpdate,
		Delete: resourceAviatrixFirewallTagDelete,

		Schema: map[string]*schema.Schema{
			"firewall_tag": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_list": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_tag_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixFirewallTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall_tag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	err := client.CreateFirewallTag(firewall_tag)
	if err != nil {
		return fmt.Errorf("Failed to create firewall tag: %s", err)
	}
	//If cidr list is present, update cidr list
	if _, ok := d.GetOk("cidr_list"); ok {
		cidr_list := d.Get("cidr_list").([]interface{})
		for _, curr_cidr := range cidr_list {
			cm := curr_cidr.(map[string]interface{})
			cidr_member := goaviatrix.CIDRMember{
				CIDRTag: cm["cidr_tag_name"].(string),
				CIDR:    cm["cidr"].(string),
			}
			firewall_tag.CIDRList = append(firewall_tag.CIDRList, cidr_member)
		}
		err := client.UpdateFirewallTag(firewall_tag)
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix FirewallTag: %s", err)
		}
	}
	d.SetId(firewall_tag.Name)
	return nil
}

func resourceAviatrixFirewallTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall_tag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	fwt, err := client.GetFirewallTag(firewall_tag)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error fetching firewall tag %s: %s", firewall_tag.Name, err)
	}
	log.Printf("[TRACE] Reading cidr list for tag %s: %#v", firewall_tag.Name, fwt)
	if fwt != nil {

		var cidr_list []map[string]interface{}
		for _, cidr_member := range fwt.CIDRList {
			cm := make(map[string]interface{})
			cm["cidr_tag_name"] = cidr_member.CIDRTag
			cm["cidr"] = cidr_member.CIDR

			cidr_list = append(cidr_list, cm)
		}
		d.Set("cidr_list", cidr_list)
	}

	return nil
}

func resourceAviatrixFirewallTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall_tag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewall_tag)
	//Update cidr list
	cidr_list := d.Get("cidr_list").([]interface{})
	for _, curr_cidr := range cidr_list {
		cm := curr_cidr.(map[string]interface{})
		cidr_member := goaviatrix.CIDRMember{
			CIDRTag: cm["cidr_tag_name"].(string),
			CIDR:    cm["cidr"].(string),
		}
		firewall_tag.CIDRList = append(firewall_tag.CIDRList, cidr_member)
	}
	err := client.UpdateFirewallTag(firewall_tag)
	if err != nil {
		return fmt.Errorf("Failed to update Aviatrix FirewallTag: %s", err)
	}
	if _, ok := d.GetOk("cidr_list"); ok {
		d.SetPartial("cidr_list")
	}
	d.Partial(false)
	return nil
}

func resourceAviatrixFirewallTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall_tag := &goaviatrix.FirewallTag{
		Name: d.Get("firewall_tag").(string),
	}
	//firewall_tag.CIDRList = make([]*goaviatrix.CIDRMember, 0)
	err := client.UpdateFirewallTag(firewall_tag)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix FirewallTag policy list: %s", err)
	}
	err = client.DeleteFirewallTag(firewall_tag)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix FirewallTag policy list: %s", err)
	}
	return nil
}
