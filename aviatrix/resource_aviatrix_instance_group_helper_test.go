package aviatrix

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setGroupUUIDFromGatewayName is the shared Read-side fix for AVX-80061: it
// re-derives group_uuid so terraform import cannot silently drop it. A blank
// group name must be a no-op (older controllers may not return one) and must
// never make a controller call — verified here without a live client.
func TestSetGroupUUIDFromGatewayName_BlankNameIsNoop(t *testing.T) {
	res := resourceAviatrixSpokeInstance()
	d := res.TestResourceData()

	// A nil client would panic if the helper tried a lookup; a blank name must
	// short-circuit before that, so reaching the assertions proves the no-op.
	err := setGroupUUIDFromGatewayName(context.Background(), nil, d, "")
	require.NoError(t, err)
	assert.Empty(t, getString(d, "group_uuid"))
}

// insertion_gateway_az has no per-gateway controller field to read back, so the
// fix marks it Computed to avoid a spurious ForceNew diff on import. Guard that
// schema property so it cannot silently regress.
func TestInsertionGatewayAzIsComputed(t *testing.T) {
	res := resourceAviatrixSpokeInstance()
	field, ok := res.Schema["insertion_gateway_az"]
	require.True(t, ok, "insertion_gateway_az must exist in the spoke instance schema")
	assert.True(t, field.Computed, "insertion_gateway_az must be Computed to survive import")
	assert.True(t, field.ForceNew, "insertion_gateway_az must remain ForceNew")
}

// group_uuid must stay Required+ForceNew on both instance resources; the fix
// relies on Read re-deriving it rather than on it being Computed.
func TestGroupUUIDSchemaContract(t *testing.T) {
	for name, res := range map[string]*schema.Resource{
		"spoke":   resourceAviatrixSpokeInstance(),
		"transit": resourceAviatrixTransitInstance(),
	} {
		field, ok := res.Schema["group_uuid"]
		require.Truef(t, ok, "%s: group_uuid must exist", name)
		assert.Truef(t, field.Required, "%s: group_uuid must be Required", name)
		assert.Truef(t, field.ForceNew, "%s: group_uuid must be ForceNew", name)
		assert.Falsef(t, field.Computed, "%s: group_uuid must not be Computed", name)
	}
}
