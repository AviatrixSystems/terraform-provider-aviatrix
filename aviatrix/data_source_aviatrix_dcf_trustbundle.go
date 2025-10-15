package aviatrix

import (
	"context"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Content of the DCF Trust Bundle as a list of strings.",
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
		return diag.Errorf("could not get DCF trust bundle: %s", err)
	}

	// fmt.Printf("trustBundle ID: %s", trustBundle.BundleID)

	if err := d.Set("bundle_id", trustBundle.BundleID); err != nil {
		return diag.Errorf("could not set bundle_id: %s", err)
	}

	if err := d.Set("bundle_content", trustBundle.BundleContent); err != nil {
		return diag.Errorf("could not set bundle_content: %s", err)
	}

	if err := d.Set("created_at", trustBundle.CreatedAt); err != nil {
		return diag.Errorf("could not set created_at: %s", err)
	}

	d.SetId(trustBundle.UUID)

	return nil
}
