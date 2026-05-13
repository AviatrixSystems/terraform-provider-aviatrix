package goaviatrix

type FirewallManagementAccess struct {
	CID                          string `form:"CID,omitempty"`
	Action                       string `form:"action,omitempty"`
	TransitFireNetGatewayName    string `form:"transit_firenet_gateway_name,omitempty" json:"gw_name,omitempty"`
	ManagementAccessResourceName string `form:"management_access,omitempty" json:"management_access,omitempty"`
}

type FirewallManagementAccessAPIResp struct {
	Return  bool                       `json:"return"`
	Results []FirewallManagementAccess `json:"results"`
	Reason  string                     `json:"reason"`
}

func (c *Client) CreateFirewallManagementAccess(firewallManagementAccess *FirewallManagementAccess) error {
	form := map[string]string{
		"CID":               c.CID,
		"action":            "edit_transit_firenet_management_access",
		"gateway_name":      firewallManagementAccess.TransitFireNetGatewayName,
		"management_access": firewallManagementAccess.ManagementAccessResourceName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetFirewallManagementAccess(firewallManagementAccess *FirewallManagementAccess) (*FirewallManagementAccess, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_transit_firenet_spoke_policies",
	}

	var data FirewallManagementAccessAPIResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		return nil, ErrNotFound
	}
	firewallManagementAccessList := data.Results
	for i := range firewallManagementAccessList {
		if firewallManagementAccessList[i].TransitFireNetGatewayName != firewallManagementAccess.TransitFireNetGatewayName {
			continue
		}
		if firewallManagementAccessList[i].ManagementAccessResourceName == "no" {
			return nil, ErrNotFound
		}
		firewallManagementAccess.ManagementAccessResourceName = firewallManagementAccessList[i].ManagementAccessResourceName
		return firewallManagementAccess, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DestroyFirewallManagementAccess(firewallManagementAccess *FirewallManagementAccess) error {
	form := map[string]string{
		"CID":               c.CID,
		"action":            "edit_transit_firenet_management_access",
		"gateway_name":      firewallManagementAccess.TransitFireNetGatewayName,
		"management_access": firewallManagementAccess.ManagementAccessResourceName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
