package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
		Name:     getString(d, "name"),
		Priority: getInt(d, "priority"),
	}

	return qosClass
}

func resourceAviatrixQosClassCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

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
	client := mustClient(meta)

	uuid := d.Id()
	mustSet(d, "uuid", uuid)

	qosClass, err := client.GetQosClass(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read qos class %s", err)
	}
	mustSet(d, "name", qosClass.Name)
	mustSet(d, "priority", qosClass.Priority)

	return nil
}

func resourceAviatrixQosClassUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

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
	client := mustClient(meta)

	uuid := d.Id()
	err := client.DeleteQosClass(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete qos class: %v", err)
	}

	return nil
}
