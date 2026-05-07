package goaviatrix

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// telixProfileEndpoint is the v2.5 base path for Telix export profile
// operations.
const telixProfileEndpoint = "telix/profile"

// Telix telemetry source enum values returned by the controller.
const (
	TelixTelemetrySourceDcfLogs                  = "TELEMETRY_SOURCE_DCF_LOGS"
	TelixTelemetrySourceNodeExporter             = "TELEMETRY_SOURCE_NODE_EXPORTER"
	TelixTelemetrySourceGatewayOperationalSyslog = "TELEMETRY_SOURCE_GATEWAY_OPERATIONAL_SYSLOG"
)

// Telix OTLP protocol enum values accepted by the controller.
const (
	TelixOtlpProtocolGRPC = "TELIX_OTLP_PROTOCOL_GRPC"
	TelixOtlpProtocolHTTP = "TELIX_OTLP_PROTOCOL_HTTP"
)

// TelixFilterConfig configures optional filtering applied to telemetry data
// before it is exported to the destination.
type TelixFilterConfig struct {
	DcfLogTypes []string `json:"dcf_log_types,omitempty"`
	DcfActions  []string `json:"dcf_actions,omitempty"`
}

// TelixAllGatewaysScope marks a profile as applying to every gateway
// compatible with its telemetry sources. The API contract defines this as an
// empty object; presence of the variant field on the parent scope is what
// carries meaning.
type TelixAllGatewaysScope struct{}

// TelixSelectedGatewaysScope restricts a profile to a specific list of
// gateways identified by name.
type TelixSelectedGatewaysScope struct {
	GatewayNames []string `json:"gateway_names"`
}

// TelixGatewayScope is the oneOf wrapper for the gateway scope of a profile.
// A well-formed payload sets exactly one variant field; the wrapper itself
// does not enforce this and callers are responsible for the invariant.
type TelixGatewayScope struct {
	AllGateways      *TelixAllGatewaysScope      `json:"all_gateways,omitempty"`
	SelectedGateways *TelixSelectedGatewaysScope `json:"selected_gateways,omitempty"`
}

// TelixTLSInputConfig is the write-side TLS configuration accepted by create
// and update requests. Values are write-only and are not returned in
// subsequent reads; the controller exposes only presence flags via
// TelixTLSConfig. Client certificate and client private key fields support
// mutual TLS (mTLS) when the destination requires it.
type TelixTLSInputConfig struct {
	CaCertificatePem     *string `json:"ca_certificate_pem,omitempty"`
	ClientCertificatePem *string `json:"client_certificate_pem,omitempty"`
	ClientPrivateKeyPem  *string `json:"client_private_key_pem,omitempty"`
	ServerNameOverride   *string `json:"server_name_override,omitempty"`
	InsecureSkipVerify   *bool   `json:"insecure_skip_verify,omitempty"`
}

// TelixOtlpDestinationInput is the write-side representation of an OTLP
// destination. Headers and TLS material may be supplied here but will be
// redacted by the controller on subsequent reads.
type TelixOtlpDestinationInput struct {
	Endpoint string               `json:"endpoint"`
	Protocol string               `json:"protocol"`
	Headers  map[string]string    `json:"headers,omitempty"`
	TLS      *TelixTLSInputConfig `json:"tls,omitempty"`
}

// TelixDestinationInput is the write-side oneOf wrapper for destination
// configurations.
type TelixDestinationInput struct {
	Otlp TelixOtlpDestinationInput `json:"otlp"`
}

// TelixProfileCreateRequest is the request body for POST /telix/profile.
type TelixProfileCreateRequest struct {
	DisplayName string                `json:"display_name"`
	Sources     []string              `json:"sources"`
	Destination TelixDestinationInput `json:"destination"`
	Scope       TelixGatewayScope     `json:"scope"`
	Enabled     *bool                 `json:"enabled,omitempty"`
	Filters     *TelixFilterConfig    `json:"filters,omitempty"`
}

// TelixProfileUpdateRequest is the request body for PATCH
// /telix/profile/{profile_id}. All fields are optional; only fields explicitly
// supplied are updated. The telemetry sources list and the destination
// protocol are immutable per the API contract and cannot be changed via this
// request.
type TelixProfileUpdateRequest struct {
	DisplayName *string                `json:"display_name,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	Destination *TelixDestinationInput `json:"destination,omitempty"`
	Scope       *TelixGatewayScope     `json:"scope,omitempty"`
	Filters     *TelixFilterConfig     `json:"filters,omitempty"`
}

// TelixProfileCreateResponse is the response body returned by both
// POST /telix/profile and PATCH /telix/profile/{profile_id}. Only the
// server-generated profile identifier is returned.
type TelixProfileCreateResponse struct {
	ProfileID string `json:"profile_id"`
}

// TelixTLSConfig is the read-side TLS configuration reported by the controller
// for an outbound destination. Sensitive material (CA certificate, client
// certificate, client private key) is never returned; only presence flags are
// reported. All fields are optional in the response and are modeled as
// pointers so absent fields can be distinguished from explicitly false /
// empty values.
type TelixTLSConfig struct {
	HasCaCertificate     *bool   `json:"has_ca_certificate,omitempty"`
	HasClientCertificate *bool   `json:"has_client_certificate,omitempty"`
	HasClientPrivateKey  *bool   `json:"has_client_private_key,omitempty"`
	ServerNameOverride   *string `json:"server_name_override,omitempty"`
	InsecureSkipVerify   *bool   `json:"insecure_skip_verify,omitempty"`
}

// TelixOtlpDestination is the read-side representation of an OTLP destination.
// Header values and TLS material are redacted by the controller; the
// HasHeaders flag indicates only whether headers are configured.
type TelixOtlpDestination struct {
	Endpoint   string          `json:"endpoint"`
	Protocol   string          `json:"protocol"`
	HasHeaders *bool           `json:"has_headers,omitempty"`
	TLS        *TelixTLSConfig `json:"tls,omitempty"`
}

// TelixDestination is the read-side oneOf wrapper for destination
// configurations.
type TelixDestination struct {
	Otlp TelixOtlpDestination `json:"otlp"`
}

// TelixProfileItem is the summary view of a Telix export profile returned by
// the list endpoint.
type TelixProfileItem struct {
	ProfileID      string            `json:"profile_id"`
	DisplayName    string            `json:"display_name"`
	Enabled        bool              `json:"enabled"`
	CreatedAt      time.Time         `json:"created_at"`
	LastModifiedAt time.Time         `json:"last_modified_at"`
	Sources        []string          `json:"sources"`
	Destination    TelixDestination  `json:"destination"`
	Scope          TelixGatewayScope `json:"scope"`
}

// TelixProfileDetail is the full detail view of a Telix export profile
// returned by the get endpoint.
type TelixProfileDetail struct {
	ProfileID      string             `json:"profile_id"`
	DisplayName    string             `json:"display_name"`
	Enabled        bool               `json:"enabled"`
	Sources        []string           `json:"sources"`
	Destination    TelixDestination   `json:"destination"`
	Scope          TelixGatewayScope  `json:"scope"`
	Filters        *TelixFilterConfig `json:"filters,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	LastModifiedAt time.Time          `json:"last_modified_at"`
}

// TelixProfileListResponse is the response body of the list endpoint.
type TelixProfileListResponse struct {
	Profiles []TelixProfileItem `json:"profiles"`
	Total    int                `json:"total"`
}

// ListTelixProfiles returns the summary view of every Telix export profile
// configured on the controller.
func (c *Client) ListTelixProfiles(ctx context.Context) (*TelixProfileListResponse, error) {
	var response TelixProfileListResponse
	if err := c.GetAPIContext25(ctx, &response, telixProfileEndpoint, nil); err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateTelixProfile creates a new Telix export profile and returns the
// server-generated profile_id.
func (c *Client) CreateTelixProfile(ctx context.Context, request *TelixProfileCreateRequest) (string, error) {
	var response TelixProfileCreateResponse
	if err := c.PostAPIContext25(ctx, &response, telixProfileEndpoint, request); err != nil {
		return "", err
	}
	return response.ProfileID, nil
}

// GetTelixProfile retrieves a Telix export profile by profile_id. It returns
// ErrNotFound if the controller reports the profile does not exist.
func (c *Client) GetTelixProfile(ctx context.Context, profileID string) (*TelixProfileDetail, error) {
	endpoint, err := url.JoinPath(telixProfileEndpoint, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to construct endpoint: %w", err)
	}

	var detail TelixProfileDetail
	if err := c.GetAPIContext25(ctx, &detail, endpoint, nil); err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &detail, nil
}

// UpdateTelixProfile applies a partial update to an existing Telix export
// profile via PATCH. Only fields explicitly set on the request are updated;
// the telemetry sources list and the destination protocol are immutable per
// the API contract and cannot be changed via this call.
func (c *Client) UpdateTelixProfile(ctx context.Context, profileID string, request *TelixProfileUpdateRequest) error {
	endpoint, err := url.JoinPath(telixProfileEndpoint, profileID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}

	var response TelixProfileCreateResponse
	return c.PatchAPIContext25(ctx, &response, endpoint, request)
}

// DeleteTelixProfile removes the Telix export profile identified by
// profile_id.
func (c *Client) DeleteTelixProfile(ctx context.Context, profileID string) error {
	endpoint, err := url.JoinPath(telixProfileEndpoint, profileID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
