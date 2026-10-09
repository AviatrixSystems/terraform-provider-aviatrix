package goaviatrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTgwName    = "tgw-east"
	testDomainName = "spoke-domain"
	testVpcID      = "vpc-0467685fc8fb28818"
)

// newAwsTgwVpcAttachmentStub answers each controller action with the given
// payload and records the form of the last attach_vpc_to_tgw submit.
func newAwsTgwVpcAttachmentStub(t *testing.T, responses map[string]any) (*Client, *map[string][]string) {
	t.Helper()
	var submitted map[string][]string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		assert.NoError(t, r.ParseForm())

		action := r.Form.Get("action")
		switch action {
		case "check_task_status":
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "done"}))
			return
		case "attach_vpc_to_tgw":
			submitted = r.Form
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": "req-1"}))
			return
		}
		results, ok := responses[action]
		if !assert.Truef(t, ok, "unexpected action %q", action) {
			return
		}
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"return": true, "results": results}))
	}))
	t.Cleanup(server.Close)
	client := &Client{
		HTTPClient: server.Client(),
		CID:        "test-cid",
		baseURL:    server.URL,
	}
	return client, &submitted
}

func TestCreateAwsTgwVpcAttachmentSendsApplianceMode(t *testing.T) {
	tests := []struct {
		name          string
		applianceMode bool
		wantSent      bool
	}{
		{name: "enabled", applianceMode: true, wantSent: true},
		{name: "disabled is omitted", applianceMode: false, wantSent: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, submitted := newAwsTgwVpcAttachmentStub(t, nil)

			require.NoError(t, client.CreateAwsTgwVpcAttachment(&AwsTgwVpcAttachment{
				TgwName:            testTgwName,
				Region:             "us-east-1",
				SecurityDomainName: testDomainName,
				VpcAccountName:     "aws-account",
				VpcID:              testVpcID,
				ApplianceMode:      tt.applianceMode,
			}))

			form := *submitted
			require.NotNil(t, form, "attach_vpc_to_tgw was not submitted")
			value, sent := form["appliance_mode"]
			assert.Equal(t, tt.wantSent, sent)
			if tt.wantSent {
				assert.Equal(t, []string{"true"}, value)
			}
		})
	}
}

func TestGetAwsTgwVpcAttachmentReadsApplianceMode(t *testing.T) {
	tests := []struct {
		name           string
		storedMode     any
		firewallDomain bool
		want           bool
	}{
		{name: "enabled", storedMode: "enable", want: true},
		{name: "disabled", storedMode: "disable", want: false},
		// Attachments created before the controller stored the value.
		{name: "not stored", storedMode: nil, want: false},
		{name: "firewall domain ignores the stored value", storedMode: "enable", firewallDomain: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attachment := map[string]any{
				"vpc_id":                       testVpcID,
				"acct_name":                    "aws-account",
				"region":                       "us-east-1",
				"associated_route_domain_name": testDomainName,
			}
			if tt.storedMode != nil {
				attachment["appliance_mode_support"] = tt.storedMode
			}
			client, _ := newAwsTgwVpcAttachmentStub(t, map[string]any{
				"get_tgw_attachment_details": []any{attachment},
				"list_route_domain_names":    []string{testDomainName},
				"view_route_domain_details": []any{map[string]any{
					"name":            testDomainName,
					"firewall_domain": tt.firewallDomain,
					"attached_vpc": []any{map[string]any{
						"vpc_id":       testVpcID,
						"account_name": "aws-account",
					}},
				}},
				"list_attachment_route_table_details": map[string]any{"vpc_id": testVpcID},
			})

			got, err := client.GetAwsTgwVpcAttachment(&AwsTgwVpcAttachment{
				TgwName:            testTgwName,
				SecurityDomainName: testDomainName,
				VpcID:              testVpcID,
			})

			require.NoError(t, err)
			assert.Equal(t, tt.want, got.ApplianceMode)
		})
	}
}
