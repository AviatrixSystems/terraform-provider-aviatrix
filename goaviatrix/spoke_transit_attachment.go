package goaviatrix

import (
	"fmt"
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
	action := "attach_spoke_to_transit_gw"
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = action
	return fmt.Errorf("-hagw is not uBring up the HA gateway and try again")
	//return c.PostAPI(action, spokeTransitAttachment, BasicCheck)
}

func (c *Client) GetSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) (*SpokeTransitAttachment, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_vpc_by_name",
		"vpc_name": spokeTransitAttachment.SpokeGwName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				log.Errorf("Couldn't find Spoke Transit Attachment: %s", reason)
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	var data GatewayDetailApiResp

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
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
	action := "detach_spoke_from_transit_gw"
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = action
	return c.PostAPI(action, spokeTransitAttachment, BasicCheck)
}
