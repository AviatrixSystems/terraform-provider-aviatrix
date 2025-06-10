package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerBgpCommunitiesGlobalConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerBgpCommunitiesGlobalConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerBgpCommunitiesGlobalConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerBgpCommunitiesGlobalConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerBgpCommunitiesGlobalConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bgp_communities_global": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "BGP communities global configuration",
			},
		},
	}
}

func resourceAviatrixControllerBgpCommunitiesGlobalConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	bgpCommunities, ok := d.Get("bgp_communities_global").(bool)
	if !ok {
		return diag.Errorf("failed to assert bgp_communities_global as bool")
	}
	if bgpCommunities {
		err := client.EnableControllerBgpCommunitiesGlobal(ctx)
		if err != nil {
			return diag.Errorf("failed to enable controller BGP communities global config: %v", err)
		}
	} else {
		err := client.DisableControllerBgpCommunitiesGlobal(ctx)
		if err != nil {
			return diag.Errorf("failed to disable controller BGP communities global config: %v", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpCommunitiesGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpCommunitiesGlobalConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	commGlobal, err := client.GetControllerBgpCommunitiesGlobal(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get controller BGP communities global config: %v", err)
	}

	err = d.Set("bgp_communities_global", commGlobal)
	if err != nil {
		return diag.Errorf("failed to set bgp_communities_global: %v", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerBgpCommunitiesGlobalConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	if d.HasChange("bgp_communities_global") {
		bgpCommunities, ok := d.Get("bgp_communities_global").(bool)
		if !ok {
			return diag.Errorf("failed to assert bgp_communities_global as bool")
		}
		if bgpCommunities {
			err := client.EnableControllerBgpCommunitiesGlobal(ctx)
			if err != nil {
				return diag.Errorf("failed to enable controller BGP communities global config: %v", err)
			}
		} else {
			err := client.DisableControllerBgpCommunitiesGlobal(ctx)
			if err != nil {
				return diag.Errorf("failed to disable controller BGP communities global config: %v", err)
			}
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerBgpCommunitiesGlobalConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerBgpCommunitiesGlobalConfigDelete(ctx context.Context, _ *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	err := client.DisableControllerBgpCommunitiesGlobal(ctx)
	if err != nil {
		return diag.Errorf("failed to delete controller BGP communities global config: %v", err)
	}

	return nil
}
