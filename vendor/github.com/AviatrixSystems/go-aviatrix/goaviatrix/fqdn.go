package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Gateway simple struct to hold fqdn details
type FQDN struct {
	FQDNTag                 string `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	Action                  string `form:"action,omitempty"`
	CID                     string `form:"CID,omitempty"`
	FQDNStatus              string `form:"status,omitempty" json:"status,omitempty"`
	FQDNMode                string `form:"wb_mode,omitempty" json:"wb_mode,omitempty"`
	GwList                  []string `form:"gw_name,omitempty" json:"members,omitempty"`
	DomainList              []string `form:"domain_names[],omitempty"`
}

type ResultListResp struct {
	Return  bool   `json:"return"`
	Results []string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateFQDN(fqdn *FQDN) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=add_nfq_tag&tag_name=%s", c.CID, fqdn.FQDNTag)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DeleteFQDN(fqdn *FQDN) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=del_nfq_tag&tag_name=%s", c.CID, fqdn.FQDNTag)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

//change state to 'enabled' or 'disabled'
func (c *Client) UpdateFQDNStatus(fqdn *FQDN) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=set_nfq_tag_state&tag_name=%s&status=%s", c.CID, fqdn.FQDNTag, fqdn.FQDNStatus)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

//Change default mode to 'white' or 'black'
func (c *Client) UpdateFQDNMode(fqdn *FQDN) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=set_nfq_tag_wbmode&tag_name=%s&wbmode=%s", c.CID, fqdn.FQDNTag, fqdn.FQDNMode)
	resp,err := c.Get(path, nil)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) UpdateDomains(fqdn *FQDN) (error) {
	fqdn.CID=c.CID
	fqdn.Action="update_nfq_tag_whitelist"
	resp,err := c.Post(c.baseURL, fqdn)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) AttachGws(fqdn *FQDN) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=add_nfq_tag_member&tag_name=%s", c.CID, fqdn.FQDNTag)
	for i := range fqdn.GwList {
		newPath := path + fmt.Sprintf("&gw_name=%s", fqdn.GwList[i])
		resp,err := c.Get(newPath, nil)
			if err != nil {
			return err
		}
		var data APIResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if(!data.Return){
			return errors.New(data.Reason)
		}
	}
	return nil
}

func (c *Client) DetachGws(fqdn *FQDN) (error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=del_nfq_tag_member&tag_name=%s", c.CID, fqdn.FQDNTag)
	for i := range fqdn.GwList {
		newPath := path + fmt.Sprintf("&gw_name=%s", fqdn.GwList[i])
		resp,err := c.Get(newPath, nil)
			if err != nil {
			return err
		}
		var data APIResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if(!data.Return){
			return errors.New(data.Reason)
		}
	}
	return nil
}

func (c *Client) GetFQDNTag(fqdn *FQDN) (*FQDN, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_nfq_tag_names&account_name=%s", c.CID)
	resp,err := c.Get(path, nil)

	if err != nil {
		return nil, err
	}

	//Output result for this query is weird. FQDN tag names have 
	//been set as keys. This cannot be unmarshalled easily as we
	//can't have a predefined structure(since tag names will be arbitrary)
	//to decode it in. So using a map of string->interface{}
	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if val, ok := data["reason"]; ok {
	    return nil, errors.New(val.(string))
	}
	if val, ok := data["results"]; ok {
		if foundTag, ok1 := val.(map[string]interface{})[fqdn.FQDNTag]; ok1 {
			tagdata := foundTag.(map[string]interface{})
			fqdn.FQDNMode = tagdata["wbmode"].(string)
			fqdn.FQDNStatus = tagdata["state"].(string)
			return fqdn, nil
		}
	}

	return nil, fmt.Errorf("FQDN Tag %s not found", fqdn.FQDNTag)
}

func (c *Client) ListDomains(fqdn *FQDN) (*FQDN, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_nfq_tag_whitelist&tag_name=%s", c.CID, fqdn.FQDNTag)
	resp,err := c.Get(path, nil)
		if err != nil {
		return nil, err
	}
	var data ResultListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	//domainList:= data.Results
	fqdn.DomainList = data.Results

	return fqdn, nil
}

func (c *Client) ListGws(fqdn *FQDN) (*FQDN, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_gateway_nfq_tag&tag_name=%s", c.CID, fqdn.FQDNTag)
	resp,err := c.Get(path, nil)
		if err != nil {
		return nil, err
	}
	var data ResultListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	fqdn.GwList = data.Results

	return fqdn, nil
}
