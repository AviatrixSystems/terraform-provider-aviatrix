package aviatrix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestSegmentationAssociationTransitScoped(t *testing.T) {
	tests := []struct {
		attachmentType string
		want           bool
	}{
		{attachmentType: "SPOKE", want: false},
		{attachmentType: "EDGESPOKE", want: false},
		{attachmentType: "", want: false},
		{attachmentType: "ARMSPOKE", want: true},
		{attachmentType: "ARMVNG", want: true},
		{attachmentType: "SITE2CLOUD", want: true},
		{attachmentType: "TGWDOMAIN", want: true},
		{attachmentType: "EDGEVLAN", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.attachmentType, func(t *testing.T) {
			assert.Equal(t, tt.want, segmentationAssociationTransitScoped(tt.attachmentType))
		})
	}
}

// domainAssociationController is the slice of the controller the association
// CustomizeDiff consults: the attachment list, plus gateway and group lookups
// used to tell a real transit move from a group-name alias.
type domainAssociationController struct {
	attachmentName string
	attachmentType string
	// apiName is the name the controller lists. Empty uses attachmentName.
	// EDGESPOKE listings are "<site>:<edge>" while config uses the site id.
	apiName        string
	domain         string
	groups         map[string]string
	gateways       map[string]bool
	listCalls      int
	gatewayLookups int
}

func (f *domainAssociationController) client(t *testing.T) *goaviatrix.Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		body, _ := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if form, err := url.ParseQuery(string(body)); err == nil {
			for k, v := range form {
				params[k] = v
			}
		}
		assert.NoError(t, json.NewEncoder(w).Encode(f.respond(t, params)))
	}))
	t.Cleanup(server.Close)
	client, err := goaviatrix.NewClient("admin", "password", strings.TrimPrefix(server.URL, "https://"), server.Client(), nil)
	require.NoError(t, err)
	return client
}

func (f *domainAssociationController) respond(t *testing.T, params url.Values) map[string]any {
	t.Helper()
	switch action := params.Get("action"); action {
	case "get_api_token":
		return map[string]any{"return": true, "results": map[string]any{"api_token": "token"}}
	case "login":
		return map[string]any{"return": true, "CID": "cid"}
	case "list_multi_cloud_domain_attachments":
		f.listCalls++
		name := f.apiName
		if name == "" {
			name = f.attachmentName
		}
		return map[string]any{"return": true, "results": map[string]any{
			"attachments": []map[string]any{{
				"name":         name,
				"domain":       f.domain,
				"transit_name": "transit-a",
				"type":         f.attachmentType,
			}},
		}}
	case "get_gateway_group_details_by_name":
		f.gatewayLookups++
		primary, ok := f.groups[params.Get("group_name")]
		if !ok {
			return map[string]any{"return": false, "reason": "not found"}
		}
		return map[string]any{"return": true, "results": map[string]any{
			"group_name":           params.Get("group_name"),
			"primary_gateway_name": primary,
		}}
	case "list_vpcs_summary":
		f.gatewayLookups++
		name := params.Get("gateway_name")
		if !f.gateways[name] {
			return map[string]any{"return": true, "results": []any{}}
		}
		return map[string]any{"return": true, "results": []map[string]any{{"vpc_name": name}}}
	default:
		t.Errorf("unexpected controller action %q", action)
		return map[string]any{"return": false, "reason": "unexpected action " + action}
	}
}

func domainAssociationConfig(attachment, transit string) map[string]any {
	return map[string]any{
		"network_domain_name":  "spoke2",
		"attachment_name":      attachment,
		"transit_gateway_name": transit,
	}
}

func TestSegmentationNetworkDomainAssociationTransitChange(t *testing.T) {
	const (
		spokeAttachment = "spoke-aws-1"
		armAttachment   = "azureAccount:vnet:rg:c0402aeb-2eba-413d-8aad-2a9d8a1c64d7"
		edgeSite        = "site-1"
	)
	groups := map[string]string{"transit-a-group": "transit-a"}
	gateways := map[string]bool{"transit-a": true, "transit-b": true}

	tests := []struct {
		name           string
		attachment     string
		attachmentType string
		apiName        string
		oldTransit     string
		newTransit     string
		wantReplace    bool
		wantList       bool
	}{
		{
			name:           "spoke gateway transit change is ignored",
			attachment:     spokeAttachment,
			attachmentType: "SPOKE",
			oldTransit:     "transit-a",
			newTransit:     "transit-b",
			wantList:       true,
		},
		{
			name:           "edge spoke transit change is ignored",
			attachment:     edgeSite,
			attachmentType: "EDGESPOKE",
			apiName:        edgeSite + ":edge-1",
			oldTransit:     "transit-a",
			newTransit:     "transit-b",
			wantList:       true,
		},
		{
			name:           "azure native spoke transit move replaces the association",
			attachment:     armAttachment,
			attachmentType: "ARMSPOKE",
			oldTransit:     "transit-a",
			newTransit:     "transit-b",
			wantReplace:    true,
			wantList:       true,
		},
		{
			name:           "azure native spoke group alias is not a move",
			attachment:     armAttachment,
			attachmentType: "ARMSPOKE",
			oldTransit:     "transit-a",
			newTransit:     "transit-a-group",
			wantList:       true,
		},
		{
			name:           "site2cloud transit move replaces the association",
			attachment:     "conn-1",
			attachmentType: "SITE2CLOUD",
			oldTransit:     "transit-a",
			newTransit:     "transit-b",
			wantReplace:    true,
			wantList:       true,
		},
		{
			name:           "unknown attachment type keeps the previous ignore behavior",
			attachment:     spokeAttachment,
			attachmentType: "",
			oldTransit:     "transit-a",
			newTransit:     "transit-b",
			wantList:       true,
		},
		{
			name:           "unchanged transit does not look up the association",
			attachment:     armAttachment,
			attachmentType: "ARMSPOKE",
			oldTransit:     "transit-a",
			newTransit:     "transit-a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &domainAssociationController{
				attachmentName: tt.attachment,
				attachmentType: tt.attachmentType,
				apiName:        tt.apiName,
				domain:         "spoke2",
				groups:         groups,
				gateways:       gateways,
			}
			client := controller.client(t)
			res := resourceAviatrixSegmentationNetworkDomainAssociation()
			state := schema.TestResourceDataRaw(t, res.Schema, domainAssociationConfig(tt.attachment, tt.oldTransit))
			state.SetId("spoke2~" + tt.attachment)

			diff, err := res.SimpleDiff(context.Background(), state.State(), terraform.NewResourceConfigRaw(domainAssociationConfig(tt.attachment, tt.newTransit)), client)
			require.NoError(t, err)
			require.NotNil(t, diff)
			assert.Equal(t, tt.wantReplace, diff.RequiresNew())
			if tt.wantList {
				assert.Equal(t, 1, controller.listCalls)
			} else {
				assert.Zero(t, controller.listCalls)
			}
			if !tt.wantReplace {
				attr, ok := diff.Attributes["transit_gateway_name"]
				if ok {
					assert.False(t, attr.RequiresNew)
					assert.Equal(t, attr.Old, attr.New)
				}
			} else {
				attr, ok := diff.Attributes["transit_gateway_name"]
				require.True(t, ok, "transit_gateway_name must be in the diff")
				assert.True(t, attr.RequiresNew)
				assert.Equal(t, tt.oldTransit, attr.Old)
				assert.Equal(t, tt.newTransit, attr.New)
			}
		})
	}
}

func TestSegmentationNetworkDomainAssociationCreateIgnoresTransitDiff(t *testing.T) {
	controller := &domainAssociationController{
		attachmentName: "spoke-aws-1",
		attachmentType: "SPOKE",
		domain:         "spoke2",
	}
	client := controller.client(t)
	res := resourceAviatrixSegmentationNetworkDomainAssociation()

	diff, err := res.SimpleDiff(context.Background(), nil, terraform.NewResourceConfigRaw(domainAssociationConfig("spoke-aws-1", "transit-a")), client)
	require.NoError(t, err)
	require.NotNil(t, diff)
	assert.Zero(t, controller.listCalls)
	assert.Zero(t, controller.gatewayLookups)
}
