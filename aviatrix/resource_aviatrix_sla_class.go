package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSLAClass() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixSLAClassCreate,
		ReadWithoutTimeout:   resourceAviatrixSLAClassRead,
		UpdateWithoutTimeout: resourceAviatrixSLAClassUpdate,
		DeleteWithoutTimeout: resourceAviatrixSLAClassDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of SLA class.",
			},
			"latency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Latency of sla class in ms.",
			},
			"jitter": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Jitter of sla class in ms.",
			},
			"packet_drop_rate": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Packet drop rate of sla class.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of SLA class.",
			},
		},
	}
}

func marshalSLAClassInput(d *schema.ResourceData) *goaviatrix.SLAClass {
	slaClass := &goaviatrix.SLAClass{
		Name:           getString(d, "name"),
		Latency:        getInt(d, "latency"),
		Jitter:         getInt(d, "jitter"),
		PacketDropRate: getFloat64(d, "packet_drop_rate"),
	}

	return slaClass
}

func resourceAviatrixSLAClassCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	slaClass := marshalSLAClassInput(d)

	flag := false
	defer resourceAviatrixSLAClassReadIfRequired(ctx, d, meta, &flag)

	uuid, err := client.CreateSLAClass(ctx, slaClass)
	if err != nil {
		return diag.Errorf("failed to create sla class: %s", err)
	}

	d.SetId(uuid)
	return resourceAviatrixSLAClassReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixSLAClassReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSLAClassRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixSLAClassRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()
	mustSet(d, "uuid", uuid)

	slaClass, err := client.GetSLAClass(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read sla class %s", err)
	}
	mustSet(d, "name", slaClass.Name)
	mustSet(d, "latency", slaClass.Latency)
	mustSet(d, "jitter", slaClass.Jitter)
	mustSet(d, "packet_drop_rate", slaClass.PacketDropRate)

	return nil
}

func resourceAviatrixSLAClassUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("name", "latency", "jitter", "package_drop_rate") {
		slaClass := marshalSLAClassInput(d)

		err := client.UpdateSLAClass(ctx, slaClass, uuid)
		if err != nil {
			return diag.Errorf("failed to update sla class: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixSLAClassRead(ctx, d, meta)
}

func resourceAviatrixSLAClassDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()
	err := client.DeleteSLAClass(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete sla: %v", err)
	}

	return nil
}
