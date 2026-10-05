package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTrafficClassifier() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixTrafficClassifierCreate,
		ReadWithoutTimeout:   resourceAviatrixTrafficClassifierRead,
		DeleteWithoutTimeout: resourceAviatrixTrafficClassifierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policies": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of traffic classifier policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Policy name.",
						},
						"source_smart_group_uuids": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							Description: "List of source smart group UUIDs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_smart_group_uuids": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							Description: "List of destination smart group UUIDs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"port_ranges": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Description: "Port ranges.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"low": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "Low port range.",
									},
									"high": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "High port range.",
									},
								},
							},
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Protocol.",
						},
						"link_hierarchy_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Link hierarchy UUID.",
						},
						"sla_class_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "SLA class UUID.",
						},
						"enable_logging": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Enable logging.",
						},
						"route_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Route type.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic classifier policy UUID.",
						},
					},
				},
			},
		},
	}
}

func marshalTrafficClassifierInput(d *schema.ResourceData) *goaviatrix.PolicyList {
	var policyList goaviatrix.PolicyList

	policies := d.Get("policies").([]interface{})
	for _, v0 := range policies {
		v1 := v0.(map[string]interface{})

		if2 := goaviatrix.TCPolicy{
			Name:          v1["name"].(string),
			Protocol:      v1["protocol"].(string),
			LinkHierarchy: v1["link_hierarchy_uuid"].(string),
			SlaClass:      v1["sla_class_uuid"].(string),
			Logging:       v1["enable_logging"].(bool),
			RouteType:     v1["route_type"].(string),
		}

		for _, ss := range v1["source_smart_group_uuids"].([]interface{}) {
			if2.SrcSgs = append(if2.SrcSgs, ss.(string))
		}

		for _, ds := range v1["destination_smart_group_uuids"].([]interface{}) {
			if2.DstSgs = append(if2.DstSgs, ds.(string))
		}

		for _, v2 := range v1["port_ranges"].(*schema.Set).List() {
			v3 := v2.(map[string]interface{})

			pr := goaviatrix.PortRange{
				Lo: v3["low"].(int),
				Hi: v3["high"].(int),
			}

			if2.PortRanges = append(if2.PortRanges, pr)
		}

		policyList.Policies = append(policyList.Policies, if2)
	}

	return &policyList
}

func resourceAviatrixTrafficClassifierCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	policyList := marshalTrafficClassifierInput(d)

	flag := false
	defer resourceAviatrixTrafficClassifierReadIfRequired(ctx, d, meta, &flag)

	err := client.CreateTrafficClassifier(ctx, policyList)
	if err != nil {
		return diag.Errorf("failed to create traffic classifier: %s", err)
	}

	d.SetId("traffic_classifier_policies")
	return resourceAviatrixTrafficClassifierReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixTrafficClassifierReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTrafficClassifierRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixTrafficClassifierRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	tcResp, err := client.GetTrafficClassifier(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read traffic classifier: %s", err)
	}

	var policies []map[string]interface{}

	for _, policyList := range *tcResp {
		for _, policy := range policyList.Policies {
			p := make(map[string]interface{})
			p["uuid"] = policy.UUID
			p["name"] = policy.Name
			p["source_smart_group_uuids"] = policy.SrcSgs
			p["destination_smart_group_uuids"] = policy.DstSgs
			p["link_hierarchy_uuid"] = policy.LinkHierarchy
			p["sla_class_uuid"] = policy.SlaClass
			p["enable_logging"] = policy.Logging
			p["route_type"] = policy.RouteType

			if policy.Protocol != "PROTOCOL_UNSPECIFIED" {
				p["protocol"] = policy.Protocol
			}

			var portRanges []map[string]interface{}
			for _, pr := range policy.PortRanges {
				p1 := make(map[string]interface{})
				p1["low"] = pr.Lo
				p1["high"] = pr.Hi
				portRanges = append(portRanges, p1)
			}
			p["port_ranges"] = portRanges

			policies = append(policies, p)
		}
	}

	if err := d.Set("policies", policies); err != nil {
		return diag.Errorf("failed to set policies: %s", err)
	}

	return nil
}

func resourceAviatrixTrafficClassifierDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteTrafficClassifier(ctx)
	if err != nil {
		return diag.Errorf("failed to delete traffic classifier: %v", err)
	}

	return nil
}
