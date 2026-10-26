package goaviatrix

import "context"

func (c *Client) EnableCopilotAssociation(ctx context.Context, addr string) error {
	form := map[string]string{
		"action":     "enable_copilot_association",
		"CID":        c.CID,
		"copilot_ip": addr,
	}
	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) DisableCopilotAssociation(ctx context.Context) error {
	form := map[string]string{
		"action": "disable_copilot_association",
		"CID":    c.CID,
	}
	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

type CopilotAssociationStatus struct {
	Status   bool   `json:"status"`
	IP       string `json:"ip"`
	PublicIp string `json:"public_ip"`
}

func (c *Client) GetCopilotAssociationStatus(ctx context.Context) (*CopilotAssociationStatus, error) {
	form := map[string]string{
		"action": "get_copilot_association_status",
		"CID":    c.CID,
	}
	var resp struct {
		APIResp
		Results CopilotAssociationStatus
	}
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	if !resp.Results.Status {
		return nil, ErrNotFound
	}
	return &resp.Results, nil
}
