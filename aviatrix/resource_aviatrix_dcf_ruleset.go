package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

//nolint:funlen
func resourceAviatrixDCFRuleset() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFRulesetCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFRulesetRead,
		UpdateWithoutTimeout: resourceAviatrixDCFRulesetUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFRulesetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the ruleset.",
			},
			"system_resource": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the ruleset is a system resource.",
			},
			"attach_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The attachment point to which the ruleset is attached.",
			},
			"rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of distributed-firewalling rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DENY", "PERMIT", "DEEP_PACKET_INSPECTION_PERMIT", "INTRUSION_DETECTION_PERMIT"}, false),
							Description: "Action for the specified source and destination Smart Groups. " +
								"Must be one of DEEP_PACKET_INSPECTION_PERMIT, INTRUSION_DETECTION_PERMIT, PERMIT or DENY.",
						},
						"decrypt_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "DECRYPT_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"DECRYPT_UNSPECIFIED", "DECRYPT_ALLOWED", "DECRYPT_NOT_ALLOWED"}, false),
							Description: "Decryption options for the rule. " +
								"Must be one of DECRYPT_UNSPECIFIED, DECRYPT_ALLOWED or DECRYPT_NOT_ALLOWED.",
						},
						"dst_smart_groups": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of destination Smart Group UUIDs for the rule.",
						},
						"exclude_sg_orchestration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If this flag is set to true, this rule will be ignored for SG orchestration.",
						},
						"flow_app_requirement": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "APP_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"APP_UNSPECIFIED", "TLS_REQUIRED", "NOT_TLS_REQUIRED"}, false),
							Description: "Flow application requirement for the rule. " +
								"Must be one of APP_UNSPECIFIED, TLS_REQUIRED or NOT_TLS_REQUIRED.",
						},
						"logging": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable logging for the rule.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the rule. Rule names must be unique.",
						},
						"port_ranges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of port ranges for the rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lo": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Lower bound of port range.",
									},
									"hi": {
										Type:             schema.TypeInt,
										Optional:         true,
										ValidateFunc:     validation.IntAtLeast(0),
										DiffSuppressFunc: DiffSuppressFuncDistributedFirewallingPolicyPortRangeHi,
										Description:      "Upper bound of port range.",
									},
								},
							},
							MaxItems: 64,
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Priority level of the rule",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP", "ANY"}, true),
							DiffSuppressFunc: func(_, oldProto, newProto string, _ *schema.ResourceData) bool {
								return strings.EqualFold(oldProto, newProto)
							},
							Description: "Protocol for the rule to filter. " +
								"Must be one of ANY, ICMP, TCP or UDP.",
						},
						"src_smart_groups": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of source Smart Group UUIDs for the rule.",
						},
						"tls_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "TLS profile UUID for the rule.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the rule.",
						},
						"watch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable watch mode for the rule.",
						},
						"web_groups": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of Web Group UUIDs for the rule.",
						},
						"log_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "def000ad-7000-0000-0000-000000000001",
							Description: "Log profile UUID for the rule. This will be set to Log at Start by default which has a UUID of def000ad-7000-0000-0000-000000000001.",
							// The log profile UUID must be one of the predefined log profile UUIDs
							// def000ad-7000-0000-0000-000000000001: DEF_LOG_PROFILE_START
							// def000ad-7000-0000-0000-000000000002: DEF_LOG_PROFILE_END
							// def000ad-7000-0000-0000-000000000003: DEF_LOG_PROFILE_ALL
							// TODO(ACK): AVX-68895@everclear-CF2, implement API+datasource
							ValidateFunc: validation.StringInSlice([]string{"def000ad-7000-0000-0000-000000000001", "def000ad-7000-0000-0000-000000000002", "def000ad-7000-0000-0000-000000000003"}, false),
						},
					},
				},
			},
		},
	}
}

func marshalDCFRulesetInput(d *schema.ResourceData) (*goaviatrix.DCFPolicyList, error) {
	policyList := &goaviatrix.DCFPolicyList{}

	name := getString(d, "name")

	policyList.Name = name

	attachTo := getString(d, "attach_to")

	policyList.AttachTo = attachTo

	policiesSet, ok := d.Get("rules").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("ruleset rules must be of type *schema.Set")
	}
	policies := []interface{}{}
	if policiesSet != nil {
		policies = policiesSet.List()
	}
	policyList.Policies = []goaviatrix.DCFPolicy{}
	for _, policyInterface := range policies {
		var ok bool

		policyMap, ok := policyInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("rules must be of type map[string]interface{}")
		}

		policy, err := marshalPolicyInput(policyMap)
		if err != nil {
			return nil, err
		}

		policyList.Policies = append(policyList.Policies, *policy)
	}

	policyList.UUID = d.Id()

	return policyList, nil
}

//nolint:funlen,cyclop
func marshalPolicyInput(policyMap map[string]interface{}) (*goaviatrix.DCFPolicy, error) {
	var ok bool
	var err error
	policy := &goaviatrix.DCFPolicy{}

	policy.Name, ok = policyMap["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name must be of type string")
	}

	policy.Action, ok = policyMap["action"].(string)
	if !ok {
		return nil, fmt.Errorf("action must be of type string")
	}

	policy.Priority, ok = policyMap["priority"].(int)
	if !ok {
		return nil, fmt.Errorf("priority must be of type int")
	}

	policy.LogProfile, ok = policyMap["log_profile"].(string)
	if !ok {
		return nil, fmt.Errorf("log_profile must be of type string")
	}

	policy.FlowAppRequirement, ok = policyMap["flow_app_requirement"].(string)
	if !ok {
		return nil, fmt.Errorf("flow_app_requirement must be of type string")
	}

	policy.DecryptPolicy, ok = policyMap["decrypt_policy"].(string)
	if !ok {
		return nil, fmt.Errorf("decrypt_policy must be of type string")
	}

	policy.ExcludeSgOrchestration, ok = policyMap["exclude_sg_orchestration"].(bool)
	if !ok {
		return nil, fmt.Errorf("exclude_sg_orchestration must be of type bool")
	}

	protocol, ok := policyMap["protocol"].(string)
	if !ok {
		return nil, fmt.Errorf("protocol must be of type string")
	}
	protocol = strings.ToUpper(protocol)

	if protocol == "ANY" {
		policy.Protocol = "PROTOCOL_UNSPECIFIED"
	} else {
		policy.Protocol = protocol
	}

	policy.SrcSmartGroups, err = marshalSmartGroupsInput(policyMap, "src_smart_groups")
	if err != nil {
		return nil, err
	}

	policy.DstSmartGroups, err = marshalSmartGroupsInput(policyMap, "dst_smart_groups")
	if err != nil {
		return nil, err
	}

	policy.WebGroups, err = marshalSmartGroupsInput(policyMap, "web_groups")
	if err != nil {
		return nil, err
	}

	policy.Logging, ok = policyMap["logging"].(bool)
	if !ok {
		return nil, fmt.Errorf("logging must be of type bool")
	}

	policy.Watch, ok = policyMap["watch"].(bool)
	if !ok {
		return nil, fmt.Errorf("watch must be of type bool")
	}

	if goaviatrix.MapContains(policyMap, "port_ranges") {
		if policy.Protocol == "ICMP" {
			return nil, fmt.Errorf("%q must not be set when %q is %q", "port_ranges", "protocol", "ICMP")
		}

		portRanges, ok := policyMap["port_ranges"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("port_ranges must be of type []interface{}")
		}

		for _, portRangeInterface := range portRanges {
			portRangeMap, ok := portRangeInterface.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("port_ranges items must be of type []interface{}")
			}

			portRange := &goaviatrix.DCFPortRange{}

			portRange.Lo, ok = portRangeMap["lo"].(int)
			if !ok {
				return nil, fmt.Errorf("port range lo must be of type int")
			}

			portRange.Hi, ok = portRangeMap["hi"].(int)
			if !ok {
				return nil, fmt.Errorf("port range hi must be of type int")
			}

			policy.PortRanges = append(policy.PortRanges, *portRange)
		}
	}

	policy.UUID, ok = policyMap["uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("uuid must be of type string")
	}

	if tlsProfileUUID, ok := policyMap["tls_profile"]; ok {
		policy.TLSProfile, ok = tlsProfileUUID.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for tls_profile, should be a string")
		}
	}

	return policy, nil
}

func marshalSmartGroupsInput(policyMap map[string]interface{}, key string) ([]string, error) {
	var smartGroups []string

	smartGroupsSet, ok := policyMap[key].(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("%s must be of type *schema.Set", key)
	}

	for _, smartGroup := range smartGroupsSet.List() {
		smartGroupStr, ok := smartGroup.(string)
		if !ok {
			return nil, fmt.Errorf("%s must be of type string", key)
		}

		smartGroups = append(smartGroups, smartGroupStr)
	}

	return smartGroups, nil
}

func resourceAviatrixDCFRulesetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	policyList, err := marshalDCFRulesetInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Ruleset during create: %s", err)
	}

	uuid, err := client.CreateDCFPolicyList(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to create DCF Ruleset: %s", err)
	}

	d.SetId(uuid)

	return nil
}

//nolint:funlen,cyclop
func resourceAviatrixDCFRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	policyList, err := client.GetDCFPolicyList(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF Ruleset: %s", err)
	}

	var policies []map[string]interface{}
	for _, policy := range policyList.Policies {
		if policy.SystemResource {
			continue
		}
		p := make(map[string]interface{})
		p["name"] = policy.Name
		p["action"] = policy.Action
		p["priority"] = policy.Priority
		p["src_smart_groups"] = policy.SrcSmartGroups
		p["dst_smart_groups"] = policy.DstSmartGroups
		p["web_groups"] = policy.WebGroups
		p["logging"] = policy.Logging
		p["watch"] = policy.Watch
		p["uuid"] = policy.UUID
		p["exclude_sg_orchestration"] = policy.ExcludeSgOrchestration
		p["log_profile"] = policy.LogProfile

		if strings.EqualFold(policy.Protocol, "PROTOCOL_UNSPECIFIED") {
			p["protocol"] = "ANY"
		} else {
			p["protocol"] = policy.Protocol
		}
		p["flow_app_requirement"] = policy.FlowAppRequirement
		p["decrypt_policy"] = policy.DecryptPolicy

		if policy.Protocol != "ICMP" {
			var portRanges []map[string]interface{}
			for _, portRange := range policy.PortRanges {
				portRangeMap := map[string]interface{}{
					"hi": portRange.Hi,
					"lo": portRange.Lo,
				}
				portRanges = append(portRanges, portRangeMap)
			}
			p["port_ranges"] = portRanges
		}
		p["tls_profile"] = policy.TLSProfile

		policies = append(policies, p)
	}

	if err := d.Set("name", policyList.Name); err != nil {
		return diag.Errorf("failed to set name during DCF Ruleset read: %s", err)
	}

	if err := d.Set("rules", policies); err != nil {
		return diag.Errorf("failed to set rules during DCF Ruleset read: %s", err)
	}

	d.SetId(policyList.UUID)

	return nil
}

func resourceAviatrixDCFRulesetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	policyList, err := marshalDCFRulesetInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Ruleset during update: %s", err)
	}

	err = client.UpdateDCFPolicyList(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to update DCF Ruleset: %s", err)
	}

	return nil
}

func resourceAviatrixDCFRulesetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	err := client.DeleteDCFPolicyList(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete DCF Ruleset: %v", err)
	}

	return nil
}
