package goaviatrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// webGroupRoundTrip sends a filter through the API encoding and back, mimicking a
// create followed by a refresh.
func webGroupRoundTrip(filter *WebGroupMatchExpression) *WebGroupMatchExpression {
	wg := newWebGroup(WebGroupResult{
		Selector: WebGroupAnyResult{
			Any: []WebGroupMatchExpressionResult{
				{All: webGroupFilterToMap(filter)},
			},
		},
	})
	return wg.Selector.Expressions[0]
}

func TestWebGroupFilterToMap_externalSNISpreadsExtArgs(t *testing.T) {
	got := webGroupFilterToMap(&WebGroupMatchExpression{
		ExternalSNI: "malicious-domain-feed",
		ExtArgs:     map[string]string{"filterA": "valueA"},
	})

	assert.Equal(t, map[string]string{
		ExternalSNIKey: "malicious-domain-feed",
		"filterA":      "valueA",
	}, got)
}

func TestWebGroupRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		filter *WebGroupMatchExpression
	}{
		{"snifilter", &WebGroupMatchExpression{SniFilter: "*.aviatrix.com"}},
		{"urlfilter", &WebGroupMatchExpression{UrlFilter: "aviatrix.com/test"}},
		{"external_sni", &WebGroupMatchExpression{ExternalSNI: "malicious-domain-feed"}},
		{"external_sni with ext_args", &WebGroupMatchExpression{
			ExternalSNI: "malicious-domain-feed",
			ExtArgs:     map[string]string{"filterA": "valueA", "filterB": "valueB"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.filter, webGroupRoundTrip(tt.filter))
		})
	}
}

// Keys outside a feed match are not feed arguments, so they must not surface as
// ext_args on a plain SNI or URL filter.
func TestNewWebGroupMatchExpression_ignoresUnknownKeysWithoutFeed(t *testing.T) {
	got := newWebGroupMatchExpression(map[string]string{
		SniFilterKey: "*.aviatrix.com",
		"unknown":    "value",
	})

	assert.Equal(t, &WebGroupMatchExpression{SniFilter: "*.aviatrix.com"}, got)
}
