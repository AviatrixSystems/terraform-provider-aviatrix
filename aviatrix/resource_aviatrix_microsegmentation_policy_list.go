package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixMicrosegmentationPolicyList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixMicrosegmentationPolicyListCreate,
		ReadWithoutTimeout:   resourceAviatrixMicrosegmentationPolicyListRead,
		UpdateWithoutTimeout: resourceAviatrixMicrosegmentationPolicyListUpdate,
		DeleteWithoutTimeout: resourceAviatrixMicrosegmentationPolicyListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of micro-segmentation policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of the policy.",
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"PERMIT", "DENY"}, false),
							Description: "Action for the specified source and destination App Domains." +
								"Must be one of PERMIT or DENY.",
						},
						"logging": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable logging for the policy.",
						},
						"dst_app_domains": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of destination App Domain UUIDs for the policy.",
						},
						"src_app_domains": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Set of source App Domain UUIDs for the policy.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
							Description:  "Protocol for the policy to filter.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Priority level of the policy",
						},
						"watch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable watch mode for the policy.",
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
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Upper bound of port range.",
									},
								},
							},
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

func marshalMicrosegmentationPolicyListInput(d *schema.ResourceData) *goaviatrix.MicrosegmentationPolicyList {
	policyList := &goaviatrix.MicrosegmentationPolicyList{}

	policies := d.Get("policies").([]interface{})
	for _, policyInterface := range policies {
		policy := policyInterface.(map[string]interface{})

		microsegmentationPolicy := &goaviatrix.MicrosegmentationPolicy{
			Name:     policy["name"].(string),
			Action:   policy["action"].(string),
			Priority: policy["priority"].(int),
			Protocol: policy["protocol"].(string),
		}

		if logging, loggingOk := policy["logging"]; loggingOk {
			microsegmentationPolicy.Logging = logging.(bool)
		}

		for _, appDomain := range policy["src_app_domains"].(*schema.Set).List() {
			microsegmentationPolicy.SrcAppDomains = append(microsegmentationPolicy.SrcAppDomains, appDomain.(string))
		}

		for _, appDomain := range policy["dst_app_domains"].(*schema.Set).List() {
			microsegmentationPolicy.DstAppDomains = append(microsegmentationPolicy.DstAppDomains, appDomain.(string))
		}

		if watch, watchOk := policy["watch"]; watchOk {
			microsegmentationPolicy.Watch = watch.(bool)
		}

		for _, portRangeInterface := range policy["port_ranges"].([]interface{}) {
			portRangeMap := portRangeInterface.(map[string]interface{})
			portRange := &goaviatrix.MicrosegmentationPortRange{
				Lo: portRangeMap["lo"].(int),
			}

			if hi, hiOk := portRangeMap["hi"]; hiOk {
				portRange.Hi = hi.(int)
			}
			microsegmentationPolicy.PortRanges = append(microsegmentationPolicy.PortRanges, *portRange)
		}

		policyList.Policies = append(policyList.Policies, *microsegmentationPolicy)
	}

	return policyList
}

func resourceAviatrixMicrosegmentationPolicyListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	policyList := marshalMicrosegmentationPolicyListInput(d)

	flag := false
	defer resourceAviatrixMicrosegmentationPolicyListReadIfRequired(ctx, d, meta, &flag)

	err := client.CreateMicrosegmentationPolicyList(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to create Micro-segmentation Policy List: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixMicrosegmentationPolicyListReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixMicrosegmentationPolicyListReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixMicrosegmentationPolicyListRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixMicrosegmentationPolicyListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	policyList, err := client.GetMicrosegmentationPolicyList(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Micro-segmentation Policy List: %s", err)
	}

	var policies []map[string]interface{}
	for _, policy := range policyList.Policies {
		p := make(map[string]interface{})
		p["name"] = policy.Name
		p["action"] = policy.Action
		p["logging"] = policy.Logging
		p["priority"] = policy.Priority
		p["protocol"] = policy.Protocol
		p["src_app_domains"] = policy.SrcAppDomains
		p["dst_app_domains"] = policy.DstAppDomains
		p["watch"] = policy.Watch
		p["uuid"] = policy.UUID

		var portRanges []map[string]interface{}
		for _, portRange := range policy.PortRanges {
			portRangeMap := map[string]interface{}{
				"hi": portRange.Hi,
				"lo": portRange.Lo,
			}
			portRanges = append(portRanges, portRangeMap)
		}
		p["port_ranges"] = portRanges

		policies = append(policies, p)
	}

	if err := d.Set("policies", policies); err != nil {
		return diag.Errorf("failed to set policies during Micro-segmentation Policy List read: %s\n", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixMicrosegmentationPolicyListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChange("policies") {
		policyList := marshalMicrosegmentationPolicyListInput(d)
		err := client.UpdateMicrosegmentationPolicyList(ctx, policyList)
		if err != nil {
			return diag.Errorf("failed to update Micro-segmentation policies: %s", err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceAviatrixMicrosegmentationPolicyListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteMicrosegmentationPolicyList(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Micro-segmentation Policy List: %v", err)
	}

	return nil
}
