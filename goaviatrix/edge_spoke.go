package goaviatrix

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EdgeSpoke struct {
	Action                             string `form:"action,omitempty"`
	CID                                string `form:"CID,omitempty"`
	Type                               string `form:"type,omitempty"`
	Caag                               bool   `form:"caag,omitempty"`
	GwName                             string `form:"gateway_name,omitempty" json:"gw_name"`
	SiteId                             string `form:"site_id,omitempty" json:"vpc_id"`
	ManagementInterfaceConfig          string
	ManagementEgressIpPrefix           string `form:"mgmt_egress_ip" json:"mgmt_egress_ip"`
	EnableManagementOverPrivateNetwork bool   `form:"mgmt_over_private_network,omitempty" json:"mgmt_over_private_network"`
	WanInterfaceIpPrefix               string `form:"wan_ip,omitempty" json:"wan_ip"`
	WanDefaultGatewayIp                string `form:"wan_default_gateway,omitempty" json:"wan_default_gateway"`
	LanInterfaceIpPrefix               string `form:"lan_ip,omitempty" json:"lan_ip"`
	ManagementInterfaceIpPrefix        string `form:"mgmt_ip,omitempty" json:"mgmt_ip"`
	ManagementDefaultGatewayIp         string `form:"mgmt_default_gateway,omitempty" json:"mgmt_default_gateway"`
	DnsServerIp                        string `form:"dns_server_ip,omitempty" json:"dns_server_ip"`
	SecondaryDnsServerIp               string `form:"dns_server_ip_secondary,omitempty" json:"dns_server_ip_secondary"`
	Dhcp                               bool   `form:"dhcp,omitempty" json:"dhcp"`
	ZtpFileType                        string `form:"ztp_file_type,omitempty"`
	ZtpFileDownloadPath                string
	ActiveStandby                      string `form:"active_standby,omitempty"`
	EnableEdgeActiveStandby            bool   `json:"edge_active_standby"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"edge_active_standby_preemptive"`
	LocalAsNumber                      string `json:"local_as_number"`
	PrependAsPath                      []string
	PrependAsPathReturn                string   `json:"prepend_as_path"`
	IncludeCidrList                    []string `json:"include_cidr_list"`
	EnableLearnedCidrsApproval         bool     `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs               []string `form:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath               bool     `json:"preserve_as_path"`
	BgpPollingTime                     int      `json:"bgp_polling_time"`
	BgpHoldTime                        int      `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing"`
	EnableJumboFrame                   bool     `json:"jumbo_frame"`
	Latitude                           string
	Longitude                          string
	LatitudeReturn                     float64 `json:"latitude"`
	LongitudeReturn                    float64 `json:"longitude"`
	WanPublicIp                        string  `form:"wan_discovery_ip" json:"public_ip"`
	PrivateIP                          string  `json:"private_ip"`
	RxQueueSize                        string  `json:"rx_queue_size"`
	State                              string  `json:"vpc_state"`
}

type EdgeSpokeListResp struct {
	Return  bool        `json:"return"`
	Results []EdgeSpoke `json:"results"`
	Reason  string      `json:"reason"`
}

func (c *Client) CreateEdgeSpoke(ctx context.Context, edgeSpoke *EdgeSpoke) error {
	edgeSpoke.Action = "create_edge_gateway"
	edgeSpoke.CID = c.CID
	edgeSpoke.Type = "spoke"
	edgeSpoke.Caag = false

	if edgeSpoke.ManagementInterfaceConfig == "DHCP" {
		edgeSpoke.Dhcp = true
	}

	resp, err := c.PostAPIDownloadContext(ctx, edgeSpoke.Action, edgeSpoke, BasicCheck)
	if err != nil {
		return err
	}

	var fileName string
	if edgeSpoke.ZtpFileType == "iso" {
		fileName = edgeSpoke.ZtpFileDownloadPath + "/" + edgeSpoke.GwName + "-" + edgeSpoke.SiteId + ".iso"
	} else {
		fileName = edgeSpoke.ZtpFileDownloadPath + "/" + edgeSpoke.GwName + "-" + edgeSpoke.SiteId + "-cloud-init.txt"
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeSpoke(ctx context.Context, gwName string) (*EdgeSpoke, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeSpokeListResp

	err := c.GetAPIContext(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeSpokeList := data.Results
	for _, edgeSpoke := range edgeSpokeList {
		if edgeSpoke.GwName == gwName {
			for _, p := range strings.Split(edgeSpoke.PrependAsPathReturn, " ") {
				if p != "" {
					edgeSpoke.PrependAsPath = append(edgeSpoke.PrependAsPath, p)
				}
			}

			return &edgeSpoke, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateEdgeSpokeIpConfigurations(ctx context.Context, edgeSpoke *EdgeSpoke) error {
	form := map[string]string{
		"action":              "update_edge_gateway",
		"CID":                 c.CID,
		"gateway_name":        edgeSpoke.GwName,
		"wan_ip":              edgeSpoke.WanInterfaceIpPrefix,
		"wan_default_gateway": edgeSpoke.WanDefaultGatewayIp,
		"lan_ip":              edgeSpoke.LanInterfaceIpPrefix,
		"mgmt_egress_ip":      edgeSpoke.ManagementEgressIpPrefix,
		"wan_discovery_ip":    edgeSpoke.WanPublicIp,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) DeleteEdgeSpoke(ctx context.Context, name string) error {
	form := map[string]string{
		"action": "delete_edge_gateway",
		"CID":    c.CID,
		"name":   name,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) EnableEdgeSpokeTransitiveRouting(ctx context.Context, name string) error {
	form := map[string]string{
		"action":       "enable_edge_transitive_routing",
		"CID":          c.CID,
		"gateway_name": name,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) DisableEdgeSpokeTransitiveRouting(ctx context.Context, name string) error {
	form := map[string]string{
		"action":       "disable_edge_transitive_routing",
		"CID":          c.CID,
		"gateway_name": name,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) UpdateEdgeSpokeGeoCoordinate(ctx context.Context, edgeSpoke *EdgeSpoke) error {
	form := map[string]string{
		"action":        "update_edge_gateway",
		"CID":           c.CID,
		"gateway_name":  edgeSpoke.GwName,
		"geo_latitude":  edgeSpoke.Latitude,
		"geo_longitude": edgeSpoke.Longitude,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func ValidateEdgeSpokeLatitude(val interface{}, key string) (warns []string, errs []error) {
	v, _ := strconv.ParseFloat(val.(string), 64)
	if v < -90 || v > 90 {
		errs = append(errs, fmt.Errorf("latitude must be between -90 and 90"))
	}
	return
}

func ValidateEdgeSpokeLongitude(val interface{}, key string) (warns []string, errs []error) {
	v, _ := strconv.ParseFloat(val.(string), 64)
	if v < -180 || v > 180 {
		errs = append(errs, fmt.Errorf("longitude must be between -180 and 180"))
	}
	return
}

func DiffSuppressFuncEdgeSpokeCoordinate(k, old, new string, d *schema.ResourceData) bool {
	o, _ := strconv.ParseFloat(old, 64)
	n, _ := strconv.ParseFloat(new, 64)
	return math.Round(o*1000000)/1000000 == math.Round(n*1000000)/1000000
}
