package goaviatrix

import (
	"fmt"
	"encoding/json"
	"errors"
	"strings"
	"log"
    "time"
)

type Version struct {
	CID                         	string `form:"CID,omitempty"`
	Action                  		string `form:"action,omitempty"`
	Version		             		string `form:"version,omitempty" json:"version,omitempty"`
}

type UpgradeResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) Upgrade(version *Version) (error) {
	path := ""
	if(version.Version == "") {
		path = c.baseURL + fmt.Sprintf("?CID=%s&action=upgrade", c.CID)
	} else {
		path = c.baseURL + fmt.Sprintf("?CID=%s&action=upgrade&version=%s", c.CID, version.Version)
	}
	for i := 0; ; i++ {
		resp,err := c.Get(path, nil)
		if err != nil {
			return err
		}
		var data UpgradeResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if(!data.Return){
			if strings.Contains(data.Reason, "Active upgrade in progress.") && i<3 {
				log.Printf("[INFO] Active upgrade is in progress. Retry after 60 secs...")
				time.Sleep(60 * time.Second)
				continue
			}
			return errors.New(data.Reason)
		}
		break
	}
	return nil
}
