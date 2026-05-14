package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// resource_aviatrix_telix_profile.go implements the lifecycle for the
// aviatrix_telix_profile managed resource.
//
// Sensitive field contract (write-only, in state):
//   - The controller redacts headers, ca_certificate_pem, client_certificate_pem,
//     and client_private_key_pem on read; only presence flags (has_*) are returned.
//   - Those four fields are marked Sensitive: true and live in state. The
//     user-supplied value is the authoritative source for the diff engine.
//   - Read populates every non-sensitive field plus the has_* mirrors but never
//     overwrites the four sensitive fields, except when has_* flips false while
//     state still holds a value. In that case Read clears the state slot so the
//     next plan re-pushes from configuration (external-deletion drift).
//   - External rotation to a different value is not detectable. Documented as a
//     known limitation in the registry doc.

const (
	telixPathOtlp                    = "destination.0.otlp.0"
	telixPathOtlpEndpoint            = telixPathOtlp + ".endpoint"
	telixPathOtlpProtocol            = telixPathOtlp + ".protocol"
	telixPathOtlpHeaders             = telixPathOtlp + ".headers"
	telixPathOtlpTLSList             = telixPathOtlp + ".tls"
	telixPathOtlpTLS                 = telixPathOtlpTLSList + ".0"
	telixPathTLSCA                   = telixPathOtlpTLS + ".ca_certificate_pem"
	telixPathTLSClientCertificatePEM = telixPathOtlpTLS + ".client_certificate_pem"
	telixPathTLSClientPrivateKeyPEM  = telixPathOtlpTLS + ".client_private_key_pem"
	telixPathTLSServerNameOverride   = telixPathOtlpTLS + ".server_name_override"
	telixPathTLSInsecureSkipVerify   = telixPathOtlpTLS + ".insecure_skip_verify"
)

func resourceAviatrixTelixProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixTelixProfileCreate,
		ReadWithoutTimeout:   resourceAviatrixTelixProfileRead,
		UpdateWithoutTimeout: resourceAviatrixTelixProfileUpdate,
		DeleteWithoutTimeout: resourceAviatrixTelixProfileDelete,
		CustomizeDiff:        validateTelixProfileFilterSources,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Display name for the Telix export profile. Must be unique across profiles.",
			},
			"sources": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						goaviatrix.TelixTelemetrySourceDcfLogs,
						goaviatrix.TelixTelemetrySourceNodeExporter,
						goaviatrix.TelixTelemetrySourceGatewayOperationalSyslog,
					}, false),
				},
				Description: "Telemetry sources exported by this profile. Immutable after creation; changing this list forces resource replacement.",
			},
			"destination": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Destination configuration for telemetry export. Exactly one variant block must be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// otlp is Required today because it is the only destination
						// variant. When a second variant (e.g. splunk) is added,
						// flip otlp to Optional and add
						//     ExactlyOneOf: []string{"destination.0.otlp", "destination.0.<new>"}
						// to both children. That transition is backward compatible:
						// existing configs with `destination { otlp { ... } }` still
						// validate. See the scope block below for the same pattern
						// applied to all_gateways / selected_gateways.
						"otlp": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "OTLP destination configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Description:  "OTLP collector or backend endpoint URL.",
									},
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											goaviatrix.TelixOtlpProtocolGRPC,
											goaviatrix.TelixOtlpProtocolHTTP,
										}, false),
										Description: "OTLP protocol used for export. Immutable after creation; changing this value forces resource replacement.",
									},
									"headers": {
										Type:        schema.TypeMap,
										Optional:    true,
										Sensitive:   true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Optional static headers sent with each OTLP request. Write-only: the controller never returns this value on read. External rotation will not be detected until the value is changed in configuration. See resource documentation.",
									},
									"has_headers": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the controller currently holds a non-empty headers value for this profile.",
									},
									"tls": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "TLS configuration for the OTLP endpoint. Omit entirely when not using TLS material; omit after use to clear TLS on the controller (presence flags nested under this block remain computed-only).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ca_certificate_pem": {
													Type:         schema.TypeString,
													Optional:     true,
													Sensitive:    true,
													ValidateFunc: goaviatrix.ValidateCertificates,
													Description:  "PEM-encoded CA certificate used to validate the remote endpoint. Write-only: the controller never returns this value on read. External rotation will not be detected until the value is changed in configuration. See resource documentation.",
												},
												"client_certificate_pem": {
													Type:         schema.TypeString,
													Optional:     true,
													Sensitive:    true,
													ValidateFunc: validation.StringIsNotWhiteSpace,
													Description:  "PEM-encoded client certificate used for mutual TLS. Write-only: the controller never returns this value on read. External rotation will not be detected until the value is changed in configuration. See resource documentation.",
												},
												"client_private_key_pem": {
													Type:         schema.TypeString,
													Optional:     true,
													Sensitive:    true,
													ValidateFunc: validation.StringIsNotWhiteSpace,
													Description:  "PEM-encoded client private key used for mutual TLS. Write-only: the controller never returns this value on read. External rotation will not be detected until the value is changed in configuration. See resource documentation.",
												},
												"server_name_override": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional server name override used during certificate validation. The controller returns the stored value on read.",
												},
												"insecure_skip_verify": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Disables TLS certificate verification. Intended for testing only. The controller returns the stored value on read.",
												},
												"has_ca_certificate": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the controller currently holds a CA certificate for this profile.",
												},
												"has_client_certificate": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the controller currently holds a client certificate for this profile.",
												},
												"has_client_private_key": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the controller currently holds a client private key for this profile.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"scope": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Gateway scope for the profile. Exactly one of all_gateways or selected_gateways must be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_gateways": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"scope.0.all_gateways", "scope.0.selected_gateways"},
							Elem:         &schema.Resource{Schema: map[string]*schema.Schema{}},
							Description:  "Marker block applying the profile to every gateway compatible with the selected sources. Set as an empty block: all_gateways {}.",
						},
						"selected_gateways": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"scope.0.all_gateways", "scope.0.selected_gateways"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_names": {
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Specific gateways to which this profile applies.",
									},
								},
							},
							Description: "Restricts the profile to a specific list of gateways.",
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether this export profile is enabled. Defaults to true.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Optional filtering for telemetry data before export.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcf_log_types": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "DCF log types to include. Requires TELEMETRY_SOURCE_DCF_LOGS in sources.",
						},
						"dcf_actions": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "DCF actions to include.",
						},
					},
				},
			},
			"profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-generated unique identifier for the profile.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RFC3339 timestamp of when the profile was created.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RFC3339 timestamp of when the profile was last modified.",
			},
		},
	}
}

// validateTelixProfileFilterSources enforces the YAML-documented rule that
// any DCF-specific filter (dcf_log_types, dcf_actions) requires
// TELEMETRY_SOURCE_DCF_LOGS to be present in the sources list. The check
// runs at plan time so misconfigurations surface before any API call.
func validateTelixProfileFilterSources(_ context.Context, d *schema.ResourceDiff, _ any) error {
	filterList, ok := d.Get("filters").([]any)
	if !ok || len(filterList) == 0 {
		return nil
	}
	filter, ok := filterList[0].(map[string]any)
	if !ok {
		return nil
	}

	dcfLogTypes, _ := filter["dcf_log_types"].([]any)
	dcfActions, _ := filter["dcf_actions"].([]any)
	if len(dcfLogTypes) == 0 && len(dcfActions) == 0 {
		return nil
	}

	sources, _ := d.Get("sources").([]any)
	for _, s := range sources {
		if s == goaviatrix.TelixTelemetrySourceDcfLogs {
			return nil
		}
	}

	return fmt.Errorf("filters.dcf_log_types or filters.dcf_actions is set but %q is not in sources; either add it to sources or remove the filter",
		goaviatrix.TelixTelemetrySourceDcfLogs)
}

func marshalTelixProfileCreateRequest(d *schema.ResourceData) *goaviatrix.TelixProfileCreateRequest {
	enabled := getBool(d, "enabled")
	return &goaviatrix.TelixProfileCreateRequest{
		DisplayName: getString(d, "display_name"),
		Sources:     getStringList(d, "sources"),
		Destination: expandTelixDestinationCreate(d),
		Scope:       expandTelixScopeValue(d),
		Enabled:     &enabled,
		Filters:     expandTelixFilters(d),
	}
}

// marshalTelixProfileUpdateRequest builds a partial PATCH payload. Most fields
// are populated only when HasChange reports a diff. For filters, HasChange
// alone is not trusted when only nested dcf_* lists change; see the inline note.
// Sensitive OTLP sub-fields are gated by their own HasChange so unrelated edits
// never re-push secrets.
func marshalTelixProfileUpdateRequest(d *schema.ResourceData) *goaviatrix.TelixProfileUpdateRequest {
	req := &goaviatrix.TelixProfileUpdateRequest{}

	if d.HasChange("display_name") {
		v := getString(d, "display_name")
		req.DisplayName = &v
	}
	if d.HasChange("enabled") {
		v := getBool(d, "enabled")
		req.Enabled = &v
	}
	if d.HasChange("destination") {
		req.Destination = expandTelixDestinationUpdate(d)
	}
	if d.HasChange("scope") {
		v := expandTelixScopeValue(d)
		req.Scope = &v
	}
	// Filters: optional MaxItems:1 block whose only attributes are nested
	// TypeLists (dcf_*). SDK v2 can leave HasChange("filters") false when
	// only those nested lists change, so we always expand filters whenever the
	// configuration includes a filters block (and still rely on HasChange when
	// the block is removed—empty config, cleared in state).
	filterList, _ := d.Get("filters").([]any)
	hasFiltersBlock := len(filterList) > 0 && filterList[0] != nil
	if hasFiltersBlock {
		req.Filters = expandTelixFilters(d)
	} else if d.HasChange("filters") {
		req.Filters = expandTelixFilters(d)
	}

	return req
}

// expandTelixDestinationCreate builds the full destination payload for create.
// Every field comes from current configuration.
func expandTelixDestinationCreate(d *schema.ResourceData) goaviatrix.TelixDestinationInput {
	headersMap := expandStringMap(d, telixPathOtlpHeaders)
	var headersPtr *map[string]string
	if len(headersMap) > 0 {
		h := headersMap
		headersPtr = &h
	}

	otlp := goaviatrix.TelixOtlpDestinationInput{
		Endpoint: getString(d, telixPathOtlpEndpoint),
		Protocol: getString(d, telixPathOtlpProtocol),
		Headers:  headersPtr,
		TLS:      expandTelixTLSCreate(d),
	}
	return goaviatrix.TelixDestinationInput{Otlp: otlp}
}

// expandTelixDestinationUpdate builds a destination payload for update. Endpoint
// and Protocol are required by the API contract and always populated from
// current config (which equals state for unchanged fields). Headers and TLS
// sub-fields are included only when their own slot changed, so unchanged
// secrets are not re-pushed by an unrelated destination edit.
func expandTelixDestinationUpdate(d *schema.ResourceData) *goaviatrix.TelixDestinationInput {
	otlp := goaviatrix.TelixOtlpDestinationInput{
		Endpoint: getString(d, telixPathOtlpEndpoint),
		Protocol: getString(d, telixPathOtlpProtocol),
	}

	if d.HasChange(telixPathOtlpHeaders) {
		m := expandStringMap(d, telixPathOtlpHeaders)
		if len(m) == 0 {
			// Non-nil pointer to {} so PATCH JSON carries "headers": {} and the
			// controller clears stored OTLP headers; nil would omit headers (no-op).
			cleared := map[string]string{}
			otlp.Headers = &cleared
		} else {
			otlp.Headers = &m
		}
	}

	if tls := buildTelixTLSPatch(d); tls != nil {
		otlp.TLS = tls
	}

	return &goaviatrix.TelixDestinationInput{Otlp: otlp}
}

// telixOtlpTLSListAbsent reports whether a ResourceData/GetChange value carries
// no OTLP tls block element (either missing list, zero length, or a nil elem).
func telixOtlpTLSListAbsent(v any) bool {
	if v == nil {
		return true
	}
	list, ok := v.([]any)
	if !ok || len(list) == 0 {
		return true
	}
	return list[0] == nil
}

// telixOtlpTLSBlockRemoval returns true when the tls list changed from
// "block present" to "absent" (user removed destination.otlp.tls { ... }).
func telixOtlpTLSBlockRemoval(hasTLSListChange bool, oldTLS, newTLS any) bool {
	if !hasTLSListChange {
		return false
	}
	return !telixOtlpTLSListAbsent(oldTLS) && telixOtlpTLSListAbsent(newTLS)
}

// telixOtlpTLSClearPayload returns a PATCH TLS body that clears every mergeable
// OTLP TLS field. The controller treats non-nil *string pointing at "" as
// "clear encrypted PEM / string field" and requires both client cert and key
// pointers when clearing mTLS material.
func telixOtlpTLSClearPayload() *goaviatrix.TelixTLSInputConfig {
	empty := ""
	insecureOff := false
	return &goaviatrix.TelixTLSInputConfig{
		CaCertificatePem:     &empty,
		ClientCertificatePem: &empty,
		ClientPrivateKeyPem:  &empty,
		ServerNameOverride:   &empty,
		InsecureSkipVerify:   &insecureOff,
	}
}

// expandTelixTLSCreate builds the TLS input config for create. Returns nil if
// the user did not supply a TLS block.
func expandTelixTLSCreate(d *schema.ResourceData) *goaviatrix.TelixTLSInputConfig {
	tlsList, ok := d.Get(telixPathOtlpTLSList).([]any)
	if !ok || len(tlsList) == 0 || tlsList[0] == nil {
		return nil
	}

	out := &goaviatrix.TelixTLSInputConfig{}
	if v := getString(d, telixPathTLSCA); v != "" {
		out.CaCertificatePem = &v
	}
	if v := getString(d, telixPathTLSClientCertificatePEM); v != "" {
		out.ClientCertificatePem = &v
	}
	if v := getString(d, telixPathTLSClientPrivateKeyPEM); v != "" {
		out.ClientPrivateKeyPem = &v
	}
	if v := getString(d, telixPathTLSServerNameOverride); v != "" {
		out.ServerNameOverride = &v
	}
	v := getBool(d, telixPathTLSInsecureSkipVerify)
	out.InsecureSkipVerify = &v
	return out
}

// buildTelixTLSPatch returns a TelixTLSInputConfig for PATCH merge. When the
// user removes destination.otlp.tls entirely, returns a teardown payload so
// the controller clears stored TLS material (omitting tls in JSON leaves it
// unchanged). Otherwise builds only PEM and related fields whose attributes
// changed.
func buildTelixTLSPatch(d *schema.ResourceData) *goaviatrix.TelixTLSInputConfig {
	oldTLS, newTLS := d.GetChange(telixPathOtlpTLSList)
	if telixOtlpTLSBlockRemoval(d.HasChange(telixPathOtlpTLSList), oldTLS, newTLS) {
		return telixOtlpTLSClearPayload()
	}

	tls := &goaviatrix.TelixTLSInputConfig{}
	changed := false

	if d.HasChange(telixPathTLSCA) {
		v := getString(d, telixPathTLSCA)
		tls.CaCertificatePem = &v
		changed = true
	}
	if d.HasChange(telixPathTLSClientCertificatePEM) {
		v := getString(d, telixPathTLSClientCertificatePEM)
		tls.ClientCertificatePem = &v
		changed = true
	}
	if d.HasChange(telixPathTLSClientPrivateKeyPEM) {
		v := getString(d, telixPathTLSClientPrivateKeyPEM)
		tls.ClientPrivateKeyPem = &v
		changed = true
	}
	if d.HasChange(telixPathTLSServerNameOverride) {
		v := getString(d, telixPathTLSServerNameOverride)
		tls.ServerNameOverride = &v
		changed = true
	}
	if d.HasChange(telixPathTLSInsecureSkipVerify) {
		v := getBool(d, telixPathTLSInsecureSkipVerify)
		tls.InsecureSkipVerify = &v
		changed = true
	}

	if !changed {
		return nil
	}
	return tls
}

func expandTelixScopeValue(d *schema.ResourceData) goaviatrix.TelixGatewayScope {
	scope := goaviatrix.TelixGatewayScope{}
	scopeList, ok := d.Get("scope").([]any)
	if !ok || len(scopeList) == 0 || scopeList[0] == nil {
		return scope
	}
	scopeMap, ok := scopeList[0].(map[string]any)
	if !ok {
		return scope
	}

	if all, ok := scopeMap["all_gateways"].([]any); ok && len(all) > 0 {
		scope.AllGateways = &goaviatrix.TelixAllGatewaysScope{}
		return scope
	}
	if sel, ok := scopeMap["selected_gateways"].([]any); ok && len(sel) > 0 {
		selMap, ok := sel[0].(map[string]any)
		if !ok {
			return scope
		}
		gws, _ := selMap["gateway_names"].([]any)
		names := make([]string, 0, len(gws))
		for _, g := range gws {
			names = append(names, mustString(g))
		}
		scope.SelectedGateways = &goaviatrix.TelixSelectedGatewaysScope{GatewayNames: names}
	}
	return scope
}

func expandTelixFilters(d *schema.ResourceData) *goaviatrix.TelixFilterConfig {
	filtersList, ok := d.Get("filters").([]any)
	if !ok || len(filtersList) == 0 || filtersList[0] == nil {
		return nil
	}
	f, ok := filtersList[0].(map[string]any)
	if !ok {
		return nil
	}

	out := &goaviatrix.TelixFilterConfig{}
	if v, ok := f["dcf_log_types"].([]any); ok {
		for _, t := range v {
			out.DcfLogTypes = append(out.DcfLogTypes, mustString(t))
		}
	}
	if v, ok := f["dcf_actions"].([]any); ok {
		for _, t := range v {
			out.DcfActions = append(out.DcfActions, mustString(t))
		}
	}
	return out
}

func expandStringMap(d *schema.ResourceData, key string) map[string]string {
	raw, ok := d.Get(key).(map[string]any)
	if !ok || len(raw) == 0 {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = mustString(v)
	}
	return out
}

// flattenTelixProfileDetail copies the controller's GET response into Terraform
// state. The single load-bearing rule of Approach 1 lives here: we do not
// overwrite the four sensitive fields with values from the response. The user-
// supplied value already in state remains authoritative. The only exception is
// external-deletion drift: when a presence flag flips to false while state
// still holds a value, we clear that state slot so the next plan re-pushes
// from configuration.
func flattenTelixProfileDetail(d *schema.ResourceData, detail *goaviatrix.TelixProfileDetail) error {
	mustSet(d, "profile_id", detail.ProfileID)
	mustSet(d, "display_name", detail.DisplayName)
	mustSet(d, "enabled", detail.Enabled)
	mustSet(d, "sources", detail.Sources)

	if !detail.CreatedAt.IsZero() {
		mustSet(d, "created_at", detail.CreatedAt.Format(time.RFC3339))
	} else {
		mustSet(d, "created_at", "")
	}
	if !detail.LastModifiedAt.IsZero() {
		mustSet(d, "last_modified_at", detail.LastModifiedAt.Format(time.RFC3339))
	} else {
		mustSet(d, "last_modified_at", "")
	}

	if err := flattenTelixDestination(d, &detail.Destination); err != nil {
		return fmt.Errorf("set destination: %w", err)
	}

	if err := flattenTelixScope(d, &detail.Scope); err != nil {
		return fmt.Errorf("set scope: %w", err)
	}

	if err := flattenTelixFilters(d, detail.Filters); err != nil {
		return fmt.Errorf("set filters: %w", err)
	}

	return nil
}

func flattenTelixDestination(d *schema.ResourceData, dest *goaviatrix.TelixDestination) error {
	hasHeaders := derefBool(dest.Otlp.HasHeaders)
	otlp := map[string]any{
		"endpoint":    dest.Otlp.Endpoint,
		"protocol":    dest.Otlp.Protocol,
		"has_headers": hasHeaders,
	}

	// Preserve user-supplied headers. Clear when the controller no longer holds
	// any (external deletion).
	statedHeaders, _ := d.Get(telixPathOtlpHeaders).(map[string]any)
	if !hasHeaders && len(statedHeaders) > 0 {
		otlp["headers"] = map[string]any{}
	} else {
		otlp["headers"] = statedHeaders
	}

	otlp["tls"] = flattenTelixTLS(d, dest.Otlp.TLS)

	return d.Set("destination", []any{
		map[string]any{
			"otlp": []any{otlp},
		},
	})
}

func flattenTelixTLS(d *schema.ResourceData, tls *goaviatrix.TelixTLSConfig) []any {
	if tls == nil {
		return []any{}
	}

	hasCa := derefBool(tls.HasCaCertificate)
	hasClientCert := derefBool(tls.HasClientCertificate)
	hasClientKey := derefBool(tls.HasClientPrivateKey)

	out := map[string]any{
		"server_name_override":   derefString(tls.ServerNameOverride),
		"insecure_skip_verify":   derefBool(tls.InsecureSkipVerify),
		"has_ca_certificate":     hasCa,
		"has_client_certificate": hasClientCert,
		"has_client_private_key": hasClientKey,
	}

	// Preserve user-supplied PEMs from current state, except when the matching
	// presence flag is false and state still holds a value (external deletion).
	statedTLS, _ := d.Get(telixPathOtlpTLS).(map[string]any)

	for _, slot := range []struct {
		field   string
		present bool
	}{
		{"ca_certificate_pem", hasCa},
		{"client_certificate_pem", hasClientCert},
		{"client_private_key_pem", hasClientKey},
	} {
		stated, _ := statedTLS[slot.field].(string)
		if !slot.present && stated != "" {
			out[slot.field] = ""
		} else {
			out[slot.field] = stated
		}
	}

	return []any{out}
}

func flattenTelixScope(d *schema.ResourceData, scope *goaviatrix.TelixGatewayScope) error {
	s := map[string]any{}
	if scope.AllGateways != nil {
		s["all_gateways"] = []any{map[string]any{}}
	}
	if scope.SelectedGateways != nil {
		s["selected_gateways"] = []any{
			map[string]any{
				"gateway_names": scope.SelectedGateways.GatewayNames,
			},
		}
	}
	return d.Set("scope", []any{s})
}

func flattenTelixFilters(d *schema.ResourceData, filters *goaviatrix.TelixFilterConfig) error {
	if filters == nil || (len(filters.DcfLogTypes) == 0 && len(filters.DcfActions) == 0) {
		return d.Set("filters", []any{})
	}
	f := map[string]any{}
	if len(filters.DcfLogTypes) > 0 {
		f["dcf_log_types"] = filters.DcfLogTypes
	}
	if len(filters.DcfActions) > 0 {
		f["dcf_actions"] = filters.DcfActions
	}
	return d.Set("filters", []any{f})
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func resourceAviatrixTelixProfileCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	profileID, err := client.CreateTelixProfile(ctx, marshalTelixProfileCreateRequest(d))
	if err != nil {
		return diag.Errorf("failed to create Telix profile: %s", err)
	}

	d.SetId(profileID)
	return resourceAviatrixTelixProfileRead(ctx, d, meta)
}

func resourceAviatrixTelixProfileRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	detail, err := client.GetTelixProfile(ctx, d.Id())
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Telix profile %s: %s", d.Id(), err)
	}

	if err := flattenTelixProfileDetail(d, detail); err != nil {
		return diag.Errorf("failed to populate Telix profile state: %s", err)
	}
	return nil
}

func resourceAviatrixTelixProfileUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if err := client.UpdateTelixProfile(ctx, d.Id(), marshalTelixProfileUpdateRequest(d)); err != nil {
		return diag.Errorf("failed to update Telix profile %s: %s", d.Id(), err)
	}

	return resourceAviatrixTelixProfileRead(ctx, d, meta)
}

func resourceAviatrixTelixProfileDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	if err := client.DeleteTelixProfile(ctx, d.Id()); err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return nil
		}
		return diag.Errorf("failed to delete Telix profile %s: %s", d.Id(), err)
	}
	return nil
}
