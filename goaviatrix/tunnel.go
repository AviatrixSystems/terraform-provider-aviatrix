package goaviatrix

// Tunnel simple struct to hold tunnel details

import (
	log "github.com/sirupsen/logrus"
)

type Tunnel struct {
	VpcName1        string `json:"vpc_name1"`
	VpcName2        string `json:"vpc_name2"`
	PeeringState    string `json:"peering_state"`
	PeeringHaStatus string `json:"peering_ha_status"`
	PeeringLink     string `json:"peering_link"`
	EnableHA        string `json:"enable_ha"`
}

type TunnelResult struct {
	PairList []Tunnel `json:"pair_list"`
}

type TunnelListResp struct {
	Return  bool         `json:"return"`
	Results TunnelResult `json:"results"`
	Reason  string       `json:"reason"`
}

func (c *Client) CreateTunnel(tunnel *Tunnel) error {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "peer_vpc_pair",
		"vpc_name1":  tunnel.VpcName1,
		"vpc_name2":  tunnel.VpcName2,
		"ha_enabled": tunnel.EnableHA,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetTunnel(tunnel *Tunnel) (*Tunnel, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_peer_vpc_pairs",
	}

	var data TunnelListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	// list_peer_vpc_pairs keys spoke-to-spoke peerings by gateway GROUP name,
	// while the config supplies the individual gateway name. Resolve each name
	// to its group name so the match works for both, mirroring GetSpokeTransitAttachment.
	grpName1 := c.resolveGroupName(tunnel.VpcName1)
	grpName2 := c.resolveGroupName(tunnel.VpcName2)

	tunList := data.Results.PairList
	for i := range tunList {
		if matchesTunnelPair(tunList[i].VpcName1, tunList[i].VpcName2, tunnel.VpcName1, tunnel.VpcName2, grpName1, grpName2) {
			log.Debugf("Found %s~%s tunnel: %#v", tunnel.VpcName1, tunnel.VpcName2, tunList[i])
			return &tunList[i], nil
		}
	}
	log.Errorf("Tunnel with gateways %s and %s not found", tunnel.VpcName1, tunnel.VpcName2)
	return nil, ErrNotFound
}

// matchesTunnelPair reports whether a listed peering pair (pair1, pair2)
// corresponds to the configured gateways (cfg1, cfg2). Each configured gateway
// may appear in the list either by its own name or by its resolved group name
// (grp1/grp2, for spoke-to-spoke peerings), and the listed pair is stored with
// gateway names sorted, so its orientation need not match the config — hence
// both orderings are accepted.
func matchesTunnelPair(pair1, pair2, cfg1, cfg2, grp1, grp2 string) bool {
	side := func(listed, cfg, grp string) bool {
		return listed == cfg || listed == grp
	}
	return (side(pair1, cfg1, grp1) && side(pair2, cfg2, grp2)) ||
		(side(pair2, cfg1, grp1) && side(pair1, cfg2, grp2))
}

// resolveGroupName maps an individual gateway name to its gateway-group name.
// Falls back to the input name when the gateway can't be looked up (e.g. it is
// already a group name, or lookup fails).
func (c *Client) resolveGroupName(gwName string) string {
	gw, err := c.GetGateway(&Gateway{GwName: gwName})
	if err == nil && gw.GroupName != "" {
		return gw.GroupName
	}
	return gwName
}

func (c *Client) UpdateTunnel(tunnel *Tunnel) error {
	return nil
}

func (c *Client) DeleteTunnel(tunnel *Tunnel) error {
	form := map[string]string{
		"CID":       c.CID,
		"action":    "unpeer_vpc_pair",
		"vpc_name1": tunnel.VpcName1,
		"vpc_name2": tunnel.VpcName2,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
