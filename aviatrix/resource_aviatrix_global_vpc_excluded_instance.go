package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixGlobalVpcExcludedInstance() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixGlobalVpcExcludedInstanceCreate,
		ReadWithoutTimeout:   resourceAviatrixGlobalVpcExcludedInstanceRead,
		UpdateWithoutTimeout: resourceAviatrixGlobalVpcExcludedInstanceUpdate,
		DeleteWithoutTimeout: resourceAviatrixGlobalVpcExcludedInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Aviatrix GCP access name.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance to be excluded for tagging.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of the instance.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the exclude list object.",
			},
		},
	}
}

func marshalGlobalVpcExcludedInstanceInput(d *schema.ResourceData) *goaviatrix.GlobalVpcExcludedInstance {
	globalVpcExcludedInstance := &goaviatrix.GlobalVpcExcludedInstance{
		AccountName:  d.Get("account_name").(string),
		InstanceName: d.Get("instance_name").(string),
		Region:       d.Get("region").(string),
	}

	return globalVpcExcludedInstance
}

func resourceAviatrixGlobalVpcExcludedInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	globalVpcExcludedInstance := marshalGlobalVpcExcludedInstanceInput(d)

	flag := false
	defer resourceAviatrixGlobalVpcExcludedInstanceReadIfRequired(ctx, d, meta, &flag)

	uuid, err := client.CreateGlobalVpcExcludedInstance(ctx, globalVpcExcludedInstance)
	if err != nil {
		return diag.Errorf("failed to create global vpc excluded instance: %s", err)
	}

	d.SetId(uuid)
	return resourceAviatrixGlobalVpcExcludedInstanceReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixGlobalVpcExcludedInstanceReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixGlobalVpcExcludedInstanceRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixGlobalVpcExcludedInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Set("uuid", uuid)

	globalVpcExcludedInstance, err := client.GetGlobalVpcExcludedInstance(ctx, uuid)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read global vpc excluded instance: %s", err)
	}

	d.Set("account_name", globalVpcExcludedInstance.AccountName)
	d.Set("instance_name", globalVpcExcludedInstance.InstanceName)
	d.Set("region", globalVpcExcludedInstance.Region)

	return nil
}

func resourceAviatrixGlobalVpcExcludedInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	d.Partial(true)
	if d.HasChanges("account_name", "instance_name", "region") {
		globalVpcExcludedInstance := marshalGlobalVpcExcludedInstanceInput(d)

		err := client.UpdateGlobalVpcExcludedInstance(ctx, globalVpcExcludedInstance, uuid)
		if err != nil {
			return diag.Errorf("failed to update global vpc excluded instance: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixGlobalVpcExcludedInstanceRead(ctx, d, meta)
}

func resourceAviatrixGlobalVpcExcludedInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	uuid := d.Id()
	err := client.DeleteGlobalVpcExcludedInstance(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete global vpc excluded instance: %v", err)
	}

	return nil
}
