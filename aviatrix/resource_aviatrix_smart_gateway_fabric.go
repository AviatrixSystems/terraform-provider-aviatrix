package aviatrix

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	smartGatewayFeatureName           = "smart_gateway"
	smartGatewayFabricFeatureNameAttr = "feature_name"
	smartGatewayFabricIsEnabledAttr   = "is_enabled"
	smartGatewayResolverAttr          = "enable_route_resolver"

	smartGatewayFabricUnderlayMesh = "underlay_mesh"
)

var smartGatewayFabricFeatures = []string{
	smartGatewayFabricUnderlayMesh,
}

// resourceAviatrixSmartGatewayFabric turns one Smart Gateway fabric feature on or off across every eligible group.
func resourceAviatrixSmartGatewayFabric() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixSmartGatewayFabricCreate,
		ReadWithoutTimeout:   resourceAviatrixSmartGatewayFabricRead,
		UpdateWithoutTimeout: resourceAviatrixSmartGatewayFabricUpdate,
		DeleteWithoutTimeout: resourceAviatrixSmartGatewayFabricDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			smartGatewayFabricFeatureNameAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(smartGatewayFabricFeatures, false),
				Description: fmt.Sprintf(
					"Which Smart Gateway fabric feature this instance owns. One of %s.",
					strings.Join(smartGatewayFabricFeatures, ", "),
				),
			},
			smartGatewayFabricIsEnabledAttr: {
				Type:     schema.TypeBool,
				Required: true,
				Description: "Turn the selected feature on or off across every eligible gateway group on the controller. " +
					"Requires the `smart_gateway` controller feature (aviatrix_config_feature).",
			},
		},
	}
}

// Create sweeps the feature to is_enabled and stores the feature name as the ID.
func resourceAviatrixSmartGatewayFabricCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)
	feature := getString(d, smartGatewayFabricFeatureNameAttr)
	enabled := getBool(d, smartGatewayFabricIsEnabledAttr)
	if err := applySmartGatewayFabricFeature(ctx, smartGatewayFabricClientOps(client), feature, enabled); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(feature)
	return resourceAviatrixSmartGatewayFabricRead(ctx, d, meta)
}

// Update sweeps the feature to the new is_enabled; on failure the prior state is kept.
func resourceAviatrixSmartGatewayFabricUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if !d.HasChange(smartGatewayFabricIsEnabledAttr) {
		return resourceAviatrixSmartGatewayFabricRead(ctx, d, meta)
	}
	d.Partial(true)
	if err := applySmartGatewayFabricFeature(ctx, smartGatewayFabricClientOps(mustClient(meta)), d.Id(), getBool(d, smartGatewayFabricIsEnabledAttr)); err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)
	return resourceAviatrixSmartGatewayFabricRead(ctx, d, meta)
}

// Read sets feature_name from the ID and is_enabled from the controller.
func resourceAviatrixSmartGatewayFabricRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	feature := d.Id()
	enabled, err := smartGatewayFabricEnabled(ctx, mustClient(meta), feature, getBool(d, smartGatewayFabricIsEnabledAttr))
	if err != nil {
		return diag.FromErr(err)
	}
	mustSet(d, smartGatewayFabricFeatureNameAttr, feature)
	mustSet(d, smartGatewayFabricIsEnabledAttr, enabled)
	return nil
}

// Delete only removes the resource from state; the controller is left unchanged.
func resourceAviatrixSmartGatewayFabricDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	d.SetId("")
	return nil
}

// smartGatewayFabricEnabled reports whether the feature is on across every eligible group.
func smartGatewayFabricEnabled(ctx context.Context, client *goaviatrix.Client, feature string, prior bool) (bool, error) {
	featureOn, err := smartGatewayFeatureEnabled(ctx, client)
	if err != nil || !featureOn {
		return false, err
	}
	groups, err := smartGatewayEligibleGroups(ctx, client)
	if err != nil {
		return false, err
	}
	if err := smartGatewayFabricCheckFeature(feature); err != nil {
		return false, err
	}
	return smartGatewayFabricEnabledFromFabric(groups, smartGatewayUnderlayOn, prior), nil
}

// smartGatewayFabricOps holds the controller calls the sweep makes, so tests can fake them.
type smartGatewayFabricOps struct {
	featureEnabled func(ctx context.Context) (bool, error)
	groups         func(ctx context.Context) ([]goaviatrix.SmartGatewayGroupStatus, error)
	setUnderlay    func(ctx context.Context, group string, enable bool) error
}

// smartGatewayFabricClientOps backs smartGatewayFabricOps with the real client.
func smartGatewayFabricClientOps(client *goaviatrix.Client) smartGatewayFabricOps {
	return smartGatewayFabricOps{
		featureEnabled: func(ctx context.Context) (bool, error) { return smartGatewayFeatureEnabled(ctx, client) },
		groups: func(ctx context.Context) ([]goaviatrix.SmartGatewayGroupStatus, error) {
			return smartGatewayEligibleGroups(ctx, client)
		},
		setUnderlay: client.SetSmartGatewayUnderlay,
	}
}

// applySmartGatewayFabricFeature turns the feature on or off on every eligible group.
func applySmartGatewayFabricFeature(ctx context.Context, ops smartGatewayFabricOps, feature string, enable bool) error {
	featureOn, err := ops.featureEnabled(ctx)
	if err != nil {
		return err
	}
	if !featureOn {
		if enable {
			return fmt.Errorf("the Smart Gateway feature is disabled on the controller: enable it first with aviatrix_config_feature (feature_name = %q)", smartGatewayFeatureName)
		}
		// Every per-group flag is already off while the feature is off.
		return nil
	}
	groups, err := ops.groups(ctx)
	if err != nil {
		return fmt.Errorf("%w; nothing was changed, run apply again once the controller responds", err)
	}
	if err := smartGatewayFabricCheckFeature(feature); err != nil {
		return err
	}
	if on := smartGatewayGroupsMatching(groups, smartGatewayResolverOn); !enable && len(on) > 0 {
		return fmt.Errorf(
			"cannot set `is_enabled = false` on the `underlay_mesh` aviatrix_smart_gateway_fabric while the route resolver is on in %d gateway group(s): %s. "+
				"Set `%s = false` on those gateways first",
			len(on), strings.Join(on, ", "), smartGatewayResolverAttr)
	}
	return smartGatewayFabricSweep(ctx, groups, "underlay mesh", smartGatewayUnderlayOn, enable, ops.setUnderlay)
}

// smartGatewayFabricCheckFeature rejects any feature_name other than underlay_mesh.
func smartGatewayFabricCheckFeature(feature string) error {
	if feature != smartGatewayFabricUnderlayMesh {
		return fmt.Errorf("unknown %s %q; expected one of %s",
			smartGatewayFabricFeatureNameAttr, feature, strings.Join(smartGatewayFabricFeatures, ", "))
	}
	return nil
}

// smartGatewayUnderlayOn reports a group's underlay flag.
func smartGatewayUnderlayOn(g goaviatrix.SmartGatewayGroupStatus) bool {
	return g.SmartGatewayUnderlayEnabled
}

// smartGatewayResolverOn reports a group's route resolver flag.
func smartGatewayResolverOn(g goaviatrix.SmartGatewayGroupStatus) bool {
	return g.SmartGatewayResolverEnabled
}

// smartGatewayFabricEnabledFromFabric keeps prior for no groups and flips it for a partly enabled fabric, so the plan shows a change.
func smartGatewayFabricEnabledFromFabric(groups []goaviatrix.SmartGatewayGroupStatus, on func(goaviatrix.SmartGatewayGroupStatus) bool, prior bool) bool {
	if len(groups) == 0 {
		return prior
	}
	all, anyOn := true, false
	for _, g := range groups {
		all = all && on(g)
		anyOn = anyOn || on(g)
	}
	switch {
	case all:
		return true
	case !anyOn:
		return false
	default:
		return !prior
	}
}

// smartGatewayGroupsMatching returns the names of the groups that match.
func smartGatewayGroupsMatching(groups []goaviatrix.SmartGatewayGroupStatus, match func(goaviatrix.SmartGatewayGroupStatus) bool) []string {
	var names []string
	for _, g := range groups {
		if match(g) {
			names = append(names, g.GatewayGroupName)
		}
	}
	return names
}

// smartGatewayFeatureEnabled reports whether the smart_gateway controller feature is on.
func smartGatewayFeatureEnabled(ctx context.Context, client *goaviatrix.Client) (bool, error) {
	feature, err := client.GetFeatureStatus(ctx, smartGatewayFeatureName)
	if err != nil {
		return false, fmt.Errorf("failed to read Smart Gateway feature status: %w", err)
	}
	return feature.Enabled, nil
}

// smartGatewayEligibleGroups returns the groups that use the Smart Gateway flags, sorted by name.
func smartGatewayEligibleGroups(ctx context.Context, client *goaviatrix.Client) ([]goaviatrix.SmartGatewayGroupStatus, error) {
	fabric, err := client.GetSmartGatewayFabricStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read Smart Gateway fabric status: %w", err)
	}
	types, err := client.ListGatewayGroupTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list gateway-group types: %w", err)
	}
	return filterSmartGatewayEligible(fabric, types), nil
}

// filterSmartGatewayEligible drops standalone, VPN and cloud spoke groups, like the controller does.
func filterSmartGatewayEligible(fabric []goaviatrix.SmartGatewayGroupStatus, types map[string]string) []goaviatrix.SmartGatewayGroupStatus {
	eligible := make([]goaviatrix.SmartGatewayGroupStatus, 0, len(fabric))
	for _, g := range fabric {
		if g.IsCloudSpokeWithoutBgpLu || !goaviatrix.IsSmartGatewayEligibleType(types[g.GatewayGroupName]) {
			continue
		}
		eligible = append(eligible, g)
	}
	sort.Slice(eligible, func(i, j int) bool { return eligible[i].GatewayGroupName < eligible[j].GatewayGroupName })
	return eligible
}

// smartGatewayFabricSweep sets one flag on every group not already at the target and reports all failures together.
func smartGatewayFabricSweep(
	ctx context.Context,
	groups []goaviatrix.SmartGatewayGroupStatus,
	flag string,
	current func(goaviatrix.SmartGatewayGroupStatus) bool,
	enable bool,
	set func(ctx context.Context, group string, enable bool) error,
) error {
	var failures []string
	for _, g := range groups {
		if current(g) == enable {
			continue
		}
		if err := set(ctx, g.GatewayGroupName, enable); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", g.GatewayGroupName, err))
		}
	}
	if len(failures) == 0 {
		return nil
	}
	verb := "disable"
	if enable {
		verb = "enable"
	}
	return fmt.Errorf("failed to %s the Smart Gateway %s on %d of %d gateway group(s):\n%s",
		verb, flag, len(failures), len(groups), strings.Join(failures, "\n"))
}

// smartGatewayResolverSchema is the enable_route_resolver attribute shared by the gateway and group resources.
func smartGatewayResolverSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
		Description: "Enable the Smart Gateway route resolver on this gateway group. Unset keeps the controller's value, " +
			"which is on for a new group when the rest of the fabric is fully rolled out. " +
			"Not supported on non-BGP cloud spokes or standalone groups.",
	}
}

// validateSmartGatewayResolverDiff rejects enable_route_resolver = true at plan time on a group without BGP-LU.
func validateSmartGatewayResolverDiff(d *schema.ResourceDiff, unsupported bool, reason string) error {
	return smartGatewayResolverUnsupportedErr(d.GetRawConfig().GetAttr(smartGatewayResolverAttr), unsupported, reason)
}

// smartGatewayResolverUnsupportedErr returns an error when the resolver is set to true on an unsupported group.
func smartGatewayResolverUnsupportedErr(raw cty.Value, unsupported bool, reason string) error {
	if unsupported && configBoolIsTrue(raw) {
		return fmt.Errorf("`%s = true` is not supported %s: the Smart Gateway route resolver needs BGP-LU", smartGatewayResolverAttr, reason)
	}
	return nil
}

// checkSmartGatewayResolverBeforeCreate checks the resolver prerequisites before a gateway or group is created.
func checkSmartGatewayResolverBeforeCreate(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, eligible bool) error {
	raw := d.GetRawConfig().GetAttr(smartGatewayResolverAttr)
	if !configBoolIsTrue(raw) {
		return nil
	}
	// Plan-time validation is skipped while enable_bgp or gw_type is unknown, so repeat it here.
	if !eligible {
		return smartGatewayResolverUnsupportedErr(raw, true, "on a gateway group that does not run BGP-LU")
	}
	featureOn, err := smartGatewayFeatureEnabled(ctx, client)
	if err != nil {
		return err
	}
	var groups []goaviatrix.SmartGatewayGroupStatus
	if featureOn {
		if groups, err = smartGatewayEligibleGroups(ctx, client); err != nil {
			return err
		}
	}
	return smartGatewayResolverCreateErr(featureOn, groups)
}

// smartGatewayResolverCreateErr requires the feature on and the underlay on in every group.
func smartGatewayResolverCreateErr(featureOn bool, groups []goaviatrix.SmartGatewayGroupStatus) error {
	if !featureOn {
		return fmt.Errorf(
			"`%s = true` requires the Smart Gateway feature: enable it first with aviatrix_config_feature (feature_name = %q)",
			smartGatewayResolverAttr, smartGatewayFeatureName)
	}
	if off := smartGatewayGroupsMatching(groups, func(g goaviatrix.SmartGatewayGroupStatus) bool { return !g.SmartGatewayUnderlayEnabled }); len(off) > 0 {
		return fmt.Errorf(
			"`%s = true` requires the Smart Gateway underlay on every gateway group, but it is off in %d: %s. "+
				"Set `is_enabled = true` on the `underlay_mesh` aviatrix_smart_gateway_fabric first",
			smartGatewayResolverAttr, len(off), strings.Join(off, ", "))
	}
	return nil
}

// smartGatewayResolverCreateWarning reports a post-create resolver failure as a warning, so the new resource is not tainted.
func smartGatewayResolverCreateWarning(err error) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("%s was not applied", smartGatewayResolverAttr),
		Detail:   fmt.Sprintf("The resource was created, but setting the route resolver failed: %v. The next plan shows the change and apply retries it.", err),
	}
}

// applySmartGatewayResolver sets the group's resolver to the configured value when it differs.
func applySmartGatewayResolver(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, group string) error {
	if d.GetRawConfig().GetAttr(smartGatewayResolverAttr).IsNull() {
		return nil
	}
	want := getBool(d, smartGatewayResolverAttr)
	current, err := client.GetSmartGatewayResolverStatus(ctx, group)
	if err != nil {
		return fmt.Errorf("failed to read %s on gateway group %q: %w", smartGatewayResolverAttr, group, err)
	}
	if current == want {
		return nil
	}
	if err := client.SetSmartGatewayResolver(ctx, group, want); err != nil {
		// Keep the controller's value in state; the group resources do not use d.Partial.
		mustSet(d, smartGatewayResolverAttr, current)
		return fmt.Errorf("failed to set %s on gateway group %q: %w", smartGatewayResolverAttr, group, err)
	}
	return nil
}

// applySmartGatewayResolverForGateway applies the resolver to the gateway's group, looking the group up on create.
func applySmartGatewayResolverForGateway(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client) error {
	if d.GetRawConfig().GetAttr(smartGatewayResolverAttr).IsNull() {
		return nil
	}
	group := getString(d, "group_name")
	if group == "" {
		gw, err := client.GetGateway(&goaviatrix.Gateway{
			AccountName: getString(d, "account_name"),
			GwName:      getString(d, "gw_name"),
		})
		if err != nil {
			return fmt.Errorf("failed to look up the gateway group of %q: %w", getString(d, "gw_name"), err)
		}
		group = gw.GroupName
	}
	return applySmartGatewayResolver(ctx, d, client, group)
}

// readSmartGatewayResolver sets enable_route_resolver from the controller, leaving state alone if the call fails.
func readSmartGatewayResolver(ctx context.Context, d *schema.ResourceData, client *goaviatrix.Client, group string) {
	if group == "" {
		return
	}
	enabled, err := client.GetSmartGatewayResolverStatus(ctx, group)
	if err != nil {
		log.Printf("[WARN] could not read Smart Gateway resolver for gateway group %q: %v", group, err)
		return
	}
	mustSet(d, smartGatewayResolverAttr, enabled)
}

// configBoolIsTrue reports whether a raw config value is a known, explicit true.
func configBoolIsTrue(v cty.Value) bool {
	return v.IsKnown() && !v.IsNull() && v.Type() == cty.Bool && v.True()
}
