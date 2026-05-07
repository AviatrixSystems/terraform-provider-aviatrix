package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeTransitAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeTransitAttachmentCreate,
		Read:   resourceAviatrixSpokeTransitAttachmentRead,
		Update: resourceAviatrixSpokeTransitAttachmentUpdate,
		Delete: resourceAviatrixSpokeTransitAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"spoke_gw_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Name of the spoke gateway or gateway group to attach to transit network.",
			},
			"transit_gw_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Name of the transit gateway or gateway group to attach the spoke gateway to.",
			},
			"route_tables": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "Learned routes will be propagated to these route tables.",
			},
			"tunnel_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 49),
				Description: "(Optional) Advanced option. Number of public tunnels. Required with both Spoke and Transit" +
					"to be insane mode enabled and max performance enabled. Type: Integer. Valid Range: 1-49." +
					"Available as of provider version R3.1.3+.",
			},
			"enable_max_performance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
				Description: "Indicates whether the maximum amount of HPE tunnels will be created. " +
					"Only valid when transit and spoke gateways are each launched in Insane Mode and in the same cloud type. " +
					"Available as of provider version R2.22.2+.",
			},
			"spoke_prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on spoke gateway.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"transit_prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit gateway.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"spoke_bgp_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the spoke gateway is BGP enabled or not.",
			},
			"transit_gateway_logical_ifnames": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Transit gateway logical interface names for edge gateways, where the peering terminates. Required for all edge gateways.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// resolveGatewayName resolves name to an actual gateway name.
// It first checks whether name is a gateway group; if so, it returns the
// primary gateway name from the group's details. If name is not a group, it
// verifies a gateway with that name exists and returns name as-is. Returns
// an error if neither a matching group nor a gateway is found.
func resolveGatewayName(ctx context.Context, client *goaviatrix.Client, name string) (string, error) {
	group, err := client.GetGatewayGroupByName(ctx, name)
	if err == nil {
		if group.PrimaryGatewayName == "" {
			return "", fmt.Errorf("gateway group %q has no primary gateway", name)
		}
		return group.PrimaryGatewayName, nil
	}
	if !errors.Is(err, goaviatrix.ErrNotFound) {
		return "", fmt.Errorf("error looking up gateway group %q: %w", name, err)
	}

	// Not a group — verify it is a direct gateway name.
	if _, err := client.GetGateway(&goaviatrix.Gateway{GwName: name}); err == nil {
		return name, nil
	} else if !errors.Is(err, goaviatrix.ErrNotFound) {
		return "", fmt.Errorf("error looking up gateway %q: %w", name, err)
	}

	return "", fmt.Errorf("no gateway group or gateway found with name %q", name)
}

/*
 * recoverSpokeTransitAttachment checks if a spoke-transit attachment was
 * created on the controller despite a timeout or connection error. Polls
 * up to 5 times with 60s intervals to allow the controller to finish.
 * Uses the provided context for cancellation and timeout control.
 * Returns true if the attachment exists.
 */
func recoverSpokeTransitAttachment(ctx context.Context, client *goaviatrix.Client, attachment *goaviatrix.SpokeTransitAttachment) bool {
	log.Printf("[WARN] Timeout/connection error attaching spoke %s to transit %s, checking if attachment was created",
		attachment.SpokeGwName, attachment.TransitGwName)
	checkAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   attachment.SpokeGwName,
		TransitGwName: attachment.TransitGwName,
	}
	// The controller may still be completing the operation, so
	// poll a few times with backoff before giving up.
	checkTries, checkMaxTries, checkBackoff := 0, 5, 60*time.Second
	for checkTries < checkMaxTries {
		checkTries++
		if _, readErr := client.GetSpokeTransitAttachmentContext(ctx, checkAttachment); readErr == nil {
			log.Printf("[INFO] Spoke %s is already attached to transit %s despite timeout, proceeding",
				attachment.SpokeGwName, attachment.TransitGwName)
			return true
		}
		log.Printf("[INFO] Attachment not found yet (attempt %d/%d), waiting %v before checking again",
			checkTries, checkMaxTries, checkBackoff)
		select {
		case <-ctx.Done():
			log.Printf("[WARN] Context cancelled while waiting for attachment recovery")
			return false
		case <-time.After(checkBackoff):
		}
	}
	// Final check after polling
	if _, readErr := client.GetSpokeTransitAttachmentContext(ctx, checkAttachment); readErr == nil {
		log.Printf("[INFO] Spoke %s is already attached to transit %s despite timeout, proceeding",
			attachment.SpokeGwName, attachment.TransitGwName)
		return true
	}
	return false
}

func resourceAviatrixSpokeTransitAttachmentCreate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	attachment := marshalSpokeTransitAttachmentInput(d)

	ctx := context.Background()
	resolvedSpokeGwName, err := resolveGatewayName(ctx, client, attachment.SpokeGwName)
	if err != nil {
		return fmt.Errorf("could not resolve spoke_gw_name %q: %w", attachment.SpokeGwName, err)
	}
	resolvedTransitGwName, err := resolveGatewayName(ctx, client, attachment.TransitGwName)
	if err != nil {
		return fmt.Errorf("could not resolve transit_gw_name %q: %w", attachment.TransitGwName, err)
	}
	attachment.SpokeGwName = resolvedSpokeGwName
	attachment.TransitGwName = resolvedTransitGwName

	spoke, err := client.GetGateway(&goaviatrix.Gateway{GwName: resolvedSpokeGwName})
	if err != nil {
		return fmt.Errorf("could not find spoke gateway: %w", err)
	}

	// get transit gateway details
	transitGatewayDetails, err := getGatewayDetails(client, resolvedTransitGwName)
	if err != nil {
		return fmt.Errorf("could not get transit gateway details for %s: %w", resolvedTransitGwName, err)
	}
	// get edge transit logical interface names
	if err := getEdgeTransitLogicalIfNames(d, transitGatewayDetails, attachment); err != nil {
		return fmt.Errorf("%w", err)
	}

	if !spoke.EnableBgp && (len(attachment.SpokePrependAsPath) != 0 || len(attachment.TransitPrependAsPath) != 0) {
		return fmt.Errorf("'spoke_prepend_as_path' and 'transit_prepend_as_path' are only valid for BGP enabled spoke gateway")
	}

	if attachment.NoMaxPerformance && attachment.InsaneModeTunnelNumber > 0 {
		return fmt.Errorf("'tunnel_count' can only be specified with max performance enabled. Please set 'enable_max_performance' to true")
	}

	d.SetId(resolvedSpokeGwName + "~" + resolvedTransitGwName)
	flag := false
	defer func() { _ = resourceAviatrixSpokeTransitAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	timeout := 15 * time.Minute
	attachCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	try, maxTries, backoff := 0, 10, 1000*time.Millisecond
	for {
		try++
		err := client.CreateSpokeTransitAttachment(attachCtx, attachment)
		if err != nil {
			if strings.Contains(err.Error(), "AVXERR-TRANSIT-0034") {
				// already joined, so we can break out of the loop
				break
			}
			if strings.Contains(err.Error(), "is not up") || strings.Contains(err.Error(), "is not ready") {
				if try == maxTries {
					return fmt.Errorf("could not attach spoke: %s to transit %s: %w", resolvedSpokeGwName, resolvedTransitGwName, err)
				}
				time.Sleep(backoff)
				// Double the backoff time after each failed try
				backoff *= 2
				continue
			}
			// On timeout or connection errors, check if the controller
			// completed the operation despite the client not receiving
			// a response (e.g. proxy timeout, context deadline exceeded).
			// Use a fresh context since the original one may have expired.
			if isTimeoutOrConnectionError(err) {
				recoveryCtx, recoveryCancel := context.WithTimeout(context.Background(), timeout)
				defer recoveryCancel()
				if recoverSpokeTransitAttachment(recoveryCtx, client, attachment) {
					break
				}
			}
			return fmt.Errorf("failed to attach spoke: %s to transit %s: %w", resolvedSpokeGwName, resolvedTransitGwName, err)
		}
		break
	}

	if len(attachment.SpokePrependAsPath) != 0 {
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: attachment.SpokeGwName,
			TransitGatewayName2: attachment.TransitGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transGwPeering, attachment.SpokePrependAsPath)
		if err != nil {
			return fmt.Errorf("could not set spoke_prepend_as_path: %w", err)
		}
	}

	if len(attachment.TransitPrependAsPath) != 0 {
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: attachment.TransitGwName,
			TransitGatewayName2: attachment.SpokeGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transGwPeering, attachment.TransitPrependAsPath)
		if err != nil {
			return fmt.Errorf("could not set transit_prepend_as_path: %w", err)
		}
	}

	return resourceAviatrixSpokeTransitAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSpokeTransitAttachmentReadIfRequired(d *schema.ResourceData, meta any, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeTransitAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixSpokeTransitAttachmentRead(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	spokeGwName := getString(d, "spoke_gw_name")
	transitGwName := getString(d, "transit_gw_name")
	isImport := spokeGwName == "" || transitGwName == ""
	if isImport {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no spoke_gw_name or transit_gw_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return fmt.Errorf("import id is invalid, expecting spoke_gw_name~transit_gw_name, but received: %s", id)
		}
		spokeGwName = parts[0]
		transitGwName = parts[1]
		mustSet(d, "spoke_gw_name", spokeGwName)
		mustSet(d, "transit_gw_name", transitGwName)
	}

	ctx := context.Background()
	resolvedSpokeGwName, err := resolveGatewayName(ctx, client, spokeGwName)
	if err != nil {
		return fmt.Errorf("could not resolve spoke_gw_name %q: %w", spokeGwName, err)
	}
	resolvedTransitGwName, err := resolveGatewayName(ctx, client, transitGwName)
	if err != nil {
		return fmt.Errorf("could not resolve transit_gw_name %q: %w", transitGwName, err)
	}

	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   resolvedSpokeGwName,
		TransitGwName: resolvedTransitGwName,
	}

	attachment, err := client.GetSpokeTransitAttachment(spokeTransitAttachment)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find spoke_transit_attachment: %w", err)
	}
	mustSet(d, "spoke_bgp_enabled", attachment.SpokeBgpEnabled)

	if attachment.RouteTables != "" {
		var routeTables []string
		for routeTable := range strings.SplitSeq(attachment.RouteTables, ",") {
			routeTables = append(routeTables, strings.Split(routeTable, "~~")[0])
		}

		err = d.Set("route_tables", routeTables)
		if err != nil {
			return fmt.Errorf("could not set route_tables: %w", err)
		}
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: resolvedSpokeGwName,
		TransitGatewayName2: resolvedTransitGwName,
	}

	transitGatewayPeering, err = client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}

	// set the transit gateway logical interface names only for edge gateways
	transitGateway, err := getGatewayDetails(client, resolvedTransitGwName)
	if err != nil {
		return fmt.Errorf("could not get transit gateway details for %s: %w", resolvedTransitGwName, err)
	}

	if goaviatrix.IsCloudType(transitGateway.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		if len(attachment.TransitGatewayLogicalIfNames) > 0 {
			logicalIfNames, err := getLogicalIfNames(transitGateway, attachment.TransitGatewayLogicalIfNames)
			if err != nil {
				return fmt.Errorf("could not get logical interface names for edge transit gateway %s: %w", resolvedTransitGwName, err)
			}
			_ = d.Set("transit_gateway_logical_ifnames", logicalIfNames)
		}
	}
	mustSet(d, "enable_max_performance", !transitGatewayPeering.NoMaxPerformance)
	if !transitGatewayPeering.NoMaxPerformance && transitGatewayPeering.InsaneModeTunnelCount > 0 {
		mustSet(d, "tunnel_count", transitGatewayPeering.InsaneModeTunnelCount)
	}

	if attachment.SpokeBgpEnabled {
		if transitGatewayPeering.PrependAsPath1 != "" {
			var prependAsPath []string
			for str := range strings.SplitSeq(transitGatewayPeering.PrependAsPath1, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("spoke_prepend_as_path", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set spoke_prepend_as_path: %w", err)
			}
		}

		if transitGatewayPeering.PrependAsPath2 != "" {
			var prependAsPath []string
			for str := range strings.SplitSeq(transitGatewayPeering.PrependAsPath2, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("transit_prepend_as_path", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set transit_prepend_as_path: %w", err)
			}
		}
	} else {
		mustSet(d, "spoke_prepend_as_path", nil)
		mustSet(d, "transit_prepend_as_path", nil)
	}

	d.SetId(resolvedSpokeGwName + "~" + resolvedTransitGwName)
	return nil
}

func resourceAviatrixSpokeTransitAttachmentUpdate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	if !getBool(d, "spoke_bgp_enabled") && d.HasChanges("spoke_prepend_as_path", "transit_prepend_as_path") {
		return fmt.Errorf("'spoke_prepend_as_path' and 'transit_prepend_as_path' are only valid for BGP enabled spoke gateway")
	}

	d.Partial(true)

	ctx := context.Background()
	spokeGwName, err := resolveGatewayName(ctx, client, getString(d, "spoke_gw_name"))
	if err != nil {
		return fmt.Errorf("could not resolve spoke_gw_name: %w", err)
	}
	transitGwName, err := resolveGatewayName(ctx, client, getString(d, "transit_gw_name"))
	if err != nil {
		return fmt.Errorf("could not resolve transit_gw_name: %w", err)
	}

	if d.HasChange("spoke_prepend_as_path") {
		transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: spokeGwName,
			TransitGatewayName2: transitGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transitGatewayPeering, getStringList(d, "spoke_prepend_as_path"))
		if err != nil {
			return fmt.Errorf("could not update spoke_prepend_as_path: %w", err)
		}
	}

	if d.HasChange("transit_prepend_as_path") {
		transitGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: transitGwName,
			TransitGatewayName2: spokeGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transitGwPeering, getStringList(d, "transit_prepend_as_path"))
		if err != nil {
			return fmt.Errorf("could not update transit_prepend_as_path: %w", err)
		}
	}

	if d.HasChange("tunnel_count") {
		transitGatewayPeering := &goaviatrix.TransitGatewayPeeringEdit{
			TransitGatewayName1: spokeGwName,
			TransitGatewayName2: transitGwName,
			TunnelCount:         getInt(d, "tunnel_count"),
		}

		if transitGatewayPeering.TunnelCount > 0 && !getBool(d, "enable_max_performance") {
			return fmt.Errorf("'tunnel_count' can't be updated with max performance disabled")
		}

		err := client.UpdateTransitGatewayPeeringTunnelCount(transitGatewayPeering)
		if err != nil {
			return fmt.Errorf("could not update tunnel_count for spoke transit attachment: %v : %w", spokeGwName+"~"+transitGwName, err)
		}
	}

	d.Partial(false)
	d.SetId(spokeGwName + "~" + transitGwName)
	return resourceAviatrixSpokeTransitAttachmentRead(d, meta)
}

func resourceAviatrixSpokeTransitAttachmentDelete(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	ctx := context.Background()
	resolvedSpokeGwName, err := resolveGatewayName(ctx, client, getString(d, "spoke_gw_name"))
	if err != nil {
		return fmt.Errorf("could not resolve spoke_gw_name: %w", err)
	}
	resolvedTransitGwName, err := resolveGatewayName(ctx, client, getString(d, "transit_gw_name"))
	if err != nil {
		return fmt.Errorf("could not resolve transit_gw_name: %w", err)
	}

	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   resolvedSpokeGwName,
		TransitGwName: resolvedTransitGwName,
	}
	if err := client.DeleteSpokeTransitAttachment(spokeTransitAttachment); err != nil {
		return fmt.Errorf("could not detach spoke: %s from transit %s: %w", resolvedSpokeGwName,
			resolvedTransitGwName, err)
	}

	return nil
}

func marshalSpokeTransitAttachmentInput(d *schema.ResourceData) *goaviatrix.SpokeTransitAttachment {
	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:            getString(d, "spoke_gw_name"),
		TransitGwName:          getString(d, "transit_gw_name"),
		RouteTables:            strings.Join(getStringSet(d, "route_tables"), ","),
		SpokePrependAsPath:     getStringList(d, "spoke_prepend_as_path"),
		TransitPrependAsPath:   getStringList(d, "transit_prepend_as_path"),
		InsaneModeTunnelNumber: getInt(d, "tunnel_count"),
		NoMaxPerformance:       !getBool(d, "enable_max_performance"),
	}

	return spokeTransitAttachment
}
