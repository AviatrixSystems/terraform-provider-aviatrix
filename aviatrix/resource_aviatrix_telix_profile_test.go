package aviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	testTelixCAPem         = "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJAOZHJC...\n-----END CERTIFICATE-----"
	testTelixClientCertPem = "-----BEGIN CERTIFICATE-----\nMIICfTCCAWWgAwIBAgI...\n-----END CERTIFICATE-----"
	testTelixClientKeyPem  = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqh...\n-----END PRIVATE KEY-----"
)

// testTelixAccClientCertPEM and testTelixAccClientKeyPEM are a matching RSA
// certificate and private key pair used only by TestAccAviatrixTelixProfile_lifecycle
// so the controller's PEM parsing and tls.X509KeyPair validation succeed.
const (
	testTelixAccClientCertPEM = `-----BEGIN CERTIFICATE-----
MIIDTzCCAjegAwIBAgIUNpHz3Nxs+LkIMu7AYm+AOK0chrIwDQYJKoZIhvcNAQEL
BQAwNzEeMBwGA1UEAwwVdGVsaXgtYWNjLXRlc3QtY2xpZW50MRUwEwYDVQQKDAxh
dmlhdHJpeC1hY2MwHhcNMjYwNTEyMjIzOTQyWhcNMzYwNTA5MjIzOTQyWjA3MR4w
HAYDVQQDDBV0ZWxpeC1hY2MtdGVzdC1jbGllbnQxFTATBgNVBAoMDGF2aWF0cml4
LWFjYzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAIdH7kIjcJ+D+/wg
ZY5r/H6MTe20TSgtcc2Fx+INpO5BI7a4/EQEgHx19PPGnFuM3gBxv+cVjIZS0ng6
JrgtEps2YRzVNkNjUPSl6eM6rRsoqa+0XvThso9ndqQhfecEX3dGPnGc+KRUFvQm
GDt03KSX9FG0oO23/W8SXeATZAm/GoZdJDvixFsnSDRQq2fdqWCds7g+ekuI01ht
ugSuBC8ChSrEpxKN/UjeCvNsTzPEFhYfVoHxIbuIUIEfHxR1FicWirdE/IUV1z5G
jQEUipvPpvHR44I8gIaSM/TvispenEA3wwdrqiC0jJKM5bOJ3Kq+pKfRXp2T1/j+
Xl8moEMCAwEAAaNTMFEwHQYDVR0OBBYEFE7X43kgJkz9jFhcCUJnc12vo1QqMB8G
A1UdIwQYMBaAFE7X43kgJkz9jFhcCUJnc12vo1QqMA8GA1UdEwEB/wQFMAMBAf8w
DQYJKoZIhvcNAQELBQADggEBAGNSq+7uf17T9XibPwDXGg9Dlo46SsP24lyzzwHt
hIIlGvUABzc0p2hfTubzjId6isDHPz0iOYjbByCPExaChqoOsrW3PfFZtlioB6ym
mBfaT2Q5qIiaiyXJVSg1cYX16Ye2ubLmy5hpRgoaHblDco/XVotVYtArlKqdW8pg
uziXceKZ+z/i42MxWRuuZP6KMv3w26qBsNM6JBFu1ZOohi+LE7dGv/uciZXzk2BK
ppsEokWqlYpr2hcPCnaMO1bmyB41VMkh8YcqLls/RElHvhVRuLAOTEwc0Ggdz+MI
uRr/SVO7666g+hS4tXG2ebjdZi0mFWaiASA/qADMaHi8eiM=
-----END CERTIFICATE-----`

	testTelixAccClientKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCHR+5CI3Cfg/v8
IGWOa/x+jE3ttE0oLXHNhcfiDaTuQSO2uPxEBIB8dfTzxpxbjN4Acb/nFYyGUtJ4
Oia4LRKbNmEc1TZDY1D0penjOq0bKKmvtF704bKPZ3akIX3nBF93Rj5xnPikVBb0
Jhg7dNykl/RRtKDtt/1vEl3gE2QJvxqGXSQ74sRbJ0g0UKtn3algnbO4PnpLiNNY
bboErgQvAoUqxKcSjf1I3grzbE8zxBYWH1aB8SG7iFCBHx8UdRYnFoq3RPyFFdc+
Ro0BFIqbz6bx0eOCPICGkjP074rKXpxAN8MHa6ogtIySjOWzidyqvqSn0V6dk9f4
/l5fJqBDAgMBAAECggEAFHwVaClyqozzUWcxyaJbc77Y0WJIfrQm8+lKmQCE6tkC
ne7JWvQC5Wsn15QfEBPheALBfYNQqfRiT3dKw//Qk0qQOQFQ0V0kIJpRCNSqlXQP
rH4CalAU2EyL9Wg8PDOi0G5avUcjisqlm9HQQIyTw/ld3dNyvPHBQW+atyK94WW4
3HqcZNXb7UOrdY8ueCjvk0U8TdxqaOIJhaTMldLTdTeanwcwWZX0zn5NMZiYGHyH
WfHoTjj1+9u48qbDE4o9VEbd5v/9Tg3hb8optUw8M2GPTZK2hjFkPSLSVGLuE0M2
RAV38xS5nR/dBuEoj/bp7JyBtucW8bUQ+AJs2Po/UQKBgQC9VS0dWLu14XFFW3kE
NjThwg6/Q3RF/oMQ2N15bDLSW6Xm+D+JhWKn6i9P8A9YMo+aFxq3kY9l32wgrwai
LYrUAOagfgg98SznDZIOqy0Igch5yNxialcL7jSCg05JBbgbGF8ylrLnxlKUC7TZ
6S5M5o/NWaK9zvZWIhM9rjR3EQKBgQC26myc/EwwKHz+wJyPq/MxN1wz5D4AhQYT
eRxBv77qgX670joDGlvBma0iiX0kGn2vof9zZ3Ffwy5bWM7eEc+O4EcnJ213VMBE
GhDp67o7KrzIf6o0xK3U1juln6uA3rgMuLRRDYM1kpgE5gLPPRep76xfCBkgXgoF
FNvKhCkqEwKBgGz4kTbK038jemZI96YM7PLjFknPMST4D8eqig5Q0A9y4FHHoAou
01GB4ClKKgrBTxWJJr9w7+/aYAmPs2m0fKr4ucS1xVihbw6tKNt4ejrjN9egW/fo
7KDZQS+9E1nECOrPZDthsSblZrH+6uBg7V0ldq7iYGCOtgltI1Xk5h2BAoGALGLA
TmvOlRUOF8dndlmUXsn/PrxQ61FcQxdtaur7ie44cZ025I/d2iHPaIUSb9NZ0meu
FDPyx/kV46auNCcARbxYp8CiiIVxTlVA63J/M2JQgxqvk7RyNiZyPON8+32QDc44
Oz7bKwHSj8W8wsshVeRJ4JmXd0o6hjckioT9dC8CgYEAm/bud7mDfNJxxrJjO5kB
+8Zn3R4c/oXNaib0wogeoNVEZb1KT14T8rYscyqefbmQhPk5FjiDcuue5u9wOSIS
ET6urH/qYodtb/UdRN3cMYB/UVQhdaOIBs2XVbNRi6WdeVdF6d/5p0QIPMJjfHv3
Db3WD9urR8PMoIjykbgEES8=
-----END PRIVATE KEY-----`
)

// telixProfileSchema returns the schema map under test. Centralized so each
// test case can build a *schema.ResourceData via schema.TestResourceDataRaw.
func telixProfileSchema() map[string]*schema.Schema {
	return resourceAviatrixTelixProfile().Schema
}

func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }

func ptrOtlpHeaders(m map[string]string) *map[string]string {
	mm := m
	return &mm
}

func TestMarshalTelixProfileCreateRequest(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]any
		want *goaviatrix.TelixProfileCreateRequest
	}{
		{
			name: "minimal_grpc_all_gateways",
			in: map[string]any{
				"display_name": "minimal",
				"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
				"enabled":      true,
				"destination": []any{
					map[string]any{
						"otlp": []any{
							map[string]any{
								"endpoint": "otel.example.com:4317",
								"protocol": goaviatrix.TelixOtlpProtocolGRPC,
							},
						},
					},
				},
				"scope": []any{
					map[string]any{
						"all_gateways": []any{map[string]any{}},
					},
				},
			},
			want: &goaviatrix.TelixProfileCreateRequest{
				DisplayName: "minimal",
				Sources:     []string{goaviatrix.TelixTelemetrySourceDcfLogs},
				Enabled:     ptrBool(true),
				Destination: goaviatrix.TelixDestinationInput{
					Otlp: goaviatrix.TelixOtlpDestinationInput{
						Endpoint: "otel.example.com:4317",
						Protocol: goaviatrix.TelixOtlpProtocolGRPC,
					},
				},
				Scope: goaviatrix.TelixGatewayScope{
					AllGateways: &goaviatrix.TelixAllGatewaysScope{},
				},
			},
		},
		{
			name: "selected_gateways_with_filters_and_tls",
			in: map[string]any{
				"display_name": "secure",
				"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
				"enabled":      true,
				"destination": []any{
					map[string]any{
						"otlp": []any{
							map[string]any{
								"endpoint": "https://collector/v1/logs",
								"protocol": goaviatrix.TelixOtlpProtocolHTTP,
								"headers": map[string]any{
									"X-Auth": "tok",
								},
								"tls": []any{
									map[string]any{
										"ca_certificate_pem":     testTelixCAPem,
										"client_certificate_pem": testTelixClientCertPem,
										"client_private_key_pem": testTelixClientKeyPem,
										"server_name_override":   "collector",
										"insecure_skip_verify":   false,
									},
								},
							},
						},
					},
				},
				"scope": []any{
					map[string]any{
						"selected_gateways": []any{
							map[string]any{
								"gateway_names": []any{"gw-a", "gw-b"},
							},
						},
					},
				},
				"filters": []any{
					map[string]any{
						"dcf_log_types": []any{"FLOW", "INTRUSION"},
						"dcf_actions":   []any{"ALLOW"},
					},
				},
			},
			want: &goaviatrix.TelixProfileCreateRequest{
				DisplayName: "secure",
				Sources:     []string{goaviatrix.TelixTelemetrySourceDcfLogs},
				Enabled:     ptrBool(true),
				Destination: goaviatrix.TelixDestinationInput{
					Otlp: goaviatrix.TelixOtlpDestinationInput{
						Endpoint: "https://collector/v1/logs",
						Protocol: goaviatrix.TelixOtlpProtocolHTTP,
						Headers:  ptrOtlpHeaders(map[string]string{"X-Auth": "tok"}),
						TLS: &goaviatrix.TelixTLSInputConfig{
							CaCertificatePem:     ptrString(testTelixCAPem),
							ClientCertificatePem: ptrString(testTelixClientCertPem),
							ClientPrivateKeyPem:  ptrString(testTelixClientKeyPem),
							ServerNameOverride:   ptrString("collector"),
							InsecureSkipVerify:   ptrBool(false),
						},
					},
				},
				Scope: goaviatrix.TelixGatewayScope{
					SelectedGateways: &goaviatrix.TelixSelectedGatewaysScope{
						GatewayNames: []string{"gw-a", "gw-b"},
					},
				},
				Filters: &goaviatrix.TelixFilterConfig{
					DcfLogTypes: []string{"FLOW", "INTRUSION"},
					DcfActions:  []string{"ALLOW"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, telixProfileSchema(), tc.in)
			got := marshalTelixProfileCreateRequest(d)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("marshalTelixProfileCreateRequest mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExpandTelixScopeValue(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]any
		want goaviatrix.TelixGatewayScope
	}{
		{
			name: "all_gateways",
			in: map[string]any{
				"display_name": "x",
				"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
				"destination": []any{map[string]any{
					"otlp": []any{map[string]any{
						"endpoint": "e",
						"protocol": goaviatrix.TelixOtlpProtocolGRPC,
					}},
				}},
				"scope": []any{map[string]any{
					"all_gateways": []any{map[string]any{}},
				}},
			},
			want: goaviatrix.TelixGatewayScope{
				AllGateways: &goaviatrix.TelixAllGatewaysScope{},
			},
		},
		{
			name: "selected_gateways",
			in: map[string]any{
				"display_name": "x",
				"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
				"destination": []any{map[string]any{
					"otlp": []any{map[string]any{
						"endpoint": "e",
						"protocol": goaviatrix.TelixOtlpProtocolGRPC,
					}},
				}},
				"scope": []any{map[string]any{
					"selected_gateways": []any{map[string]any{
						"gateway_names": []any{"a", "b", "c"},
					}},
				}},
			},
			want: goaviatrix.TelixGatewayScope{
				SelectedGateways: &goaviatrix.TelixSelectedGatewaysScope{
					GatewayNames: []string{"a", "b", "c"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, telixProfileSchema(), tc.in)
			got := expandTelixScopeValue(d)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("expandTelixScopeValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExpandTelixFilters(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]any
		want *goaviatrix.TelixFilterConfig
	}{
		{
			name: "absent",
			in: map[string]any{
				"display_name": "x",
				"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
				"destination": []any{map[string]any{
					"otlp": []any{map[string]any{
						"endpoint": "e",
						"protocol": goaviatrix.TelixOtlpProtocolGRPC,
					}},
				}},
				"scope": []any{map[string]any{
					"all_gateways": []any{map[string]any{}},
				}},
			},
			want: nil,
		},
		{
			name: "log_types_only",
			in: map[string]any{
				"display_name": "x",
				"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
				"destination": []any{map[string]any{
					"otlp": []any{map[string]any{
						"endpoint": "e",
						"protocol": goaviatrix.TelixOtlpProtocolGRPC,
					}},
				}},
				"scope": []any{map[string]any{
					"all_gateways": []any{map[string]any{}},
				}},
				"filters": []any{map[string]any{
					"dcf_log_types": []any{"FLOW"},
				}},
			},
			want: &goaviatrix.TelixFilterConfig{
				DcfLogTypes: []string{"FLOW"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, telixProfileSchema(), tc.in)
			got := expandTelixFilters(d)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("expandTelixFilters mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTelixOtlpTLSListAbsent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    any
		want bool
	}{
		{name: "nil", v: nil, want: true},
		{name: "empty_slice", v: []any{}, want: true},
		{name: "slice_nil_elem", v: []any{nil}, want: true},
		{name: "block_present_empty_map", v: []any{map[string]any{}}, want: false},
		{name: "block_present_ca", v: []any{map[string]any{"ca_certificate_pem": "x"}}, want: false},
		{name: "wrong_slice_elm_type", v: []string{}, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := telixOtlpTLSListAbsent(tc.v)
			if got != tc.want {
				t.Fatalf("telixOtlpTLSListAbsent(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTelixOtlpTLSBlockRemoval(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		hasTLSChange   bool
		oldV, newV     any
		wantRemovalTag bool // short name avoids "want" clash
	}{
		{
			name:           "no_change_flag",
			hasTLSChange:   false,
			oldV:           []any{map[string]any{}},
			newV:           []any{},
			wantRemovalTag: false,
		},
		{
			name:           "removed_tls_block",
			hasTLSChange:   true,
			oldV:           []any{map[string]any{"ca_certificate_pem": "pem"}},
			newV:           []any{},
			wantRemovalTag: true,
		},
		{
			name:           "added_tls_block",
			hasTLSChange:   true,
			oldV:           []any{},
			newV:           []any{map[string]any{}},
			wantRemovalTag: false,
		},
		{
			name:           "both_absent_always_false",
			hasTLSChange:   true,
			oldV:           []any{},
			newV:           []any{},
			wantRemovalTag: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := telixOtlpTLSBlockRemoval(tc.hasTLSChange, tc.oldV, tc.newV)
			if got != tc.wantRemovalTag {
				t.Fatalf("telixOtlpTLSBlockRemoval(...) = %v, want %v", got, tc.wantRemovalTag)
			}
		})
	}
}

func TestTelixOtlpTLSClearPayloadMarshals(t *testing.T) {
	payload := telixOtlpTLSClearPayload()
	dest := goaviatrix.TelixDestinationInput{
		Otlp: goaviatrix.TelixOtlpDestinationInput{
			Endpoint: "otel.example.com:4317",
			Protocol: goaviatrix.TelixOtlpProtocolGRPC,
			TLS:      payload,
		},
	}

	b, err := json.Marshal(dest)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, sub := range []string{
		`"ca_certificate_pem":""`,
		`"client_certificate_pem":""`,
		`"client_private_key_pem":""`,
		`"server_name_override":""`,
		`"insecure_skip_verify":false`,
	} {
		if !strings.Contains(s, sub) {
			t.Fatalf("marshaled PATCH TLS clear missing %q\nfull: %s", sub, s)
		}
	}
}

// TestFlattenTelixProfileDetail_PreservesSecretsInState is the most important
// test in this file. If it ever fails, Approach 1's contract is broken: Read
// would overwrite secrets in state with redacted/empty values from the API,
// which would either cause perpetual diffs or silently lose user-supplied
// material on every refresh.
func TestFlattenTelixProfileDetail_PreservesSecretsInState(t *testing.T) {
	initial := map[string]any{
		"display_name": "preserve",
		"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
		"destination": []any{map[string]any{
			"otlp": []any{map[string]any{
				"endpoint": "otel.example.com:4317",
				"protocol": goaviatrix.TelixOtlpProtocolGRPC,
				"headers":  map[string]any{"X-Auth": "stored-token"},
				"tls": []any{map[string]any{
					"ca_certificate_pem":     testTelixCAPem,
					"client_certificate_pem": testTelixClientCertPem,
					"client_private_key_pem": testTelixClientKeyPem,
				}},
			}},
		}},
		"scope": []any{map[string]any{
			"all_gateways": []any{map[string]any{}},
		}},
	}

	d := schema.TestResourceDataRaw(t, telixProfileSchema(), initial)

	// API response: presence flags are true, no actual material returned.
	detail := &goaviatrix.TelixProfileDetail{
		ProfileID:   "abc-123",
		DisplayName: "preserve",
		Enabled:     true,
		Sources:     []string{goaviatrix.TelixTelemetrySourceDcfLogs},
		Destination: goaviatrix.TelixDestination{
			Otlp: goaviatrix.TelixOtlpDestination{
				Endpoint:   "otel.example.com:4317",
				Protocol:   goaviatrix.TelixOtlpProtocolGRPC,
				HasHeaders: ptrBool(true),
				TLS: &goaviatrix.TelixTLSConfig{
					HasCaCertificate:     ptrBool(true),
					HasClientCertificate: ptrBool(true),
					HasClientPrivateKey:  ptrBool(true),
				},
			},
		},
		Scope: goaviatrix.TelixGatewayScope{
			AllGateways: &goaviatrix.TelixAllGatewaysScope{},
		},
	}

	if err := flattenTelixProfileDetail(d, detail); err != nil {
		t.Fatalf("flattenTelixProfileDetail: %v", err)
	}

	// Headers must be preserved.
	v := d.Get("destination.0.otlp.0.headers")
	gotHeaders, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("destination.0.otlp.0.headers: expected map in state, got %T (%v)", v, v)
	}
	if gotHeaders["X-Auth"] != "stored-token" {
		t.Fatalf("headers were overwritten: got %v, want preserved 'stored-token'", gotHeaders)
	}

	// PEMs must be preserved.
	for path, want := range map[string]string{
		"destination.0.otlp.0.tls.0.ca_certificate_pem":     testTelixCAPem,
		"destination.0.otlp.0.tls.0.client_certificate_pem": testTelixClientCertPem,
		"destination.0.otlp.0.tls.0.client_private_key_pem": testTelixClientKeyPem,
	} {
		v := d.Get(path)
		got, ok := v.(string)
		if !ok {
			t.Fatalf("%s: expected string in state, got %T (%v)", path, v, v)
		}
		if got != want {
			t.Fatalf("%s was overwritten: got %q, want preserved value", path, got)
		}
	}

	// has_* mirrors must be set to the API values.
	for _, path := range []string{
		"destination.0.otlp.0.has_headers",
		"destination.0.otlp.0.tls.0.has_ca_certificate",
		"destination.0.otlp.0.tls.0.has_client_certificate",
		"destination.0.otlp.0.tls.0.has_client_private_key",
	} {
		v := d.Get(path)
		got, ok := v.(bool)
		if !ok {
			t.Errorf("%s: expected bool, got %T (%v)", path, v, v)
			continue
		}
		if !got {
			t.Errorf("%s should be true (API said present), got false", path)
		}
	}
}

// TestFlattenTelixProfileDetail_ExternalDeletion verifies that when the
// controller no longer holds a sensitive value (has_* flipped to false)
// while state still has the user-supplied value, Read clears the state slot
// so the next plan re-pushes from configuration.
func TestFlattenTelixProfileDetail_ExternalDeletion(t *testing.T) {
	initial := map[string]any{
		"display_name": "deleted-externally",
		"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
		"destination": []any{map[string]any{
			"otlp": []any{map[string]any{
				"endpoint": "otel.example.com:4317",
				"protocol": goaviatrix.TelixOtlpProtocolGRPC,
				"headers":  map[string]any{"X-Auth": "stored-token"},
				"tls": []any{map[string]any{
					"ca_certificate_pem":     testTelixCAPem,
					"client_certificate_pem": testTelixClientCertPem,
					"client_private_key_pem": testTelixClientKeyPem,
				}},
			}},
		}},
		"scope": []any{map[string]any{
			"all_gateways": []any{map[string]any{}},
		}},
	}

	d := schema.TestResourceDataRaw(t, telixProfileSchema(), initial)

	// API now reports all has_* as false (admin cleared the values via UI).
	detail := &goaviatrix.TelixProfileDetail{
		ProfileID:   "abc-123",
		DisplayName: "deleted-externally",
		Enabled:     true,
		Sources:     []string{goaviatrix.TelixTelemetrySourceDcfLogs},
		Destination: goaviatrix.TelixDestination{
			Otlp: goaviatrix.TelixOtlpDestination{
				Endpoint:   "otel.example.com:4317",
				Protocol:   goaviatrix.TelixOtlpProtocolGRPC,
				HasHeaders: ptrBool(false),
				TLS: &goaviatrix.TelixTLSConfig{
					HasCaCertificate:     ptrBool(false),
					HasClientCertificate: ptrBool(false),
					HasClientPrivateKey:  ptrBool(false),
				},
			},
		},
		Scope: goaviatrix.TelixGatewayScope{
			AllGateways: &goaviatrix.TelixAllGatewaysScope{},
		},
	}

	if err := flattenTelixProfileDetail(d, detail); err != nil {
		t.Fatalf("flattenTelixProfileDetail: %v", err)
	}

	// Headers must be cleared.
	hv := d.Get("destination.0.otlp.0.headers")
	gotHeaders, ok := hv.(map[string]any)
	if !ok {
		t.Fatalf("destination.0.otlp.0.headers: expected map in state, got %T (%v)", hv, hv)
	}
	if len(gotHeaders) != 0 {
		t.Fatalf("headers should have been cleared, got %v", gotHeaders)
	}

	// PEMs must be cleared.
	for _, path := range []string{
		"destination.0.otlp.0.tls.0.ca_certificate_pem",
		"destination.0.otlp.0.tls.0.client_certificate_pem",
		"destination.0.otlp.0.tls.0.client_private_key_pem",
	} {
		v := d.Get(path)
		got, ok := v.(string)
		if !ok {
			t.Fatalf("%s: expected string in state, got %T (%v)", path, v, v)
		}
		if got != "" {
			t.Errorf("%s should have been cleared, got %q", path, got)
		}
	}
}

// TestFlattenTelixProfileDetail_NonSensitiveDriftDetected ensures normal
// drift detection keeps working for non-sensitive fields. If an admin
// changes display_name in the controller, Read writes the new value into
// state so the next plan shows a diff against config.
func TestFlattenTelixProfileDetail_NonSensitiveDriftDetected(t *testing.T) {
	initial := map[string]any{
		"display_name": "old-name",
		"sources":      []any{goaviatrix.TelixTelemetrySourceDcfLogs},
		"destination": []any{map[string]any{
			"otlp": []any{map[string]any{
				"endpoint": "otel.example.com:4317",
				"protocol": goaviatrix.TelixOtlpProtocolGRPC,
			}},
		}},
		"scope": []any{map[string]any{
			"all_gateways": []any{map[string]any{}},
		}},
	}
	d := schema.TestResourceDataRaw(t, telixProfileSchema(), initial)

	detail := &goaviatrix.TelixProfileDetail{
		ProfileID:   "abc",
		DisplayName: "renamed-via-ui",
		Enabled:     true,
		Sources:     []string{goaviatrix.TelixTelemetrySourceDcfLogs},
		Destination: goaviatrix.TelixDestination{
			Otlp: goaviatrix.TelixOtlpDestination{
				Endpoint: "otel.example.com:4317",
				Protocol: goaviatrix.TelixOtlpProtocolGRPC,
			},
		},
		Scope: goaviatrix.TelixGatewayScope{
			AllGateways: &goaviatrix.TelixAllGatewaysScope{},
		},
	}

	if err := flattenTelixProfileDetail(d, detail); err != nil {
		t.Fatalf("flattenTelixProfileDetail: %v", err)
	}

	v := d.Get("display_name")
	got, ok := v.(string)
	if !ok {
		t.Fatalf("display_name: expected string in state, got %T (%v)", v, v)
	}
	if got != "renamed-via-ui" {
		t.Fatalf("display_name drift not reflected: got %q want renamed-via-ui", got)
	}
}

// TestAccAviatrixTelixProfile_lifecycle exercises create, in-place update,
// non-sensitive TLS field toggle, and destroy for an aviatrix_telix_profile
// resource. The test framework's implicit "no changes after apply" check
// also enforces the load-bearing rule that Read does not produce a spurious
// diff against the same configuration.
func TestAccAviatrixTelixProfile_lifecycle(t *testing.T) {
	if os.Getenv("SKIP_TELIX_PROFILE") == "yes" {
		t.Skip("Skipping Telix profile test as SKIP_TELIX_PROFILE is set")
	}
	resourceName := "aviatrix_telix_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckTelixProfileDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with full TLS material and headers.
			{
				Config: testAccTelixProfileWithTLS("test-telix-profile", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTelixProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-telix-profile"),
					resource.TestCheckResourceAttrSet(resourceName, "profile_id"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.has_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_ca_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_client_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_client_private_key", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.insecure_skip_verify", "false"),
				),
			},
			// Step 2: Rename the profile without changing any TLS field.
			// Update path should send only display_name and leave secrets in
			// state untouched. has_* flags must stay true (server still has
			// the material).
			{
				Config: testAccTelixProfileWithTLS("test-telix-profile-renamed", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTelixProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-telix-profile-renamed"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.has_headers", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_ca_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_client_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_client_private_key", "true"),
				),
			},
			// Step 3: Toggle a non-sensitive TLS sub-field (insecure_skip_verify).
			// Tests that a TLS update does not require re-pushing secret
			// material from config; the patch only contains the changed sub-
			// field. has_* flags remain true.
			{
				Config: testAccTelixProfileWithTLS("test-telix-profile-renamed", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTelixProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.insecure_skip_verify", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_ca_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_client_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination.0.otlp.0.tls.0.has_client_private_key", "true"),
				),
			},
		},
	})
}

func testAccTelixProfileWithTLS(displayName string, insecure bool) string {
	return fmt.Sprintf(`
resource "aviatrix_telix_profile" "test" {
  display_name = %q
  sources      = ["TELEMETRY_SOURCE_DCF_LOGS"]
  enabled      = true

  destination {
    otlp {
      endpoint = "otel-collector.example.com:4317"
      protocol = "TELIX_OTLP_PROTOCOL_GRPC"

      headers = {
        "X-Auth-Token" = "tok-1"
      }

      tls {
        ca_certificate_pem     = %q
        client_certificate_pem = %q
        client_private_key_pem = %q
        server_name_override   = "collector.example.com"
        insecure_skip_verify   = %t
      }
    }
  }

  scope {
    all_gateways {}
  }
}
`, displayName, testCertificateContent(), testTelixAccClientCertPEM, testTelixAccClientKeyPEM, insecure)
}

func testAccCheckTelixProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Telix profile resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Telix profile ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		detail, err := client.GetTelixProfile(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch Telix profile %s: %w", rs.Primary.ID, err)
		}
		if detail.ProfileID != rs.Primary.ID {
			return fmt.Errorf("Telix profile ID mismatch: %s vs %s", detail.ProfileID, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckTelixProfileDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_telix_profile" {
			continue
		}
		_, err := client.GetTelixProfile(context.Background(), rs.Primary.ID)
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("Telix profile %s still exists when it should be destroyed", rs.Primary.ID)
		}
	}
	return nil
}
