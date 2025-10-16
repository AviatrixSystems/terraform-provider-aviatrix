package goaviatrix

import (
	"context"
	"fmt"
)

// DCFTrustBundle represents a DCF trust bundle structure for GET responses
type DCFTrustBundle struct {
	DisplayName   string   `json:"display_name"`
	BundleID      string   `json:"bundle_id"`
	BundleContent []string `json:"bundle_content"`
	CreatedAt     string   `json:"created_at"`
	UUID          string   `json:"uuid,omitempty"`
}

// DCFTrustBundleCreateRequest represents the request structure for creating a trust bundle
type DCFTrustBundleCreateRequest struct {
	BundleContent string `json:"bundle_content"`
	DisplayName   string `json:"display_name"`
}

// DCFTrustBundleCreateResponse represents the response structure for creating a trust bundle
type DCFTrustBundleCreateResponse struct {
	BundleID string `json:"bundle_id"`
}

// DCFTrustbundleErrorResponse represents the error response structure for DCF trust bundle operations
type DCFTrustbundleErrorResponse struct {
	Message string `json:"message"`
}

// CreateDCFTrustBundle creates a new DCF trust bundle
func (c *Client) CreateDCFTrustBundle(ctx context.Context, bundleContent, displayName string) (string, error) {
	endpoint := "dcf/trustbundle"

	request := DCFTrustBundleCreateRequest{
		BundleContent: bundleContent,
		DisplayName:   displayName,
	}

	var response DCFTrustBundleCreateResponse
	err := c.PostAPIContext25(ctx, &response, endpoint, request)
	if err != nil {
		return "", err
	}

	return response.BundleID, nil
}

// GetDCFTrustBundleByID retrieves a DCF trust bundle by UUID
func (c *Client) GetDCFTrustBundleByID(ctx context.Context, bundleUUID string) (*DCFTrustBundle, error) {
	endpoint := fmt.Sprintf("dcf/trustbundle/%s", bundleUUID)

	var trustBundle DCFTrustBundle
	err := c.GetAPIContext25(ctx, &trustBundle, endpoint, nil)
	if err != nil {
		return nil, err
	}
	// Set the UUID from the parameter for consistency
	trustBundle.UUID = bundleUUID
	return &trustBundle, nil
}

// GetDCFTrustBundleByName retrieves a DCF trust bundle by name
func (c *Client) GetDCFTrustBundleByName(ctx context.Context, bundleName string) (*DCFTrustBundle, error) {
	endpoint := fmt.Sprintf("dcf/trustbundle/name/%s", bundleName)
	var trustBundle DCFTrustBundle
	err := c.GetAPIContext25(ctx, &trustBundle, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if trustBundle.BundleID == "" {
		return nil, ErrNotFound
	}
	trustBundle.UUID = trustBundle.BundleID
	return &trustBundle, nil
}

// DeleteDCFTrustBundle deletes a DCF trust bundle by UUID
func (c *Client) DeleteDCFTrustBundle(ctx context.Context, bundleUUID string) error {
	endpoint := fmt.Sprintf("dcf/trustbundle/%s", bundleUUID)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
