package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EdgeSpoke struct {
	Action                             string `json:"action,omitempty"`
	CID                                string `json:"CID,omitempty"`
	Type                               string `json:"type,omitempty"`
	GwName                             string `json:"gateway_name,omitempty"`
	SiteId                             string `json:"site_id,omitempty"`
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip,omitempty"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network,omitempty"`
	DnsServerIp                        string `json:"dns_server_ip,omitempty"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary,omitempty"`
	ZtpFileType                        string `json:"ztp_file_type,omitempty"`
	ZtpFileDownloadPath                string
	ActiveStandby                      string `json:"active_standby,omitempty"`
	EnableEdgeActiveStandby            bool   `json:"enable_active_standby,omitempty"`
	DisableEdgeActiveStandby           bool   `json:"disable_active_standby,omitempty"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"enable_active_standby_preemptive,omitempty"`
	DisableEdgeActiveStandbyPreemptive bool   `json:"disable_active_standby_preemptive,omitempty"`
	LocalAsNumber                      string `json:"local_as_number,omitempty"`
	PrependAsPath                      []string
	PrependAsPathReturn                string   `json:"prepend_as_path,omitempty"`
	IncludeCidrList                    []string `json:"include_cidr_list,omitempty"`
	EnableLearnedCidrsApproval         bool     `json:"enable_learned_cidrs_approval,omitempty"`
	ApprovedLearnedCidrs               []string `form:"approved_learned_cidrs,omitempty,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string `json:"bgp_manual_spoke_advertise_cidrs,omitempty"`
	EnablePreserveAsPath               bool     `json:"preserve_as_path,omitempty"`
	BgpPollingTime                     int      `json:"bgp_polling_time,omitempty"`
	BgpBfdPollingTime                  int      `json:"bgp_neighbor_status_polling_time,omitempty"`
	BgpHoldTime                        int      `json:"bgp_hold_time,omitempty"`
	EnableEdgeTransitiveRouting        bool     `json:"edge_transitive_routing,omitempty"`
	EnableJumboFrame                   bool     `json:"jumbo_frame,omitempty"`
	Latitude                           string
	Longitude                          string
	LatitudeReturn                     float64 `json:"latitude,omitempty"`
	LongitudeReturn                    float64 `json:"longitude,omitempty"`
	RxQueueSize                        string  `json:"rx_queue_size,omitempty"`
	State                              string  `json:"vpc_state,omitempty"`
	InterfaceList                      []*EdgeSpokeInterface
	Interfaces                         string `json:"interfaces,omitempty"`
	VlanList                           []*EdgeSpokeVlan
	Vlan                               string                        `json:"vlan,omitempty"`
	CustomInterfaceMapping             map[string]CustomInterfaceMap `json:"custom_interface_mapping,omitempty"`
	AdvertisedCidrList                 []string                      `json:"advertise_cidr_list,omitempty"`
	TunnelEncryptionCipher             string                        `json:"ph2_encryption_policy,omitempty"`
	TunnelForwardSecrecy               string                        `json:"ph2_pfs_policy,omitempty"`
}

type EdgeSpokeInterface struct {
	IfName        string           `json:"ifname"`
	Type          string           `json:"type"`
	Dhcp          bool             `json:"dhcp,omitempty"`
	PublicIp      string           `json:"public_ip,omitempty"`
	IpAddr        string           `json:"ipaddr,omitempty"`
	GatewayIp     string           `json:"gateway_ip,omitempty"`
	DNSPrimary    string           `json:"dns_primary,omitempty"`
	DNSSecondary  string           `json:"dns_secondary,omitempty"`
	SubInterfaces []*EdgeSpokeVlan `json:"subinterfaces,omitempty"`
	VrrpState     bool             `json:"vrrp_state,omitempty"`
	VirtualIp     string           `json:"virtual_ip,omitempty"`
	Tag           string           `json:"tag,omitempty"`
	IPv6Addr      string           `json:"v6_ipaddr,omitempty"`
	GatewayIPv6   string           `json:"gateway_ipv6_ip,omitempty"`
}

type CustomInterfaceMap struct {
	IdentifierType  string `json:"identifier_type"`
	IdentifierValue string `json:"identifier_value"`
}

type EdgeSpokeVlan struct {
	ParentInterface string `json:"parent_interface"`
	VlanId          string `json:"vlan_id,omitempty"`
	IpAddr          string `json:"ipaddr,omitempty"`
	GatewayIp       string `json:"gateway_ip,omitempty"`
	PeerIpAddr      string `json:"peer_ipaddr,omitempty"`
	PeerGatewayIp   string `json:"peer_gateway_ip,omitempty"`
	VirtualIp       string `json:"virtual_ip,omitempty"`
	Tag             string `json:"tag,omitempty"`
}

type EdgeSpokeResp struct {
	GwName                             string `json:"gw_name"`
	SiteId                             string `json:"vpc_id"`
	CloudType                          int    `json:"cloud_type"`
	ManagementEgressIpPrefix           string `json:"mgmt_egress_ip"`
	EnableManagementOverPrivateNetwork bool   `json:"mgmt_over_private_network"`
	DnsServerIp                        string `json:"dns_server_ip"`
	SecondaryDnsServerIp               string `json:"dns_server_ip_secondary"`
	ZtpFileType                        string `json:"ztp_file_type"`
	EnableEdgeActiveStandby            bool   `json:"edge_active_standby"`
	EnableEdgeActiveStandbyPreemptive  bool   `json:"edge_active_standby_preemptive"`
	LocalAsNumber                      string `json:"local_as_number"`
	PrependAsPath                      []string
	PrependAsPathReturn                string                        `json:"prepend_as_path"`
	IncludeCidrList                    []string                      `json:"include_cidr_list"`
	EnableLearnedCidrsApproval         bool                          `json:"enable_learned_cidrs_approval"`
	ApprovedLearnedCidrs               []string                      `form:"approved_learned_cidrs,omitempty"`
	SpokeBgpManualAdvertisedCidrs      []string                      `json:"bgp_manual_spoke_advertise_cidrs"`
	EnablePreserveAsPath               bool                          `json:"preserve_as_path"`
	BgpPollingTime                     int                           `json:"bgp_polling_time"`
	BgpBfdPollingTime                  int                           `json:"bgp_neighbor_status_polling_time"`
	BgpHoldTime                        int                           `json:"bgp_hold_time"`
	EnableEdgeTransitiveRouting        bool                          `json:"edge_transitive_routing"`
	EnableJumboFrame                   bool                          `json:"jumbo_frame"`
	Latitude                           float64                       `json:"latitude"`
	Longitude                          float64                       `json:"longitude"`
	RxQueueSize                        string                        `json:"rx_queue_size"`
	State                              string                        `json:"vpc_state"`
	InterfaceList                      []*EdgeSpokeInterface         `json:"interfaces"`
	CustomInterfaceMapping             map[string]CustomInterfaceMap `json:"custom_interface_mapping,omitempty"`
	AdvertisedCidrList                 []string                      `json:"advertise_cidr_list,omitempty"`
	TunnelEncryptionCipher             string                        `form:"ph2_encryption_policy,omitempty"`
	TunnelForwardSecrecy               string                        `form:"ph2_pfs_policy,omitempty"`
}

type EdgeSpokeListResp struct {
	Return  bool            `json:"return"`
	Results []EdgeSpokeResp `json:"results"`
	Reason  string          `json:"reason"`
}

func (c *Client) CreateEdgeSpoke(ctx context.Context, edgeSpoke *EdgeSpoke) error {
	edgeSpoke.Action = "create_edge_gateway"
	edgeSpoke.CID = c.CID
	edgeSpoke.Type = "spoke"

	interfaces, err := json.Marshal(edgeSpoke.InterfaceList)
	if err != nil {
		return err
	}

	edgeSpoke.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if len(edgeSpoke.VlanList) == 0 {
		edgeSpoke.VlanList = []*EdgeSpokeVlan{}
	}

	vlan, err := json.Marshal(edgeSpoke.VlanList)
	if err != nil {
		return err
	}

	edgeSpoke.Vlan = b64.StdEncoding.EncodeToString(vlan)

	resp, err := c.PostAPIContext2Download(ctx, edgeSpoke.Action, edgeSpoke, BasicCheck)
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

// EdgeSpokeHa represents an HA edge spoke gateway
type EdgeSpokeHa struct {
	Action                   string `form:"action" json:"action"`
	CID                      string `form:"CID" json:"CID"`
	GroupUUID                string `form:"group_uuid,omitempty" json:"group_uuid,omitempty"`
	PrimaryGwName            string `form:"primary_gw_name,omitempty" json:"primary_gw_name,omitempty"`
	SiteID                   string `form:"site_id,omitempty" json:"site_id,omitempty"`
	ZtpFileType              string `form:"ztp_file_type,omitempty" json:"ztp_file_type,omitempty"`
	ZtpFileDownloadPath      string `form:"-"`
	InterfaceList            []*EdgeSpokeInterface
	Interfaces               string `form:"interfaces,omitempty" json:"interfaces"`
	NoProgressBar            bool   `form:"no_progress_bar,omitempty" json:"no_progress_bar,omitempty"`
	ManagementEgressIPPrefix string `form:"mgmt_egress_ip,omitempty" json:"mgmt_egress_ip,omitempty"`
	CloudInit                bool   `form:"cloud_init,omitempty" json:"cloud_init,omitempty"`
}

// CreateEdgeSpokeHa creates an HA edge spoke gateway
func (c *Client) CreateEdgeSpokeHa(ctx context.Context, edgeSpokeHa *EdgeSpokeHa) (string, error) {
	edgeSpokeHa.CID = c.CID
	edgeSpokeHa.Action = "create_multicloud_ha_gateway"
	edgeSpokeHa.NoProgressBar = true

	if edgeSpokeHa.ZtpFileType == "iso" {
		edgeSpokeHa.CloudInit = false
	} else {
		edgeSpokeHa.CloudInit = true
	}

	interfaces, err := json.Marshal(edgeSpokeHa.InterfaceList)
	if err != nil {
		return "", err
	}

	edgeSpokeHa.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	type CreateEdgeSpokeHaResp struct {
		Return bool   `json:"return"`
		Result string `json:"results"`
		Reason string `json:"reason"`
	}

	var data CreateEdgeSpokeHaResp

	gwName, err := c.PostAPIContext2HaGw(ctx, &data, edgeSpokeHa.Action, edgeSpokeHa, BasicCheck)
	if err != nil {
		return "", err
	}

	// Write ZTP file if download path is provided
	if edgeSpokeHa.ZtpFileDownloadPath != "" {
		var fileName string
		if edgeSpokeHa.ZtpFileType == "iso" {
			fileName = edgeSpokeHa.ZtpFileDownloadPath + "/" + edgeSpokeHa.SiteID + "-hagw.iso"
		} else {
			fileName = edgeSpokeHa.ZtpFileDownloadPath + "/" + edgeSpokeHa.SiteID + "-hagw-cloud-init.txt"
		}

		outFile, err := os.Create(fileName)
		if err != nil {
			return gwName, err
		}
		defer outFile.Close()

		if edgeSpokeHa.ZtpFileType == "iso" {
			decodedResult, err := b64.StdEncoding.DecodeString(data.Result)
			if err != nil {
				return gwName, err
			}
			_, err = outFile.Write(decodedResult)
			if err != nil {
				return gwName, err
			}
		} else {
			_, err = outFile.WriteString(data.Result)
			if err != nil {
				return gwName, err
			}
		}
	}

	return gwName, nil
}

func (c *Client) GetEdgeSpoke(ctx context.Context, gwName string) (*EdgeSpokeResp, error) {
	form := map[string]string{
		"action":       "list_vpcs_summary",
		"CID":          c.CID,
		"gateway_name": gwName,
	}

	var data EdgeSpokeListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
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

func (c *Client) UpdateEdgeSpoke(ctx context.Context, edgeSpoke *EdgeSpoke) error {
	edgeSpoke.Action = "update_edge_gateway"
	edgeSpoke.CID = c.CID

	interfaces, err := json.Marshal(edgeSpoke.InterfaceList)
	if err != nil {
		return err
	}

	edgeSpoke.Interfaces = b64.StdEncoding.EncodeToString(interfaces)

	if len(edgeSpoke.VlanList) == 0 {
		edgeSpoke.VlanList = []*EdgeSpokeVlan{}
	}

	vlan, err := json.Marshal(edgeSpoke.VlanList)
	if err != nil {
		return err
	}

	edgeSpoke.Vlan = b64.StdEncoding.EncodeToString(vlan)

	return c.PostAPIContext2(ctx, nil, edgeSpoke.Action, edgeSpoke, BasicCheck)
}

func (c *Client) DeleteEdgeSpoke(ctx context.Context, name string) error {
	form := map[string]string{
		"action": "delete_edge_gateway",
		"CID":    c.CID,
		"name":   name,
	}

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
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
	s, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%s must be a string", key))
		return
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s must be a valid number", key))
		return
	}

	if v < -90 || v > 90 {
		errs = append(errs, fmt.Errorf("%s must be between -90 and 90", key))
	}
	return
}

func ValidateEdgeSpokeLongitude(val interface{}, key string) (warns []string, errs []error) {
	s, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%s must be a string", key))
		return
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s must be a valid number", key))
		return
	}

	if v < -180 || v > 180 {
		errs = append(errs, fmt.Errorf("%s must be between -180 and 180", key))
	}
	return
}

func DiffSuppressFuncEdgeSpokeCoordinate(k, old, new string, d *schema.ResourceData) bool {
	o, _ := strconv.ParseFloat(old, 64)
	n, _ := strconv.ParseFloat(new, 64)
	return math.Round(o*1000000)/1000000 == math.Round(n*1000000)/1000000
}
