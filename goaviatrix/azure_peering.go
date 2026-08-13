package goaviatrix

type AzurePeer struct {
	Action       string `form:"action,omitempty"`
	CID          string `form:"CID,omitempty"`
	AccountName1 string `form:"req_account_name,omitempty"`
	AccountName2 string `form:"acc_account_name,omitempty"`
	VNet1        string `form:"req_vpc_id,omitempty"`
	VNet2        string `form:"acc_vpc_id,omitempty"`
	Region1      string `form:"req_region,omitempty"`
	Region2      string `form:"acc_region,omitempty"`
	VNetCidr1    []string
	VNetCidr2    []string
}

type AzurePeerAPIResp struct {
	Return  bool   `json:"return"`
	Reason  string `json:"reason"`
	Results string `json:"results"`
}

func (c *Client) CreateAzurePeer(azurePeer *AzurePeer) error {
	azurePeer.CID = c.CID
	azurePeer.Action = "arm_peer_vnet_pair"
	return c.PostAPI(azurePeer.Action, azurePeer, BasicCheck)
}

func (c *Client) GetAzurePeer(azurePeer *AzurePeer) (*AzurePeer, error) {
	var data map[string]interface{}
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_arm_peer_vnet_pairs",
	}

	if err := c.GetAPI(&data, form["action"], form, BasicCheck); err != nil {
		return nil, err
	}

	val, ok := data["results"]
	if !ok || val == nil {
		return nil, ErrNotFound
	}

	pairList, ok := val.([]interface{})
	if !ok {
		return nil, ErrNotFound
	}

	for _, item := range pairList {
		pair, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		requester, ok := pair["requester"].(map[string]interface{})
		if !ok {
			continue
		}
		accepter, ok := pair["accepter"].(map[string]interface{})
		if !ok {
			continue
		}

		reqVPC, ok := requester["vpc_id"].(string)
		if !ok {
			continue
		}
		accVPC, ok := accepter["vpc_id"].(string)
		if !ok {
			continue
		}

		if reqVPC != azurePeer.VNet1 || accVPC != azurePeer.VNet2 {
			continue
		}

		out := &AzurePeer{
			VNet1:        reqVPC,
			VNet2:        accVPC,
			AccountName1: stringOrEmpty(requester["account_name"]),
			AccountName2: stringOrEmpty(accepter["account_name"]),
			Region1:      stringOrEmpty(requester["region"]),
			Region2:      stringOrEmpty(accepter["region"]),
			VNetCidr1:    toStringSlice(requester["vpc_cidr"]),
			VNetCidr2:    toStringSlice(accepter["vpc_cidr"]),
		}

		return out, nil
	}

	return nil, ErrNotFound
}

func stringOrEmpty(v interface{}) string {
	s, _ := v.(string)
	return s
}

func toStringSlice(v interface{}) []string {
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, x := range list {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func (c *Client) DeleteAzurePeer(azurePeer *AzurePeer) error {
	form := map[string]string{
		"CID":       c.CID,
		"action":    "arm_unpeer_vnet_pair",
		"vpc_name1": azurePeer.VNet1,
		"vpc_name2": azurePeer.VNet2,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}
