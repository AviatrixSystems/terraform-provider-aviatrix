package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type ArmSpokeNativePeering struct {
	CID                string `form:"CID,omitempty"`
	Action             string `form:"action,omitempty"`
	TransitGatewayName string `form:"transit_gateway_name,omitempty"`
	SpokeAccountName   string `form:"account_name,omitempty"`
	SpokeRegion        string `form:"region,omitempty"`
	SpokeVpcID         string `form:"vpc_id,omitempty"`
}

type ArmSpokeNativePeeringAPIResp struct {
	Return  bool                        `json:"return"`
	Results []ArmSpokeNativePeeringEdit `json:"results"`
	Reason  string                      `json:"reason"`
}

type ArmSpokeNativePeeringEdit struct {
	Region string `json:"region"`
	Name   string `json:"name"`
}

func (c *Client) CreateArmSpokeNativePeering(armSpokeNativePeering *ArmSpokeNativePeering) error {
	armSpokeNativePeering.CID = c.CID
	armSpokeNativePeering.Action = "attach_arm_native_spoke_to_transit"
	resp, err := c.Post(c.baseURL, armSpokeNativePeering)
	if err != nil {
		return errors.New("HTTP Post 'attach_arm_native_spoke_to_transit' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'attach_arm_native_spoke_to_transit' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'attach_arm_native_spoke_to_transit' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetArmSpokeNativePeering(armSpokeNativePeering *ArmSpokeNativePeering) (*ArmSpokeNativePeering, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'GetArmSpokeNativePeering': ") + err.Error())
	}
	listArmNativeSpokes := url.Values{}
	listArmNativeSpokes.Add("CID", c.CID)
	listArmNativeSpokes.Add("action", "list_arm_native_spokes")
	listArmNativeSpokes.Add("transit_gateway_name", armSpokeNativePeering.TransitGatewayName)
	listArmNativeSpokes.Add("details", "true")
	Url.RawQuery = listArmNativeSpokes.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_arm_native_spokes' failed: " + err.Error())
	}
	var data ArmSpokeNativePeeringAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_arm_native_spokes' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'list_arm_native_spokes' Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		return nil, ErrNotFound
	}
	peeringList := data.Results
	for i := range peeringList {
		if peeringList[i].Name == "" || len(strings.Split(peeringList[i].Name, ":")) != 3 {
			continue
		}
		spokeArray := strings.Split(peeringList[i].Name, ":")
		spokeAccountName := spokeArray[0]
		spokeVpcID := "" + spokeArray[1] + ":" + spokeArray[2]
		if armSpokeNativePeering.SpokeAccountName != spokeAccountName || armSpokeNativePeering.SpokeVpcID != spokeVpcID {
			continue
		}
		armSpokeNativePeering.SpokeRegion = peeringList[i].Region
		return armSpokeNativePeering, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DeleteArmSpokeNativePeering(armSpokeNativePeering *ArmSpokeNativePeering) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'DeleteArmSpokeNativePeering': ") + err.Error())
	}
	detachArmNativeSpokeToTransit := url.Values{}
	detachArmNativeSpokeToTransit.Add("CID", c.CID)
	detachArmNativeSpokeToTransit.Add("action", "detach_arm_native_spoke_to_transit")
	detachArmNativeSpokeToTransit.Add("transit_gateway_name", armSpokeNativePeering.TransitGatewayName)
	detachArmNativeSpokeToTransit.Add("spoke_name", ""+armSpokeNativePeering.SpokeAccountName+":"+armSpokeNativePeering.SpokeVpcID)
	Url.RawQuery = detachArmNativeSpokeToTransit.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'detach_arm_native_spoke_to_transit' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'detach_arm_native_spoke_to_transit' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'detach_arm_native_spoke_to_transit' Post failed: " + data.Reason)
	}
	return nil
}
