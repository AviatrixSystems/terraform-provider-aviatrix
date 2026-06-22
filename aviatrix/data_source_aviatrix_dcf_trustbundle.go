package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixDcfTrustbundle() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDcfTrustbundleRead,
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display Name of the DCF Trust Bundle.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the DCF Trust Bundle.",
			},
			"bundle_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content of the DCF Trust Bundle as a string.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the DCF Trust Bundle was created.",
			},
		},
	}
}

func dataSourceAviatrixDcfTrustbundleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	name, ok := d.Get("display_name").(string)
	if !ok {
		return diag.Errorf("display_name must be of type string")
	}
	if name == "" {
		return diag.Errorf("display_name must be specified")
	}

	trustBundle, err := client.GetDCFTrustBundleByName(ctx, name)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return diag.Errorf("DCF Trust Bundle not found: %s", err)
		}
		return diag.Errorf("failed to read DCF Trust Bundle: %s", err)
	}

	d.Set("bundle_id", trustBundle.BundleID)
	d.Set("display_name", trustBundle.DisplayName)
	if !trustBundle.CreatedAt.IsZero() {
		d.Set("created_at", trustBundle.CreatedAt.Format("2006-01-02T15:04:05Z"))
	} else {
		d.Set("created_at", "")
	}
	bundleContent := strings.TrimSpace(strings.Join(trustBundle.BundleContent, "\n"))
	d.Set("bundle_content", bundleContent)
	d.SetId(trustBundle.UUID)
	return nil
}
