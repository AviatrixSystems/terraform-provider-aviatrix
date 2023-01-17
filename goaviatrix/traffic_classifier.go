package goaviatrix

import (
	"context"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PortRange struct {
	Lo int `json:"lo"`
	Hi int `json:"hi"`
}

type TCPolicy struct {
	UUID          string      `json:"uuid,omitempty"`
	Name          string      `json:"name,omitempty"`
	SrcSgs        []string    `json:"src_sgs,omitempty"`
	DstSgs        []string    `json:"dst_sgs,omitempty"`
	PortRanges    []PortRange `json:"port_ranges,omitempty"`
	Protocol      string      `json:"protocol,omitempty"`
	LinkHierarchy string      `json:"link_hierarchy,omitempty"`
	SlaClass      string      `json:"sla_class,omitempty"`
	Logging       bool        `json:"logging,omitempty"`
	RouteType     string      `json:"route_type,omitempty"`
}

type PolicyList struct {
	Policies []TCPolicy `json:"policies"`
}

type TrafficClassifierResp struct {
	TrafficClassifier []PolicyList `json:"traffic_classifier_policies"`
}

func (c *Client) CreateTrafficClassifier(ctx context.Context, policyList *PolicyList) error {
	endpoint := "ipsla/traffic-classifier"

	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) GetTrafficClassifier(ctx context.Context) (*[]PolicyList, error) {
	endpoint := "ipsla/traffic-classifier"

	var data TrafficClassifierResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	if len(data.TrafficClassifier) == 0 {
		return nil, ErrNotFound
	}

	return &data.TrafficClassifier, nil
}

func (c *Client) DeleteTrafficClassifier(ctx context.Context) error {
	endpoint := "ipsla/traffic-classifier"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func DiffSuppressFuncTrafficClassifier(k, old, new string, d *schema.ResourceData) bool {
	pOld, pNew := d.GetChange("policies")
	var policiesOld []map[string]interface{}

	for _, p0 := range pOld.([]interface{}) {
		p1 := p0.(map[string]interface{})
		p2 := make(map[string]interface{})

		for k, v := range p1 {
			if k != "uuid" && k != "port_ranges" {
				p2[k] = v
			}
		}

		pr := p1["port_ranges"].(*schema.Set).List()
		var portRanges []map[string]interface{}
		for _, v := range pr {
			temp := make(map[string]interface{})
			temp["lo"] = v.(map[string]interface{})["lo"]
			temp["hi"] = v.(map[string]interface{})["hi"]
			portRanges = append(portRanges, temp)
		}

		sort.Slice(portRanges, func(i, j int) bool {
			return portRanges[i]["lo"].(string) < portRanges[j]["lo"].(string)
		})

		p2["port_ranges"] = portRanges

		policiesOld = append(policiesOld, p2)
	}

	var policiesNew []map[string]interface{}

	for _, p0 := range pNew.([]interface{}) {
		p1 := p0.(map[string]interface{})
		p2 := make(map[string]interface{})

		for k, v := range p1 {
			if k != "uuid" && k != "port_ranges" {
				p2[k] = v
			}
		}

		pr := p1["port_ranges"].(*schema.Set).List()
		var portRanges []map[string]interface{}
		for _, v := range pr {
			temp := make(map[string]interface{})
			temp["lo"] = v.(map[string]interface{})["lo"]
			temp["hi"] = v.(map[string]interface{})["hi"]
			portRanges = append(portRanges, temp)
		}

		sort.Slice(portRanges, func(i, j int) bool {
			return portRanges[i]["lo"].(string) < portRanges[j]["lo"].(string)
		})

		p2["port_ranges"] = portRanges

		policiesNew = append(policiesNew, p2)
	}

	sort.Slice(policiesOld, func(i, j int) bool {
		return policiesOld[i]["name"].(string) < policiesOld[j]["name"].(string)
	})

	sort.Slice(policiesNew, func(i, j int) bool {
		return policiesNew[i]["name"].(string) < policiesNew[j]["name"].(string)
	})

	return reflect.DeepEqual(policiesOld, policiesNew)
}
