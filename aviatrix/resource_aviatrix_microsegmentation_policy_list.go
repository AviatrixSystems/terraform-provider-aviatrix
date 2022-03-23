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
				Type:     schema.TypeList,
				Required: true,
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
							Description:  "Microsegmentation action. Must be one of PERMIT or DENY.",
						},
						"dst_app_domains": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of destination App Domain UUIDs for this policy.",
						},
						"src_app_domains": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of source App Domain UUIDs for this policy.",
						},
						"port_ranges": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hi": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Upper bound of port range.",
									},
									"lo": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Lower bound of port range.",
									},
								},
							},
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Priority level of this policy",
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

		for _, appDomain := range policy["src_app_domains"].([]interface{}) {
			microsegmentationPolicy.SrcAppDomains = append(microsegmentationPolicy.SrcAppDomains, appDomain.(string))
		}

		for _, appDomain := range policy["dst_app_domains"].([]interface{}) {
			microsegmentationPolicy.DstAppDomains = append(microsegmentationPolicy.DstAppDomains, appDomain.(string))
		}

		for _, portRangeInterface := range policy["port_ranges"].([]interface{}) {
			portRangeMap := portRangeInterface.(map[string]interface{})
			portRange := &goaviatrix.MicrosegmentationPortRange{
				Hi: portRangeMap["hi"].(int),
				Lo: portRangeMap["lo"].(int),
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
		return diag.Errorf("failed to create Microsegmentation Policy List: %s", err)
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
		return diag.Errorf("failed to read Microsegmentation Policy List: %s", err)
	}

	var policies []map[string]interface{}
	for _, policy := range policyList.Policies {
		p := make(map[string]interface{})
		p["name"] = policy.Name
		p["action"] = policy.Action
		p["priority"] = policy.Priority
		p["protocol"] = policy.Protocol
		p["src_app_domains"] = policy.SrcAppDomains
		p["dst_app_domains"] = policy.DstAppDomains

		var portRanges []map[string]interface{}
		for _, portRange := range policy.PortRanges {
			portRangeMap := map[string]interface{}{
				"hi": portRange.Hi,
				"lo": portRange.Lo,
			}
			portRanges = append(portRanges, portRangeMap)
		}

		policies = append(policies, p)
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
			return diag.Errorf("failed to update Microsegmentation Policies: %s", err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceAviatrixMicrosegmentationPolicyListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteMicrosegmentationPolicyList(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Microsegmentation Policy List: %v", err)
	}

	return nil
}
