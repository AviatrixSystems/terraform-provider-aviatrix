package aviatrix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixFirewall() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFirewallRead,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the gateway the firewall is associated with.",
			},
			"base_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The firewall's base policy.",
			},
			"base_log_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether logging is enabled or not.",
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of policies associated with the firewall.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDRs separated by a comma or tag names such 'HR' or 'marketing' etc.",
						},
						"dst_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDRs separated by a comma or tag names such 'HR' or 'marketing' etc.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A single port or a range of port numbers.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid values: 'allow', 'deny' or 'force-drop'(in stateful firewall rule to allow immediate packet dropping on established sessions).",
						},
						"log_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Valid values: true or false.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of this firewall policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)

	firewall := &goaviatrix.Firewall{
		GwName: gwName,
	}

	fw, err := client.GetPolicy(firewall)

	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error fetching firewall policy for gateway %s: %s", firewall.GwName, err)
	}

	d.Set("gw_name", gwName)

	d.Set("base_policy", "deny-all")
	if fw.BasePolicy == "allow-all" {
		d.Set("base_policy", "allow-all")
	}

	d.Set("base_log_enabled", false)
	if fw.BaseLogEnabled == "on" {
		d.Set("base_log_enabled", true)
	}

	var policies []map[string]interface{}
	for _, p := range fw.PolicyList {
		policies = append(policies, goaviatrix.PolicyToMap(p))
	}
	if err = d.Set("policies", policies); err != nil {
		return fmt.Errorf("error setting firewall policies for gateway %s: %s", firewall.GwName, err)
	}

	d.SetId(gwName)

	return nil
}
