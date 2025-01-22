package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type EdgeMegaportHa struct {
	Action                   string `json:"action"`
	CID                      string `json:"CID"`
	PrimaryGwName            string `json:"primary_gw_name"`
	ZtpFileDownloadPath      string
	InterfaceList            []*EdgeMegaportInterface
	Interfaces               string `json:"interfaces"`
	NoProgressBar            bool   `json:"no_progress_bar,omitempty"`
	ManagementEgressIPPrefix string `json:"mgmt_egress_ip,omitempty"`
}

type EdgeMegaportHaResp struct {
	AccountName              string                   `json:"account_name"`
	PrimaryGwName            string                   `json:"primary_gw_name"`
	GwName                   string                   `json:"gw_name"`
	InterfaceList            []*EdgeMegaportInterface `json:"interfaces"`
	ManagementEgressIPPrefix string                   `json:"mgmt_egress_ip"`
}

type EdgeMegaportHaListResp struct {
	Return  bool                 `json:"return"`
	Results []EdgeMegaportHaResp `json:"results"`
	Reason  string               `json:"reason"`
}

func (c *Client) CreateEdgeMegaportHa(ctx context.Context, edgeMegaportHa *EdgeMegaportHa) (string, error) {
	edgeMegaportHa.CID = c.CID
	edgeMegaportHa.Action = "create_multicloud_ha_gateway"
	edgeMegaportHa.NoProgressBar = true

	interfaces, err := json.Marshal(edgeMegaportHa.InterfaceList)
	if err != nil {
		return "", err
	}

	edgeMegaportHa.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	var data CreateEdgeMegaportResp

	gwName, err := c.PostAPIContext2HaGw(ctx, &data, edgeMegaportHa.Action, edgeMegaportHa, BasicCheck)
	if err != nil {
		return "", err
	}

	fileName := edgeMegaportHa.ZtpFileDownloadPath + "/" + edgeMegaportHa.PrimaryGwName + "-hagw-cloud-init.txt"

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

func (c *Client) GetEdgeMegaportHa(ctx context.Context, gwName string) (*EdgeMegaportHaResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeMegaportHaListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeMegaportHaList := data.Results
	for _, edgeMegaportHa := range edgeMegaportHaList {
		if edgeMegaportHa.GwName == gwName {
			return &edgeMegaportHa, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeMegaportHa(ctx context.Context, edgeMegaport *EdgeMegaport) error {
	form := map[string]string{
		"action":         "update_edge_gateway",
		"CID":            c.CID,
		"name":           edgeMegaport.GwName,
		"mgmt_egress_ip": edgeMegaport.ManagementEgressIPPrefix,
		"cloud_type":     fmt.Sprintf("%v", EDGEMEGAPORT),
	}

	interfaces, err := json.Marshal(edgeMegaport.InterfaceList)
	if err != nil {
		return err
	}

	form["interfaces"] = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
