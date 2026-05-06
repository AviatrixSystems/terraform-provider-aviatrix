package aviatrix

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// telixProfileNotImplemented is the diagnostic returned by every CRUD hook
// in this PR. The schema lands first so the user-facing surface can be
// reviewed in isolation; the provider.go registration and the lifecycle
// implementation (marshal / flatten helpers, drift suppression for
// write-only fields) will land together in a follow-up PR.
const telixProfileNotImplemented = "aviatrix_telix_profile %s is not implemented yet; lifecycle handlers will land in a follow-up PR"

// Keep resourceAviatrixTelixProfile and its helpers reachable from the
// `unused` linter's perspective until provider.go registers the resource
// in a follow-up PR. Remove this declaration when the registration lands.
var _ = resourceAviatrixTelixProfile

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
										Description: "Optional static headers sent with each OTLP request. Values are write-only and redacted by the controller on read.",
									},
									"tls": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "TLS configuration for the OTLP endpoint.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ca_certificate_pem": {
													Type:         schema.TypeString,
													Optional:     true,
													Sensitive:    true,
													ValidateFunc: goaviatrix.ValidateCertificates,
													Description:  "PEM-encoded CA certificate used to validate the remote endpoint. Write-only.",
												},
												"client_certificate_pem": {
													Type:         schema.TypeString,
													Optional:     true,
													Sensitive:    true,
													ValidateFunc: validation.StringIsNotWhiteSpace,
													Description:  "PEM-encoded client certificate used for mutual TLS. Write-only.",
												},
												"client_private_key_pem": {
													Type:         schema.TypeString,
													Optional:     true,
													Sensitive:    true,
													ValidateFunc: validation.StringIsNotWhiteSpace,
													Description:  "PEM-encoded client private key used for mutual TLS. Write-only.",
												},
												"server_name_override": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional server name override used during certificate validation.",
												},
												"insecure_skip_verify": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Disables TLS certificate verification. Intended for testing only.",
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

func resourceAviatrixTelixProfileCreate(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return diag.Errorf(telixProfileNotImplemented, "create")
}

func resourceAviatrixTelixProfileRead(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return diag.Errorf(telixProfileNotImplemented, "read")
}

func resourceAviatrixTelixProfileUpdate(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return diag.Errorf(telixProfileNotImplemented, "update")
}

func resourceAviatrixTelixProfileDelete(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return diag.Errorf(telixProfileNotImplemented, "delete")
}
