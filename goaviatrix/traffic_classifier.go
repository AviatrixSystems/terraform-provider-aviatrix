package goaviatrix

import (
	"context"
	"fmt"
)

type PortRange struct {
	Lo int `json:"lo"`
	Hi int `json:"hi"`
}

type TCPolicy struct {
	UUID          string      `json:"uuid"`
	Name          string      `json:"name"`
	SrcSgs        []string    `json:"src_sgs"`
	DstSgs        []string    `json:"dst_sgs"`
	PortRanges    []PortRange `json:"port_ranges,omitempty"`
	Protocol      string      `json:"protocol"`
	LinkHierarchy string      `json:"link_hierarchy"`
	SlaClass      string      `json:"sla_class"`
	Logging       bool        `json:"logging"`
	RouteType     string      `json:"route_type"`
}

type PolicyList struct {
	Policies []TCPolicy `json:"policies"`
}

type TrafficClassifierResp struct {
	TrafficClassifier []PolicyList `json:"traffic_classifier_policies"`
}

func (c *Client) CreateTrafficClassifier(ctx context.Context, policyList *PolicyList) (string, error) {
	endpoint := "ipsla/traffic-classifier"

	type resp struct {
		UUID string `json:"uuid"`
	}

	var data resp
	//err := c.PostAPIContext25(ctx, &data, endpoint, policyList)
	err := c.PutAPIContext25(ctx, endpoint, policyList)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetTrafficClassifier(ctx context.Context, uuid string) (*PolicyList, error) {
	//endpoint := fmt.Sprintf("ipsla/traffic-classifier/%s", uuid)
	//
	//var data TrafficClassifierResp
	//err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, trafficClassifier := range data.TrafficClassifier {
	//	for _, policy := range trafficClassifier.Policies {
	//		if policy.UUID == uuid {
	//			return &policy, nil
	//		}
	//	}
	//}

	return nil, ErrNotFound
}

func (c *Client) DeleteTrafficClassifier(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("ipsla/traffic-classifier/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

//func DiffSuppressFuncLinkHierarchy(k, old, new string, d *schema.ResourceData) bool {
//	lOld, lNew := d.GetChange("links")
//	var linksOld []map[string]interface{}
//
//	for _, l0 := range lOld.([]interface{}) {
//		l1 := l0.(map[string]interface{})
//		linksOld = append(linksOld, l1)
//	}
//
//	var linksNew []map[string]interface{}
//
//	for _, l0 := range lNew.([]interface{}) {
//		l1 := l0.(map[string]interface{})
//		linksNew = append(linksNew, l1)
//	}
//
//	sort.Slice(linksOld, func(i, j int) bool {
//		return linksOld[i]["name"].(string) < linksOld[j]["name"].(string)
//	})
//
//	sort.Slice(linksNew, func(i, j int) bool {
//		return linksNew[i]["name"].(string) < linksNew[j]["name"].(string)
//	})
//
//	return reflect.DeepEqual(linksOld, linksNew)
//}
