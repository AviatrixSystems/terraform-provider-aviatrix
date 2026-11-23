package goaviatrix

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const dcfPolicyScopesEndpoint = "dcf/policy-scopes"

type DCFPolicyScopeSelector map[string][]string

type DCFPolicyScopeSelectors struct {
	Selectors []DCFPolicyScopeSelector `json:"selectors"`
}

type DCFPolicyScopeInput struct {
	Name     string                   `json:"name"`
	Included DCFPolicyScopeSelectors  `json:"included"`
	Excluded *DCFPolicyScopeSelectors `json:"excluded,omitempty"`
}

type DCFPolicyScope struct {
	Name     string                   `json:"name"`
	UUID     string                   `json:"uuid"`
	Included DCFPolicyScopeSelectors  `json:"included"`
	Excluded *DCFPolicyScopeSelectors `json:"excluded,omitempty"`
}

// CreateDCFPolicyScope creates a DCF policy scope.
func (c *Client) CreateDCFPolicyScope(ctx context.Context, scope *DCFPolicyScopeInput) (*DCFPolicyScope, error) {
	var created DCFPolicyScope
	if err := c.PostAPIContext25(ctx, &created, dcfPolicyScopesEndpoint, scope); err != nil {
		return nil, fmt.Errorf("creating DCF policy scope: %w", err)
	}
	return &created, nil
}

// GetDCFPolicyScope retrieves a DCF policy scope by UUID.
func (c *Client) GetDCFPolicyScope(ctx context.Context, scopeUUID string) (*DCFPolicyScope, error) {
	endpoint, err := url.JoinPath(dcfPolicyScopesEndpoint, scopeUUID)
	if err != nil {
		return nil, fmt.Errorf("constructing DCF policy scope endpoint: %w", err)
	}

	var scope DCFPolicyScope
	if err := c.GetAPIContext25(ctx, &scope, endpoint, nil); err != nil {
		return nil, normalizeDCFPolicyScopeNotFound(err)
	}
	return &scope, nil
}

type dcfPolicyScopeList struct {
	Scopes []DCFPolicyScope `json:"scopes"`
}

// ListDCFPolicyScopes retrieves all DCF policy scopes.
func (c *Client) ListDCFPolicyScopes(ctx context.Context) ([]DCFPolicyScope, error) {
	var list dcfPolicyScopeList
	if err := c.GetAPIContext25(ctx, &list, dcfPolicyScopesEndpoint, nil); err != nil {
		return nil, fmt.Errorf("listing DCF policy scopes: %w", err)
	}
	return list.Scopes, nil
}

// GetDCFPolicyScopeByName retrieves the DCF policy scope with the given name.
func (c *Client) GetDCFPolicyScopeByName(ctx context.Context, name string) (*DCFPolicyScope, error) {
	scopes, err := c.ListDCFPolicyScopes(ctx)
	if err != nil {
		return nil, err
	}
	for i := range scopes {
		if scopes[i].Name == name {
			return &scopes[i], nil
		}
	}
	return nil, ErrNotFound
}

// UpdateDCFPolicyScope replaces a DCF policy scope by UUID.
func (c *Client) UpdateDCFPolicyScope(ctx context.Context, scopeUUID string, scope *DCFPolicyScopeInput) error {
	endpoint, err := url.JoinPath(dcfPolicyScopesEndpoint, scopeUUID)
	if err != nil {
		return fmt.Errorf("constructing DCF policy scope endpoint: %w", err)
	}
	if err := c.PutAPIContext25(ctx, endpoint, scope); err != nil {
		return fmt.Errorf("updating DCF policy scope %q: %w", scopeUUID, err)
	}
	return nil
}

// DeleteDCFPolicyScope deletes a DCF policy scope by UUID.
func (c *Client) DeleteDCFPolicyScope(ctx context.Context, scopeUUID string) error {
	endpoint, err := url.JoinPath(dcfPolicyScopesEndpoint, scopeUUID)
	if err != nil {
		return fmt.Errorf("constructing DCF policy scope endpoint: %w", err)
	}
	if err := c.DeleteAPIContext25(ctx, endpoint, nil); err != nil {
		return normalizeDCFPolicyScopeNotFound(err)
	}
	return nil
}

// normalizeDCFPolicyScopeNotFound maps an API 404 to the provider's not-found sentinel.
func normalizeDCFPolicyScopeNotFound(err error) error {
	var statusErr StatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode() == http.StatusNotFound {
		return ErrNotFound
	}
	return err
}
