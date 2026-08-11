package aviatrix

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// setGroupUUIDFromGatewayName resolves the owning gateway group's UUID from its
// name and persists it into state as group_uuid.
//
// group_uuid is Required+ForceNew (not Computed) on both aviatrix_spoke_instance
// and aviatrix_transit_instance, so the SDK never re-derives it. Without this,
// terraform import — which populates state purely via Read — leaves group_uuid
// empty, which both breaks Update/Delete (they pass it to the controller) and
// forces a disruptive replacement on the next plan. GetGateway already returns
// the group name, and GetGatewayGroupByName maps a name back to a UUID, so this
// is a pure Read-side fix over existing client calls.
//
// A blank groupName (e.g. an older controller that does not return it) is a
// no-op rather than an error, so Read never fails on gateways that legitimately
// have no owning group in the response.
func setGroupUUIDFromGatewayName(ctx context.Context, client *goaviatrix.Client, d *schema.ResourceData, groupName string) error {
	if groupName == "" {
		return nil
	}
	group, err := client.GetGatewayGroupByName(ctx, groupName)
	if err != nil {
		return fmt.Errorf("could not look up gateway group %q: %w", groupName, err)
	}
	mustSet(d, "group_uuid", group.GroupUUID)
	return nil
}
