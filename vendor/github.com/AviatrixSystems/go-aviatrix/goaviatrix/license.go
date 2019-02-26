package goaviatrix

import (
	"encoding/json"
	"errors"
)

type CustomerRequest struct {
	APIRequest
	CustomerID string `form:"customer_id,omitempty" json:"CustomerID" url:"customer_id"`
}
type License struct {
	Verified   int    `json:"Verified"`
	Type       string `json:"Type"`
	Expiration string `json:"Expiration"`
	Allocated  int    `json:"Allocated"`
	IssueDate  string `json:"IssueDate"`
	Quantity   int    `json:"Quantity"`
}

type ViewLicenseList struct {
	LicenseList []License `json:"license_list"`
}

type SetLicenseList struct {
	LicenseList []map[string]License `json:"license_list"`
}

type DeleteLicenseList struct {
	LicenseList []map[string]string `json:"license_list"`
}

type ViewLicenseResponse struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}
type SetLicenseResponse struct {
	Return     bool           `json:"return"`
	Results    SetLicenseList `json:"results"`
	Reason     string         `json:"reason"`
	CustomerID string         `json:"CustomerID"`
}

type DeleteLicenseResponse struct {
	Return  bool              `json:"return"`
	Results DeleteLicenseList `json:"results"`
	Reason  string            `json:"reason"`
}

func (c *Client) SetCustomerID(customerID string) (*SetLicenseList, error) {
	customerNew := new(CustomerRequest)
	customerNew.CustomerID = customerID
	customerNew.CID = c.CID
	customerNew.Action = "setup_customer_id"
	var response SetLicenseResponse
	_, body, err := c.Do("GET", customerNew)
	if err != nil {
		return nil, errors.New("HTTP Get setup_customer_id failed: " + err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("Json Decode setup_customer_id failed: " + err.Error())
	}
	return &response.Results, err
}

func (c *Client) DeleteCustomerID() (*DeleteLicenseList, error) {
	customerNew := new(CustomerRequest)
	customerNew.CustomerID = " "
	customerNew.CID = c.CID
	customerNew.Action = "setup_customer_id"
	var response DeleteLicenseResponse
	_, body, err := c.Do("GET", customerNew)
	if err != nil {
		return nil, errors.New("HTTP Get setup_customer_id failed: " + err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("Json Decode setup_customer_id failed: " + err.Error())
	}
	return &response.Results, err
}

func (c *Client) GetCustomerID() (string, error) {
	customerNew := new(CustomerRequest)
	customerNew.CID = c.CID
	customerNew.Action = "list_customer_id"
	var response ViewLicenseResponse
	_, body, err := c.Do("GET", customerNew)
	if err != nil {
		return "", errors.New("HTTP Get list_customer_id failed: " + err.Error())
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return "", errors.New("Json Decode list_customer_id failed: " + err.Error())
	}
	return response.Results, nil
}
