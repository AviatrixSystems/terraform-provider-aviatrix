package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

// Gateway simple struct to hold gateway details
type TransitVpc struct {
	AccountName            string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                 string `form:"action,omitempty"`
	CID                    string `form:"CID,omitempty"`
	CloudType              int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer              string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                 string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                 string `form:"gw_size,omitempty"`
	VpcID                  string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup  string `form:"vnet_and_resource_group_names,omitempty" json:"vpc_id,omitempty"`
	Subnet                 string `form:"public_subnet,omitempty" json:"vpc_net,omitempty"`
	HASubnet               string `form:"ha_subnet,omitempty"`
	PeeringHASubnet        string `json:"public_subnet,omitempty"`
	VpcRegion              string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                string `form:"gw_size,omitempty" json:"gw_size,omitempty"`
	EnableNAT              string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	TagList                string `form:"tags,omitempty"`
	EnableHybridConnection bool   `form:"enable_hybrid_connection" json:"tgw_enabled,omitempty"`
	ConnectedTransit       string `form:"connected_transit" json:"connected_transit,omitempty"`
	InsaneMode             string `form:"insane_mode,omitempty"`
	ReuseEip               string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead     bool   `json:"newly_allocated_eip,omitempty"`
	Eip                    string `form:"eip,omitempty"`
	EnableActiveMesh       string `form:"enable_activemesh,omitempty" json:"enable_activemesh,omitempty"`
}

type TransitGwFireNetInterfaces struct {
	VpcName                  string `json:"vpc_name"`
	VpcRegion                string `json:"vpc_region"`
	TransitVpc               string `json:"transit_vpc"`
	FireNetInterfacesEnabled bool   `json:"dmz_enabled"`
	Name                     string `json:"name"`
}

type TransitGwFireNetInterfacesResp struct {
	Return  bool                       `json:"return"`
	Results TransitGwFireNetInterfaces `json:"results"`
	Reason  string                     `json:"reason"`
}

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_transit_gw"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post create_transit_gw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode create_transit_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API create_transit_gw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableHaTransitVpc(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for enable_transit_ha " + err.Error())
	}
	enableTransitHa := url.Values{}
	enableTransitHa.Add("CID", c.CID)
	enableTransitHa.Add("action", "enable_transit_ha")
	enableTransitHa.Add("gw_name", gateway.GwName)
	enableTransitHa.Add("public_subnet", gateway.HASubnet)
	enableTransitHa.Add("eip", gateway.Eip)

	Url.RawQuery = enableTransitHa.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_transit_ha failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_transit_ha failed: " + err.Error())
	}
	if !data.Return {
		if strings.Contains(data.Reason, "HA GW already exists") {
			log.Printf("[INFO] HA is already enabled %s", data.Reason)
			return nil
		}
		log.Printf("[ERROR] Enabling HA failed with error %s", data.Reason)

		return errors.New("Rest API enable_transit_ha Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachTransitGWForHybrid(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_transit_gateway_interface_to_aws_tgw ") + err.Error())
	}
	enableTransitGatewayInterfaceToAwsTgw := url.Values{}
	enableTransitGatewayInterfaceToAwsTgw.Add("CID", c.CID)
	enableTransitGatewayInterfaceToAwsTgw.Add("action", "enable_transit_gateway_interface_to_aws_tgw")
	enableTransitGatewayInterfaceToAwsTgw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableTransitGatewayInterfaceToAwsTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		if strings.Contains(err.Error(), "already enabled tgw interface") {
			return nil
		}
		return errors.New("Json Decode enable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API enable_transit_gateway_interface_to_aws_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DetachTransitGWForHybrid(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_transit_gateway_interface_to_aws_tgw") + err.Error())
	}
	disableTransitGatewayInterfaceToAwsTgw := url.Values{}
	disableTransitGatewayInterfaceToAwsTgw.Add("CID", c.CID)
	disableTransitGatewayInterfaceToAwsTgw.Add("action", "disable_transit_gateway_interface_to_aws_tgw")
	disableTransitGatewayInterfaceToAwsTgw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableTransitGatewayInterfaceToAwsTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_transit_gateway_interface_to_aws_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableConnectedTransit(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_connected_transit_on_gateway") + err.Error())
	}
	enableConnectedTransitOnGateway := url.Values{}
	enableConnectedTransitOnGateway.Add("CID", c.CID)
	enableConnectedTransitOnGateway.Add("action", "enable_connected_transit_on_gateway")
	enableConnectedTransitOnGateway.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_connected_transit_on_gateway failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_connected_transit_on_gateway failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API enable_connected_transit_on_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableConnectedTransit(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_connected_transit_on_gateway") + err.Error())
	}
	disableConnectedTransitOnGateway := url.Values{}
	disableConnectedTransitOnGateway.Add("CID", c.CID)
	disableConnectedTransitOnGateway.Add("action", "disable_connected_transit_on_gateway")
	disableConnectedTransitOnGateway.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_connected_transit_on_gateway failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_connected_transit_on_gateway failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_connected_transit_on_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableGatewayFireNetInterfaces(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_gateway_firenet_interfaces") + err.Error())
	}
	disableConnectedTransitOnGateway := url.Values{}
	disableConnectedTransitOnGateway.Add("CID", c.CID)
	disableConnectedTransitOnGateway.Add("action", "enable_gateway_firenet_interfaces")
	disableConnectedTransitOnGateway.Add("gateway_name", gateway.GwName)

	Url.RawQuery = disableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Post enable_gateway_firenet_interfaces failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_gateway_firenet_interfaces failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API enable_gateway_firenet_interfaces Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableGatewayFireNetInterfaces(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_gateway_firenet_interfaces") + err.Error())
	}
	disableConnectedTransitOnGateway := url.Values{}
	disableConnectedTransitOnGateway.Add("CID", c.CID)
	disableConnectedTransitOnGateway.Add("action", "disable_gateway_firenet_interfaces")
	disableConnectedTransitOnGateway.Add("gateway", gateway.GwName)

	Url.RawQuery = disableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Post disable_gateway_firenet_interfaces failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_gateway_firenet_interfaces failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_gateway_firenet_interfaces Post failed: " + data.Reason)
	}
	return nil
}
