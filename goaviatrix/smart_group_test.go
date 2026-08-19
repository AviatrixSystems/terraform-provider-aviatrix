package goaviatrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// roundTrip sends a filter through the API encoding and back, mimicking a
// create followed by a refresh.
func roundTrip(filter *SmartGroupMatchExpression) *SmartGroupMatchExpression {
	sg := createSmartGroup(SmartGroupResult{
		Selector: SmartGroupAnyResult{
			Any: []SmartGroupMatchExpressionResult{
				{All: SmartGroupFilterToAPIMap(filter)},
			},
		},
	})
	return sg.Selector.Expressions[0]
}

// Both maps are populated so that either one leaking into the other, or either
// losing its separating dot, shows up here.
func TestSmartGroupRoundTrip_typedRowKeepsTagsAndNamespaceTags(t *testing.T) {
	got := roundTrip(&SmartGroupMatchExpression{
		Type:          "k8s",
		K8sClusterID:  "cluster1",
		Tags:          map[string]string{"app": "web"},
		NamespaceTags: map[string]string{"dept": "310"},
	})

	assert.Equal(t, "k8s", got.Type)
	assert.Equal(t, "cluster1", got.K8sClusterID)
	assert.Equal(t, map[string]string{"app": "web"}, got.Tags)
	assert.Equal(t, map[string]string{"dept": "310"}, got.NamespaceTags)
}

// The resource and data source read the maps whole instead of flattening them to
// dotted keys, so that direction needs its own case.
func TestSmartGroupFilterToResource_keepsBothMapsWhole(t *testing.T) {
	got := SmartGroupFilterToResource(&SmartGroupMatchExpression{
		Type:          "k8s",
		Tags:          map[string]string{"app": "web"},
		NamespaceTags: map[string]string{"dept": "310"},
	})

	assert.Equal(t, map[string]string{"app": "web"}, got[TagsPrefix])
	assert.Equal(t, map[string]string{"dept": "310"}, got[NamespaceTagsPrefix])
}

// An external feed defines its own argument names, so external takes precedence
// over type: every other key stays literal in ext_args instead of being
// un-prefixed into a map.
func TestSmartGroupRoundTrip_externalTakesPrecedenceOverType(t *testing.T) {
	got := roundTrip(&SmartGroupMatchExpression{
		External:      "geo",
		Type:          "k8s",
		ExtArgs:       map[string]string{"country_iso_code": "AU"},
		NamespaceTags: map[string]string{"dept": "310"},
	})

	assert.Equal(t, "geo", got.External)
	assert.Empty(t, got.NamespaceTags)
	assert.Empty(t, got.Type)
	assert.Equal(t, map[string]string{
		"country_iso_code":   "AU",
		"namespacetags.dept": "310",
		"type":               "k8s",
	}, got.ExtArgs)
}
