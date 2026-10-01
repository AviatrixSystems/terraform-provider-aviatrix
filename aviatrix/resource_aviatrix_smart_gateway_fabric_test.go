package aviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestResourceAviatrixSmartGatewayFabricSchema(t *testing.T) {
	s := resourceAviatrixSmartGatewayFabric().Schema
	// Same attribute names as aviatrix_config_feature.
	assert.Equal(t, "feature_name", smartGatewayFabricFeatureNameAttr)
	assert.Equal(t, "is_enabled", smartGatewayFabricIsEnabledAttr)

	feature := s[smartGatewayFabricFeatureNameAttr]
	require.NotNil(t, feature)
	assert.True(t, feature.Required)
	assert.True(t, feature.ForceNew, "changing the feature means a different fabric instance")

	enabled := s[smartGatewayFabricIsEnabledAttr]
	require.NotNil(t, enabled)
	assert.True(t, enabled.Required)
	assert.False(t, enabled.ForceNew, "flipping is_enabled must be an update, not a re-create")

	for _, gone := range []string{"feature", "enabled", "underlay_mesh_enabled", "route_resolver_enabled", "rollout_complete"} {
		_, present := s[gone]
		assert.False(t, present, "%s must not be on the schema", gone)
	}
}

func TestSmartGatewayResolverAttributeOnGatewayResources(t *testing.T) {
	assert.Equal(t, "enable_route_resolver", smartGatewayResolverAttr)
	for name, res := range map[string]*schema.Resource{
		"spoke_gateway":   resourceAviatrixSpokeGateway(),
		"transit_gateway": resourceAviatrixTransitGateway(),
		"spoke_group":     resourceAviatrixSpokeGroup(),
		"transit_group":   resourceAviatrixTransitGroup(),
	} {
		attr, ok := res.Schema[smartGatewayResolverAttr]
		require.True(t, ok, "%s: missing %s", name, smartGatewayResolverAttr)
		assert.True(t, attr.Optional, "%s: resolver must be Optional", name)
		assert.True(t, attr.Computed, "%s: resolver must be Computed so an unset value keeps the controller's", name)
		assert.NotNil(t, res.CustomizeDiff, "%s: resolver needs plan-time validation", name)

		// The old names must not linger anywhere on these resources.
		for _, dead := range []string{"enable_smart_gw_resolver", "enable_smart_gateway_resolver"} {
			_, present := res.Schema[dead]
			assert.False(t, present, "%s: %s must have been renamed to enable_route_resolver", name, dead)
		}
	}
}

func TestSmartGatewayFabricFeatures(t *testing.T) {
	assert.Equal(t, []string{"underlay_mesh"}, smartGatewayFabricFeatures)

	validate := resourceAviatrixSmartGatewayFabric().Schema[smartGatewayFabricFeatureNameAttr].ValidateFunc
	_, errs := validate("underlay_mesh", smartGatewayFabricFeatureNameAttr)
	assert.Empty(t, errs)
	// The route resolver is set per gateway with enable_route_resolver, not fabric-wide.
	_, errs = validate("route_resolver", smartGatewayFabricFeatureNameAttr)
	assert.NotEmpty(t, errs)
}

func TestSmartGatewayFabricCheckFeature(t *testing.T) {
	assert.NoError(t, smartGatewayFabricCheckFeature(smartGatewayFabricUnderlayMesh))
	assert.Error(t, smartGatewayFabricCheckFeature("route_resolver"))
	assert.ErrorContains(t, smartGatewayFabricCheckFeature("nope"), "unknown feature")
}

func TestFilterSmartGatewayEligible(t *testing.T) {
	fabric := []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "transit-b"},
		{GatewayGroupName: "spoke-bgp"},
		{GatewayGroupName: "spoke-csp", IsCloudSpokeWithoutBgpLu: true},
		{GatewayGroupName: "edge-spoke"},
		{GatewayGroupName: "edge-transit"},
		{GatewayGroupName: "standalone"},
		{GatewayGroupName: "vpn"},
		{GatewayGroupName: "untyped"},
	}
	types := map[string]string{
		"transit-b":    "GwGroupType.TRANSIT",
		"spoke-bgp":    "SPOKE",
		"spoke-csp":    "SPOKE",
		"edge-spoke":   "GwGroupType.EDGESPOKE",
		"edge-transit": "EDGETRANSIT",
		"standalone":   "GwGroupType.STANDALONE",
		"vpn":          "VPN",
	}

	got := filterSmartGatewayEligible(fabric, types)

	names := make([]string, 0, len(got))
	for _, g := range got {
		names = append(names, g.GatewayGroupName)
	}
	assert.Equal(t, []string{"edge-spoke", "edge-transit", "spoke-bgp", "transit-b"}, names)
}

func TestSmartGatewayFabricEnabledFromFabric(t *testing.T) {
	on := goaviatrix.SmartGatewayGroupStatus{SmartGatewayUnderlayEnabled: true}
	off := goaviatrix.SmartGatewayGroupStatus{}
	all := []goaviatrix.SmartGatewayGroupStatus{on, on}
	none := []goaviatrix.SmartGatewayGroupStatus{off, off}
	mixed := []goaviatrix.SmartGatewayGroupStatus{on, off}

	assert.True(t, smartGatewayFabricEnabledFromFabric(nil, smartGatewayUnderlayOn, true), "no groups keeps the prior value")
	assert.True(t, smartGatewayFabricEnabledFromFabric(all, smartGatewayUnderlayOn, false))
	assert.False(t, smartGatewayFabricEnabledFromFabric(none, smartGatewayUnderlayOn, true))
	// A mixed fabric must never match the configured value, or the drift is silent.
	assert.True(t, smartGatewayFabricEnabledFromFabric(mixed, smartGatewayUnderlayOn, false), "config false, one group on: must diff")
	assert.False(t, smartGatewayFabricEnabledFromFabric(mixed, smartGatewayUnderlayOn, true), "config true, one group off: must diff")
}

func TestSmartGatewayFabricSweepSkipsSettledGroups(t *testing.T) {
	groups := []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "a", SmartGatewayUnderlayEnabled: true},
		{GatewayGroupName: "b"},
		{GatewayGroupName: "c"},
	}
	var called []string
	err := smartGatewayFabricSweep(context.Background(), groups, "underlay", smartGatewayUnderlayOn, true, func(_ context.Context, g string, enable bool) error {
		assert.True(t, enable)
		called = append(called, g)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"b", "c"}, called)
}

func TestSmartGatewayFabricSweepReportsEveryFailure(t *testing.T) {
	groups := []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "a", SmartGatewayUnderlayEnabled: true},
		{GatewayGroupName: "b", SmartGatewayUnderlayEnabled: true},
		{GatewayGroupName: "c", SmartGatewayUnderlayEnabled: true},
	}
	var called []string
	err := smartGatewayFabricSweep(context.Background(), groups, "underlay", smartGatewayUnderlayOn, false, func(_ context.Context, g string, _ bool) error {
		called = append(called, g)
		if g == "c" {
			return nil
		}
		return errors.New("resolver still on")
	})
	require.Error(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, called, "a failure must not stop the sweep")
	assert.Contains(t, err.Error(), "disable the Smart Gateway underlay")
	assert.Contains(t, err.Error(), "2 of 3")
	assert.Contains(t, err.Error(), "a: resolver still on")
	assert.Contains(t, err.Error(), "b: resolver still on")
	assert.NotContains(t, err.Error(), "c:")
}

// fakeSmartGatewayFabric records the underlay calls a sweep makes and applies them to its groups.
type fakeSmartGatewayFabric struct {
	featureOn bool
	groups    []goaviatrix.SmartGatewayGroupStatus
	calls     []string
}

func (f *fakeSmartGatewayFabric) ops() smartGatewayFabricOps {
	return smartGatewayFabricOps{
		featureEnabled: func(context.Context) (bool, error) { return f.featureOn, nil },
		groups: func(context.Context) ([]goaviatrix.SmartGatewayGroupStatus, error) {
			return append([]goaviatrix.SmartGatewayGroupStatus(nil), f.groups...), nil
		},
		setUnderlay: func(_ context.Context, group string, enable bool) error {
			verb := "off"
			if enable {
				verb = "on"
			}
			f.calls = append(f.calls, "underlay "+verb+" "+group)
			for i := range f.groups {
				if f.groups[i].GatewayGroupName == group {
					f.groups[i].SmartGatewayUnderlayEnabled = enable
				}
			}
			return nil
		},
	}
}

func TestApplySmartGatewayFabricFeatureUnderlaySweep(t *testing.T) {
	f := &fakeSmartGatewayFabric{featureOn: true, groups: []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "t1"},
		{GatewayGroupName: "t2"},
	}}
	require.NoError(t, applySmartGatewayFabricFeature(context.Background(), f.ops(), smartGatewayFabricUnderlayMesh, true))
	assert.Equal(t, []string{"underlay on t1", "underlay on t2"}, f.calls)
}

func TestApplySmartGatewayFabricFeatureRefusesUnderlayOffWhileResolverOn(t *testing.T) {
	f := &fakeSmartGatewayFabric{featureOn: true, groups: []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "t1", SmartGatewayUnderlayEnabled: true, SmartGatewayResolverEnabled: true},
		{GatewayGroupName: "t2", SmartGatewayUnderlayEnabled: true},
	}}
	err := applySmartGatewayFabricFeature(context.Background(), f.ops(), smartGatewayFabricUnderlayMesh, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "on in 1 gateway group(s): t1")
	assert.Contains(t, err.Error(), "enable_route_resolver = false")
	assert.Empty(t, f.calls)
}

func TestApplySmartGatewayFabricFeatureRefusesEnableWhileFeatureOff(t *testing.T) {
	f := &fakeSmartGatewayFabric{featureOn: false}
	err := applySmartGatewayFabricFeature(context.Background(), f.ops(), smartGatewayFabricUnderlayMesh, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "aviatrix_config_feature")

	assert.NoError(t, applySmartGatewayFabricFeature(context.Background(), f.ops(), smartGatewayFabricUnderlayMesh, false))
	assert.Empty(t, f.calls)
}

func TestApplySmartGatewayFabricFeatureUnknownFeature(t *testing.T) {
	f := &fakeSmartGatewayFabric{featureOn: true}
	err := applySmartGatewayFabricFeature(context.Background(), f.ops(), "bogus", true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown feature")
}

func TestConfigBool(t *testing.T) {
	assert.True(t, configBoolIsTrue(cty.True))
	assert.False(t, configBoolIsTrue(cty.False))
	assert.False(t, configBoolIsTrue(cty.NullVal(cty.Bool)))
	assert.False(t, configBoolIsTrue(cty.UnknownVal(cty.Bool)))
	assert.False(t, configBoolIsTrue(cty.StringVal("true")))
}

func TestSmartGatewayResolverUnsupportedErr(t *testing.T) {
	assert.Error(t, smartGatewayResolverUnsupportedErr(cty.True, true, "on a non-BGP cloud spoke gateway"))
	assert.NoError(t, smartGatewayResolverUnsupportedErr(cty.False, true, "x"), "explicit false is valid on any group")
	assert.NoError(t, smartGatewayResolverUnsupportedErr(cty.NullVal(cty.Bool), true, "x"), "unset is valid on any group")
	assert.NoError(t, smartGatewayResolverUnsupportedErr(cty.True, false, "x"))
}

func TestSpokeGroupWithoutBgpLu(t *testing.T) {
	assert.True(t, spokeGroupWithoutBgpLu("SPOKE", false))
	assert.False(t, spokeGroupWithoutBgpLu("SPOKE", true))
	assert.False(t, spokeGroupWithoutBgpLu("EDGESPOKE", false))
	assert.True(t, spokeGroupWithoutBgpLu("STANDALONE", true))
}

func TestSmartGatewayResolverCreateErr(t *testing.T) {
	midRollout := []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "a", SmartGatewayUnderlayEnabled: true, SmartGatewayResolverEnabled: true},
		{GatewayGroupName: "b", SmartGatewayUnderlayEnabled: true},
	}
	underlayGap := []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "a", SmartGatewayUnderlayEnabled: true},
		{GatewayGroupName: "b"},
	}

	err := smartGatewayResolverCreateErr(false, nil)
	require.Error(t, err, "resolver on needs the feature on")
	assert.Contains(t, err.Error(), "aviatrix_config_feature")

	err = smartGatewayResolverCreateErr(true, underlayGap)
	require.Error(t, err, "resolver on needs the underlay everywhere")
	assert.Contains(t, err.Error(), "off in 1: b")
	assert.Contains(t, err.Error(), "underlay_mesh")

	assert.NoError(t, smartGatewayResolverCreateErr(true, midRollout))
	assert.NoError(t, smartGatewayResolverCreateErr(true, nil), "first group in an empty fabric")
}

// fakeController answers the Smart Gateway calls; actions in refuse fail with that reason.
type fakeController struct {
	groups []goaviatrix.SmartGatewayGroupStatus
	refuse map[string]string
}

func (f *fakeController) client(t *testing.T) *goaviatrix.Client {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		body, _ := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		var fields map[string]any
		if json.Unmarshal(body, &fields) == nil {
			for k, v := range fields {
				if s, ok := v.(string); ok {
					params.Set(k, s)
				}
			}
		} else if form, err := url.ParseQuery(string(body)); err == nil {
			for k, v := range form {
				params[k] = v
			}
		}
		assert.NoError(t, json.NewEncoder(w).Encode(f.respond(params.Get("action"), params.Get("gateway_group_name"))))
	}))
	t.Cleanup(server.Close)
	client, err := goaviatrix.NewClient("admin", "password", strings.TrimPrefix(server.URL, "https://"), server.Client(), nil)
	require.NoError(t, err)
	return client
}

func (f *fakeController) respond(action, group string) map[string]any {
	if reason, ok := f.refuse[action]; ok {
		return map[string]any{"return": false, "reason": reason}
	}
	switch action {
	case "get_api_token":
		return map[string]any{"return": true, "results": map[string]any{"api_token": "token"}}
	case "login":
		return map[string]any{"return": true, "CID": "cid"}
	case "get_controller_feature":
		return map[string]any{"return": true, "results": map[string]any{"enabled": true}}
	case "get_smart_gw_fabric_status":
		return map[string]any{"return": true, "results": f.groups}
	case "get_smart_gw_resolver_status":
		resolver := "disabled"
		for _, g := range f.groups {
			if g.GatewayGroupName == group && g.SmartGatewayResolverEnabled {
				resolver = "enabled"
			}
		}
		return map[string]any{"return": true, "results": "Gateway group " + group + " smart_gw_resolver (resolver=" + resolver + ")"}
	case "list_gateway_groups":
		var results []map[string]string
		for _, g := range f.groups {
			results = append(results, map[string]string{"name": g.GatewayGroupName, "gw_type": "TRANSIT"})
		}
		return map[string]any{"return": true, "results": results}
	}
	enable := strings.HasPrefix(action, "enable_")
	for i := range f.groups {
		if f.groups[i].GatewayGroupName != group {
			continue
		}
		if strings.HasSuffix(action, "_underlay") {
			f.groups[i].SmartGatewayUnderlayEnabled = enable
		} else {
			f.groups[i].SmartGatewayResolverEnabled = enable
		}
	}
	return map[string]any{"return": true}
}

func TestSmartGatewayFabricResourceKeepsStateWhenRefused(t *testing.T) {
	f := &fakeController{groups: []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "transit-1", SmartGatewayResolverEnabled: true},
		{GatewayGroupName: "transit-2"},
	}}
	client := f.client(t)
	resource := resourceAviatrixSmartGatewayFabric()
	apply := func(prior *terraform.InstanceState, enabled bool) *terraform.InstanceState {
		config := terraform.NewResourceConfigRaw(map[string]any{"feature_name": "underlay_mesh", "is_enabled": enabled})
		instanceDiff, err := resource.Diff(context.Background(), prior, config, client)
		require.NoError(t, err)
		state, _ := resource.Apply(context.Background(), prior, instanceDiff, client)
		return state
	}

	created := apply(nil, true)
	assert.Equal(t, "true", created.Attributes["is_enabled"])

	// The resolver is still on in transit-1, so turning the underlay off is refused.
	refused := apply(created, false)
	assert.Equal(t, "true", refused.Attributes["is_enabled"])

	f.groups[0].SmartGatewayResolverEnabled = false
	disabled := apply(created, false)
	assert.Equal(t, "false", disabled.Attributes["is_enabled"])
	assert.False(t, f.groups[0].SmartGatewayUnderlayEnabled)
}

func TestApplySmartGatewayResolverKeepsControllerValueWhenRefused(t *testing.T) {
	f := &fakeController{
		groups: []goaviatrix.SmartGatewayGroupStatus{{GatewayGroupName: "transit-1"}},
		refuse: map[string]string{"enable_smart_gw_resolver": "smart_gw_underlay is disabled"},
	}
	resource := &schema.Resource{Schema: map[string]*schema.Schema{smartGatewayResolverAttr: smartGatewayResolverSchema()}}
	d := resource.Data(&terraform.InstanceState{
		ID:         "transit-1",
		Attributes: map[string]string{smartGatewayResolverAttr: "true"},
		RawConfig:  cty.ObjectVal(map[string]cty.Value{smartGatewayResolverAttr: cty.True}),
	})

	err := applySmartGatewayResolver(context.Background(), d, f.client(t), "transit-1")
	require.Error(t, err)
	assert.False(t, getBool(d, smartGatewayResolverAttr))
}

func TestCheckSmartGatewayResolverBeforeCreate(t *testing.T) {
	f := &fakeController{groups: []goaviatrix.SmartGatewayGroupStatus{
		{GatewayGroupName: "transit-1", SmartGatewayUnderlayEnabled: true},
		{GatewayGroupName: "transit-2"},
	}}
	client := f.client(t)
	resource := &schema.Resource{Schema: map[string]*schema.Schema{smartGatewayResolverAttr: smartGatewayResolverSchema()}}
	d := resource.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{smartGatewayResolverAttr: cty.True}),
	})

	err := checkSmartGatewayResolverBeforeCreate(context.Background(), d, client, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "off in 1: transit-2")

	f.groups[1].SmartGatewayUnderlayEnabled = true
	assert.NoError(t, checkSmartGatewayResolverBeforeCreate(context.Background(), d, client, true))
}
