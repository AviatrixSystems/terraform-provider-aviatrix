package aviatrix

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDCFTLSProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixDCFTLSProfileCreate,
		ReadContext:   resourceAviatrixDCFTLSProfileRead,
		UpdateContext: resourceAviatrixDCFTLSProfileUpdate,
		DeleteContext: resourceAviatrixDCFTLSProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name for the TLS profile.",
			},
			"certificate_validation": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"CERTIFICATE_VALIDATION_NONE", "CERTIFICATE_VALIDATION_LOG_ONLY", "CERTIFICATE_VALIDATION_ENFORCE"}, false),
				Description:  "Certificate validation mode. Must be one of CERTIFICATE_VALIDATION_NONE, CERTIFICATE_VALIDATION_LOG_ONLY, or CERTIFICATE_VALIDATION_ENFORCE.",
			},
			"verify_sni": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Toggle to enable advanced SNI verification of client connections.",
			},
			"ca_bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UUID of the CA bundle that should be used for origin certificate validation. If not populated the default bundle would be used. The aviatrix_dcf_trustbundle data source can be used to get the UUID from the bundle name.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the TLS profile.",
			},
		},
	}
}

func marshalDCFTLSProfileInput(d *schema.ResourceData) (*goaviatrix.TLSProfile, error) {
	tlsProfile := &goaviatrix.TLSProfile{}

	displayName := getString(d, "display_name")

	tlsProfile.DisplayName = displayName

	certificateValidation := getString(d, "certificate_validation")

	tlsProfile.CertificateValidation = certificateValidation

	verifySni := getBool(d, "verify_sni")

	tlsProfile.VerifySni = verifySni

	if caBundleID := getString(d, "ca_bundle_id"); caBundleID != "" {
		if _, err := uuid.Parse(caBundleID); err != nil {
			return nil, fmt.Errorf("ca_bundle_id must be a valid UUID: %w", err)
		}
		tlsProfile.CABundleID = &caBundleID
	}

	return tlsProfile, nil
}

func resourceAviatrixDCFTLSProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	tlsProfile, err := marshalDCFTLSProfileInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF TLS Profile during create: %s", err)
	}

	uuid, err := client.CreateTLSProfile(ctx, tlsProfile)
	if err != nil {
		return diag.Errorf("failed to create DCF TLS Profile: %s", err)
	}
	d.SetId(uuid)
	mustSet(d, "uuid", uuid)
	return nil
}

func resourceAviatrixDCFTLSProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	tlsProfile, err := client.GetTLSProfile(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF TLS Profile: %s", err)
	}

	if err := d.Set("display_name", tlsProfile.DisplayName); err != nil {
		return diag.Errorf("failed to set display_name during DCF TLS Profile read: %s", err)
	}

	if err := d.Set("certificate_validation", tlsProfile.CertificateValidation); err != nil {
		return diag.Errorf("failed to set certificate_validation during DCF TLS Profile read: %s", err)
	}

	if err := d.Set("verify_sni", tlsProfile.VerifySni); err != nil {
		return diag.Errorf("failed to set verify_sni during DCF TLS Profile read: %s", err)
	}

	if tlsProfile.CABundleID != nil {
		if err := d.Set("ca_bundle_id", *tlsProfile.CABundleID); err != nil {
			return diag.Errorf("failed to set ca_bundle_id during DCF TLS Profile read: %s", err)
		}
	}

	if err := d.Set("uuid", uuid); err != nil {
		return diag.Errorf("failed to set uuid during DCF TLS Profile read: %s", err)
	}

	return nil
}

func resourceAviatrixDCFTLSProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	tlsProfile, err := marshalDCFTLSProfileInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF TLS Profile during update: %s", err)
	}

	uuid := d.Id()
	err = client.UpdateTLSProfile(ctx, uuid, tlsProfile)
	if err != nil {
		return diag.Errorf("failed to update DCF TLS Profile: %s", err)
	}

	return nil
}

func resourceAviatrixDCFTLSProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	err := client.DeleteTLSProfile(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete DCF TLS Profile: %v", err)
	}

	return nil
}
