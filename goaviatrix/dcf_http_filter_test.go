package goaviatrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDCFHTTPFilterSendsRequest(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusCreated, map[string]any{
		"uuid": "11111111-2222-3333-4444-555555555555",
	})

	desc := "write methods"
	filterUUID, err := client.CreateDCFHTTPFilter(context.Background(), &DCFHTTPFilter{
		Name: "write-methods",
		Desc: &desc,
		Filters: []DCFHTTPFilterAndClause{
			{All: []DCFHTTPFilterElement{{Hdr: ":method", Values: []string{"POST", "PUT"}}}},
			{All: []DCFHTTPFilterElement{{Hdr: "x-present"}}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "11111111-2222-3333-4444-555555555555", filterUUID)
	assert.Equal(t, http.MethodPost, rec.method)
	assert.Equal(t, "/v2.5/api/dcf/http-filters", rec.path)
	assert.Equal(t, map[string]any{
		"name": "write-methods",
		"desc": "write methods",
		"filters": []any{
			map[string]any{"all": []any{map[string]any{"hdr": ":method", "values": []any{"POST", "PUT"}}}},
			map[string]any{"all": []any{map[string]any{"hdr": "x-present"}}},
		},
	}, rec.body)
}

func TestCreateDCFHTTPFilterOmitsEmptyDescription(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusCreated, map[string]any{"uuid": "u"})

	_, err := client.CreateDCFHTTPFilter(context.Background(), &DCFHTTPFilter{
		Name:    "f",
		Filters: []DCFHTTPFilterAndClause{{All: []DCFHTTPFilterElement{{Hdr: "app"}}}},
	})
	require.NoError(t, err)
	assert.NotContains(t, rec.body, "desc")
	assert.NotContains(t, rec.body, "uuid")
}

func TestGetDCFHTTPFilterParsesResponse(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusOK, map[string]any{
		"uuid": "abc",
		"name": "f",
		"desc": "d",
		"filters": []any{
			map[string]any{"all": []any{map[string]any{"hdr": "app", "values": []any{"web"}}}},
		},
	})

	filter, err := client.GetDCFHTTPFilter(context.Background(), "abc")
	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/v2.5/api/dcf/http-filters/abc", rec.path)
	require.NotNil(t, filter.Desc)
	assert.Equal(t, "d", *filter.Desc)
	assert.Equal(t, &DCFHTTPFilter{
		UUID: "abc",
		Name: "f",
		Desc: filter.Desc,
		Filters: []DCFHTTPFilterAndClause{
			{All: []DCFHTTPFilterElement{{Hdr: "app", Values: []string{"web"}}}},
		},
	}, filter)
}

func TestGetDCFHTTPFilterMapsNotFound(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusNotFound, map[string]any{
		"message": "http filter not found with uuid: abc",
	})

	_, err := client.GetDCFHTTPFilter(context.Background(), "abc")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetDCFHTTPFilterPropagatesOtherErrors(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusInternalServerError, map[string]any{
		"message": "boom",
	})

	_, err := client.GetDCFHTTPFilter(context.Background(), "abc")
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNotFound)
}

func TestGetDCFHTTPFilterByNameUsesNameEndpoint(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusNotFound, map[string]any{
		"message": "http filter not found with name: f",
	})

	_, err := client.GetDCFHTTPFilterByName(context.Background(), "f")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Equal(t, "/v2.5/api/dcf/http-filters/name/f", rec.path)
}

func TestGetDCFHTTPFilterByNameEscapesName(t *testing.T) {
	var rawPath string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath = r.URL.EscapedPath()
		_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "abc", "name": "../a b/c"})
	}))
	t.Cleanup(server.Close)
	client := &Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
	}

	filter, err := client.GetDCFHTTPFilterByName(context.Background(), "../a b/c")
	require.NoError(t, err)
	assert.Equal(t, "abc", filter.UUID)
	assert.Equal(t, "/v2.5/api/dcf/http-filters/name/..%2Fa%20b%2Fc", rawPath)
}

func TestUpdateDCFHTTPFilterSendsPut(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusOK, map[string]any{"uuid": "abc"})

	err := client.UpdateDCFHTTPFilter(context.Background(), "abc", &DCFHTTPFilter{
		Name:    "f",
		Filters: []DCFHTTPFilterAndClause{{All: []DCFHTTPFilterElement{{Hdr: "app"}}}},
	})
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, rec.method)
	assert.Equal(t, "/v2.5/api/dcf/http-filters/abc", rec.path)
}

func TestDeleteDCFHTTPFilterReportsConflict(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusConflict, map[string]any{
		"message": "http filter is currently in use and cannot be deleted",
	})

	err := client.DeleteDCFHTTPFilter(context.Background(), "abc")
	require.ErrorContains(t, err, "in use")
	assert.Equal(t, http.MethodDelete, rec.method)
	assert.Equal(t, "/v2.5/api/dcf/http-filters/abc", rec.path)
}
