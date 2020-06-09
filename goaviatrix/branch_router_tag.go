package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

type BranchRouterTag struct {
	Name           string `form:"tag_name,omitempty"`
	Config         string `form:"custom_cfg,omitempty"`
	Branches       []string
	BranchesString string `form:"include_branch_list,omitempty"`
	Commit         bool
	CID            string `form:"CID"`
	Action         string `form:"action"`
}

func (c *Client) CreateBranchRouterTag(brt *BranchRouterTag) error {
	// Create the tag
	brt.CID = c.CID
	brt.Action = "add_cloudwan_configtag"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post add_cloudwan_configtag failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body add_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode add_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API add_cloudwan_configtag Post failed: " + data.Reason)
	}

	// Set the tag config
	if err := c.UpdateBranchRouterTagConfig(brt); err != nil {
		return err
	}

	// Attach the branches to the tag
	if err := c.UpdateBranchRouterTagBranches(brt); err != nil {
		return err
	}

	// Commit the tag config to the branches if 'commit' == true
	if !brt.Commit {
		return nil
	}
	if err := c.CommitBranchRouterTag(brt); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetBranchRouterTag(brt *BranchRouterTag) (*BranchRouterTag, error) {
	// Check if a tag exists with the given name
	brt.CID = c.CID
	brt.Action = "list_cloudwan_configtag_names"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return nil, errors.New("HTTP Post get_cloudwan_configtag_details failed: " + err.Error())
	}

	type Resp struct {
		Return  bool     `json:"return,omitempty"`
		Results []string `json:"results,omitempty"`
		Reason  string   `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body list_cloudwan_configtag_names failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_cloudwan_configtag_names failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_cloudwan_configtag_names Post failed: " + data.Reason)
	}

	if !Contains(data.Results, brt.Name) {
		return nil, ErrNotFound
	}

	// Get the details for the tag
	brt.Action = "get_cloudwan_configtag_details"
	resp, err = c.Post(c.baseURL, brt)

	if err != nil {
		return nil, errors.New("HTTP Post get_cloudwan_configtag_details failed: " + err.Error())
	}

	type DetailsResults struct {
		TagName          string   `json:"gtag_name"`
		AttachedBranches []string `json:"rgw_name"`
		Config           string   `json:"custom_cfg"`
	}
	type DetailsResp struct {
		Return  bool           `json:"return,omitempty"`
		Results DetailsResults `json:"results,omitempty"`
		Reason  string         `json:"reason,omitempty"`
	}
	var detailsData DetailsResp
	b = bytes.Buffer{}
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_cloudwan_configtag_details failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&detailsData); err != nil {
		return nil, errors.New("Json Decode get_cloudwan_configtag_details failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !detailsData.Return {
		return nil, errors.New("Rest API get_cloudwan_configtag_details Post failed: " + detailsData.Reason)
	}

	brt.Branches = detailsData.Results.AttachedBranches
	brt.Config = detailsData.Results.Config
	return brt, nil
}

func (c *Client) UpdateBranchRouterTagConfig(brt *BranchRouterTag) error {
	brt.CID = c.CID
	brt.Action = "edit_cloudwan_configtag"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post edit_cloudwan_configtag failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body edit_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode edit_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API edit_cloudwan_configtag Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) UpdateBranchRouterTagBranches(brt *BranchRouterTag) error {
	brt.CID = c.CID
	brt.Action = "attach_branches_to_cloudwan_configtag"
	brt.BranchesString = strings.Join(brt.Branches, ", ")
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post attach_branches_to_cloudwan_configtag failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}

	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body attach_branches_to_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode attach_branches_to_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API attach_branches_to_cloudwan_configtag Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) CommitBranchRouterTag(brt *BranchRouterTag) error {
	brt.CID = c.CID
	brt.Action = "commit_cloudwan_configtag_to_branches"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post commit_cloudwan_configtag_to_branches failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}

	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body commit_cloudwan_configtag_to_branches failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode commit_cloudwan_configtag_to_branches failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API commit_cloudwan_configtag_to_branches Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DeleteBranchRouterTag(brt *BranchRouterTag) error {
	brt.CID = c.CID
	brt.Action = "delete_cloudwan_configtag"
	resp, err := c.Post(c.baseURL, brt)

	if err != nil {
		return errors.New("HTTP Post delete_cloudwan_configtag failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}

	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body delete_cloudwan_configtag failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode delete_cloudwan_configtag failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API delete_cloudwan_configtag Post failed: " + data.Reason)
	}

	return nil
}
