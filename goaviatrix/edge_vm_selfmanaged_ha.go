package goaviatrix

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"os"
)

type EdgeVmSelfmanagedHa struct {
	Action                   string `json:"action"`
	CID                      string `json:"CID"`
	PrimaryGwName            string `json:"primary_gw_name"`
	SiteId                   string
	ZtpFileType              string
	ZtpFileDownloadPath      string
	DnsServerIp              string `json:"dns_server_ip,omitempty"`
	SecondaryDnsServerIp     string `json:"dns_server_ip_secondary,omitempty"`
	InterfaceList            []*EdgeSpokeInterface
	Interfaces               string `json:"interfaces"`
	NoProgressBar            bool   `json:"no_progress_bar,omitempty"`
	ManagementEgressIPPrefix string `json:"mgmt_egress_ip,omitempty"`
	CloudInit                bool   `json:"cloud_init"`
}

type EdgeVmSelfmanagedHaResp struct {
	PrimaryGwName            string                `json:"primary_gw_name"`
	GwName                   string                `json:"gw_name"`
	SiteID                   string                `json:"vpc_id"`
	ZtpFileType              string                `json:"ztp_file_type"`
	InterfaceList            []*EdgeSpokeInterface `json:"interfaces"`
	ManagementEgressIPPrefix string                `json:"mgmt_egress_ip"`
	DNSServerIP              string                `json:"dns_server_ip,omitempty"`
	SecondaryDNSServerIP     string                `json:"dns_server_ip_secondary,omitempty"`
}

type EdgeVmSelfmanagedHaListResp struct {
	Return  bool                      `json:"return"`
	Results []EdgeVmSelfmanagedHaResp `json:"results"`
	Reason  string                    `json:"reason"`
}

type CreateEdgeVmSelfmanagedHaResp struct {
	Return bool   `json:"return"`
	Result string `json:"results"`
	Reason string `json:"reason"`
}

func (c *Client) CreateEdgeVmSelfmanagedHa(ctx context.Context, edgeVmSelfmanagedHa *EdgeVmSelfmanagedHa) (string, error) {
	edgeVmSelfmanagedHa.CID = c.CID
	edgeVmSelfmanagedHa.Action = "create_multicloud_ha_gateway"
	edgeVmSelfmanagedHa.NoProgressBar = true

	if edgeVmSelfmanagedHa.ZtpFileType == "iso" {
		edgeVmSelfmanagedHa.CloudInit = false
	} else {
		edgeVmSelfmanagedHa.CloudInit = true
	}

	interfaces, err := json.Marshal(edgeVmSelfmanagedHa.InterfaceList)
	if err != nil {
		return "", err
	}

	edgeVmSelfmanagedHa.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	var data CreateEdgeVmSelfmanagedHaResp

	gwName, err := c.PostAPIContext2HaGw(ctx, &data, edgeVmSelfmanagedHa.Action, edgeVmSelfmanagedHa, BasicCheck)
	if err != nil {
		return "", err
	}

	var fileName string
	if edgeVmSelfmanagedHa.ZtpFileType == "iso" {
		fileName = edgeVmSelfmanagedHa.ZtpFileDownloadPath + "/" + edgeVmSelfmanagedHa.PrimaryGwName + "-" + edgeVmSelfmanagedHa.SiteId + "-ha.iso"
	} else {
		fileName = edgeVmSelfmanagedHa.ZtpFileDownloadPath + "/" + edgeVmSelfmanagedHa.PrimaryGwName + "-" + edgeVmSelfmanagedHa.SiteId + "-ha-cloud-init.txt"
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	if edgeVmSelfmanagedHa.ZtpFileType == "iso" {
		decodedResult, err := b64.StdEncoding.DecodeString(data.Result)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, bytes.NewReader(decodedResult))
		if err != nil {
			return "", err
		}
	} else {
		_, err = outFile.WriteString(data.Result)
		if err != nil {
			return "", err
		}
	}

	return gwName, nil
}

func (c *Client) GetEdgeVmSelfmanagedHa(ctx context.Context, gwName string) (*EdgeVmSelfmanagedHaResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeVmSelfmanagedHaListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeVmSelfmanagedHaList := data.Results
	for _, edgeVmSelfmanagedHa := range edgeVmSelfmanagedHaList {
		if edgeVmSelfmanagedHa.GwName == gwName {
			return &edgeVmSelfmanagedHa, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeVmSelfmanagedHa(ctx context.Context, edgeVmSelfmanaged *EdgeSpoke) error {
	form := map[string]string{
		"action":         "update_edge_gateway",
		"CID":            c.CID,
		"name":           edgeVmSelfmanaged.GwName,
		"mgmt_egress_ip": edgeVmSelfmanaged.ManagementEgressIpPrefix,
	}

	interfaces, err := json.Marshal(edgeVmSelfmanaged.InterfaceList)
	if err != nil {
		return err
	}

	form["interfaces"] = b64.StdEncoding.EncodeToString(interfaces)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
