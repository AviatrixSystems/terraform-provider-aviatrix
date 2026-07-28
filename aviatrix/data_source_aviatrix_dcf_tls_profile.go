package aviatrix

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixDCFTLSProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAviatrixDCFTLSProfileRead,
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the TLS profile.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the TLS profile.",
			},
			"certificate_validation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate validation mode.",
			},
			"verify_sni": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Toggle to enable advanced SNI verification of client connections.",
			},
			"ca_bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the CA bundle used for origin certificate validation.",
			},
		},
	}
}

func dataSourceAviatrixDCFTLSProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	displayName := getString(d, "display_name")

	if displayName == "" {
		return diag.Errorf("display_name must be specified")
	}

	tlsProfile, err := client.GetTLSProfileByName(ctx, displayName)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return diag.Errorf("DCF TLS Profile with display_name '%s' not found", displayName)
		}
		return diag.Errorf("failed to read DCF TLS Profile: %s", err)
	}

	mustSet(d, "display_name", tlsProfile.DisplayName)
	mustSet(d, "uuid", tlsProfile.UUID)
	mustSet(d, "certificate_validation", tlsProfile.CertificateValidation)
	mustSet(d, "verify_sni", tlsProfile.VerifySni)

	if tlsProfile.CABundleID != nil {
		mustSet(d, "ca_bundle_id", *tlsProfile.CABundleID)
	}

	d.SetId(tlsProfile.UUID)

	return nil
}
