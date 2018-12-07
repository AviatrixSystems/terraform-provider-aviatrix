package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AwsTGW simple struct to hold aws_tgw details
type AWSTgw struct {
	Action           string `form:"action,omitempty"`
	CID              string `form:"CID,omitempty"`
	Name 			 string `form:"tgw_name,omitempty"`
	AccountName      string `form:"account_name,omitempty"`
	Region           string `form:"region,omitempty"`
	AwsSideAsNumber  string `form:"aws_side_asn,omitempty"`
}

type AWSTgwAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type AWSTgwList struct {
	Return  bool   	 `json:"return"`
	Results []AWSTgw `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateAWSTgw(awsTgw *AWSTgw) (error) {
	awsTgw.CID = c.CID
	awsTgw.Action = "add_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) GetAWSTgw(awsTgw *AWSTgw) (string, error) {
	awsTgw.CID = c.CID
	awsTgw.Action = "list_tgw_names"
	path := c.baseURL + fmt.Sprintf("?action=%s&CID=%s&account_name=%s&region=%s", awsTgw.Action, awsTgw.CID,
		awsTgw.AccountName, awsTgw.Region)

	resp, err := c.Get(path, nil)
	if err != nil {
		return "", err
	}

	data := AWSTgwAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if !data.Return {
		return "", errors.New(data.Reason)
	}

	awsTgwList := data.Results
	for i := range awsTgwList {
		if awsTgwList[i] == awsTgw.Name {
			return awsTgwList[i], nil
		}
	}

	return "", ErrNotFound
}

func (c *Client) UpdateAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DeleteAWSTgw(awsTgw *AWSTgw) (error) {
	awsTgw.CID = c.CID
	awsTgw.Action = "delete_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if !data.Return {
 		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) AttachVpcToAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DetachVpcFromAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) AttachAviatrixTgwToAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DetachAviatrixTgwFromAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}