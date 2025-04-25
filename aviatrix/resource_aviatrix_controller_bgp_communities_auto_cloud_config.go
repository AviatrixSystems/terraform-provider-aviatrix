package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerBgpCommunitiesAutoCloudConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerBgpCommunitiesAutoCloudConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerBgpCommunitiesAutoCloudConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerBgpCommunitiesAutoCloudConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"auto_cloud_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "BGP communities auto cloud configuration",
			},
			"community_prefix": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Community prefix for auto cloud BGP communities",
			},
		},
	}
}

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	autoCloud := d.Get("auto_cloud_enabled").(bool)
	if autoCloud {
		commPrefix := d.Get("community_prefix").(int)
		err := client.SetControllerBgpCommunitiesAutoCloud(ctx, commPrefix)
		if err != nil {
			return diag.Errorf("failed to enable controller BGP communities auto cloud config: %v", err)
		}
	} else {
		err := client.DisableControllerBgpCommunitiesAutoCloud(ctx)
		if err != nil {
			return diag.Errorf("failed to disable controller BGP communities auto cloud config: %v", err)
		}
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpCommunitiesAutoCloudConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	commPrefix, err := client.GetControllerBgpCommunitiesAutoCloud(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get controller BGP communities auto cloud config: %v", err)
	}

	d.Set("community_prefix", commPrefix)
	if commPrefix > 0 {
		d.Set("auto_cloud_enabled", true)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("community_prefix") {
		autoCloud := d.Get("auto_cloud_enabled").(bool)
		if autoCloud {
			commPrefix := d.Get("community_prefix").(int)
			err := client.SetControllerBgpCommunitiesAutoCloud(ctx, commPrefix)
			if err != nil {
				return diag.Errorf("failed to enable controller BGP communities auto cloud config: %v", err)
			}
		} else {
			err := client.DisableControllerBgpCommunitiesAutoCloud(ctx)
			if err != nil {
				return diag.Errorf("failed to disable controller BGP communities auto cloud config: %v", err)
			}
		}
	}

	return nil
}

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableControllerBgpCommunitiesAutoCloud(ctx)
	if err != nil {
		return diag.Errorf("failed to disable controller BGP communities auto cloud config: %v", err)
	}

	return nil
}
