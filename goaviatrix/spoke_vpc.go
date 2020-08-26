package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Spoke gateway simple struct to hold spoke details
type SpokeVpc struct {
	AccountName           string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                string `form:"action,omitempty"`
	CID                   string `form:"CID,omitempty"`
	CloudType             int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer             string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                string `form:"gw_size,omitempty"`
	VpcID                 string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup string `form:"vnet_and_resource_group_names,omitempty"`
	Subnet                string `form:"public_subnet,omitempty" json:"public_subnet,omitempty"`
	VpcRegion             string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize               string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	EnableNat             string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	EnableVpcDnsServer    string `json:"use_vpc_dns,omitempty"`
	HASubnet              string `form:"ha_subnet,omitempty"`
	HAZone                string `form:"new_zone,omitempty"`
	HASubnetGCP           string `form:"new_subnet,omitempty"`
	SingleAzHa            string `form:"single_az_ha,omitempty"`
	TransitGateway        string `form:"transit_gw,omitempty"`
	TagList               string `form:"tags,omitempty"`
	ReuseEip              string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead    bool   `json:"newly_allocated_eip,omitempty"`
	Eip                   string `form:"eip,omitempty" json:"eip,omitempty"`
	InsaneMode            string `form:"insane_mode,omitempty"`
	EnableActiveMesh      string `form:"enable_activemesh,omitempty" json:"enable_activemesh,omitempty"`
	Zone                  string `form:"zone,omitempty" json:"zone,omitempty"`
	EncVolume             string `form:"enc_volume,omitempty"`
}

func (c *Client) LaunchSpokeVpc(spoke *SpokeVpc) error {
	spoke.CID = c.CID
	spoke.Action = "create_spoke_gw"
	resp, err := c.Post(c.baseURL, spoke)
	if err != nil {
		return errors.New("HTTP Post create_spoke_gw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_spoke_gw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_spoke_gw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SpokeJoinTransit(spoke *SpokeVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_spoke_to_transit_gw") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "attach_spoke_to_transit_gw")
	attachSpokeToTransitGw.Add("spoke_gw", spoke.GwName)
	attachSpokeToTransitGw.Add("transit_gw", spoke.TransitGateway)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get attach_spoke_to_transit_gw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode attach_spoke_to_transit_gw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API attach_spoke_to_transit_gw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SpokeLeaveAllTransit(spoke *SpokeVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for detach_spoke_from_transit_gw") + err.Error())
	}
	detachSpokeFromTransitGw := url.Values{}
	detachSpokeFromTransitGw.Add("CID", c.CID)
	detachSpokeFromTransitGw.Add("action", "detach_spoke_from_transit_gw")
	detachSpokeFromTransitGw.Add("spoke_gw", spoke.GwName)
	Url.RawQuery = detachSpokeFromTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get detach_spoke_from_transit_gw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode detach_spoke_from_transit_gw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "has not joined to any transit") {
			log.Errorf("spoke VPC is already left from transit VPC %s", data.Reason)
			return nil
		}
		return errors.New("Rest API detach_spoke_from_transit_gw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SpokeLeaveTransit(spoke *SpokeVpc) error {
	action := "detach_spoke_from_transit_gw"
	data := map[string]interface{}{
		"CID":        c.CID,
		"action":     action,
		"spoke_gw":   spoke.GwName,
		"transit_gw": spoke.TransitGateway,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) EnableHaSpokeVpc(spoke *SpokeVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_spoke_ha") + err.Error())
	}
	enableSpokeHa := url.Values{}
	enableSpokeHa.Add("CID", c.CID)
	enableSpokeHa.Add("action", "enable_spoke_ha")
	enableSpokeHa.Add("gw_name", spoke.GwName)
	enableSpokeHa.Add("eip", spoke.Eip)

	if spoke.CloudType == AWS || spoke.CloudType == AZURE || spoke.CloudType == OCI || spoke.CloudType == AWSGOV {
		enableSpokeHa.Add("public_subnet", spoke.HASubnet)
	} else if spoke.CloudType == GCP {
		enableSpokeHa.Add("new_zone", spoke.HAZone)
	} else {
		return errors.New("invalid cloud type")
	}
	Url.RawQuery = enableSpokeHa.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_spoke_ha failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_spoke_ha failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "HA GW already exists") {
			log.Infof("HA is already enabled %s", data.Reason)
			return nil
		}
		log.Errorf("Enabling HA failed with error %s", data.Reason)
		return errors.New("Rest API enable_spoke_ha Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableHaSpokeGateway(gateway *SpokeVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post create_peering_ha_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_peering_ha_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_peering_ha_gateway Post failed: " + data.Reason)
	}
	return nil
}
