package goaviatrix

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type GeoVPN struct {
	Action      string `json:"action,omitempty"`
	CID         string `json:"CID,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	CloudType   int    `json:"cloud_type,omitempty"`
	DomainName  string `json:"domain_name,omitempty"`
	ElbDNSName  string `json:"elb_dns_name,omitempty"`
	ServiceName string `json:"cname,omitempty"`
}

type GeoVPNEdit struct {
	AccountName string         `json:"account_name,omitempty"`
	CloudType   int            `json:"cloud_type,omitempty"`
	DomainName  string         `json:"domain_name,omitempty"`
	DnsName     string         `json:"service_name,omitempty"`
	ElbDNSNames []GeoVPNPolicy `json:"geo_vpn_policy,omitempty"`
	ServiceName string         `json:"cname,omitempty"`
}

type GeoVPNPolicy struct {
	ElbDNSName string `json:"elb_dns,omitempty"`
	Region     string `json:"region,omitempty"`
}

type GetGeoVPNInfoResp struct {
	Return  bool       `json:"return"`
	Results GeoVPNEdit `json:"results"`
	Reason  string     `json:"reason"`
}

func (c *Client) EnableGeoVPN(ctx context.Context, geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "enable_geo_vpn"

	return c.PostAPIContext2(ctx, nil, geoVPN.Action, geoVPN, BasicCheck)
}

func (c *Client) GetGeoVPNInfo(geoVPN *GeoVPN) (*GeoVPN, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "get_geo_vpn_info",
		"cloud_type": strconv.Itoa(geoVPN.CloudType),
	}

	var data GetGeoVPNInfoResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Geo VPN is not enbled") || strings.Contains(reason, "Geo VPN is not enabled") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	if data.Results.ServiceName == geoVPN.ServiceName && data.Results.DomainName == geoVPN.DomainName {
		geoVPN.AccountName = data.Results.AccountName
		geoVPN.ServiceName = data.Results.ServiceName
		geoVPN.DomainName = data.Results.DomainName
		elbDNSNameList := make([]string, 0)
		for i := 0; i < len(data.Results.ElbDNSNames); i++ {
			elbDNS := data.Results.ElbDNSNames[i]
			elbDNSNameList = append(elbDNSNameList, elbDNS.ElbDNSName)
		}
		geoVPN.ElbDNSName = strings.Join(elbDNSNameList, ",")
		return geoVPN, nil
	}

	log.Errorf("Couldn't find Aviatrix Geo VPN: %v", geoVPN)
	return nil, ErrNotFound
}

func (c *Client) AddElbToGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "add_elb_to_geo_vpn"

	return c.PostAPI(geoVPN.Action, geoVPN, BasicCheck)
}

func (c *Client) DeleteElbFromGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "delete_elb_from_geo_vpn"

	return c.PostAPI(geoVPN.Action, geoVPN, BasicCheck)
}

func (c *Client) DisableGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "disable_geo_vpn"

	return c.PostAPI(geoVPN.Action, geoVPN, BasicCheck)
}

func (c *Client) GetGeoVPNName(gateway *Gateway) (*GeoVPN, error) {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "get_geo_vpn_info",
		"cloud_type": strconv.Itoa(gateway.CloudType),
	}

	var data GetGeoVPNInfoResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Geo VPN is not enbled") || strings.Contains(reason, "Geo VPN is not enabled") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	policyList := data.Results.ElbDNSNames
	for i := range policyList {
		if policyList[i].ElbDNSName == gateway.ElbDNSName {
			geoVpn := &GeoVPN{
				CloudType:   data.Results.CloudType,
				AccountName: data.Results.AccountName,
				ServiceName: data.Results.DnsName,
				DomainName:  data.Results.DomainName,
			}
			return geoVpn, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix Geo VPN")
	return nil, ErrNotFound
}
