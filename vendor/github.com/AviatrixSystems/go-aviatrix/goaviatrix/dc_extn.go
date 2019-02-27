package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"

	//"log"
	//"github.com/davecgh/go-spew/spew"
)

// DCExtn simple struct

type DCExtn struct {
	CID            string `form:"CID,omitempty"`
	Action         string `form:"action,omitempty"`
	CloudType      int    `form:"cloud_type" json:"cloud_type"`
	AccountName    string `form:"account_name" json:"account_name"`
	GwName         string `form:"vpc_name" json:"vpc_name"`
	VpcRegion      string `form:"vpc_reg" json:"vpc_reg"`
	GwSize         string `form:"vpc_size" json:"vpc_size"`
	SubnetCIDR     string `form:"vpc_net" json:"vpc_net"`
	InternetAccess string `form:"internet_access" json:"internet_access"`
	PublicSubnet   string `form:"public_subnet" json:"public_subnet"`
	TunnelType     string `form:"tunnel_type" json:"tunnel_type"`
}

type DCExtnListResp struct {
	Return  bool     `json:"return"`
	Results []DCExtn `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateDCExtn(dc_extn *DCExtn) error {
	dc_extn.CID = c.CID
	dc_extn.Action = "create_container"
	resp, err := c.Post(c.baseURL, dc_extn)
	if err != nil {
		return errors.New("HTTP Post create_container failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode create_container failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API create_container Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetDCExtn(dc_extn *DCExtn) (*DCExtn, error) {
	dc_extn.CID = c.CID
	dc_extn.Action = "list_extended_vpc_peer"
	resp, err := c.Post(c.baseURL, dc_extn)
	if err != nil {
		return nil, errors.New("HTTP Post list_extended_vpc_peer failed: " + err.Error())
	}
	var data DCExtnListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_extended_vpc_peer failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_extended_vpc_peer Post failed: " + data.Reason)
	}
	// dc_extnList:= data.Results
	// for i := range dc_extnList {
	// 	if dc_extnList[i].Source == dc_extn.Source && dc_extnList[i].Nexthop == dc_extn.Nexthop {
	// 		return &dc_extnList[i], nil
	// 	}
	// }
	// log.Printf("Transitive peering with gateways %s and %s with subnet %s not found", dc_extn.Source, dc_extn.Nexthop, dc_extn.ReachableCidr)
	return nil, ErrNotFound
}

func (c *Client) UpdateDCExtn(dcx *DCExtn) error {
	dcx.CID = c.CID
	dcx.Action = "list_cidr_of_available_vpcs"
	resp, err := c.Post(c.baseURL, dcx)
	if err != nil {
		return errors.New("HTTP Post list_cidr_of_available_vpcs failed: " + err.Error())
	}
	var data DCExtnListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode list_cidr_of_available_vpcs failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API list_cidr_of_available_vpcs Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteDCExtn(dcx *DCExtn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_container") + err.Error())
	}
	deleteContainer := url.Values{}
	deleteContainer.Add("CID", c.CID)
	deleteContainer.Add("action", "delete_container")
	deleteContainer.Add("cloud_type", string(dcx.CloudType))
	deleteContainer.Add("gw_name", dcx.GwName)
	Url.RawQuery = deleteContainer.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get delete_container failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_container failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_container Get failed: " + data.Reason)
	}
	return nil
}
