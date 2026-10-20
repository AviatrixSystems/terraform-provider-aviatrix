package goaviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	log "github.com/sirupsen/logrus"
)

type AwsTgwPeering struct {
	Action        string `form:"action,omitempty"`
	CID           string `form:"CID,omitempty"`
	TgwName1      string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	TgwName2      string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
	PeerTgwID     string `form:"peer_tgw_id,omitempty" json:"peer_tgw_id,omitempty"`
	PeerAccountID string `form:"peer_account_id,omitempty" json:"peer_account_id,omitempty"`
	PeerRegion    string `form:"peer_region,omitempty" json:"peer_region,omitempty"`
	Async         bool   `form:"async,omitempty"`
}

type AwsTgwPeeringAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type awsTgwPeeringAttachment struct {
	Type          string `json:"type"`
	PeerTgwName   string `json:"peer_tgw_name"`
	ExternalPeer  bool   `json:"external_peer"`
	PeerAccountID string `json:"peer_account_id"`
	PeerRegion    string `json:"peer_region"`
}

type awsTgwPeeringDetailsResp struct {
	Return  bool `json:"return"`
	Results struct {
		Attachments map[string]awsTgwPeeringAttachment `json:"attachments"`
	} `json:"results"`
	Reason string `json:"reason"`
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

// GetAwsTgwExternalPeering fills the peer TGW ID, account ID and region of a
// peering whose tgw_name2 is a label for a TGW that TGW-O does not manage. It
// returns false, leaving the fields alone, when tgw_name2 is a managed TGW.
// A PeerTgwID already set picks between peerings that share a label.
func (c *Client) GetAwsTgwExternalPeering(awsTgwPeering *AwsTgwPeering) (bool, error) {
	var data awsTgwPeeringDetailsResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_tgw_details",
		"tgw_name": awsTgwPeering.TgwName1,
	}
	if err := c.GetAPI(&data, form["action"], form, BasicCheck); err != nil {
		return false, err
	}

	var matches []string
	for peerTgwID, attm := range data.Results.Attachments {
		if attm.Type != "peering" || !attm.ExternalPeer || attm.PeerTgwName != awsTgwPeering.TgwName2 {
			continue
		}
		if awsTgwPeering.PeerTgwID != "" && peerTgwID != awsTgwPeering.PeerTgwID {
			continue
		}
		matches = append(matches, peerTgwID)
	}
	if len(matches) == 0 {
		return false, nil
	}
	if len(matches) > 1 {
		return false, fmt.Errorf("TGW %s has more than one external peering labeled %s; set peer_tgw_id to choose one", awsTgwPeering.TgwName1, awsTgwPeering.TgwName2)
	}

	attm := data.Results.Attachments[matches[0]]
	awsTgwPeering.PeerTgwID = matches[0]
	awsTgwPeering.PeerAccountID = attm.PeerAccountID
	awsTgwPeering.PeerRegion = attm.PeerRegion
	return true, nil
}

func (c *Client) DeleteAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	form := &AwsTgwPeering{
		CID:       c.CID,
		Action:    "delete_tgw_peering",
		TgwName1:  awsTgwPeering.TgwName1,
		TgwName2:  awsTgwPeering.TgwName2,
		PeerTgwID: awsTgwPeering.PeerTgwID,
		Async:     true,
	}
	return c.PostAsyncAPI(form.Action, form, BasicCheck)
}

func DiffSuppressFuncAwsTgwPeeringTgwName1(k, old, new string, d *schema.ResourceData) bool {
	tgwName2Old, _ := d.GetChange("tgw_name2")

	tgw2Cur, ok := d.Get("tgw_name2").(string)
	if !ok {
		return false
	}
	tgw2Old, ok := tgwName2Old.(string)
	if !ok {
		return false
	}

	return old == tgw2Cur && new == tgw2Old
}

func DiffSuppressFuncAwsTgwPeeringTgwName2(k, old, new string, d *schema.ResourceData) bool {
	tgwName1Old, _ := d.GetChange("tgw_name1")

	tgw1Cur, ok := d.Get("tgw_name1").(string)
	if !ok {
		return false
	}
	tgw1Old, ok := tgwName1Old.(string)
	if !ok {
		return false
	}

	return old == tgw1Cur && new == tgw1Old
}
