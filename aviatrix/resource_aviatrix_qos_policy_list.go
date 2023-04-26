package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixQosPolicyList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixQosPolicyListCreate,
		ReadWithoutTimeout:   resourceAviatrixQosPolicyListRead,
		UpdateWithoutTimeout: resourceAviatrixQosPolicyListUpdate,
		DeleteWithoutTimeout: resourceAviatrixQosPolicyListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of QoS policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "QoS policy name.",
						},
						"dscp_values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of DSCP values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"qos_class_uuid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "QoS class UUID.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "QoS policy UUID.",
						},
					},
				},
			},
		},
	}
}

func marshalQosPolicyListInput(d *schema.ResourceData) *goaviatrix.QosPolicyList {
	var qosPolicyList goaviatrix.QosPolicyList

	policies := d.Get("policies").([]interface{})
	for _, v0 := range policies {
		v1 := v0.(map[string]interface{})

		v2 := goaviatrix.QosPolicy{
			Name:         v1["name"].(string),
			QosClassUuid: v1["qos_class_uuid"].(string),
		}

		for _, ss := range v1["dscp_values"].([]interface{}) {
			v2.DscpValues = append(v2.DscpValues, ss.(string))
		}

		qosPolicyList.Policies = append(qosPolicyList.Policies, v2)
	}

	return &qosPolicyList
}

func resourceAviatrixQosPolicyListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	qosPolicyList := marshalQosPolicyListInput(d)

	flag := false
	defer resourceAviatrixQosPolicyListReadIfRequired(ctx, d, meta, &flag)

	err := client.UpdateQosPolicyList(ctx, qosPolicyList)
	if err != nil {
		return diag.Errorf("failed to create qos policy list: %s", err)
	}

	d.SetId("qos_policy_list")
	return resourceAviatrixQosPolicyListReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixQosPolicyListReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixQosPolicyListRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixQosPolicyListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	qosPolicyResp, err := client.GetQosPolicyList(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read qos policy list: %s", err)
	}

	var policies []map[string]interface{}

	for _, policyList := range *qosPolicyResp {
		for _, policy := range policyList.Policies {
			p := make(map[string]interface{})
			p["uuid"] = policy.UUID
			p["name"] = policy.Name
			p["dscp_values"] = policy.DscpValues
			p["qos_class_uuid"] = policy.QosClassUuid

			policies = append(policies, p)
		}
	}

	if err := d.Set("policies", policies); err != nil {
		return diag.Errorf("failed to set policies: %s", err)
	}

	return nil
}

func resourceAviatrixQosPolicyListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChanges("policies") {
		qosPolicyList := marshalQosPolicyListInput(d)

		err := client.UpdateQosPolicyList(ctx, qosPolicyList)
		if err != nil {
			return diag.Errorf("failed to update qos policy list: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixQosPolicyListRead(ctx, d, meta)
}

func resourceAviatrixQosPolicyListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteQosPolicyList(ctx)
	if err != nil {
		return diag.Errorf("failed to delete qos policy list: %v", err)
	}

	return nil
}
