package goaviatrix

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// TrustBundleItemRequest represents the request body for creating/updating trust bundles
type TrustBundleItemRequest struct {
	BundleContent string `json:"bundle_content"`
	DisplayName   string `json:"display_name"`
}

// DCFTrustBundle represents the full trust bundle response from GET operations
type DCFTrustBundle struct {
	BundleID      string    `json:"bundle_id"`
	DisplayName   string    `json:"display_name"`
	BundleContent []string  `json:"bundle_content"`
	CreatedAt     time.Time `json:"created_at"`
	UUID          string    `json:"uuid,omitempty"`
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
func (c *Client) GetDCFTrustBundleByID(ctx context.Context, bundleUUID string) (*DCFTrustBundle, error) {
	endpoint, err := url.JoinPath("dcf/trustbundle", bundleUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to construct endpoint: %w", err)
	}

	var trustBundle DCFTrustBundle
	err = c.GetAPIContext25(ctx, &trustBundle, endpoint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// Set the UUID from the parameter for consistency
	trustBundle.UUID = bundleUUID
	return &trustBundle, nil
}

// GetDCFTrustBundleByName retrieves a DCF trust bundle by name
func (c *Client) GetDCFTrustBundleByName(ctx context.Context, bundleName string) (*DCFTrustBundle, error) {
	endpoint, err := url.JoinPath("dcf/trustbundle/name", bundleName)
	if err != nil {
		return nil, fmt.Errorf("failed to construct endpoint: %w", err)
	}
	var trustBundle DCFTrustBundle
	err = c.GetAPIContext25(ctx, &trustBundle, endpoint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}
	trustBundle.UUID = trustBundle.BundleID
	return &trustBundle, nil
}

// UpdateDCFTrustBundle updates an existing DCF trust bundle
func (c *Client) UpdateDCFTrustBundle(ctx context.Context, bundleUUID string, trustBundle *TrustBundleItemRequest) error {
	endpoint, err := url.JoinPath("dcf/trustbundle", bundleUUID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}
	err = c.PutAPIContext25(ctx, endpoint, trustBundle)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDCFTrustBundle deletes a DCF trust bundle by UUID
func (c *Client) DeleteDCFTrustBundle(ctx context.Context, bundleUUID string) error {
	endpoint, err := url.JoinPath("dcf/trustbundle", bundleUUID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func ValidateTrustbundle(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	certs, err := ParseCertificates([]byte(v))
	if err != nil {
		return nil, []error{err}
	}
	if len(certs) == 0 {
		return nil, []error{fmt.Errorf("no certificates found in bundle")}
	}
	for _, cert := range certs {
		if !cert.IsCA {
			return nil, []error{fmt.Errorf("certificate %q is not a CA", cert.Subject.CommonName)}
		}
	}
	return nil, nil
}

func ParseCertificates(remain []byte) ([]*x509.Certificate, error) {
	// Remove UTF-8 BOM if present using standard library
	utf8BOM := []byte{0xEF, 0xBB, 0xBF}
	remain = bytes.TrimPrefix(remain, utf8BOM)
	return parseCertificatesNoBom(remain)
}

func parseCertificatesNoBom(remain []byte) ([]*x509.Certificate, error) {
	var chain []*x509.Certificate

	for {
		var block *pem.Block

		block, remain = pem.Decode(remain)
		if block == nil {
			break
		}

		// We ignore non-certificate PEM blocks because
		// that's what (*CertPool).AppendCertsFromPEM()
		// does.
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %q object: %w", block.Type, err)
			}

			chain = append(chain, cert)
		}
	}

	return chain, nil
}
