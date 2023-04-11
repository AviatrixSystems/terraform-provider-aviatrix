package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixQosClass() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixQosClassCreate,
		ReadWithoutTimeout:   resourceAviatrixQosClassRead,
		UpdateWithoutTimeout: resourceAviatrixQosClassUpdate,
		DeleteWithoutTimeout: resourceAviatrixQosClassDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "QoS class name.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "QoS class priority.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "QoS class UUID.",
			},
		},
	}
}

func marshalQosClassInput(d *schema.ResourceData) *goaviatrix.QosClass {
	qosClass := &goaviatrix.QosClass{
		Name:     d.Get("name").(string),
		Priority: d.Get("priority").(int),
	}

	return qosClass
}

func resourceAviatrixQosClassCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	qosClass := marshalQosClassInput(d)

	flag := false
	defer resourceAviatrixQosClassReadIfRequired(ctx, d, meta, &flag)

	uuid, err := client.CreateQosClass(ctx, qosClass)
	if err != nil {
		return diag.Errorf("failed to create qos class: %s", err)
	}

	d.SetId(uuid)
	return resourceAviatrixQosClassReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixQosClassReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixQosClassRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixQosClassRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Set("uuid", uuid)

	qosClass, err := client.GetQosClass(ctx, uuid)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read qos class %s", err)
	}

	d.Set("name", qosClass.Name)
	d.Set("priority", qosClass.Priority)

	return nil
}

func resourceAviatrixQosClassUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("name", "priority") {
		qosClass := marshalQosClassInput(d)

		err := client.UpdateQosClass(ctx, qosClass, uuid)
		if err != nil {
			return diag.Errorf("failed to update qos class: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixQosClassRead(ctx, d, meta)
}

func resourceAviatrixQosClassDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	err := client.DeleteQosClass(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete qos class: %v", err)
	}

	return nil
}
