package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SpokeTransitAttachment struct {
	Action        string `form:"action,omitempty"`
	CID           string `form:"CID,omitempty"`
	SpokeGwName   string `form:"spoke_gw,omitempty"`
	TransitGwName string `form:"transit_gw,omitempty"`
	RouteTables   string `form:"route_table_list,omitempty"`
}

func (c *Client) CreateSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) error {
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = "attach_spoke_to_transit_gw"
	resp, err := c.Post(c.baseURL, spokeTransitAttachment)
	if err != nil {
		return errors.New("HTTP Post attach_spoke_to_transit_gw failed: " + err.Error())
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
		return errors.New("Rest API attach_spoke_to_transit_gw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) (*SpokeTransitAttachment, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_vpc_by_name") + err.Error())
	}
	listVpcByName := url.Values{}
	listVpcByName.Add("CID", c.CID)
	listVpcByName.Add("action", "list_vpc_by_name")
	listVpcByName.Add("vpc_name", spokeTransitAttachment.SpokeGwName)
	Url.RawQuery = listVpcByName.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_vpc_by_name failed: " + err.Error())
	}
	var data GatewayDetailApiResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpc_by_name failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpc_by_name Get failed: " + data.Reason)
	}
	if data.Results.GwName == spokeTransitAttachment.SpokeGwName {
		if data.Results.TransitGwName == spokeTransitAttachment.TransitGwName || data.Results.EgressTransitGwName == spokeTransitAttachment.TransitGwName {
			spokeTransitAttachment.RouteTables = strings.Join(data.Results.RouteTables, ",")
			return spokeTransitAttachment, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix gateway %s", spokeTransitAttachment.SpokeGwName)
	return nil, ErrNotFound
}

func (c *Client) DeleteSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) error {
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = "detach_spoke_from_transit_gw"
	resp, err := c.Post(c.baseURL, spokeTransitAttachment)
	if err != nil {
		return errors.New("HTTP Post detach_spoke_from_transit_gw failed: " + err.Error())
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
		return errors.New("Rest API detach_spoke_from_transit_gw Post failed: " + data.Reason)
	}
	return nil
}
