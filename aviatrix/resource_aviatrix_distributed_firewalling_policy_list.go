package aviatrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixDistributedFirewallingPolicyList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingPolicyListCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingPolicyListRead,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingPolicyListUpdate,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingPolicyListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of distributed-firewalling policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the policy.",
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DENY", "PERMIT", "DEEP_PACKET_INSPECTION_PERMIT", "INTRUSION_DETECTION_PERMIT"}, false),
							Description: "Action for the specified source and destination Smart Groups." +
								"Must be one of PERMIT or DENY.",
						},
						"dst_smart_groups": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of destination Smart Group UUIDs for the policy.",
						},
						"src_smart_groups": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of source Smart Group UUIDs for the policy.",
						},
						"web_groups": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of Web Group UUIDs for the policy.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP", "ANY"}, true),
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return strings.EqualFold(old, new)
							},
							Description: "Protocol for the policy to filter.",
						},
						"flow_app_requirement": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "APP_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"APP_UNSPECIFIED", "TLS_REQUIRED", "NOT_TLS_REQUIRED"}, false),
							Description:  "Flow application requirement for the policy.",
						},
						"decrypt_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "DECRYPT_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"DECRYPT_UNSPECIFIED", "DECRYPT_ALLOWED", "DECRYPT_NOT_ALLOWED"}, false),
							Description:  "Decryption options for the policy.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Priority level of the policy",
						},
						"logging": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable logging for the policy.",
						},
						"watch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable watch mode for the policy.",
						},
						"exclude_sg_orchestration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If this flag is set to true, this policy will be ignored for SG orchestration.",
						},
						"port_ranges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of port ranges for the policy.",
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
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the policy.",
						},
					},
				},
			},
		},
	}
}

func marshalDistributedFirewallingPolicyListInput(d *schema.ResourceData) (*goaviatrix.DistributedFirewallingPolicyList, error) {
	policyList := &goaviatrix.DistributedFirewallingPolicyList{}

	policies := d.Get("policies").([]interface{})
	for _, policyInterface := range policies {
		policy := policyInterface.(map[string]interface{})

		distributedFirewallingPolicy := &goaviatrix.DistributedFirewallingPolicy{
			Name:                   policy["name"].(string),
			Action:                 policy["action"].(string),
			Priority:               policy["priority"].(int),
			FlowAppRequirement:     policy["flow_app_requirement"].(string),
			DecryptPolicy:          policy["decrypt_policy"].(string),
			ExcludeSgOrchestration: policy["exclude_sg_orchestration"].(bool),
		}

		protocol := strings.ToUpper(policy["protocol"].(string))
		if protocol == "ANY" {
			distributedFirewallingPolicy.Protocol = "PROTOCOL_UNSPECIFIED"
		} else {
			distributedFirewallingPolicy.Protocol = protocol
		}

		for _, smartGroup := range policy["src_smart_groups"].(*schema.Set).List() {
			distributedFirewallingPolicy.SrcSmartGroups = append(distributedFirewallingPolicy.SrcSmartGroups, smartGroup.(string))
		}

		for _, smartGroup := range policy["dst_smart_groups"].(*schema.Set).List() {
			distributedFirewallingPolicy.DstSmartGroups = append(distributedFirewallingPolicy.DstSmartGroups, smartGroup.(string))
		}

		for _, webGroup := range policy["web_groups"].(*schema.Set).List() {
			distributedFirewallingPolicy.WebGroups = append(distributedFirewallingPolicy.WebGroups, webGroup.(string))
		}

		if logging, loggingOk := policy["logging"]; loggingOk {
			distributedFirewallingPolicy.Logging = logging.(bool)
		}

		if watch, watchOk := policy["watch"]; watchOk {
			distributedFirewallingPolicy.Watch = watch.(bool)
		}

		if mapContains(policy, "port_ranges") {
			if distributedFirewallingPolicy.Protocol == "ICMP" {
				return nil, fmt.Errorf("%q must not be set when %q is %q", "port_ranges", "protocol", "ICMP")
			}
			for _, portRangeInterface := range policy["port_ranges"].([]interface{}) {
				portRangeMap := portRangeInterface.(map[string]interface{})
				portRange := &goaviatrix.DistributedFirewallingPortRange{
					Lo: portRangeMap["lo"].(int),
				}

				if hi, hiOk := portRangeMap["hi"]; hiOk {
					portRange.Hi = hi.(int)
				}

				distributedFirewallingPolicy.PortRanges = append(distributedFirewallingPolicy.PortRanges, *portRange)
			}
		}

		if uuid, uuidOk := policy["uuid"]; uuidOk {
			distributedFirewallingPolicy.UUID = uuid.(string)
		}

		policyList.Policies = append(policyList.Policies, *distributedFirewallingPolicy)
	}

	return policyList, nil
}

func resourceAviatrixDistributedFirewallingPolicyListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	policyList, err := marshalDistributedFirewallingPolicyListInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for Distributed-firewalling Policy during create: %s\n", err)
	}

	flag := false
	defer resourceAviatrixDistributedFirewallingPolicyListReadIfRequired(ctx, d, meta, &flag)

	err = client.CreateDistributedFirewallingPolicyList(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to create Distributed-firewalling Policy List: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingPolicyListReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixDistributedFirewallingPolicyListReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixDistributedFirewallingPolicyListRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixDistributedFirewallingPolicyListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	policyList, err := client.GetDistributedFirewallingPolicyList(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Distributed-firewalling Policy List: %s", err)
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

		policies = append(policies, p)
	}

	if err := d.Set("policies", policies); err != nil {
		return diag.Errorf("failed to set policies during Distributed-firewalling Policy List read: %s\n", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingPolicyListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChange("policies") {
		policyList, err := marshalDistributedFirewallingPolicyListInput(d)
		if err != nil {
			return diag.Errorf("invalid inputs for Distributed-firewalling Policy during update: %s\n", err)
		}
		err = client.UpdateDistributedFirewallingPolicyList(ctx, policyList)
		if err != nil {
			return diag.Errorf("failed to update Distributed-firewalling policies: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixDistributedFirewallingPolicyListRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingPolicyListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteDistributedFirewallingPolicyList(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling Policy List: %v", err)
	}

	return nil
}
