package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
	client := mustClient(meta)

	autoCloud := getBool(d, "auto_cloud_enabled")

	if autoCloud {
		commPrefix := getInt(d, "community_prefix")

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
	client := mustClient(meta)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	commPrefix, err := client.GetControllerBgpCommunitiesAutoCloud(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get controller BGP communities auto cloud config: %v", err)
	}

	err = d.Set("community_prefix", commPrefix)
	if err != nil {
		return diag.Errorf("failed to set community prefix: %v", err)
	}
	if commPrefix > 0 {
		err = d.Set("auto_cloud_enabled", true)
		if err != nil {
			return diag.Errorf("failed to set auto cloud enabled: %v", err)
		}
	} else {
		err = d.Set("auto_cloud_enabled", false)
		if err != nil {
			return diag.Errorf("failed to set auto cloud enabled: %v", err)
		}
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.HasChange("auto_cloud_enabled") || d.HasChange("community_prefix") {
		autoCloud := getBool(d, "auto_cloud_enabled")

		if autoCloud {
			commPrefix := getInt(d, "community_prefix")

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

func resourceAviatrixControllerBgpCommunitiesAutoCloudConfigDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DisableControllerBgpCommunitiesAutoCloud(ctx)
	if err != nil {
		return diag.Errorf("failed to disable controller BGP communities auto cloud config: %v", err)
	}

	return nil
}
