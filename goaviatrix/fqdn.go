package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Filters struct {
	FQDN     string `form:"fqdn,omitempty" json:"fqdn,omitempty"`
	Protocol string `form:"proto,omitempty" json:"proto,omitempty"`
	Port     string `form:"port,omitempty" json:"port,omitempty"`
}

// Gateway simple struct to hold fqdn details
type FQDN struct {
	FQDNTag         string        `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	Action          string        `form:"action,omitempty"`
	CID             string        `form:"CID,omitempty"`
	FQDNStatus      string        `form:"status,omitempty" json:"status,omitempty"`
	FQDNMode        string        `form:"color,omitempty" json:"color,omitempty"`
	GwFilterTagList []GwFilterTag `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	DomainList      []*Filters    `form:"domain_names[],omitempty" json:"domain_names,omitempty"`
}

type GwFilterTag struct {
	Name         string   `json:"gw_name, omitempty"`
	SourceIPList []string `json:"source_ip_list, omitempty"`
}

type ResultListResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type GetFqdnExceptionRuleResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

type ResultListSourceIPResp struct {
	Return  bool       `json:"return"`
	Results GwSourceIP `json:"results"`
	Reason  string     `json:"reason"`
}

type GwSourceIP struct {
	ConfiguredIPs []string `json:"configured_ips"`
	VpcSubnets    []string `json:"vpc_subnets"`
}

func (c *Client) CreateFQDN(fqdn *FQDN) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for add_fqdn_filter_tag " + err.Error())
	}
	addFQDNFilterTag := url.Values{}
	addFQDNFilterTag.Add("CID", c.CID)
	addFQDNFilterTag.Add("action", "add_fqdn_filter_tag")
	addFQDNFilterTag.Add("tag_name", fqdn.FQDNTag)
	Url.RawQuery = addFQDNFilterTag.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get add_fqdn_filter_tag failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_fqdn_filter_tag failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_fqdn_filter_tag Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteFQDN(fqdn *FQDN) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for del_fqdn_filter_tag " + err.Error())
	}
	delFQDNFilterTag := url.Values{}
	delFQDNFilterTag.Add("CID", c.CID)
	delFQDNFilterTag.Add("action", "del_fqdn_filter_tag")
	delFQDNFilterTag.Add("tag_name", fqdn.FQDNTag)
	Url.RawQuery = delFQDNFilterTag.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get del_fqdn_filter_tag failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode del_fqdn_filter_tag failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API del_fqdn_filter_tag Get failed: " + data.Reason)
	}
	return nil
}

//change state to 'enabled' or 'disabled'
func (c *Client) UpdateFQDNStatus(fqdn *FQDN) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for set_fqdn_filter_tag_state " + err.Error())
	}
	setFQDNFilterTagState := url.Values{}
	setFQDNFilterTagState.Add("CID", c.CID)
	setFQDNFilterTagState.Add("action", "set_fqdn_filter_tag_state")
	setFQDNFilterTagState.Add("tag_name", fqdn.FQDNTag)
	setFQDNFilterTagState.Add("status", fqdn.FQDNStatus)
	Url.RawQuery = setFQDNFilterTagState.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get set_fqdn_filter_tag_state failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode set_fqdn_filter_tag_state failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API set_fqdn_filter_tag_state Get failed: " + data.Reason)
	}
	return nil
}

//Change default mode to 'white' or 'black'
func (c *Client) UpdateFQDNMode(fqdn *FQDN) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for set_fqdn_filter_tag_color " + err.Error())
	}
	setFQDNFilterTagColor := url.Values{}
	setFQDNFilterTagColor.Add("CID", c.CID)
	setFQDNFilterTagColor.Add("action", "set_fqdn_filter_tag_color")
	setFQDNFilterTagColor.Add("tag_name", fqdn.FQDNTag)
	setFQDNFilterTagColor.Add("color", fqdn.FQDNMode)
	Url.RawQuery = setFQDNFilterTagColor.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get set_fqdn_filter_tag_color failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode set_fqdn_filter_tag_color failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API set_fqdn_filter_tag_color Get failed: " + data.Reason)
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
		return errors.New("HTTP NewRequest set_fqdn_filter_tag_domain_names failed: " + err.Error())
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.New("HTTP Post set_fqdn_filter_tag_domain_names failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode set_fqdn_filter_tag_domain_names failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API set_fqdn_filter_tag_domain_names Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachGws(fqdn *FQDN) error {
	//log.Printf("[TRACE] inside AttachGWs ------------------------------------------------%#v", fqdn)
	//Url, err := url.Parse(c.baseURL)
	//if err != nil {
	//	return errors.New("url Parsing failed for attach_fqdn_filter_tag_to_gw " + err.Error())
	//}
	//attachFQDNFilterTagToGw := url.Values{}
	//attachFQDNFilterTagToGw.Add("CID", c.CID)
	//attachFQDNFilterTagToGw.Add("action", "attach_fqdn_filter_tag_to_gw")
	//attachFQDNFilterTagToGw.Add("tag_name", fqdn.FQDNTag)
	//
	//for i := range fqdn.GwList {
	//	attachFQDNFilterTagToGw.Add("gw_name", fqdn.GwList[i])
	//	Url.RawQuery = attachFQDNFilterTagToGw.Encode()
	//	resp, err := c.Get(Url.String(), nil)
	//	if err != nil {
	//		return errors.New("HTTP Get attach_fqdn_filter_tag_to_gw failed: " + err.Error())
	//	}
	//	var data APIResp
	//	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
	//		return errors.New("Json Decode attach_fqdn_filter_tag_to_gw failed: " + err.Error())
	//	}
	//	if !data.Return {
	//		return errors.New("Rest API attach_fqdn_filter_tag_to_gw Get failed: " + data.Reason)
	//	}
	//}
	return nil
}

func (c *Client) DetachGws(fqdn *FQDN, gwList []string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for detach_fqdn_filter_tag_from_gw " + err.Error())
	}
	detachFQDNFilterTagToGw := url.Values{}
	detachFQDNFilterTagToGw.Add("CID", c.CID)
	detachFQDNFilterTagToGw.Add("action", "detach_fqdn_filter_tag_from_gw")
	detachFQDNFilterTagToGw.Add("tag_name", fqdn.FQDNTag)

	for i := range gwList {
		detachFQDNFilterTagToGw.Add("gw_name", gwList[i])
		Url.RawQuery = detachFQDNFilterTagToGw.Encode()
		resp, err := c.Get(Url.String(), nil)
		if err != nil {
			return errors.New("HTTP Get detach_fqdn_filter_tag_from_gw failed: " + err.Error())
		}
		var data APIResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.New("Json Decode detach_fqdn_filter_tag_from_gw failed: " + err.Error())
		}
		if !data.Return {
			return errors.New("Rest API detach_fqdn_filter_tag_from_gw Get failed: " + data.Reason)
		}
	}
	return nil
}

func (c *Client) ListFQDNTags() ([]*FQDN, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("url Parsing failed for list_fqdn_filter_tags " + err.Error())
	}
	listFQDNFilterTags := url.Values{}
	listFQDNFilterTags.Add("CID", c.CID)
	listFQDNFilterTags.Add("action", "list_fqdn_filter_tags")
	Url.RawQuery = listFQDNFilterTags.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_fqdn_filter_tags failed: " + err.Error())
	}

	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_fqdn_filter_tags failed: " + err.Error())
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("url Parsing failed for list_fqdn_filter_tag_domain_names " + err.Error())
	}
	listFQDNFilterTagDomainNames := url.Values{}
	listFQDNFilterTagDomainNames.Add("CID", c.CID)
	listFQDNFilterTagDomainNames.Add("action", "list_fqdn_filter_tag_domain_names")
	listFQDNFilterTagDomainNames.Add("tag_name", fqdn.FQDNTag)
	Url.RawQuery = listFQDNFilterTagDomainNames.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_fqdn_filter_tag_domain_names failed: " + err.Error())
	}
	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_fqdn_filter_tag_domain_names failed: " + err.Error())
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

func (c *Client) ListGws(fqdn *FQDN) ([]string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("url Parsing failed for list_fqdn_filter_tag_attached_gws " + err.Error())
	}
	listFQDNFilterTagAttachedGws := url.Values{}
	listFQDNFilterTagAttachedGws.Add("CID", c.CID)
	listFQDNFilterTagAttachedGws.Add("action", "list_fqdn_filter_tag_attached_gws")
	listFQDNFilterTagAttachedGws.Add("tag_name", fqdn.FQDNTag)
	Url.RawQuery = listFQDNFilterTagAttachedGws.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_fqdn_filter_tag_attached_gws failed: " + err.Error())
	}
	var data ResultListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_fqdn_filter_tag_attached_gws failed: " + err.Error())
	}
	if !data.Return {
		log.Printf("[INFO] Couldn't find Aviatrix FQDN tag names: %s , Reason: %s", fqdn.FQDNTag,
			data.Reason)
		return nil, errors.New("Rest API list_fqdn_filter_tag_attached_gws Get failed: " + data.Reason)
	}

	return data.Results, nil
}

func (c *Client) EnableExceptionRule(fqdn *FQDN) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for enable_fqdn_exception_rule " + err.Error())
	}
	enableFqdnExceptionRule := url.Values{}
	enableFqdnExceptionRule.Add("CID", c.CID)
	enableFqdnExceptionRule.Add("action", "enable_fqdn_exception_rule")
	Url.RawQuery = enableFqdnExceptionRule.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_fqdn_exception_rule failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode detach_fqdn_filter_tag_from_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API detach_fqdn_filter_tag_from_gw Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DisableExceptionRule(fqdn *FQDN) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for disable_fqdn_exception_rule " + err.Error())
	}
	disableFqdnExceptionRule := url.Values{}
	disableFqdnExceptionRule.Add("CID", c.CID)
	disableFqdnExceptionRule.Add("action", "disable_fqdn_exception_rule")
	Url.RawQuery = disableFqdnExceptionRule.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_fqdn_exception_rule failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_fqdn_exception_rule failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_fqdn_exception_rule Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) GetExceptionRuleStatus(fqdn *FQDN) (bool, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return false, errors.New("url Parsing failed for get_fqdn_exception_rule_status " + err.Error())
	}
	getFqdnExceptionRuleStatus := url.Values{}
	getFqdnExceptionRuleStatus.Add("CID", c.CID)
	getFqdnExceptionRuleStatus.Add("action", "get_fqdn_exception_rule_status")
	Url.RawQuery = getFqdnExceptionRuleStatus.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return false, errors.New("HTTP Get get_fqdn_exception_rule_status failed: " + err.Error())
	}
	data := GetFqdnExceptionRuleResp{
		Return:  false,
		Results: "",
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, errors.New("Json Decode get_fqdn_exception_rule_status failed: " + err.Error())
	}
	if !data.Return {
		return false, errors.New("Rest API get_fqdn_exception_rule_status Get failed: " + data.Reason)
	}

	if data.Results == "disabled" {
		return false, nil
	}
	return true, nil
}

func (c *Client) AttachTagToGw(fqdn *FQDN, gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for attach_fqdn_filter_tag_to_gw " + err.Error())
	}
	attachFQDNFilterTagToGw := url.Values{}
	attachFQDNFilterTagToGw.Add("CID", c.CID)
	attachFQDNFilterTagToGw.Add("action", "attach_fqdn_filter_tag_to_gw")
	attachFQDNFilterTagToGw.Add("tag_name", fqdn.FQDNTag)
	attachFQDNFilterTagToGw.Add("gw_name", gateway.GwName)
	Url.RawQuery = attachFQDNFilterTagToGw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get attach_fqdn_filter_tag_to_gw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode attach_fqdn_filter_tag_to_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API attach_fqdn_filter_tag_to_gw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateSourceIPFilters(fqdn *FQDN, gateway *Gateway, sourceIPs []string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for update_fqdn_filter_tag_source_ip_filters " + err.Error())
	}
	attachFQDNFilterTagToGw := url.Values{}
	attachFQDNFilterTagToGw.Add("CID", c.CID)
	attachFQDNFilterTagToGw.Add("action", "update_fqdn_filter_tag_source_ip_filters")
	attachFQDNFilterTagToGw.Add("tag_name", fqdn.FQDNTag)
	attachFQDNFilterTagToGw.Add("gateway_name", gateway.GwName)
	if len(sourceIPs) != 0 {
		for i := range sourceIPs {
			attachFQDNFilterTagToGw.Add("source_ips["+strconv.Itoa(i)+"]", sourceIPs[i])
		}
	}

	Url.RawQuery = attachFQDNFilterTagToGw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get update_fqdn_filter_tag_source_ip_filters failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode update_fqdn_filter_tag_source_ip_filters failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API update_fqdn_filter_tag_source_ip_filters Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetGwFilterTagList(fqdn *FQDN) (*FQDN, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("url Parsing failed for list_fqdn_filter_tag_source_ip_filters: " + err.Error())
	}

	listGws, err := c.ListGws(fqdn)
	if err != nil {
		return nil, errors.New("failed for list_fqdn_filter_tag_source_ip_filters: " + err.Error())
	}

	var gwFilterTagList []GwFilterTag

	for i := range listGws {
		listFQDNFilterTagSourceIpFilters := url.Values{}
		listFQDNFilterTagSourceIpFilters.Add("CID", c.CID)
		listFQDNFilterTagSourceIpFilters.Add("action", "list_fqdn_filter_tag_source_ip_filters")
		listFQDNFilterTagSourceIpFilters.Add("tag_name", fqdn.FQDNTag)
		listFQDNFilterTagSourceIpFilters.Add("gateway_name", listGws[i])

		Url.RawQuery = listFQDNFilterTagSourceIpFilters.Encode()
		resp, err := c.Get(Url.String(), nil)
		if err != nil {
			return nil, errors.New("HTTP Get list_fqdn_filter_tag_source_ip_filters failed: " + err.Error())
		}
		var data ResultListSourceIPResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.New("Json Decode list_fqdn_filter_tag_source_ip_filters failed: " + err.Error())
		}
		if !data.Return {
			return nil, errors.New("Rest API list_fqdn_filter_tag_source_ip_filters Get failed: " + data.Reason)
		}

		var gwFilterTag GwFilterTag
		gwFilterTag.Name = listGws[i]
		sourceIPs := make([]string, 0)
		for j := range data.Results.ConfiguredIPs {
			sourceIPs = append(sourceIPs, strings.Split(data.Results.ConfiguredIPs[j], "~~")[0])
		}
		gwFilterTag.SourceIPList = sourceIPs
		gwFilterTagList = append(gwFilterTagList, gwFilterTag)
	}

	fqdn.GwFilterTagList = gwFilterTagList
	return fqdn, nil
}
