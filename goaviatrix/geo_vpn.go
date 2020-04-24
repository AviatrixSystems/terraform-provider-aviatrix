package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GeoVPN struct {
	AccountName string   `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action      string   `form:"action,omitempty"`
	CID         string   `form:"CID,omitempty"`
	CloudType   int      `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	ServiceName string   `form:"cname,omitempty" json:"cname,omitempty"`
	DomainName  string   `form:"domain_name,omitempty" json:"domain_name,omitempty"`
	ElbDNSName  string   `form:"elb_dns_name,omitempty"`
	ElbDNSNames []string `json:"elb_dns_name,omitempty"`
}

type GeoVPNEdit struct {
	AccountName string         `json:"account_name,omitempty"`
	CloudType   int            `json:"cloud_type,omitempty"`
	ServiceName string         `json:"cname,omitempty"`
	DomainName  string         `json:"domain_name,omitempty"`
	ElbDNSNames []GeoVPNPolicy `json:"geo_vpn_policy,omitempty"`
	DnsName     string         `json:"service_name,omitempty"`
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

func (c *Client) EnableGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "enable_geo_vpn"
	resp, err := c.Post(c.baseURL, geoVPN)
	if err != nil {
		return errors.New("HTTP Post 'enable_geo_vpn' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'enable_geo_vpn' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'enable_geo_vpn' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetGeoVPNInfo(geoVPN *GeoVPN) (*GeoVPN, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'get_geo_vpn_info' ") + err.Error())
	}
	getGeoVPNInfo := url.Values{}
	getGeoVPNInfo.Add("CID", c.CID)
	getGeoVPNInfo.Add("action", "get_geo_vpn_info")
	getGeoVPNInfo.Add("cloud_type", strconv.Itoa(geoVPN.CloudType))
	Url.RawQuery = getGeoVPNInfo.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'get_geo_vpn_info' failed: " + err.Error())
	}
	var data GetGeoVPNInfoResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_geo_vpn_info failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "Geo VPN is not enbled") || strings.Contains(data.Reason, "Geo VPN is not enabled") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API get_geo_vpn_info Get failed: " + data.Reason)
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
		geoVPN.ElbDNSNames = elbDNSNameList
		return geoVPN, nil
	}

	log.Errorf("Couldn't find Aviatrix Geo VPN: %v", geoVPN)
	return nil, ErrNotFound
}

func (c *Client) AddElbToGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "add_elb_to_geo_vpn"
	resp, err := c.Post(c.baseURL, geoVPN)
	if err != nil {
		return errors.New("HTTP Post 'add_elb_to_geo_vpn' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_elb_to_geo_vpn' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_elb_to_geo_vpn' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteElbFromGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "delete_elb_from_geo_vpn"
	resp, err := c.Post(c.baseURL, geoVPN)
	if err != nil {
		return errors.New("HTTP Post 'delete_elb_from_geo_vpn' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_elb_from_geo_vpn' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_elb_from_geo_vpn' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableGeoVPN(geoVPN *GeoVPN) error {
	geoVPN.CID = c.CID
	geoVPN.Action = "disable_geo_vpn"
	resp, err := c.Post(c.baseURL, geoVPN)
	if err != nil {
		return errors.New("HTTP Post 'disable_geo_vpn' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disable_geo_vpn' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disable_geo_vpn' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetGeoVPNName(gateway *Gateway) (*GeoVPN, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'get_geo_vpn_info': ") + err.Error())
	}
	getGeoVPNInfo := url.Values{}
	getGeoVPNInfo.Add("CID", c.CID)
	getGeoVPNInfo.Add("action", "get_geo_vpn_info")
	getGeoVPNInfo.Add("cloud_type", strconv.Itoa(gateway.CloudType))
	Url.RawQuery = getGeoVPNInfo.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'get_geo_vpn_info' failed: " + err.Error())
	}
	var data GetGeoVPNInfoResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json 'Decode get_geo_vpn_info' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "Geo VPN is not enbled") || strings.Contains(data.Reason, "Geo VPN is not enabled") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API 'get_geo_vpn_info Get' failed: " + data.Reason)
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
