package goaviatrix

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type TLSProfile struct {
	// TLSProfile defines model for TLSProfile.
	// CABundleID is the UUID of the CA bundle that should be used for origin certificate validation. If not populated the default bundle would be used.
	CABundleID *string `json:"CA_bundle_id,omitempty"`

	// CertificateValidation Certificate validation mode
	CertificateValidation string `json:"certificate_validation"`

	// DisplayName Display name for the TLS profile
	DisplayName string `json:"display_name"`

	// VerifySni Configuration to verify SNI of client
	VerifySni bool `json:"verify_sni"`
}

type TLSProfileWithID struct {
	// TLSProfile defines model for TLSProfile.
	// CABundleID is the UUID of the CA bundle that should be used for origin certificate validation. If not populated the default bundle would be used.
	CABundleID *string `json:"CA_bundle_id,omitempty"`

	// CertificateValidation Certificate validation mode
	CertificateValidation string `json:"certificate_validation"`

	// DisplayName Display name for the TLS profile
	DisplayName string `json:"display_name"`

	// UUID The unique identifier for the TLS profile
	UUID string `json:"uuid"`

	// VerifySni Configuration to verify SNI of client
	VerifySni bool `json:"verify_sni"`
}

type TLSProfileResponse struct {
	// TLSProfileResponse defines response model for TLSProfile.
	UUID string `json:"uuid"`
}

type TLSProfilesListResponse struct {
	// Profiles List of all TLS profiles
	Profiles []TLSProfileWithID `json:"profiles"`
}

func isValidUUID(u string) error {
	if _, err := uuid.Parse(u); err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	return nil
}

func (c *Client) GetTLSProfile(ctx context.Context, uuidStr string) (*TLSProfileWithID, error) {
	endpoint := fmt.Sprintf("dcf/tls-profile/%s", uuidStr)
	if err := isValidUUID(uuidStr); err != nil {
		return nil, err
	}
	var tlsProfile TLSProfileWithID
	err := c.GetAPIContext25(ctx, &tlsProfile, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return &tlsProfile, nil
}

func (c *Client) CreateTLSProfile(ctx context.Context, tlsProfile *TLSProfile) (string, error) {
	endpoint := "dcf/tls-profile"

	var data TLSProfileResponse
	if err := c.PostAPIContext25(ctx, &data, endpoint, tlsProfile); err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) ListTLSProfiles(ctx context.Context) (*TLSProfilesListResponse, error) {
	endpoint := "dcf/tls-profile"

	var listResponse TLSProfilesListResponse
	err := c.GetAPIContext25(ctx, &listResponse, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &listResponse, nil
}

func (c *Client) UpdateTLSProfile(ctx context.Context, uuid string, tlsProfile *TLSProfile) error {
	if err := isValidUUID(uuid); err != nil {
		return err
	}
	endpoint := fmt.Sprintf("dcf/tls-profile/%s", uuid)
	err := c.PutAPIContext25(ctx, endpoint, tlsProfile)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteTLSProfile(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("dcf/tls-profile/%s", uuid)
	if err := isValidUUID(uuid); err != nil {
		return err
	}
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
