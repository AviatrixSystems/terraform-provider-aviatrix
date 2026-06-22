package goaviatrix

import "context"

type AllowIp struct {
	IpAddress   string `json:"addr"`
	Description string `json:"desc"`
}

type AllowList struct {
	AllowList []AllowIp `json:"allow_list"`
	Enforce   bool      `json:"enforce"`
	Enable    bool      `json:"enable"`
}

func (c *Client) CreateControllerAccessAllowList(ctx context.Context, allowList *AllowList) error {
	endpoint := "controller/allow-list"
	allowList.Enable = true

	err := c.PutAPIContext25(ctx, endpoint, allowList)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetControllerAccessAllowList(ctx context.Context) (*AllowList, error) {
	endpoint := "controller/allow-list"

	var data AllowList
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	} else if !data.Enable {
		return nil, ErrNotFound
	}

	return &data, nil
}

func (c *Client) UpdateControllerAccessAllowList(ctx context.Context, allowList *AllowList) error {
	endpoint := "controller/allow-list"
	allowList.Enable = true

	err := c.PutAPIContext25(ctx, endpoint, allowList)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteControllerAccessAllowList(ctx context.Context) error {
	endpoint := "controller/allow-list"

	var allowList AllowList
	allowList.AllowList = []AllowIp{}
	allowList.Enable = false

	err := c.PutAPIContext25(ctx, endpoint, allowList)
	if err != nil {
		return err
	}

	return nil
}
