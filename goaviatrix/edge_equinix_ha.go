package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"os"
)

type EdgeEquinixHa struct {
	Action                   string `json:"action"`
	CID                      string `json:"CID"`
	PrimaryGwName            string `json:"primary_gw_name"`
	ZtpFileDownloadPath      string
	InterfaceList            []*EdgeEquinixInterface
	Interfaces               string `json:"interfaces"`
	NoProgressBar            bool   `json:"no_progress_bar,omitempty"`
	ManagementEgressIpPrefix string `json:"mgmt_egress_ip,omitempty"`
}

type EdgeEquinixHaResp struct {
	AccountName              string                  `json:"account_name"`
	PrimaryGwName            string                  `json:"primary_gw_name"`
	GwName                   string                  `json:"gw_name"`
	InterfaceList            []*EdgeEquinixInterface `json:"interfaces"`
	ManagementEgressIpPrefix string                  `json:"mgmt_egress_ip"`
}

type EdgeEquinixHaListResp struct {
	Return  bool                `json:"return"`
	Results []EdgeEquinixHaResp `json:"results"`
	Reason  string              `json:"reason"`
}

func (c *Client) CreateEdgeEquinixHa(ctx context.Context, edgeEquinixHa *EdgeEquinixHa) (string, error) {
	edgeEquinixHa.CID = c.CID
	edgeEquinixHa.Action = "create_multicloud_ha_gateway"
	edgeEquinixHa.NoProgressBar = true

	interfaces, err := json.Marshal(edgeEquinixHa.InterfaceList)
	if err != nil {
		return "", err
	}

	edgeEquinixHa.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	var data CreateEdgeEquinixResp

	gwName, err := c.PostAPIContext2HaGw(ctx, &data, edgeEquinixHa.Action, edgeEquinixHa, BasicCheck)
	if err != nil {
		return "", err
	}

	fileName := edgeEquinixHa.ZtpFileDownloadPath + "/" + edgeEquinixHa.PrimaryGwName + "-hagw-cloud-init.txt"

	outFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	_, err = outFile.WriteString(data.Result)
	if err != nil {
		return "", err
	}

	return gwName, nil
}

func (c *Client) GetEdgeEquinixHa(ctx context.Context, gwName string) (*EdgeEquinixHaResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeEquinixHaListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeEquinixHaList := data.Results
	for _, edgeEquinixHa := range edgeEquinixHaList {
		if edgeEquinixHa.GwName == gwName {
			return &edgeEquinixHa, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeEquinixHa(ctx context.Context, edgeEquinix *EdgeEquinix) error {
	form := map[string]string{
		"action":         "update_edge_gateway",
		"CID":            c.CID,
		"name":           edgeEquinix.GwName,
		"mgmt_egress_ip": edgeEquinix.ManagementEgressIpPrefix,
	}

	interfaces, err := json.Marshal(edgeEquinix.InterfaceList)
	if err != nil {
		return err
	}

	form["interfaces"] = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
