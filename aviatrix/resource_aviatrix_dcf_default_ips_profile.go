package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const defaultIPSProfileUUID = "defa11a1-0000-4000-8000-000000000000"

func resourceAviatrixDCFDefaultIpsProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFDefaultIpsProfileCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFDefaultIpsProfileRead,
		UpdateWithoutTimeout: resourceAviatrixDCFDefaultIpsProfileUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFDefaultIpsProfileDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"default_ips_profile": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "List of default IPS profile UUIDs. Only one profile is supported at this time.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAviatrixDCFDefaultIpsProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	profiles := expandStringList(getSet(d, "default_ips_profile").List())

	_, err := client.SetDefaultIpsProfile(ctx, profiles)
	if err != nil {
		return diag.Errorf("failed to set default IPS profile: %v", err)
	}

	d.SetId("dcf_default_ips_profile")
	return resourceAviatrixDCFDefaultIpsProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFDefaultIpsProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	if d.Id() != "dcf_default_ips_profile" {
		return diag.Errorf("ID: %s does not match expected ID \"dcf_default_ips_profile\": please provide correct ID for importing", d.Id())
	}

	profile, err := client.GetDefaultIpsProfile(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read default IPS profile: %v", err)
	}

	mustSet(d, "default_ips_profile", profile.DefaultIpsProfile)

	d.SetId("dcf_default_ips_profile")
	return nil
}

func resourceAviatrixDCFDefaultIpsProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	profiles := expandStringList(getSet(d, "default_ips_profile").List())

	_, err := client.SetDefaultIpsProfile(ctx, profiles)
	if err != nil {
		return diag.Errorf("failed to update default IPS profile: %v", err)
	}

	return resourceAviatrixDCFDefaultIpsProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFDefaultIpsProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	// Need to reset the default IPS profile to the system default, there must be a default profile
	_, err := client.SetDefaultIpsProfile(ctx, []string{defaultIPSProfileUUID})
	if err != nil {
		return diag.Errorf("failed to clear default IPS profile: %v", err)
	}

	d.SetId("")
	return nil
}
