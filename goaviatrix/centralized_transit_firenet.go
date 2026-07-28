package goaviatrix

import "context"

type CentralizedTransitFirenet struct {
	CID             string `form:"CID,omitempty"`
	Action          string `form:"action,omitempty"`
	PrimaryGwName   string `form:"primary_gw_name,omitempty"`
	SecondaryGwName string `form:"secondary_gw_name,omitempty"`
}

func (c *Client) GetPrimaryFireNet(ctx context.Context) ([]string, error) {
	form := map[string]string{
		"action": "list_primary_firenet",
		"CID":    c.CID,
	}

	type PrimaryFirenetList struct {
		GwName string `json:"gw_name"`
	}

	type Resp struct {
		Return  bool                 `json:"return"`
		Results []PrimaryFirenetList `json:"results"`
	}
	var resp Resp
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	var primaryFirenetList []string

	for _, primaryFirenet := range resp.Results {
		primaryFirenetList = append(primaryFirenetList, primaryFirenet.GwName)
	}

	return primaryFirenetList, nil
}

func (c *Client) GetSecondaryFireNet(ctx context.Context) ([]string, error) {
	form := map[string]string{
		"action": "list_secondary_firenet",
		"CID":    c.CID,
	}

	type SecondaryFirenetList struct {
		GwName string `json:"gw_name"`
	}

	type Resp struct {
		Return  bool                   `json:"return"`
		Results []SecondaryFirenetList `json:"results"`
	}

	var resp Resp
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	var secondaryFirenetList []string

	for _, secondaryFirenet := range resp.Results {
		secondaryFirenetList = append(secondaryFirenetList, secondaryFirenet.GwName)
	}

	return secondaryFirenetList, nil
}

func (c *Client) CreateCentralizedTransitFireNet(ctx context.Context, firenetAttachment *CentralizedTransitFirenet) error {
	firenetAttachment.Action = "attach_centralized_firenet"
	firenetAttachment.CID = c.CID

	return c.PostAPIContext(ctx, firenetAttachment.Action, firenetAttachment, BasicCheck)
}

func (c *Client) GetCentralizedTransitFireNet(ctx context.Context, centralizedTransitFirenet *CentralizedTransitFirenet) error {
	form := map[string]string{
		"action":      "list_transit_firenet",
		"CID":         c.CID,
		"centralized": "true",
	}

	type PF struct {
		GwName string `json:"gw_name"`
	}

	type SF struct {
		GwName string `json:"gw_name"`
	}

	type CentralizedTransitFirenetList struct {
		PrimaryFirenet       PF   `json:"primary_firenet"`
		SecondaryFirenetList []SF `json:"secondary_firenet_list"`
	}

	type Resp struct {
		Return  bool                            `json:"return"`
		Results []CentralizedTransitFirenetList `json:"results"`
	}

	var resp Resp
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return err
	}

	for _, cf := range resp.Results {
		if centralizedTransitFirenet.PrimaryGwName == cf.PrimaryFirenet.GwName {
			var SFList []string
			for _, sf := range cf.SecondaryFirenetList {
				SFList = append(SFList, sf.GwName)
			}

			if !Contains(SFList, centralizedTransitFirenet.SecondaryGwName) {
				return ErrNotFound
			}
		}
	}

	return nil
}

func (c *Client) DeleteCentralizedTransitFireNet(ctx context.Context, firenetAttachment *CentralizedTransitFirenet) error {
	firenetAttachment.Action = "detach_centralized_firenet"
	firenetAttachment.CID = c.CID

	return c.PostAPIContext(ctx, firenetAttachment.Action, firenetAttachment, BasicCheck)
}
