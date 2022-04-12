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
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	if val, ok := data["results"]; ok {
		pairList := val.([]interface{})
		for i := range pairList {
			if pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string) == azurePeer.VNet1 &&
				pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string) == azurePeer.VNet2 {
				azurePeer := &AzurePeer{
					VNet1:        pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string),
					VNet2:        pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string),
					AccountName1: pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["account_name"].(string),
					AccountName2: pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["account_name"].(string),
					Region1:      pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["region"].(string),
					Region2:      pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["region"].(string),
				}

				vnetCidrList1 := pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_cidr"].([]interface{})
				var vnetCidr1 []string
				for i := range vnetCidrList1 {
					vnetCidr1 = append(vnetCidr1, vnetCidrList1[i].(interface{}).(string))
				}
				azurePeer.VNetCidr1 = vnetCidr1

				vnetCidrList2 := pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_cidr"].([]interface{})
				var vnetCidr2 []string
				for i := range vnetCidrList2 {
					vnetCidr2 = append(vnetCidr2, vnetCidrList2[i].(string))
				}
				azurePeer.VNetCidr2 = vnetCidr2

				return azurePeer, nil
			}
		}
	}
	return nil, ErrNotFound
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
