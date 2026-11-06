package goaviatrix

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const dcfHTTPFilterEndpoint = "dcf/http-filters"

// DCFHTTPFilterElement matches a single HTTP header against zero or more values.
type DCFHTTPFilterElement struct {
	Hdr    string   `json:"hdr"`
	Values []string `json:"values,omitempty"`
}

// DCFHTTPFilterAndClause is a set of header matches that must all hold.
type DCFHTTPFilterAndClause struct {
	All []DCFHTTPFilterElement `json:"all"`
}

// DCFHTTPFilter is an HTTP filter expressed in disjunctive normal form: it
// matches when any of its clauses matches.
type DCFHTTPFilter struct {
	Name    string                   `json:"name"`
	Desc    *string                  `json:"desc,omitempty"`
	UUID    string                   `json:"uuid,omitempty"`
	Filters []DCFHTTPFilterAndClause `json:"filters"`
}

type dcfHTTPFilterUUIDResponse struct {
	UUID string `json:"uuid"`
}

// mapDCFHTTPFilterNotFound converts an HTTP 404 from the http-filters API into ErrNotFound.
func mapDCFHTTPFilterNotFound(err error) error {
	var statusErr StatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}
	return err
}

// CreateDCFHTTPFilter creates a DCF HTTP filter and returns its UUID.
func (c *Client) CreateDCFHTTPFilter(ctx context.Context, filter *DCFHTTPFilter) (string, error) {
	var resp dcfHTTPFilterUUIDResponse
	if err := c.PostAPIContext25(ctx, &resp, dcfHTTPFilterEndpoint, filter); err != nil {
		return "", fmt.Errorf("failed to create DCF HTTP filter %q: %w", filter.Name, err)
	}
	return resp.UUID, nil
}

// GetDCFHTTPFilter returns the DCF HTTP filter with the given UUID.
func (c *Client) GetDCFHTTPFilter(ctx context.Context, filterUUID string) (*DCFHTTPFilter, error) {
	endpoint, err := url.JoinPath(dcfHTTPFilterEndpoint, filterUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to construct endpoint: %w", err)
	}
	var filter DCFHTTPFilter
	if err := c.GetAPIContext25(ctx, &filter, endpoint, nil); err != nil {
		return nil, mapDCFHTTPFilterNotFound(err)
	}
	return &filter, nil
}

// GetDCFHTTPFilterByName returns the DCF HTTP filter with the given unique name.
func (c *Client) GetDCFHTTPFilterByName(ctx context.Context, name string) (*DCFHTTPFilter, error) {
	// Names are free-form, so escape the name as a single path segment rather
	// than letting "/" or ".." in it change the request path.
	endpoint := dcfHTTPFilterEndpoint + "/name/" + url.PathEscape(name)
	var filter DCFHTTPFilter
	if err := c.GetAPIContext25(ctx, &filter, endpoint, nil); err != nil {
		return nil, mapDCFHTTPFilterNotFound(err)
	}
	return &filter, nil
}

// UpdateDCFHTTPFilter replaces the DCF HTTP filter with the given UUID.
func (c *Client) UpdateDCFHTTPFilter(ctx context.Context, filterUUID string, filter *DCFHTTPFilter) error {
	endpoint, err := url.JoinPath(dcfHTTPFilterEndpoint, filterUUID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}
	if err := c.PutAPIContext25(ctx, endpoint, filter); err != nil {
		return fmt.Errorf("failed to update DCF HTTP filter %s: %w", filterUUID, err)
	}
	return nil
}

// DeleteDCFHTTPFilter deletes the DCF HTTP filter with the given UUID.
func (c *Client) DeleteDCFHTTPFilter(ctx context.Context, filterUUID string) error {
	endpoint, err := url.JoinPath(dcfHTTPFilterEndpoint, filterUUID)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint: %w", err)
	}
	if err := c.DeleteAPIContext25(ctx, endpoint, nil); err != nil {
		return fmt.Errorf("failed to delete DCF HTTP filter %s: %w", filterUUID, err)
	}
	return nil
}
