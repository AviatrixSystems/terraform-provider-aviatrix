package goaviatrix

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExternalConnectionReturnsConnection(t *testing.T) {
	client, rec := newGatewayInsertionTestClient(t, http.StatusOK, map[string]any{
		"uuid": "abc",
		"name": "conn",
	})

	conn, err := client.GetExternalConnection(context.Background(), "abc")
	require.NoError(t, err)
	assert.Equal(t, "abc", conn.UUID)
	assert.Equal(t, "conn", conn.Name)
	assert.Equal(t, http.MethodGet, rec.method)
	assert.Equal(t, "/v2.5/api/external-connections/abc", rec.path)
}

func TestGetExternalConnectionMapsNotFound(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusNotFound, map[string]any{
		"message": "[AVXERR-EXTCONN-0001] Invalid request: Connection abc not found",
	})

	_, err := client.GetExternalConnection(context.Background(), "abc")
	require.ErrorIs(t, err, ErrNotFound)
	assert.ErrorContains(t, err, "Connection abc not found")
}

func TestGetExternalConnectionKeepsBadRequest(t *testing.T) {
	client, _ := newGatewayInsertionTestClient(t, http.StatusBadRequest, map[string]any{
		"message": "bad request",
	})

	_, err := client.GetExternalConnection(context.Background(), "abc")
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNotFound)
}
