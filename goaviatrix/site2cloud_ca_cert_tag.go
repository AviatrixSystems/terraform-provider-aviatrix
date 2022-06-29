package goaviatrix

type S2CCaCertTag struct {
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	TagName        string `form:"name,omitempty"`
	CaCertificates []CaCertInstance
}

type S2CCaCert struct {
	Action          string `form:"action,omitempty"`
	CID             string `form:"CID,omitempty"`
	TagName         string `form:"name,omitempty"`
	CaCertificate   string `form:"ca_cert"`
	CaCertInstances []CaCertInstance
}

type S2CCaCertListResp struct {
	Return  bool             `json:"return,omitempty"`
	Results []CaCertInstance `json:"results,omitempty"`
	Reason  string           `json:"reason,omitempty"`
}

type CaCertInstance struct {
	ID             string `json:"_id,omitempty"`
	TagName        string `json:"s2c_cacert_tag_name,omitempty"`
	SerialNumber   string `json:"serial_number,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	CommonName     string `json:"common_name,omitempty"`
	ExpirationDate string `json:"expire_date,omitempty"`
	CertContent    string `json:"cert_content,omitempty"`
}

func (c *Client) CreateS2CCaCert(s2cCaCert *S2CCaCert) error {
	action := "add_s2c_ca_cert"
	params := map[string]string{
		"CID":                 c.CID,
		"action":              action,
		"s2c_cacert_tag_name": s2cCaCert.TagName,
		"only_one_content":    "true",
	}

	if s2cCaCert.CaCertificate != "" {
		var files []File
		ca := File{
			ParamName:      "ca_cert",
			UseFileContent: true,
			FileName:       "ca.pem", // fake name for ca
			FileContent:    s2cCaCert.CaCertificate,
		}
		files = append(files, ca)
		return c.PostFileAPI(params, files, BasicCheck)
	} else {
		return c.PostAPI(action, params, BasicCheck)
	}

}

func (c *Client) GetS2CCaCertTag(s2cCaCertTag *S2CCaCertTag) (*S2CCaCertTag, error) {
	formData := map[string]string{
		"action":              "get_s2c_ca_cert_list_by_name",
		"CID":                 c.CID,
		"s2c_cacert_tag_name": s2cCaCertTag.TagName,
	}
	var data S2CCaCertListResp
	err := c.GetAPI(&data, formData["action"], formData, BasicCheck)
	if err != nil {
		return nil, err
	}

	s2cCaCertTagReturn := &S2CCaCertTag{
		TagName: s2cCaCertTag.TagName,
	}
	if len(data.Results) != 0 {
		for _, certInstance := range data.Results {
			s2cCaCertTagReturn.CaCertificates = append(s2cCaCertTagReturn.CaCertificates, certInstance)
		}
		return s2cCaCertTagReturn, nil
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
