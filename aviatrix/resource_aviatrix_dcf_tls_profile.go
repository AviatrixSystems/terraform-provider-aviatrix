package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"slices"

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
		CustomizeDiff: resourceAviatrixDCFTLSProfileCustomizeDiff,
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
			"tls_log_fields_extension": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(append([]string{tlsLogAllField}, allTLSLogFields...), false),
				},
				Description: "TLS log fields to include in the logs. Use TLS_LOG_FIELDS_ALL on its own to include every field. Omitting this argument includes no fields.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the TLS profile.",
			},
		},
	}
}

// tlsLogAllField is a provider-only shorthand for every value in allTLSLogFields.
const tlsLogAllField = "TLS_LOG_FIELDS_ALL"

// allTLSLogFields lists every TLS log field the controller API accepts.
var allTLSLogFields = []string{
	"TLS_CLIENT_VERSIONS",
	"TLS_CLIENT_CIPHERS",
	"TLS_CLIENT_KEY_SHARE_GROUPS",
	"TLS_CLIENT_SIGNATURE_ALGORITHMS",
	"TLS_CLIENT_JA4_FINGERPRINT",
	"TLS_CLIENT_SNI",
	"TLS_CLIENT_SUPPORTED_GROUPS",
	"TLS_CLIENT_ALPNS",
	"TLS_CLIENT_ECH_PRESENT",
}

// resourceAviatrixDCFTLSProfileCustomizeDiff rejects TLS_LOG_FIELDS_ALL combined with other log fields at plan time.
func resourceAviatrixDCFTLSProfileCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}
	raw := rawConfig.GetAttr("tls_log_fields_extension")
	if raw.IsNull() || !raw.IsWhollyKnown() {
		return nil
	}
	var fields []string
	for it := raw.ElementIterator(); it.Next(); {
		_, v := it.Element()
		fields = append(fields, v.AsString())
	}
	return validateTLSLogFields(fields)
}

// validateTLSLogFields returns an error when TLS_LOG_FIELDS_ALL is not the only configured value.
func validateTLSLogFields(fields []string) error {
	if slices.Contains(fields, tlsLogAllField) && len(fields) > 1 {
		return fmt.Errorf("tls_log_fields_extension: %s must be the only value when set", tlsLogAllField)
	}
	return nil
}

// expandTLSLogFields converts configured log fields into the API request list.
func expandTLSLogFields(fields []string) []string {
	// The API has no ALL value; it recognizes the complete list of fields instead.
	if slices.Contains(fields, tlsLogAllField) {
		return slices.Clone(allTLSLogFields)
	}
	return fields
}

// flattenTLSLogFields converts the API's log fields into the state value.
func flattenTLSLogFields(apiFields, priorFields []string) []string {
	// The API reports ALL as the complete list of fields; keep the shorthand in
	// state when the prior state used it so configs using it do not plan a diff.
	if slices.Contains(priorFields, tlsLogAllField) && isAllTLSLogFields(apiFields) {
		return []string{tlsLogAllField}
	}
	return apiFields
}

// isAllTLSLogFields reports whether fields contains exactly every TLS log field.
func isAllTLSLogFields(fields []string) bool {
	if len(fields) != len(allTLSLogFields) {
		return false
	}
	for _, f := range allTLSLogFields {
		if !slices.Contains(fields, f) {
			return false
		}
	}
	return true
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

	tlsLogFields := getStringSet(d, "tls_log_fields_extension")
	if err := validateTLSLogFields(tlsLogFields); err != nil {
		return nil, err
	}
	tlsProfile.TLSLogFieldsExtension = expandTLSLogFields(tlsLogFields)

	return tlsProfile, nil
}

func resourceAviatrixDCFTLSProfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	return resourceAviatrixDCFTLSProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFTLSProfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	mustSet(d, "tls_log_fields_extension", flattenTLSLogFields(tlsProfile.TLSLogFieldsExtension, getStringSet(d, "tls_log_fields_extension")))

	if err := d.Set("uuid", uuid); err != nil {
		return diag.Errorf("failed to set uuid during DCF TLS Profile read: %s", err)
	}

	return nil
}

func resourceAviatrixDCFTLSProfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	return resourceAviatrixDCFTLSProfileRead(ctx, d, meta)
}

func resourceAviatrixDCFTLSProfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	err := client.DeleteTLSProfile(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete DCF TLS Profile: %v", err)
	}

	return nil
}
