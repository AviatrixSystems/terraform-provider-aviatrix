package goaviatrix

type S2CCaCert struct {
	Action          string `form:"action,omitempty"`
	CID             string `form:"CID,omitempty"`
	TagName         string `form:"name,omitempty"`
	CaCertificate   string `form:"ca_cert"`
	CaCertInstances []CaCertInstance
}

type S2CCaCertResp struct {
	Return  bool             `json:"return,omitempty"`
	Results []CaCertInstance `json:"results,omitempty"`
	Reason  string           `json:"reason,omitempty"`
}

type CaCertInstance struct {
	ID             string `json:"_id,omitempty"`
	CertName       string `json:"cert_name,omitempty"`
	SerialNumber   string `json:"serial_number,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	CommonName     string `json:"common_name,omitempty"`
	ExpirationDate string `json:"expire_date,omitempty"`
}

func (c *Client) CreateS2CCaCert(s2cCaCert *S2CCaCert) error {
	action := "add_s2c_ca_cert"
	params := map[string]string{
		"CID":    c.CID,
		"action": action,
		"name":   s2cCaCert.TagName,
	}

	files := []File{
		{
			Path:      s2cCaCert.CaCertificate,
			ParamName: "ca_cert",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) GetS2CCaCert(s2cCaCert *S2CCaCert) (*S2CCaCert, error) {
	formData := map[string]string{
		"action": "get_s2c_ca_cert_list",
		"CID":    c.CID,
	}
	var data S2CCaCertResp
	err := c.GetAPI(&data, formData["action"], formData, BasicCheck)
	if err != nil {
		return nil, err
	}

	s2cCaCertReturn := &S2CCaCert{
		TagName:       s2cCaCert.TagName,
		CaCertificate: s2cCaCert.CaCertificate,
	}
	for _, certInstance := range data.Results {
		if certInstance.CertName == s2cCaCert.TagName {
			s2cCaCertReturn.CaCertInstances = append(s2cCaCertReturn.CaCertInstances, certInstance)
		}
	}

	if s2cCaCertReturn.CaCertInstances != nil {
		return s2cCaCertReturn, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DeleteCertInstance(caCertInstance *CaCertInstance) error {
	action := "delete_s2c_ca_cert"
	data := map[string]interface{}{
		"action": action,
		"CID":    c.CID,
		"id":     caCertInstance.ID,
	}
	return c.PostAPI(action, data, BasicCheck)
}
