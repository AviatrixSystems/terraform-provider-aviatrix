package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Filters struct {
	FQDN     string `form:"fqdn,omitempty" json:"fqdn,omitempty"`
	Protocol string `form:"proto,omitempty" json:"proto,omitempty"`
	Port     string `form:"port,omitempty" json:"port,omitempty"`
}

// Gateway simple struct to hold fqdn details
type FQDN struct {
	FQDNTag    string     `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	Action     string     `form:"action,omitempty"`
	CID        string     `form:"CID,omitempty"`
	FQDNStatus string     `form:"status,omitempty" json:"status,omitempty"`
	FQDNMode   string     `form:"color,omitempty" json:"color,omitempty"`
	GwList     []string   `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	DomainList []*Filters `form:"domain_names[],omitempty" json:"domain_names,omitempty"`
}

type ResultListResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateFQDN(fqdn *FQDN) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=add_fqdn_filter_tag&tag_name=%s", c.CID, fqdn.FQDNTag)
	resp, err := c.Get(path, nil)
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

func (c *Client) DeleteFQDN(fqdn *FQDN) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=del_fqdn_filter_tag&tag_name=%s", c.CID, fqdn.FQDNTag)
	resp, err := c.Get(path, nil)
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

//change state to 'enabled' or 'disabled'
func (c *Client) UpdateFQDNStatus(fqdn *FQDN) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=set_fqdn_filter_tag_state&tag_name=%s&status=%s",
		c.CID, fqdn.FQDNTag, fqdn.FQDNStatus)
	resp, err := c.Get(path, nil)
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

//Change default mode to 'white' or 'black'
func (c *Client) UpdateFQDNMode(fqdn *FQDN) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=set_fqdn_filter_tag_color&tag_name=%s&color=%s",
		c.CID, fqdn.FQDNTag, fqdn.FQDNMode)
	resp, err := c.Get(path, nil)
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

func (c *Client) UpdateDomains(fqdn *FQDN) error {
	fqdn.CID = c.CID
	fqdn.Action = "set_fqdn_filter_tag_domain_names"
	log.Printf("[INFO] Update domains: %#v", fqdn)

	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&tag_name=%s", c.CID, fqdn.Action, fqdn.FQDNTag)
	for i, dn := range fqdn.DomainList {
		body = body + fmt.Sprintf("&domain_names[%d][fqdn]=%s&domain_names[%d]"+
			"[proto]=%s&domain_names[%d][port]=%s", i, dn.FQDN, i, dn.Protocol, i, dn.Port)
	}
	log.Printf("[TRACE] %s %s Body: %s", verb, c.baseURL, body)
	req, err := http.NewRequest(verb, c.baseURL, strings.NewReader(body))
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
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

func (c *Client) AttachGws(fqdn *FQDN) error {
	log.Printf("[TRACE] inside AttachGWs ------------------------------------------------%#v", fqdn)
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=attach_fqdn_filter_tag_to_gw&tag_name=%s", c.CID,
		fqdn.FQDNTag)
	for i := range fqdn.GwList {
		newPath := path + fmt.Sprintf("&gw_name=%s", fqdn.GwList[i])
		resp, err := c.Get(newPath, nil)
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
	}
	return nil
}

func (c *Client) DetachGws(fqdn *FQDN) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=detach_fqdn_filter_tag_from_gw&tag_name=%s", c.CID,
		fqdn.FQDNTag)
	for i := range fqdn.GwList {
		newPath := path + fmt.Sprintf("&gw_name=%s", fqdn.GwList[i])
		resp, err := c.Get(newPath, nil)
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
	}
	return nil
}

func (c *Client) ListFQDNTags() ([]*FQDN, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_fqdn_filter_tags", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if _, ok := data["reason"]; ok {
		log.Printf("[INFO] Couldn't find Aviatrix FQDN tags: %s", data["reason"])
		return nil, ErrNotFound
	}
	tags := make([]*FQDN, 0)
	if val, ok := data["results"]; ok {
		for tag, data := range val.(map[string]interface{}) {
			tagData := data.(map[string]interface{})
			fqdn := &FQDN{
				FQDNTag:    tag,
				FQDNMode:   tagData["wbmode"].(string),
				FQDNStatus: tagData["state"].(string),
			}
			tags = append(tags, fqdn)
		}
	}

	return tags, nil
}

func (c *Client) GetFQDNTag(fqdn *FQDN) (*FQDN, error) {
	tags, err := c.ListFQDNTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.FQDNTag == fqdn.FQDNTag {
			fqdn.FQDNMode = tag.FQDNMode
			fqdn.FQDNStatus = tag.FQDNStatus
			return fqdn, nil
		}
	}
	log.Printf("[INFO] Couldn't find Aviatrix FQDN tag %s", fqdn.FQDNTag)
	return nil, ErrNotFound
}

func (c *Client) ListDomains(fqdn *FQDN) (*FQDN, error) {
	fqdn.CID = c.CID
	fqdn.Action = "list_fqdn_filter_tag_domain_names"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_fqdn_filter_tag_domain_names&tag_name=%s",
		c.CID, fqdn.FQDNTag)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	dn := data
	names := dn["results"].([]interface{})
	for _, domain := range names {
		dn := domain.(map[string]interface{})
		//log.Printf("[TRACE] domain ------------------------->>>>>>>>>>>>: %#v", dn["fqdn"])
		//log.Printf("[TRACE] domain ------------------------->>>>>>>>>>>>: %#v", dn["protocol"])
		//log.Printf("[TRACE] domain ------------------------->>>>>>>>>>>>: %#v", dn["port"])
		fqdnFilter := Filters{
			FQDN:     dn["fqdn"].(string),
			Protocol: dn["proto"].(string),
			Port:     dn["port"].(string),
		}
		//log.Printf("[TRACE] DOMAIN key FOUND ------------------------>>>>>>>>>>>>: %#v",fqdnFilter)
		fqdn.DomainList = append(fqdn.DomainList, &fqdnFilter)
	}
	//value, ok := dn["results"].([]interface{})
	//if ok {
	//    log.Printf("[TRACE] ListDomains FOUND ------------------------------->>>>>>>>>>>>: %#v", value)
	//} else {
	//    log.Printf("[TRACE] ListDomains NOT_FOUND --------------------------->>>>>>>>>>>>: %#v", value)
	//}
	// error when passing value or when passing fqdnFilter
	return fqdn, nil
}

func (c *Client) ListGws(fqdn *FQDN) (*FQDN, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_fqdn_filter_tag_attached_gws&tag_name=%s", c.CID,
		fqdn.FQDNTag)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data ResultListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find Aviatrix FQDN tag names: %s , Reason: %s", fqdn.FQDNTag,
			data.Reason)
		return nil, errors.New(data.Reason)
	}
	fqdn.GwList = data.Results

	return fqdn, nil
}
