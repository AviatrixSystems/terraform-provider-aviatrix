package goaviatrix

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// TrustBundleItemRequest represents the request body for creating/updating trust bundles
type TrustBundleItemRequest struct {
	BundleContent string `json:"bundle_content"`
	DisplayName   string `json:"display_name"`
}

// TrustBundle represents the full trust bundle response from GET operations
type TrustBundle struct {
	BundleID      string    `json:"bundle_id"`
	DisplayName   string    `json:"display_name"`
	BundleContent []string  `json:"bundle_content"`
	CreatedAt     time.Time `json:"created_at"`
}

// TrustBundleCreateResponse represents the response from creating/updating trust bundles
type TrustBundleCreateResponse struct {
	BundleID string `json:"bundle_id"`
}

// CreateDCFTrustBundle creates a new DCF trust bundle
func (c *Client) CreateDCFTrustBundle(ctx context.Context, trustBundle *TrustBundleItemRequest) (string, error) {
	endpoint := "dcf/trustbundle"

	var response TrustBundleCreateResponse
	err := c.PostAPIContext25(ctx, &response, endpoint, trustBundle)
	if err != nil {
		return "", err
	}

	return response.BundleID, nil
}

// GetDCFTrustBundle retrieves a DCF trust bundle by UUID
func (c *Client) GetDCFTrustBundle(ctx context.Context, bundleUUID string) (*TrustBundle, error) {
	endpoint := fmt.Sprintf("dcf/trustbundle/%s", bundleUUID)

	var trustBundle TrustBundle
	err := c.GetAPIContext25(ctx, &trustBundle, endpoint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &trustBundle, nil
}

// UpdateDCFTrustBundle updates an existing DCF trust bundle
func (c *Client) UpdateDCFTrustBundle(ctx context.Context, bundleUUID string, trustBundle *TrustBundleItemRequest) error {
	endpoint := fmt.Sprintf("dcf/trustbundle/%s", bundleUUID)
	err := c.PutAPIContext25(ctx, endpoint, trustBundle)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDCFTrustBundle deletes a DCF trust bundle by UUID
func (c *Client) DeleteDCFTrustBundle(ctx context.Context, bundleUUID string) error {
	endpoint := fmt.Sprintf("dcf/trustbundle/%s", bundleUUID)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
