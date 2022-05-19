package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (c *Client) PostAPIContext2(ctx context.Context, action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	return c.DoAPIContext2(ctx, "POST", action, d, checkFunc)
}

func (c *Client) DoAPIContext2(ctx context.Context, verb string, action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	Url := fmt.Sprintf("https://%s/v2/api", c.ControllerIP)
	resp, err := c.RequestContext2(ctx, verb, Url, d)
	if err != nil {
		return fmt.Errorf("HTTP %s %q failed: %v", verb, Url, err)
	}

	return checkAPIResp(resp, action, checkFunc)
}

func (c *Client) RequestContext2(ctx context.Context, verb string, path string, i interface{}) (*http.Response, error) {
	log.Tracef("%s %s", verb, path)

	try, maxTries, backoff := 0, 2, 500*time.Millisecond
	var req *http.Request
	var err error
	var data *APIResp
	var resp *http.Response

	for {
		try++

		if i != nil {
			body, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			log.Tracef("%s %s Body: %s", verb, path, body)
			reader := bytes.NewReader(body)

			req, err = http.NewRequestWithContext(ctx, verb, path, reader)
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
			}
		} else {
			req, err = http.NewRequestWithContext(ctx, verb, path, nil)
		}

		if err != nil {
			return nil, err
		}

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return resp, err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		resp.Body.Close()

		// Replace resp.Body with new ReadCloser so that other methods can read the buffer again
		resp.Body = io.NopCloser(buf)

		if strings.Contains(resp.Header.Get("Content-Type"), "json") {
			bodyString := buf.String()
			data = new(APIResp)
			if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(data); err != nil {
				return resp, fmt.Errorf("Json Decode into standard format failed: %v\n Body: %s", err, bodyString)
			}

			// Any error not related to expired CID should return
			if !(strings.Contains(data.Reason, fmt.Sprintf("Session %s expired", c.CID))) {
				return resp, err
			}

			log.Tracef("CID invalid or expired. Trying to login again")
			if err = c.Login(); err != nil {
				return resp, err
			}

			// update the CID value in the object passed
			if i != nil {
				// Update CID in POST body
				v := reflect.ValueOf(i)
				if v.Kind() == reflect.Map {
					v.SetMapIndex(reflect.ValueOf("CID"), reflect.ValueOf(c.CID))
				} else {
					s := v.Elem()
					f := s.FieldByName("CID")
					if f.IsValid() && f.CanSet() {
						f.SetString(c.CID)
					}
				}
			} else {
				// Update CID in GET URL
				Url, err := url.Parse(path)
				if err != nil {
					return resp, fmt.Errorf("failed to parse url: %v", err)
				}
				query := Url.Query()
				query["CID"] = []string{c.CID}
				Url.RawQuery = query.Encode()
				path = Url.String()
			}

			log.WithFields(log.Fields{
				"try": try,
				"err": "CID is invalid",
			}).Warnf("HTTP request failed with expired CID")

			if try == maxTries {
				return resp, fmt.Errorf("%v", data.Reason)
			}
			time.Sleep(backoff)
			// Double the backoff time after each failed try
			backoff *= 2
		} else {
			return resp, nil
		}
	}
}
