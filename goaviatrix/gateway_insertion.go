package goaviatrix

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const (
	GatewayInsertionEgressModeNonEgress   = "non-egress"
	GatewayInsertionEgressModeLocalEgress = "local-egress"
)

const gatewayInsertionEndpoint = "gateway-insertion/vpc"

type GatewayInsertionVpcConfig struct {
	VpcID       string   `json:"vpc_id"`
	RouteTables []string `json:"route_tables"`
	EgressMode  string   `json:"egress_mode"`
}

type GatewayInsertionRequest struct {
	VpcConfig GatewayInsertionVpcConfig `json:"vpc_config"`
}

type GatewayInsertionResult struct {
	VpcID      string `json:"vpc_id"`
	Status     string `json:"status"`
	ErrMessage string `json:"err_message"`
}

type GatewayInsertionStatus struct {
	VpcID            string   `json:"vpc_id"`
	Enabled          bool     `json:"enabled"`
	RouteTableConfig []string `json:"route_table_config"`
	InsertionMode    string   `json:"insertion_mode"`
}

func checkGatewayInsertionResult(action string, result *GatewayInsertionResult) error {
	if result.Status != "success" {
		return fmt.Errorf("%s gateway insertion for VPC %s failed: %s", action, result.VpcID, result.ErrMessage)
	}
	if result.ErrMessage != "" {
		log.Warnf("%s gateway insertion for VPC %s succeeded with message: %s", action, result.VpcID, result.ErrMessage)
	}
	return nil
}

// EnableGatewayInsertion enables gateway insertion for the VPC, or reapplies the
// route table and egress configuration when it is already enabled. An empty
// RouteTables list lets the controller select the VPC's private route tables.
func (c *Client) EnableGatewayInsertion(ctx context.Context, config *GatewayInsertionVpcConfig) error {
	req := GatewayInsertionRequest{VpcConfig: *config}
	if req.VpcConfig.RouteTables == nil {
		req.VpcConfig.RouteTables = []string{}
	}
	var result GatewayInsertionResult
	if err := c.DoAPIContext25(ctx, "PUT", &result, gatewayInsertionEndpoint, req); err != nil {
		return err
	}
	return checkGatewayInsertionResult("enabling", &result)
}

func (c *Client) GetGatewayInsertionStatus(ctx context.Context, vpcID string) (*GatewayInsertionStatus, error) {
	var statuses []GatewayInsertionStatus
	if err := c.GetAPIContext25(ctx, &statuses, gatewayInsertionEndpoint, map[string]string{"vpc_id": vpcID}); err != nil {
		return nil, err
	}
	for i := range statuses {
		if statuses[i].VpcID == vpcID {
			return &statuses[i], nil
		}
	}
	return nil, fmt.Errorf("gateway insertion status for VPC %s missing from controller response", vpcID)
}

// IsGatewayInsertionEnabled reports whether gateway insertion is enabled for the
// VPC. A controller that has the feature disabled (HTTP 501) or does not serve
// the endpoint (HTTP 404) cannot have it enabled, so both report false.
func (c *Client) IsGatewayInsertionEnabled(ctx context.Context, vpcID string) (bool, error) {
	status, err := c.GetGatewayInsertionStatus(ctx, vpcID)
	if err != nil {
		var statusErr StatusError
		if errors.As(err, &statusErr) &&
			(statusErr.StatusCode() == http.StatusNotImplemented || statusErr.StatusCode() == http.StatusNotFound) {
			return false, nil
		}
		return false, err
	}
	return status.Enabled, nil
}

func (c *Client) DisableGatewayInsertion(ctx context.Context, vpcID string) error {
	path := gatewayInsertionEndpoint + "?" + url.Values{"vpc_id": {vpcID}}.Encode()
	var result GatewayInsertionResult
	if err := c.DoAPIContext25(ctx, "DELETE", &result, path, nil); err != nil {
		return err
	}
	return checkGatewayInsertionResult("disabling", &result)
}
