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
func TestSmartGroupRoundTrip_typedRowKeepsTagsAndK8sNamespaceTags(t *testing.T) {
	got := roundTrip(&SmartGroupMatchExpression{
		Type:             "k8s",
		K8sClusterID:     "cluster1",
		Tags:             map[string]string{"app": "web"},
		K8sNamespaceTags: map[string]string{"dept": "310"},
	})

	assert.Equal(t, "k8s", got.Type)
	assert.Equal(t, "cluster1", got.K8sClusterID)
	assert.Equal(t, map[string]string{"app": "web"}, got.Tags)
	assert.Equal(t, map[string]string{"dept": "310"}, got.K8sNamespaceTags)
}

// The resource and data source read the maps whole instead of flattening them to
// dotted keys, so that direction needs its own case.
func TestSmartGroupFilterToResource_keepsBothMapsWhole(t *testing.T) {
	got := SmartGroupFilterToResource(&SmartGroupMatchExpression{
		Type:             "k8s",
		Tags:             map[string]string{"app": "web"},
		K8sNamespaceTags: map[string]string{"dept": "310"},
	})

	assert.Equal(t, map[string]string{"app": "web"}, got[TagsPrefix])
	assert.Equal(t, map[string]string{"dept": "310"}, got[K8sNamespaceTagsPrefix])
}

func TestSmartGroupRoundTrip_resourceGroupAndVpcEndpointSelectors(t *testing.T) {
	tests := []struct {
		name   string
		filter *SmartGroupMatchExpression
	}{
		{
			name: "AWS",
			filter: &SmartGroupMatchExpression{
				Type:          "vpc_endpoint",
				ServiceName:   "com.amazonaws.us-east-1.s3",
				ServiceRegion: "us-east-1",
			},
		},
		{
			name: "Azure resource",
			filter: &SmartGroupMatchExpression{
				Type:          "vm",
				ResourceGroup: "my-resource-group",
			},
		},
		{
			name: "Azure private endpoint",
			filter: &SmartGroupMatchExpression{
				Type:          "vpc_endpoint",
				ResourceGroup: "my-private-endpoint-resource-group",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.filter, roundTrip(tt.filter))
		})
	}
}

// An external feed defines its own argument names, so external takes precedence
// over type: every other key stays literal in ext_args instead of being
// un-prefixed into a map.
func TestSmartGroupRoundTrip_externalTakesPrecedenceOverType(t *testing.T) {
	got := roundTrip(&SmartGroupMatchExpression{
		External:         "geo",
		Type:             "k8s",
		ExtArgs:          map[string]string{"country_iso_code": "AU"},
		K8sNamespaceTags: map[string]string{"dept": "310"},
	})

	assert.Equal(t, "geo", got.External)
	assert.Empty(t, got.K8sNamespaceTags)
	assert.Empty(t, got.Type)
	assert.Equal(t, map[string]string{
		"country_iso_code":        "AU",
		"k8s_namespace_tags.dept": "310",
		"type":                    "k8s",
	}, got.ExtArgs)
}
