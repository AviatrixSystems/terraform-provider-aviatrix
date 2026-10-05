package goaviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	log "github.com/sirupsen/logrus"
)

type AwsTgwPeering struct {
	Action   string `form:"action,omitempty"`
	CID      string `form:"CID,omitempty"`
	TgwName1 string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	TgwName2 string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
	Async    bool   `form:"async,omitempty"`
}

type AwsTgwPeeringAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	awsTgwPeering.CID = c.CID
	awsTgwPeering.Action = "add_tgw_peering"
	awsTgwPeering.Async = true
	return c.PostAsyncAPI(awsTgwPeering.Action, awsTgwPeering, BasicCheck)
}

func (c *Client) GetAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	var data AwsTgwPeeringAPIResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_peered_tgw_names",
		"tgw_name": awsTgwPeering.TgwName1,
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
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
	awsTgwPeering.Async = true
	return c.PostAsyncAPI(awsTgwPeering.Action, awsTgwPeering, BasicCheck)
}

func DiffSuppressFuncAwsTgwPeeringTgwName1(k, old, new string, d *schema.ResourceData) bool {
	tgwName2Old, _ := d.GetChange("tgw_name2")
	return old == d.Get("tgw_name2").(string) && new == tgwName2Old.(string)
}

func DiffSuppressFuncAwsTgwPeeringTgwName2(k, old, new string, d *schema.ResourceData) bool {
	tgwName1Old, _ := d.GetChange("tgw_name1")
	return old == d.Get("tgw_name1").(string) && new == tgwName1Old.(string)
}
