package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	EgressPathDefault = "EGRESS_PATH_DEFAULT"
	EgressPathLocal   = "EGRESS_PATH_LOCAL"
)

var dcfRuleElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"DENY", "PERMIT", "DEEP_PACKET_INSPECTION_PERMIT", "INTRUSION_DETECTION_PERMIT"}, false),
			Description: "Action for the specified source and destination Smart Groups. " +
				"Must be one of INTRUSION_DETECTION_PERMIT, PERMIT or DENY.",
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
			Required:    true,
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
		"egress_path": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      EgressPathDefault,
			ValidateFunc: validation.StringInSlice([]string{EgressPathDefault, EgressPathLocal}, false),
			Description: "Egress path for this rule. Must be one of EGRESS_PATH_DEFAULT or EGRESS_PATH_LOCAL." +
				"EGRESS_PATH_DEFAULT routes traffic through the spoke's configured egress transit (FireNet, TGW, etc.). " +
				"EGRESS_PATH_LOCAL routes traffic out through the spoke gateway directly. " +
				"Example: `egress_path = \"EGRESS_PATH_LOCAL\"`.",
		},
	},
}

// dcfRuleSetHash normalizes the "protocol" field to uppercase before hashing
// so that case-only differences (e.g. "tcp" vs "TCP") don't cause Terraform to
// detect a spurious diff and trigger an unnecessary update.
func dcfRuleSetHash(v interface{}) int {
	raw, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}
	normalized := make(map[string]interface{}, len(raw))
	for k, val := range raw {
		normalized[k] = val
	}
	if protocol, ok := normalized["protocol"].(string); ok {
		normalized["protocol"] = strings.ToUpper(protocol)
	}
	return schema.HashResource(dcfRuleElem)(normalized)
}

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
				Elem:        dcfRuleElem,
				Set:         dcfRuleSetHash,
			},
		},
	}
}

func getUUIDSFromRules(rules []goaviatrix.DCFPolicy) (map[string]struct{}, error) {
	setUUID := make(map[string]struct{})
	for _, rule := range rules {
		uuid := rule.UUID
		if uuid != "" {
			setUUID[uuid] = struct{}{}
		}
	}
	return setUUID, nil
}

func buildNameToUUIDMapForUUIDNotInGivenSet(rules []goaviatrix.DCFPolicy, givenSet map[string]struct{}) (map[string]string, error) {
	nameToUUIDMap := make(map[string]string)
	for _, rule := range rules {
		// If the rule is not in the new set, then it might be an old rule that had attributes changed or a rule that was deleted, so we will add it to the map.
		if _, exists := givenSet[rule.UUID]; !exists {
			nameToUUIDMap[rule.Name] = rule.UUID
		}
	}
	return nameToUUIDMap, nil
}

func sortRules(rules []goaviatrix.DCFPolicy) {
	slices.SortFunc(rules, func(a goaviatrix.DCFPolicy, b goaviatrix.DCFPolicy) int {
		return a.Priority - b.Priority
	})
}

func marshalRulesList(rulesSet *schema.Set) ([]goaviatrix.DCFPolicy, error) {
	if rulesSet == nil {
		return []goaviatrix.DCFPolicy{}, nil
	}
	rules := []goaviatrix.DCFPolicy{}
	for _, rule := range rulesSet.List() {
		ruleMap, ok := rule.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("rules must be of type map[string]any")
		}
		rule, err := marshalPolicyInput(ruleMap)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

func marshalDCFRulesetInput(d *schema.ResourceData) (*goaviatrix.DCFPolicyList, error) {
	policyList := &goaviatrix.DCFPolicyList{}

	name := getString(d, "name")

	policyList.Name = name

	attachTo := getString(d, "attach_to")

	policyList.AttachTo = attachTo

	// Set diff in Terraform does not compare each object attributes in details, but instead it just compares the object hash
	// If the hash is different, the object is removed and a new object matching user input is added, so here oldPoliciesSet and newPoliciesSet are the old and new sets of rules, respectively.
	// The oldPoliciesSet contains a set of rules with UUIDs which match the user state before the update, and newPoliciesSet contains a set of rules along with new rules that had any attributes changed.
	// These new rules do not have UUIDs as it is a computed field, and if the newPoliciesSet is sent as is then the backend will generate new UUIDs for the rules which had attibutes changed.
	// This is not what we want, as we want to update the rules that had attributes changed, and not update the rule UUIDs.
	// To solve this, we will build a map of name -> UUID for the old rules, and then use this map to set the UUIDs for the new rules which match the name of the old rules.
	// This way, we can update the rules that had attributes changed, and not update the rule UUIDs as backend will not create new UUIDs if we pass in a UUID with our request.
	// The algorithm is as follows:
	// 1. Build a set of UUIDs for the new rules.
	// 2. Build a map of name => UUID for the old rules.
	// 3. Iterate over the new rules and set the UUID to the old rule UUID if the name exists in the mapped old rule name => UUID map and the new rule does not have a UUID.
	// 4. Send the new rules to the backend.

	// The drawback to this solution is that now name and UUID are tagged together, so if the user flips names of rules in an update then the UUIDs are also flipped accordingly.
	// This flip might result in log history having inconsistencies and hit counters for rules showing wrong values, and this is a case we will release note to the customer.

	oldRulesSchemaSet, newRulesSchemaSet := d.GetChange("rules")
	oldRulesSet, ok := oldRulesSchemaSet.(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("ruleset rules must be of type *schema.Set")
	}
	newRulesSet, ok := newRulesSchemaSet.(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("ruleset rules must be of type *schema.Set")
	}
	newRules, err := marshalRulesList(newRulesSet)
	if err != nil {
		return nil, err
	}
	oldRules, err := marshalRulesList(oldRulesSet)
	if err != nil {
		return nil, err
	}
	sortRules(newRules)

	// Build a set of UUIDs for the new rules.
	newUUIDs, err := getUUIDSFromRules(newRules)
	if err != nil {
		return nil, err
	}
	// While iterating over the old rules, we will build a map of name => UUID for the old rules.
	oldRuleNameToUUIDMap, err := buildNameToUUIDMapForUUIDNotInGivenSet(oldRules, newUUIDs)
	if err != nil {
		return nil, err
	}

	policyList.Policies = []goaviatrix.DCFPolicy{}
	for _, policy := range newRules {
		// If the rule name exists in the mapped old rule name => UUID map, and that the new rule does not have a UUID, then we will set the UUID to the old rule UUID.
		if policyUUID, exists := oldRuleNameToUUIDMap[policy.Name]; exists && policyUUID != "" {
			policy.UUID = policyUUID
			// delete the mapping so that we do not reuse it in case of a duplicate name
			delete(oldRuleNameToUUIDMap, policy.Name)
		}
		policyList.Policies = append(policyList.Policies, policy)
	}

	policyList.UUID = d.Id()

	return policyList, nil
}

//nolint:funlen,cyclop
func marshalPolicyInput(policyMap map[string]any) (*goaviatrix.DCFPolicy, error) {
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

	policy.EgressPath, ok = policyMap["egress_path"].(string)
	if !ok {
		return nil, fmt.Errorf("egress_path must be of type string")
	}

	if goaviatrix.MapContains(policyMap, "port_ranges") {
		if policy.Protocol == "ICMP" {
			return nil, fmt.Errorf("%q must not be set when %q is %q", "port_ranges", "protocol", "ICMP")
		}

		portRanges, ok := policyMap["port_ranges"].([]any)
		if !ok {
			return nil, fmt.Errorf("port_ranges must be of type []interface{}")
		}

		for _, portRangeInterface := range portRanges {
			portRangeMap, ok := portRangeInterface.(map[string]any)
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

func marshalSmartGroupsInput(policyMap map[string]any, key string) ([]string, error) {
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

func resourceAviatrixDCFRulesetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	policyList, err := marshalDCFRulesetInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Ruleset during create: %s", err)
	}

	var returnDiag diag.Diagnostics

	for _, policy := range policyList.Policies {
		if policy.Action == "DEEP_PACKET_INSPECTION_PERMIT" {
			returnDiag = diag.Errorf("DEEP_PACKET_INSPECTION_PERMIT will no longer be a valid Action value in the next major release. Use INTRUSION_DETECTION_PERMIT and DECRYPT_ALLOWED instead")
			returnDiag[0].Severity = diag.Warning
			break
		}
	}

	uuid, err := client.CreateDCFPolicyList(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to create DCF Ruleset: %s", err)
	}

	d.SetId(uuid)
	// todo: consider refactoring client.CreateDCFPolicyList or using the create result to populate the readDiags to prevent 2 API calls
	// ticket tracking this: https://aviatrix.atlassian.net/browse/AVX-76199
	readDiags := resourceAviatrixDCFRulesetRead(ctx, d, meta)
	return append(returnDiag, readDiags...)
}

//nolint:funlen,cyclop
func resourceAviatrixDCFRulesetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var policies []map[string]any
	for _, policy := range policyList.Policies {
		if policy.SystemResource {
			continue
		}
		p := make(map[string]any)
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
			p["protocol"] = strings.ToUpper(policy.Protocol)
		}
		p["flow_app_requirement"] = policy.FlowAppRequirement
		p["decrypt_policy"] = policy.DecryptPolicy

		if policy.Protocol != "ICMP" {
			var portRanges []map[string]any
			for _, portRange := range policy.PortRanges {
				portRangeMap := map[string]any{
					"hi": portRange.Hi,
					"lo": portRange.Lo,
				}
				portRanges = append(portRanges, portRangeMap)
			}
			p["port_ranges"] = portRanges
		}
		p["tls_profile"] = policy.TLSProfile
		p["egress_path"] = policy.EgressPath

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

func resourceAviatrixDCFRulesetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	policyList, err := marshalDCFRulesetInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Ruleset during update: %s", err)
	}

	err = client.UpdateDCFPolicyList(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to update DCF Ruleset: %s", err)
	}

	return resourceAviatrixDCFRulesetRead(ctx, d, meta)
}

func resourceAviatrixDCFRulesetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	err := client.DeleteDCFPolicyList(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete DCF Ruleset: %v", err)
	}

	return nil
}
