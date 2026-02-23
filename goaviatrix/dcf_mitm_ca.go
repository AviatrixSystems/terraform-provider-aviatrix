package goaviatrix

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	DCFMITMSystemCAID    = "def000ca-6000-1000-8000-000000000001"
	DCFMitmCaStateActive = "active"
)

// MitmCaItemRequest represents the request body for creating a MITM CA
type MitmCaItemRequest struct {
	Name             string `json:"name"`
	Key              string `json:"key"`
	CertificateChain string `json:"certificate_chain"`
}

// MitmCaPatchRequest represents the request body for updating a MITM CA
type MitmCaPatchRequest struct {
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
}

// MitmCaResponse represents the MITM CA response from GET operations
type MitmCaResponse struct {
	Name             string    `json:"name"`
	CaID             string    `json:"ca_id"`
	CaHash           string    `json:"ca_hash"`
	CertificateChain string    `json:"certificate_chain"`
	State            string    `json:"state"`
	Origin           string    `json:"origin"`
	CreatedAt        time.Time `json:"created_at"`
}

// MitmCaCreateResponse represents the response from creating a MITM CA
type MitmCaCreateResponse struct {
	CaID string `json:"ca_id"`
}

// MitmCaListResponse represents the response from listing DCF MITM CAs
type MitmCaListResponse struct {
	Cas   []MitmCaResponse `json:"cas"`
	Total int              `json:"total"`
}

// ListDCFMitmCa lists all DCF MITM CAs
func (c *Client) ListDCFMitmCa(ctx context.Context) (*MitmCaListResponse, error) {
	endpoint := "dcf/mitm-ca"

	var response MitmCaListResponse
	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateDCFMitmCa creates a new DCF MITM CA
func (c *Client) CreateDCFMitmCa(ctx context.Context, mitmCa *MitmCaItemRequest) (string, error) {
	endpoint := "dcf/mitm-ca"

	var response MitmCaCreateResponse
	err := c.PostAPIContext25(ctx, &response, endpoint, mitmCa)
	if err != nil {
		return "", err
	}

	return response.CaID, nil
}

// GetDCFMitmCa retrieves a DCF MITM CA by ID
func (c *Client) GetDCFMitmCa(ctx context.Context, caID string) (*MitmCaResponse, error) {
	endpoint, err := url.JoinPath("dcf/mitm-ca", caID)
	if err != nil {
		return nil, fmt.Errorf("failed to construct endpoint: %w", err)
	}

	var mitmCa MitmCaResponse
	err = c.GetAPIContext25(ctx, &mitmCa, endpoint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &mitmCa, nil
}

// UpdateDCFMitmCa updates an existing DCF MITM CA (name and/or state)
func (c *Client) UpdateDCFMitmCa(ctx context.Context, caID string, patchRequest *MitmCaPatchRequest) (*MitmCaResponse, error) {
	endpoint, err := url.JoinPath("dcf/mitm-ca", caID)
	if err != nil {
		return nil, fmt.Errorf("failed to construct endpoint: %w", err)
	}

	var response MitmCaResponse
	err = c.PatchAPIContext25(ctx, &response, endpoint, patchRequest)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteDCFMitmCa deletes a DCF MITM CA by ID
func (c *Client) DeleteDCFMitmCa(ctx context.Context, caID string) error {
	endpoint, err := url.JoinPath("dcf/mitm-ca", caID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}

	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func (c *Client) RefreshDCFMitmSysatemCA(ctx context.Context) error {
	endpoint, err := url.JoinPath("dcf/mitm-ca", DCFMITMSystemCAID, "regenerate")
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}
	err = c.PostAPIContext25(ctx, nil, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to refresh DCF MITM system CA: %w", err)
	}
	return nil
}
