package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	log "github.com/sirupsen/logrus"
)

// AWSPeer simple struct to hold aws_peer details
type AWSPeer struct {
	Action       string `form:"action,omitempty"`
	CID          string `form:"CID,omitempty"`
	AccountName1 string `form:"peer1_account_name,omitempty"`
	AccountName2 string `form:"peer2_account_name,omitempty"`
	VpcID1       string `form:"peer1_vpc_id,omitempty"`
	VpcID2       string `form:"peer2_vpc_id,omitempty"`
	Region1      string `form:"peer1_region,omitempty"`
	Region2      string `form:"peer2_region,omitempty"`
	RtbList1     string `form:"peer1_rtb_id,omitempty"`
	RtbList2     string `form:"peer2_rtb_id,omitempty"`
}

type AwsPeerAPIResp struct {
	Return  bool              `json:"return"`
	Reason  string            `json:"reason"`
	Results map[string]string `json:"results"`
}

type AwsPeerGetAPIResp struct {
	Return  bool     `json:"return"`
	Reason  string   `json:"reason"`
	Results PairList `json:"results"`
}

type PairList struct {
	PairLists []AWSPeerEdit `json:"pair_list,omitempty"`
}

type AWSPeerEdit struct {
	AWSVpc1 AwsVpcInfo `json:"requester,omitempty"`
	AWSVpc2 AwsVpcInfo `json:"accepter,omitempty"`
}

type AwsVpcInfo struct {
	AccountName   string   `json:"account_name,omitempty"`
	VpcID         string   `json:"vpc_id,omitempty"`
	Region        string   `json:"region,omitempty"`
	RoutingTables []string `json:"peering_route_tables,omitempty"`
}

func (c *Client) CreateAWSPeer(awsPeer *AWSPeer) (string, error) {
	awsPeer.CID = c.CID
	awsPeer.Action = "create_aws_peering"
	resp, err := c.Post(c.baseURL, awsPeer)
	if err != nil {
		return "", errors.New("HTTP Post create_aws_peering failed: " + err.Error())
	}
	defer resp.Body.Close()
	var data AwsPeerAPIResp
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", errors.New("ReadFrom create_aws_peering failed: " + err.Error())
	}
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return "", errors.New("Json Decode create_aws_peering failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return "", errors.New("Rest API create_aws_peering Post failed: " + data.Reason)
	}
	r, _ := regexp.Compile(`pcx-\w+`)
	id := r.FindString(data.Results["text"])
	return id, nil
}

func (c *Client) GetAWSPeer(awsPeer *AWSPeer) (*AWSPeer, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_aws_peerings ") + err.Error())
	}
	listAwsPeering := url.Values{}
	listAwsPeering.Add("CID", c.CID)
	listAwsPeering.Add("action", "list_aws_peerings")
	Url.RawQuery = listAwsPeering.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_aws_peerings failed: " + err.Error())
	}

	var data AwsPeerGetAPIResp
	buf := new(bytes.Buffer)
	defer resp.Body.Close()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("ReadFrom list_aws_peerings failed: " + err.Error())
	}
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_aws_peerings failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		log.Errorf("Couldn't find AWS peering between VPCs %s and %s: %s", awsPeer.VpcID1, awsPeer.VpcID2, data.Reason)
		return nil, errors.New("Rest API list_aws_peerings Get failed: " + data.Reason)
	}
	for i := range data.Results.PairLists {
		if data.Results.PairLists[i].AWSVpc1.VpcID == awsPeer.VpcID1 && data.Results.PairLists[i].AWSVpc2.VpcID == awsPeer.VpcID2 {
			awsPeer.AccountName1 = data.Results.PairLists[i].AWSVpc1.AccountName
			awsPeer.AccountName2 = data.Results.PairLists[i].AWSVpc2.AccountName
			awsPeer.Region1 = data.Results.PairLists[i].AWSVpc1.Region
			awsPeer.Region2 = data.Results.PairLists[i].AWSVpc2.Region
			if len(data.Results.PairLists[i].AWSVpc1.RoutingTables) != 0 {
				awsPeer.RtbList1 = strings.Join(data.Results.PairLists[i].AWSVpc1.RoutingTables, ",")
			}
			if len(data.Results.PairLists[i].AWSVpc2.RoutingTables) != 0 {
				awsPeer.RtbList2 = strings.Join(data.Results.PairLists[i].AWSVpc2.RoutingTables, ",")
			}
			return awsPeer, nil
		}
	}
	log.Errorf("No AWS peering between VPC %s and %s is present.", awsPeer.VpcID1, awsPeer.VpcID2)
	return nil, ErrNotFound
}

func (c *Client) UpdateAWSPeer(awsPeer *AWSPeer) error {
	return nil
}

func (c *Client) DeleteAWSPeer(awsPeer *AWSPeer) error {
	awsPeer.CID = c.CID
	awsPeer.Action = "delete_aws_peering"

	resp, err := c.Post(c.baseURL, awsPeer)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return fmt.Errorf("HTTP Post delete_aws_peering failed: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return fmt.Errorf("HTTP Post delete_aws_peering returned nil response/body")
	}
	defer resp.Body.Close()

	const maxBody = 256 << 10 // 256 KiB
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if readErr != nil {
		return fmt.Errorf("read delete_aws_peering response body failed: %w", readErr)
	}

	if len(body) > maxBody {
		return fmt.Errorf("delete_aws_peering: response body too large (>%d bytes)", maxBody)
	}

	var data APIResp
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("delete_aws_peering: failed to decode json response: %w (body: %s)", err, body)
	}
	if !data.Return {
		return fmt.Errorf("rest api delete_aws_peering post failed: %s", data.Reason)
	}
	return nil
}

func DiffSuppressFuncRtbList1(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("rtb_list1")

	oList, ok := o.([]interface{})
	if !ok {
		return false
	}
	nList, ok := n.([]interface{})
	if !ok {
		return false
	}

	rtbListOld := ExpandStringList(oList)
	rtbListNew := ExpandStringList(nList)

	return Equivalent(rtbListOld, rtbListNew)
}

func DiffSuppressFuncRtbList2(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("rtb_list2")

	oList, ok := o.([]interface{})
	if !ok {
		return false
	}
	nList, ok := n.([]interface{})
	if !ok {
		return false
	}

	rtbListOld := ExpandStringList(oList)
	rtbListNew := ExpandStringList(nList)

	return Equivalent(rtbListOld, rtbListNew)
}
