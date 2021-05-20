package goaviatrix

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AwsTgwPeering struct {
	Action   string `form:"action,omitempty"`
	CID      string `form:"CID,omitempty"`
	TgwName1 string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	TgwName2 string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
}

type AwsTgwPeeringAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	awsTgwPeering.CID = c.CID
	awsTgwPeering.Action = "add_tgw_peering"
	return c.PostAPI(awsTgwPeering.Action, awsTgwPeering, BasicCheck)
}

func (c *Client) GetAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	var data AwsTgwPeeringAPIResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_peered_tgw_names",
		"tgw_name": awsTgwPeering.TgwName1,
	}
	check := func(action, reason string, ret bool) error {
		if !ret {
			if strings.Contains(data.Reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return err
	}
	if len(data.Results) == 0 {
		log.Errorf("Aws tgw peering with tgw: %s and tgw: %s not found", awsTgwPeering.TgwName1, awsTgwPeering.TgwName2)
		return ErrNotFound
	}
	peeringList := data.Results
	for i := range peeringList {
		if peeringList[i] == awsTgwPeering.TgwName2 {
			return nil
		}
	}
	return ErrNotFound
}

func (c *Client) DeleteAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	awsTgwPeering.CID = c.CID
	awsTgwPeering.Action = "delete_tgw_peering"
	return c.PostAPI(awsTgwPeering.Action, awsTgwPeering, BasicCheck)
}
