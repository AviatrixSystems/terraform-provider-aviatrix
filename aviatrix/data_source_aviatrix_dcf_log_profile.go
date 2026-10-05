package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixDcfLogProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDcfLogProfileRead,
		Schema: map[string]*schema.Schema{
			"profile_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Log Profile.",
			},
			"profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the Log Profile which can be referenced in a DCF rule.",
			},
			"session_end": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Tells us if the logging of session end is enabled.",
			},
			"session_start": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Tells us if the logging of session start is enabled.",
			},
		},
	}
}

func dataSourceAviatrixDcfLogProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	profileName, ok := d.Get("profile_name").(string)
	if !ok {
		return diag.Errorf("profile_name must be of type string")
	}
	if profileName == "" {
		return diag.Errorf("profile_name must be specified")
	}

	logProfile, err := client.GetLogProfileByName(ctx, profileName)
	if err != nil {
		return diag.Errorf("could not get DCF log profile: %s", err)
	}

	err = d.Set("profile_name", logProfile.ProfileName)
	if err != nil {
		return diag.Errorf("could not set profile_name: %s", err)
	}
	err = d.Set("profile_id", logProfile.ProfileID)
	if err != nil {
		return diag.Errorf("could not set profile_id: %s", err)
	}
	err = d.Set("session_end", logProfile.SessionEnd)
	if err != nil {
		return diag.Errorf("could not set session_end: %s", err)
	}
	err = d.Set("session_start", logProfile.SessionStart)
	if err != nil {
		return diag.Errorf("could not set session_start: %s", err)
	}

	d.SetId(logProfile.ProfileID)

	return nil
}
