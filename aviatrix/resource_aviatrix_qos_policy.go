package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixQosPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixQosPolicyCreate,
		ReadWithoutTimeout:   resourceAviatrixQosPolicyRead,
		UpdateWithoutTimeout: resourceAviatrixQosPolicyUpdate,
		DeleteWithoutTimeout: resourceAviatrixQosPolicyDelete,
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

func marshalQosPolicyInput(d *schema.ResourceData) *goaviatrix.QosPolicyList {
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

func resourceAviatrixQosPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	qosPolicyList := marshalQosPolicyInput(d)

	flag := false
	defer resourceAviatrixQosPolicyReadIfRequired(ctx, d, meta, &flag)

	err := client.UpdateQosPolicy(ctx, qosPolicyList)
	if err != nil {
		return diag.Errorf("failed to create qos policy: %s", err)
	}

	d.SetId("qos_policy")
	return resourceAviatrixQosPolicyReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixQosPolicyReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixQosPolicyRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixQosPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	qosPolicyResp, err := client.GetQosPolicy(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read qos policy: %s", err)
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

func resourceAviatrixQosPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChanges("policies") {
		qosPolicyList := marshalQosPolicyInput(d)

		err := client.UpdateQosPolicy(ctx, qosPolicyList)
		if err != nil {
			return diag.Errorf("failed to update qos policy: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixQosPolicyRead(ctx, d, meta)
}

func resourceAviatrixQosPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteQosPolicy(ctx)
	if err != nil {
		return diag.Errorf("failed to delete qos policy: %v", err)
	}

	return nil
}
