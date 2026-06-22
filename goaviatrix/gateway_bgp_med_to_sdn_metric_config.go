package goaviatrix

import (
	"fmt"
	"strings"
)

// This file contains the API calls related to gateway BGP MED to SDN metric global config
// Details on the feature and related APIs can be found here:
// Design doc: https://docs.google.com/document/d/1h-gxgwZ6OxNuLNLFgKymldqdoRU3kY1zEshbZs13C0A/edit?usp=sharing
// APIs: https://aviatrix.atlassian.net/wiki/spaces/AVXENG/pages/3136946177/BGP+MED+to+SDN+Metric+APIs

func (c *Client) SetGatewayBgpMedToSdnMetric(gwName string, override bool) error {
	data := map[string]string{
		"action":           "set_gateway_accept_bgp_med_to_sdn_metric_override",
		"CID":              c.CID,
		"gateway_name":     gwName,
		"gateway_override": fmt.Sprint(override),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) GetGatewayBgpMedToSdnMetric(gwName string) (bool, error) {
	data := map[string]string{
		"action":       "show_bgp_med_to_sdn_metric_gateway_override",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	type BgpMedToSdnMetricGatewayResults struct {
		BgpMedToSdnMetricGateway     string `json:"set_bgp_med_to_sdn_metric"`
		BgpMedToSdnMetricGatewayText string `json:"text"`
	}

	type BgpMedToSdnMetricGatewayResponse struct {
		Results BgpMedToSdnMetricGatewayResults
	}

	var resp BgpMedToSdnMetricGatewayResponse
	err := c.GetAPI(&resp, data["action"], data, BasicCheck)
	if err != nil {
		return false, err
	}

	// Somehow the API returns "true" or "false" as strings, so we need to convert them to bool
	var override bool
	switch strings.ToLower(resp.Results.BgpMedToSdnMetricGateway) {
	case "true":
		override = true
	case "false":
		override = false
	}
	return override, nil
}
