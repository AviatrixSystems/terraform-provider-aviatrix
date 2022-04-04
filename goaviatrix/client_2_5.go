package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type APIError struct {
	Message string
}

func checkAndReturnAPIResp25(resp *http.Response, v interface{}, method, path string) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body %q failed: %v", path, err)
	}
	bodyString := buf.String()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		var apiError APIError
		if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&apiError); err != nil {
			return fmt.Errorf("Json Decode failed: %v\n Body: %s", err, bodyString)
		}
		return fmt.Errorf("HTTP %s %q failed: %v\n", method, path, apiError.Message)
	}

	if v != nil {
		if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&v); err != nil {
			return fmt.Errorf("Json Decode failed: %v\n Body: %s", err, bodyString)
		}
	}

	return nil
}

func (c *Client) urlencode25(d map[string]string, path string) (string, error) {
	link := fmt.Sprintf("https://%s/v2.5/api/%s", c.ControllerIP, path)
	Url, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("parsing url: %v", err)
	}
	v := url.Values{}
	for key, val := range d {
		v.Add(key, val)
	}
	Url.RawQuery = v.Encode()
	return Url.String(), nil
}

func (c *Client) GetAPIContext25(ctx context.Context, v interface{}, path string, d map[string]string) error {
	Url, err := c.urlencode25(d, path)
	if err != nil {
		return fmt.Errorf("could not url encode values for path %q: %v", path, err)
	}

	try, maxTries, backoff := 0, 5, 500*time.Millisecond
	var resp *http.Response
	for {
		try++
		resp, err = c.RequestContext25(ctx, "GET", Url, nil)
		if err == nil {
			break
		}

		log.WithFields(log.Fields{
			"try":    try,
			"action": path,
			"err":    err.Error(),
		}).Warnf("HTTP GET request failed")

		if try == maxTries {
			return fmt.Errorf("HTTP Get %s failed: %v", path, err)
		}
		time.Sleep(backoff)
		// Double the backoff time after each failed try
		backoff *= 2
	}

	return checkAndReturnAPIResp25(resp, v, "GET", path)
}

func (c *Client) PostAPIContext25(ctx context.Context, v interface{}, path string, d interface{}) error {
	return c.DoAPIContext25(ctx, "POST", v, path, d)
}

func (c *Client) PutAPIContext25(ctx context.Context, path string, d interface{}) error {
	return c.DoAPIContext25(ctx, "PUT", nil, path, d)
}

func (c *Client) DeleteAPIContext25(ctx context.Context, path string, d interface{}) error {
	return c.DoAPIContext25(ctx, "DELETE", nil, path, d)
}

func (c *Client) DoAPIContext25(ctx context.Context, verb string, v interface{}, path string, d interface{}) error {
	Url := fmt.Sprintf("https://%s/v2.5/api/%s", c.ControllerIP, path)
	resp, err := c.RequestContext25(ctx, verb, Url, d)
	if err != nil {
		return fmt.Errorf("HTTP %s %q failed: %v", verb, path, err)
	}

	return checkAndReturnAPIResp25(resp, v, verb, path)
}

func (c *Client) RequestContext25(ctx context.Context, verb string, Url string, i interface{}) (*http.Response, error) {
	log.Tracef("%s %s", verb, Url)

	try, maxTries, backoff := 0, 2, 500*time.Millisecond
	var req *http.Request
	var err error
	var apiError *APIError
	var resp *http.Response

	if i != nil {
		body, err := json.Marshal(i)
		if err != nil {
			return nil, err
		}
		log.Tracef("%s %s Body: %s", verb, Url, body)
		reader := bytes.NewReader(body)

		req, err = http.NewRequestWithContext(ctx, verb, Url, reader)
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, verb, Url, nil)
	}

	if err != nil {
		return nil, err
	}

	for {
		try++

		// Set CID as Authorization header for v2.5
		req.Header.Set("Authorization", fmt.Sprintf("cid %s", c.CID))

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode == 403 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			resp.Body.Close()

			// Replace resp.Body with new ReadCloser so that other methods can read the buffer again
			resp.Body = io.NopCloser(buf)

			bodyString := buf.String()
			apiError = new(APIError)
			if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(apiError); err != nil {
				return resp, fmt.Errorf("Json Decode into error message failed: %v\n Body: %s", err, bodyString)
			}

			if !strings.Contains(apiError.Message, "Invalid CID") {
				log.Debugf("[DEBUG] API Response Error: %s\n", apiError.Message)
				return resp, err
			}

			log.Tracef("CID invalid or expired. Trying to login again")
			if err = c.Login(); err != nil {
				return resp, err
			}
			log.WithFields(log.Fields{
				"try": try,
				"err": "CID is invalid",
			}).Warnf("HTTP request failed with expired CID")

			if try == maxTries {
				return resp, fmt.Errorf("%v", apiError.Message)
			}
			time.Sleep(backoff)
			// Double the backoff time after each failed try
			backoff *= 2
		} else {
			log.Debugf("[DEBUG] HTTP Response: %v\n", resp)
			return resp, err
		}
	}
}
