package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallCreate,
		Read:   resourceAviatrixFirewallRead,
		Update: resourceAviatrixFirewallUpdate,
		Delete: resourceAviatrixFirewallDelete,

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_allow_deny": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"base_log_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dst_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"allow_deny": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_enable": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall := &goaviatrix.Firewall{
		GwName: d.Get("gw_name").(string),
	}
	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewall)
	if _, ok := d.GetOk("base_allow_deny"); ok {
		firewall.BaseAllowDeny = d.Get("base_allow_deny").(string)
		if firewall.BaseAllowDeny == "allow" {
			firewall.BaseAllowDeny = "allow-all"
		}
		if firewall.BaseAllowDeny == "deny" {
			firewall.BaseAllowDeny = "deny-all"
		}
	}
	if _, ok := d.GetOk("base_log_enable"); ok {
		firewall.BaseLogEnable = d.Get("base_log_enable").(string)
	}
	//If base_allow_deny or base_log enable is present, set base policy
	if firewall.BaseAllowDeny != "" || firewall.BaseLogEnable != "" {

		err := client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base firewall policies for GW %s: %s", firewall.GwName, err)
		}
	}
	//If policy list is present, update policy list
	if _, ok := d.GetOk("policy"); ok {
		policies := d.Get("policy").([]interface{})
		for _, policy := range policies {
			pl := policy.(map[string]interface{})
			firewallPolicy := &goaviatrix.Policy{
				SrcIP:     pl["src_ip"].(string),
				DstIP:     pl["dst_ip"].(string),
				Protocol:  pl["protocol"].(string),
				Port:      pl["port"].(string),
				AllowDeny: pl["allow_deny"].(string),
				LogEnable: pl["log_enable"].(string),
			}
			firewall.PolicyList = append(firewall.PolicyList, firewallPolicy)
		}
		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Firewall: %s", err)
		}
	}
	d.SetId(firewall.GwName)
	return nil
}

func resourceAviatrixFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall := &goaviatrix.Firewall{
		GwName: d.Get("gw_name").(string),
	}
	fw, err := client.GetPolicy(firewall)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error fetching policy for gateway %s: %s", firewall.GwName, err)
	}
	log.Printf("[TRACE] Reading policy for gateway %s: %#v",
		firewall.GwName, fw)
	if fw != nil {
		if fw.BaseAllowDeny == "allow-all" {
			d.Set("base_allow_deny", "allow")
		}
		if fw.BaseAllowDeny == "deny-all" {
			d.Set("base_allow_deny", "deny")
		}

		var policies []map[string]interface{}
		for _, policy := range fw.PolicyList {
			pl := make(map[string]interface{})
			pl["src_ip"] = policy.SrcIP
			pl["dst_ip"] = policy.DstIP
			pl["protocol"] = policy.Protocol
			pl["port"] = policy.Port
			pl["log_enable"] = policy.LogEnable
			pl["allow_deny"] = policy.AllowDeny
			policies = append(policies, pl)
		}
		d.Set("policy", policies)
	}

	return nil
}

func resourceAviatrixFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall := &goaviatrix.Firewall{
		GwName: d.Get("gw_name").(string),
	}
	d.Partial(true)
	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewall)
	if ok := d.HasChange("base_allow_deny"); ok {
		firewall.BaseAllowDeny = d.Get("base_allow_deny").(string)
		if firewall.BaseAllowDeny == "allow" {
			firewall.BaseAllowDeny = "allow-all"
		}
		if firewall.BaseAllowDeny == "deny" {
			firewall.BaseAllowDeny = "deny-all"
		}
		firewall.BaseLogEnable = d.Get("base_log_enable").(string)
	}
	if ok := d.HasChange("base_log_enable"); ok {
		firewall.BaseAllowDeny = d.Get("base_allow_deny").(string)
		firewall.BaseLogEnable = d.Get("base_log_enable").(string)
	}
	//If base_allow_deny or base_log enable is present, first delete
	//existing policies, set base policy, and then reapply deleted policies.
	if firewall.BaseAllowDeny != "" || firewall.BaseLogEnable != "" {
		firewall.PolicyList = make([]*goaviatrix.Policy, 0)
		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Firewall: %s", err)
		}
		err = client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base firewall policies for GW %s: %s", firewall.GwName, err)
		}
		if d.HasChange("base_allow_deny") {
			d.SetPartial("base_allow_deny")
		}
		if d.HasChange("base_log_enable") {
			d.SetPartial("base_log_enable")
		}
	}
	//If policy list is present, update policy list
	if _, ok := d.GetOk("policy"); ok {
		policies := d.Get("policy").([]interface{})
		for _, policy := range policies {
			pl := policy.(map[string]interface{})
			firewallPolicy := &goaviatrix.Policy{
				SrcIP:     pl["src_ip"].(string),
				DstIP:     pl["dst_ip"].(string),
				Protocol:  pl["protocol"].(string),
				Port:      pl["port"].(string),
				AllowDeny: pl["allow_deny"].(string),
				LogEnable: pl["log_enable"].(string),
			}
			firewall.PolicyList = append(firewall.PolicyList, firewallPolicy)
		}
		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Firewall: %s", err)
		}
		d.SetPartial("policy")
	}
	d.Partial(false)
	return nil
}

func resourceAviatrixFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	firewall := &goaviatrix.Firewall{
		GwName: d.Get("gw_name").(string),
	}
	firewall.PolicyList = make([]*goaviatrix.Policy, 0)
	err := client.UpdatePolicy(firewall)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Firewall policy list: %s", err)
	}
	//FIXME: Need to reset base policy rules and base logging too to
	//allow-all and on(default values).
	//There is a bug in API set_vpc_base_policy, in which changing
	//both base_policy and base_policy_log_enable together to the
	//opposite of their current values gives error.
	//Add base policy resetting after the bug gets fixed
	return nil
}
