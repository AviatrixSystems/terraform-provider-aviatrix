package goaviatrix

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
)

// LoginResp represents the response object from the `login` action
type LoginResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
	CID     string `json:"CID"`
}

// APIResp represents the basic response from any action
type APIResp struct {
	Return bool   `json:"return"`
	Reason string `json:"reason"`
}

// APIRequest represents the basic fields for any request
type APIRequest struct {
	CID    string `form:"CID,omitempty" json:"CID" url:"CID"`
	Action string `form:"action,omitempty" json:"action" url:"action"`
}

// Client for accessing the Aviatrix Controller
type Client struct {
	HTTPClient   *http.Client
	Username     string
	Password     string
	CID          string
	ControllerIP string
	baseURL      string
}

// Login to the Aviatrix controller with the username/password provided in
// the client structure.
// Arguments:
//    None
// Returns:
//    error - if any
func (c *Client) Login() error {
	account := make(map[string]interface{})
	account["action"] = "login"
	account["username"] = c.Username
	account["password"] = c.Password

	log.Infof("Parsed Aviatrix login: %s", account["username"])
	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return err
	}
	var data LoginResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	log.Tracef("CID is '%s'.", data.CID)
	c.CID = data.CID
	return nil
}

// NewClient creates a Client object using the arguments provided.
// Arguments:
//   username - the controller username
//   password - the controller password
//   controllerIP - the controller IP/host
//   HTTPClient - the http client object
// Returns:
//   Client - the newly created client
//   error - if any
// See Also:
//   init()
func NewClient(username string, password string, controllerIP string, HTTPClient *http.Client) (*Client, error) {
	client := &Client{Username: username, Password: password, HTTPClient: HTTPClient, ControllerIP: controllerIP}
	return client.init(controllerIP)
}

// init initializes the new client with the given controller IP/host.  Logs
// in to the controller and sets up the http client.
// Arguments:
//    controllerIP - the controller host/IP
// Returns:
//   Client - the updated client object
//   error - if any
func (c *Client) init(controllerIP string) (*Client, error) {
	if len(controllerIP) == 0 {
		return nil, fmt.Errorf("Aviatrix: Client: Controller IP is not set")
	}

	c.baseURL = "https://" + controllerIP + "/v1/api"

	if c.HTTPClient == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		c.HTTPClient = &http.Client{Transport: tr}
	}
	if err := c.Login(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Get(path string, i interface{}) (*http.Response, error) {
	return c.Request("GET", path, i)
}

// Post issues an HTTP POST request with the given interface form-encoded.
func (c *Client) Post(path string, i interface{}) (*http.Response, error) {
	return c.Request("POST", path, i)
}

// Put issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) Put(path string, i interface{}) (*http.Response, error) {
	return c.Request("PUT", path, i)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(path string, i interface{}) (*http.Response, error) {
	return c.Request("GET", path, i)
}

type File struct {
	Path      string
	ParamName string
}

// PostFile will encode the files and parameters with multipart form encoding.
func (c *Client) PostFile(path string, params map[string]string, files []File) (*http.Response, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Encode the files
	for _, f := range files {
		if f.Path == "" {
			continue
		}
		file, err := os.Open(f.Path)
		if err != nil {
			return nil, err
		}
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}
		_ = file.Close()
		part, err := createFormFile(f.ParamName, fi.Name(), http.DetectContentType(fileContents), writer)
		if err != nil {
			return nil, err
		}
		_, _ = part.Write(fileContents)
	}

	// Encode the other params
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", path, body)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.HTTPClient.Do(req)
}

func createFormFile(fieldname, filename, fileContentType string, w *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", fileContentType)
	return w.CreatePart(h)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// Do performs the HTTP request.
// Arguments:
//   verb - GET, PUT, POST, DELETE, etc
//   req  - the query string (for GET) or body for others
// Returns:
//   http.Response - the HTTP response object (body is closed)
//   []byte - the body string as a byte array
//   error - if any
func (c *Client) Do(verb string, req interface{}) (*http.Response, []byte, error) {
	var err error
	var resp *http.Response
	var url string
	var body []byte
	respdata := new(APIResp)

	// do request
	var loop int
	for {
		url = c.baseURL
		loop = loop + 1
		if verb == "GET" {
			// prepare query string
			v, _ := query.Values(req)
			url = url + "?" + v.Encode()
			resp, err = c.Request(verb, url, nil)
		} else {
			resp, err = c.Request(verb, url, req)
		}

		// check response for error
		if err != nil {
			if loop > 2 {
				return resp, nil, err
			} else {
				continue // try again
			}
		}

		log.Tracef("%s %s: %d", verb, url, resp.StatusCode)
		// decode the json response and look for errors to retry
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ = ioutil.ReadAll(resp.Body)
			if err = json.Unmarshal(body, respdata); err != nil {
				return resp, body, err
			}
			// Check if the CID has expired; if so re-login
			if respdata.Reason == "CID is invalid or expired." && loop < 2 {
				log.Tracef("re-login (expired CID)")
				time.Sleep(500 * time.Millisecond)
				if err = c.Login(); err != nil {
					return resp, body, err
				}
				// update the CID value in the object passed
				s := reflect.ValueOf(req).Elem()
				f := s.FieldByName("CID")
				if f.IsValid() && f.CanSet() {
					f.SetString(c.CID)
				}
				// loop around again using new CID
			} else if !respdata.Return {
				return resp, body, errors.New(respdata.Reason)
			} else {
				return resp, body, nil
			}
		} else {
			return resp, body, errors.New("Status code")
		}
	}

	return resp, body, err
}

// Request makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) Request(verb string, path string, i interface{}) (*http.Response, error) {
	log.Tracef("%s %s", verb, path)
	var req *http.Request
	var err error
	if i != nil {
		buf := new(bytes.Buffer)
		if err = form.NewEncoder(buf).Encode(i); err != nil {
			return nil, err
		}
		body := buf.String()
		log.Tracef("%s %s Body: %s", verb, path, body)
		reader := strings.NewReader(body)
		req, err = http.NewRequest(verb, path, reader)
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = http.NewRequest(verb, path, nil)
	}

	if err != nil {
		return nil, err
	}
	return c.HTTPClient.Do(req)
}
