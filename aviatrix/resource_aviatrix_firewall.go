package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "deny-all",
				ValidateFunc: validation.StringInSlice([]string{"deny-all", "allow-all"}, false),
				Description:  "New base policy.",
			},
			"base_log_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether enable logging or not. Valid values: true, false. Default value: false.",
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
				Description: "New access policy for the gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source address, a valid IPv4 address or tag name.",
						},
						"dst_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination address, a valid IPv4 address or tag name.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "all",
							ValidateFunc: validation.StringInSlice([]string{"all", "tcp", "udp", "icmp", "sctp", "rdp", "dccp"}, false),
							Description:  "Valid values: 'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'.",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A single port or a range of port numbers.",
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"allow", "deny", "force-drop"}, false),
							Description:  "Valid values: 'allow', 'deny' or 'force-drop'(in stateful firewall rule to allow immediate packet dropping on established sessions).",
						},
						"log_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Valid values: true, false. Default value: false.",
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

	var mapPolicyKey = make(map[string]bool)
	var previousActionIsForceDrop = true
	// If policies are present and manage_firewall_policies is set to true, update policies
	if hasSetPolicies && enabledInlinePolicies {
		policies := d.Get("policy").([]interface{})
		for index, policy := range policies {
			pl := policy.(map[string]interface{})
			firewallPolicy := &goaviatrix.Policy{
				SrcIP:       pl["src_ip"].(string),
				DstIP:       pl["dst_ip"].(string),
				Protocol:    pl["protocol"].(string),
				Port:        pl["port"].(string),
				Action:      pl["action"].(string),
				Description: pl["description"].(string),
			}

			if !previousActionIsForceDrop && firewallPolicy.Action == "force-drop" {
				return fmt.Errorf("validation on policy rules failed: rule no. %v in policy list is a 'force-drop' rule. It should be ahead of other type rules in Policy list", index+1)
			}
			if previousActionIsForceDrop && firewallPolicy.Action != "force-drop" {
				previousActionIsForceDrop = false
			}

			key := firewallPolicy.SrcIP + "~" + firewallPolicy.DstIP + "~" + firewallPolicy.Protocol + "~" + firewallPolicy.Port
			if mapPolicyKey[key] {
				return fmt.Errorf("validation on policy rules failed: rule no. %v in policy list is a duplicate rule", index+1)
			}
			mapPolicyKey[key] = true

			if firewallPolicy.Protocol == "all" && firewallPolicy.Port != "0:65535" {
				return fmt.Errorf("validation on policy rules failed: rule no. %v's port should be '0:65535' for protocol 'all'", index+1)
			} else if firewallPolicy.Protocol == "all" {
				firewallPolicy.Port = ""
			}
			if firewallPolicy.Protocol == "icmp" && (firewallPolicy.Port != "") {
				return fmt.Errorf("validation on policy rules failed: rule no. %v's port should be empty for protocol 'icmp'", index+1)
			}

			logEnabled := pl["log_enabled"].(bool)
			if logEnabled {
				firewallPolicy.LogEnabled = "on"
			} else {
				firewallPolicy.LogEnabled = "off"
			}

			firewall.PolicyList = append(firewall.PolicyList, firewallPolicy)
		}
	}

	log.Printf("[INFO] Creating Aviatrix firewall: %#v", firewall)

	d.SetId(firewall.GwName)
	flag := false
	defer resourceAviatrixFirewallReadIfRequired(d, meta, &flag)

	//If base_policy or base_log enable is present, set base policy
	if firewall.BasePolicy == "allow-all" {
		firewall.BaseLogEnabled = "off"
		err := client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base firewall policy for GW %s: %s", firewall.GwName, err)
		}
	}

	baseLogEnabled := d.Get("base_log_enabled").(bool)
	if baseLogEnabled {
		firewall.BaseLogEnabled = "on"
		err := client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to enable base logging for GW %s: %s", firewall.GwName, err)
		}
	}

	if hasSetPolicies && enabledInlinePolicies {
		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set Aviatrix firewall policies for GW %s: %s", firewall.GwName, err)
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

	var policiesFromFile []map[string]interface{}
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
			policiesFromFile = append(policiesFromFile, goaviatrix.PolicyToMap(policy))
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
		} else if firewall.BasePolicy == "deny" {
			firewall.BasePolicy = "deny-all"
		}

		if d.HasChange("base_log_enabled") {
			old, _ := d.GetChange("base_log_enabled")
			if old.(bool) {
				firewall.BaseLogEnabled = "on"
			} else {
				firewall.BaseLogEnabled = "off"
			}
		} else {
			if d.Get("base_log_enabled").(bool) {
				firewall.BaseLogEnabled = "on"
			} else {
				firewall.BaseLogEnabled = "off"
			}
		}

		err := client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to update base firewall policies for GW %s: %s", firewall.GwName, err)
		}
	}

	if ok := d.HasChange("base_log_enabled"); ok {
		firewall.BasePolicy = d.Get("base_policy").(string)
		if d.Get("base_log_enabled").(bool) {
			firewall.BaseLogEnabled = "on"
		} else {
			firewall.BaseLogEnabled = "off"
		}

		err := client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to update base logging for GW %s: %s", firewall.GwName, err)
		}
	}

	if ok := d.HasChange("policy"); ok && enabledInlinePolicies {
		var mapPolicyKey = make(map[string]bool)
		var previousActionIsForceDrop = true
		policies := d.Get("policy").([]interface{})
		for index, policy := range policies {
			pl := policy.(map[string]interface{})
			firewallPolicy := &goaviatrix.Policy{
				SrcIP:       pl["src_ip"].(string),
				DstIP:       pl["dst_ip"].(string),
				Protocol:    pl["protocol"].(string),
				Port:        pl["port"].(string),
				Action:      pl["action"].(string),
				Description: pl["description"].(string),
			}

			if !previousActionIsForceDrop && firewallPolicy.Action == "force-drop" {
				return fmt.Errorf("validation on policy rules failed: rule no. %v in policy list is a 'force-drop' rule. It should be ahead of other type rules in Policy list", index+1)
			}
			if previousActionIsForceDrop && firewallPolicy.Action != "force-drop" {
				previousActionIsForceDrop = false
			}

			key := firewallPolicy.SrcIP + "~" + firewallPolicy.DstIP + "~" + firewallPolicy.Protocol + "~" + firewallPolicy.Port
			if mapPolicyKey[key] {
				return fmt.Errorf("validation on policy rules failed: rule no. %v in policy list is a duplicate rule", index+1)
			}
			mapPolicyKey[key] = true

			if firewallPolicy.Protocol == "all" && firewallPolicy.Port != "0:65535" {
				return fmt.Errorf("validation on policy rules failed: rule no. %v's port should be '0:65535' for protocol 'all'", index+1)
			} else if firewallPolicy.Protocol == "all" {
				firewallPolicy.Port = ""
			}
			if firewallPolicy.Protocol == "icmp" && (firewallPolicy.Port != "") {
				return fmt.Errorf("validation on policy rules failed: rule no. %v's port should be empty for protocol 'icmp'", index+1)
			}

			logEnabled := pl["log_enabled"].(bool)
			if logEnabled {
				firewallPolicy.LogEnabled = "on"
			} else {
				firewallPolicy.LogEnabled = "off"
			}

			firewall.PolicyList = append(firewall.PolicyList, firewallPolicy)
		}

		err := client.UpdatePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Firewall policy: %s", err)
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

	if d.Get("base_policy").(string) != "deny-all" {
		firewall.BasePolicy = "deny-all"
		baseLogEnabled := d.Get("base_log_enabled").(bool)
		if baseLogEnabled {
			firewall.BaseLogEnabled = "on"
		} else {
			firewall.BaseLogEnabled = "off"
		}
		err = client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base firewall policies to default: %s", err)
		}
	}

	if d.Get("base_log_enabled").(bool) {
		firewall.BasePolicy = "deny-all"

		firewall.BaseLogEnabled = "off"
		err = client.SetBasePolicy(firewall)
		if err != nil {
			return fmt.Errorf("failed to set base logging to default: %s", err)
		}
	}

	return nil
}
