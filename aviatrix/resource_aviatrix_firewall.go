package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallCreate,
		Read:   resourceAviatrixFirewallRead,
		Update: resourceAviatrixFirewallUpdate,
		Delete: resourceAviatrixFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixFirewallResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixFirewallStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of gateway.",
			},
			"base_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "deny-all",
				Description: "New base policy.",
			},
			"base_log_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether enable logging or not. Valid values: true or false.",
			},
			"manage_firewall_policies": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Enable to manage firewall policies via in-line rules. If false, policies must be managed " +
					"using `aviatrix_firewall_policy` resources.",
			},
			"policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Deprecated:  "Please set `manage_firewall_policies` to false, and use the standalone aviatrix_firewall_policy resource instead.",
				Description: "New access policy for the gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CIDRs separated by comma or tag names such 'HR' or 'marketing' etc.",
						},
						"dst_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CIDRs separated by comma or tag names such 'HR' or 'marketing' etc.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "all",
							Description: "'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'.",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A single port or a range of port numbers.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Valid values: 'allow', 'deny' or 'force-drop'(in stateful firewall rule to allow immediate packet dropping on established sessions).",
						},
						"log_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Valid values: true or false.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Description of this firewall policy.",
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
		GwName:     d.Get("gw_name").(string),
		BasePolicy: d.Get("base_policy").(string),
	}

	_, hasSetPolicies := d.GetOk("policy")
	enabledInlinePolicies := d.Get("manage_firewall_policies").(bool)
	if hasSetPolicies && !enabledInlinePolicies {
		return fmt.Errorf("manage_firewall_policies must be set to true to set in-line policies")
	}

	if firewall.BasePolicy != "allow-all" && firewall.BasePolicy != "deny-all" {
		return fmt.Errorf("base_policy can only be 'allow-all', or 'deny-all'")
	}

	baseLogEnabled := d.Get("base_log_enabled").(bool)
	if baseLogEnabled {
		firewall.BaseLogEnabled = "on"
	} else {
		firewall.BaseLogEnabled = "off"
	}

	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewall)

	d.SetId(firewall.GwName)
	flag := false
	defer resourceAviatrixFirewallReadIfRequired(d, meta, &flag)

	//If base_policy or base_log enable is present, set base policy
	if firewall.BasePolicy != "" {
		err := client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base firewall policies for GW %s: %s", firewall.GwName, err)
		}
	}

	// If policies are present and manage_firewall_policies is set to true, update policies
	if hasSetPolicies && enabledInlinePolicies {
		policies := d.Get("policy").([]interface{})
		for _, policy := range policies {
			pl := policy.(map[string]interface{})
			firewallPolicy := &goaviatrix.Policy{
				SrcIP:       pl["src_ip"].(string),
				DstIP:       pl["dst_ip"].(string),
				Protocol:    pl["protocol"].(string),
				Port:        pl["port"].(string),
				Action:      pl["action"].(string),
				Description: pl["description"].(string),
			}

			logEnabled := pl["log_enabled"].(bool)
			if logEnabled {
				firewallPolicy.LogEnabled = "on"
			} else {
				firewallPolicy.LogEnabled = "off"
			}

			err := client.ValidatePolicy(firewallPolicy)
			if err != nil {
				return fmt.Errorf("policy validation failed: %v", err)
			}
			if firewallPolicy.Protocol == "all" {
				firewallPolicy.Port = ""
			}

			firewall.PolicyList = append(firewall.PolicyList, firewallPolicy)
		}

		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Firewall: %s", err)
		}
	}

	return resourceAviatrixFirewallReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFirewallReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFirewallRead(d, meta)
	}
	return nil
}

func resourceAviatrixFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.Set("manage_firewall_policies", true)
		d.SetId(id)
	}

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

	log.Printf("[TRACE] Reading policy for gateway %s: %#v", firewall.GwName, fw)

	policyMap := make(map[string]map[string]interface{})
	var policyKeyArray []string
	if fw != nil {
		if fw.BasePolicy == "allow-all" {
			d.Set("base_policy", "allow-all")
		} else {
			d.Set("base_policy", "deny-all")
		}
		if fw.BaseLogEnabled == "on" {
			d.Set("base_log_enabled", true)
		} else {
			d.Set("base_log_enabled", false)
		}

		for _, policy := range fw.PolicyList {
			pl := goaviatrix.PolicyToMap(policy)
			key := policy.SrcIP + "~" + policy.DstIP + "~" + policy.Protocol + "~" + policy.Port
			policyMap[key] = pl
			policyKeyArray = append(policyKeyArray, key)
		}
	}

	var policiesFromFile []map[string]interface{}
	policies := d.Get("policy").([]interface{})
	for _, policy := range policies {
		pl := policy.(map[string]interface{})
		firewallPolicy := &goaviatrix.Policy{
			SrcIP:       pl["src_ip"].(string),
			DstIP:       pl["dst_ip"].(string),
			Protocol:    pl["protocol"].(string),
			Port:        pl["port"].(string),
			Action:      pl["action"].(string),
			Description: pl["description"].(string),
		}
		logEnabled := pl["log_enabled"].(bool)
		if logEnabled {
			firewallPolicy.LogEnabled = "on"
		} else {
			firewallPolicy.LogEnabled = "off"
		}

		key := firewallPolicy.SrcIP + "~" + firewallPolicy.DstIP + "~" + firewallPolicy.Protocol + "~" + firewallPolicy.Port
		if val, ok := policyMap[key]; ok {
			if goaviatrix.CompareMapOfInterface(pl, val) {
				policiesFromFile = append(policiesFromFile, pl)
				delete(policyMap, key)
			}
		}
	}
	if len(policyKeyArray) != 0 {
		for i := 0; i < len(policyKeyArray); i++ {
			if policyMap[policyKeyArray[i]] != nil {
				policiesFromFile = append(policiesFromFile, policyMap[policyKeyArray[i]])
			}
		}
	}

	// Only write policies to state if the user has enabled in-line policies.
	if d.Get("manage_firewall_policies").(bool) {
		if err := d.Set("policy", policiesFromFile); err != nil {
			log.Printf("[WARN] Error setting policy for (%s): %s", d.Id(), err)
		}
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

	_, hasSetPolicies := d.GetOk("policy")
	enabledInlinePolicies := d.Get("manage_firewall_policies").(bool)
	if hasSetPolicies && !enabledInlinePolicies {
		return fmt.Errorf("manage_firewall_policies must be set to true to set in-line policies")
	}

	if ok := d.HasChange("base_policy"); ok {
		firewall.BasePolicy = d.Get("base_policy").(string)
		if firewall.BasePolicy == "allow" {
			firewall.BasePolicy = "allow-all"
		}
		if firewall.BasePolicy == "deny" {
			firewall.BasePolicy = "deny-all"
		}
		if d.Get("base_log_enabled").(bool) {
			firewall.BaseLogEnabled = "on"
		} else {
			firewall.BaseLogEnabled = "off"
		}
	}

	if ok := d.HasChange("base_log_enabled"); ok {
		firewall.BasePolicy = d.Get("base_policy").(string)

		baseLogEnabled := d.Get("base_log_enabled").(bool)
		if baseLogEnabled {
			firewall.BaseLogEnabled = "on"
		} else {
			firewall.BaseLogEnabled = "off"
		}
	}

	//If base_policy or base_log enable is present, first delete
	//existing policies, set base policy, and then reapply deleted policies.
	if firewall.BasePolicy != "" || firewall.BaseLogEnabled != "" {
		firewall.PolicyList = make([]*goaviatrix.Policy, 0)
		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Firewall: %s", err)
		}
		err = client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base firewall policies for GW %s: %s", firewall.GwName, err)
		}
	}

	//If policy list is present, update policy list
	if ok := d.HasChange("policy"); ok && enabledInlinePolicies {
		policies := d.Get("policy").([]interface{})
		for _, policy := range policies {
			pl := policy.(map[string]interface{})
			firewallPolicy := &goaviatrix.Policy{
				SrcIP:       pl["src_ip"].(string),
				DstIP:       pl["dst_ip"].(string),
				Protocol:    pl["protocol"].(string),
				Port:        pl["port"].(string),
				Action:      pl["action"].(string),
				Description: pl["description"].(string),
			}

			if pl["log_enabled"].(bool) {
				firewallPolicy.LogEnabled = "on"
			} else {
				firewallPolicy.LogEnabled = "off"
			}

			err := client.ValidatePolicy(firewallPolicy)
			if err != nil {
				return fmt.Errorf("policy validation failed: %v", err)
			}
			if firewallPolicy.Protocol == "all" {
				firewallPolicy.Port = ""
			}

			firewall.PolicyList = append(firewall.PolicyList, firewallPolicy)
		}

		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Firewall: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixFirewallRead(d, meta)
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
