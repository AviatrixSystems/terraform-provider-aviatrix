package goaviatrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudTypesToString(t *testing.T) {
	tt := []struct {
		Name     string
		Mask     int
		Expected string
	}{
		{"empty", 0, ""},
		{"single AWS", AWS, "AWS (1)"},
		{
			"transit firenet allowlist",
			AWSRelatedCloudTypes | GCPRelatedCloudTypes | AzureArmRelatedCloudTypes | OCIRelatedCloudTypes,
			"AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048)",
		},
		{"china edit restriction", AWSChina | AzureChina, "AWSChina (1024), AzureChina (2048)"},
		{
			"gcp-azure edit restriction",
			GCPRelatedCloudTypes | AzureArmRelatedCloudTypes,
			"GCP (4), Azure (8), AzureGov (32), AzureChina (2048)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, CloudTypesToString(tc.Mask))
		})
	}
}
