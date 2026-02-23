package goaviatrix

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ajg/form"
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

//go:generate moq -rm -out client_mock.go . ClientInterface
type ClientInterface interface {
	DeleteAccount(account *Account) error
	GetAccount(account *Account) (Account, error)
	AuditAccount(ctx context.Context, account *Account) error
	InvalidateCache()
}

// Client for accessing the Aviatrix Controller
type Client struct {
	ClientInterface
	HTTPClient       *http.Client
	Username         string
	Password         string
	CID              string
	ControllerIP     string
	baseURL          string
	IgnoreTagsConfig *IgnoreTagsConfig
	cachedAccounts   []Account
	cacheMutex       sync.Mutex
}

type GetApiTokenResp struct {
	Return  bool         `json:"return"`
	Results ApiTokenInfo `json:"results"`
	Reason  string       `json:"reason"`
}

type ApiTokenInfo struct {
	ApiToken   string `json:"api_token"`
	LegalTerms string `json:"legal_terms"`
	Reason     string `json:"reason"`
}

func (c *Client) GetApiToken() (string, error) {
	apiToken := make(map[string]interface{})
	apiToken["action"] = "get_api_token"
	apiToken["log_enable"] = true

	log.Infof("Getting API token...")
	Url := fmt.Sprintf("https://%s/v2/api", c.ControllerIP)
	resp, err := c.Get(Url, apiToken)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data GetApiTokenResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if !data.Return {
		return "", errors.New(data.Reason)
	}
	log.Tracef("Token is '%s'.", data.Results.ApiToken)
	return data.Results.ApiToken, nil
}

// Login to the Aviatrix controller with the username/password provided in
// the client structure.
// Arguments:
//
//	None
//
// Returns:
//
//	error - if any
func (c *Client) Login() error {
	ApiToken, err := c.GetApiToken()
	if err != nil {
		return err
	}

	account := make(map[string]interface{})
	account["action"] = "login"
	account["username"] = c.Username
	account["password"] = c.Password

	Url := fmt.Sprintf("https://%s/v2/api", c.ControllerIP)
	log.Infof("Parsed Aviatrix login: %s", account["username"])
	resp, err := c.RequestContextLogin(context.Background(), "POST", Url, account, ApiToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data LoginResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	c.CID = data.CID
	return nil
}

func (c *Client) LoginForCloudn() error {
	account := make(map[string]interface{})
	account["action"] = "login"
	account["username"] = c.Username
	account["password"] = c.Password

	log.Infof("Parsed Aviatrix login: %s", account["username"])
	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
//
//	username - the controller username
//	password - the controller password
//	controllerIP - the controller IP/host
//	HTTPClient - the http client object
//
// Returns:
//
//	Client - the newly created client
//	error - if any
//
// See Also:
//
//	init()
func NewClient(username string, password string, controllerIP string, HTTPClient *http.Client, ignoreTagsConfig *IgnoreTagsConfig) (*Client, error) {
	client := &Client{Username: username, Password: password, HTTPClient: HTTPClient, ControllerIP: controllerIP, IgnoreTagsConfig: ignoreTagsConfig}
	return client.init(controllerIP)
}

// init initializes the new client with the given controller IP/host.  Logs
// in to the controller and sets up the http client.
// Arguments:
//
//	controllerIP - the controller host/IP
//
// Returns:
//
//	Client - the updated client object
//	error - if any
func (c *Client) init(controllerIP string) (*Client, error) {
	if len(controllerIP) == 0 {
		return nil, fmt.Errorf("controller ip is not set")
	}
	c.baseURL = "https://" + controllerIP + "/v2/api"
	if c.HTTPClient == nil {
		tr := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			// whynosemgrep: we need to support insecure TLS
			// nosemgrep: problem-based-packs.insecure-transport.go-stdlib.bypass-tls-verification.bypass-tls-verification
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
				MinVersion:         tls.VersionTLS12,
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

func (c *Client) GetContext(ctx context.Context, path string, i interface{}) (*http.Response, error) {
	return c.RequestContext(ctx, "GET", path, i)
}

// PostContext issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostContext(ctx context.Context, path string, i interface{}) (*http.Response, error) {
	return c.RequestContext(ctx, "POST", path, i)
}

// CheckAPIResponseFunc looks at the Reason and Return fields from an API response
// and returns an error
type CheckAPIResponseFunc func(action, method, reason string, ret bool) error

// BasicCheck will only verify that the Return field was set to true
var BasicCheck CheckAPIResponseFunc = func(action, method, reason string, ret bool) error {
	if !ret {
		return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
	}
	return nil
}

// DuplicateBasicCheck will verify that the Return field was set to true
// If the Return is false and Reason contains "already exists", it will return a DuplicateError
var DuplicateBasicCheck CheckAPIResponseFunc = func(action, method, reason string, ret bool) error {
	if !ret {
		err := fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		if strings.Contains(strings.ToLower(reason), "already exists") {
			return DuplicateError{
				Err: err,
			}
		}
		return err
	}
	return nil
}

// PostAPI makes a post request to the Aviatrix API, decodes the response and checks for any errors
func (c *Client) PostAPI(action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	return c.PostAPIContext(context.Background(), action, d, checkFunc)
}

// PostAPIContext makes a post request to the Aviatrix API, decodes the response and checks for any errors
func (c *Client) PostAPIContext(ctx context.Context, action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	resp, err := c.PostContext(ctx, c.baseURL, d)
	if err != nil {
		return fmt.Errorf("HTTP POST %q failed: %w", action, err)
	}
	defer resp.Body.Close()
	return checkAPIResp(resp, action, checkFunc)
}

// PostAPIDownloadContext makes a post request to the Aviatrix API, checks for errors and returns the response body
func (c *Client) PostAPIDownloadContext(ctx context.Context, action string, d interface{}, checkFunc CheckAPIResponseFunc) (io.ReadCloser, error) {
	resp, err := c.PostContext(ctx, c.baseURL, d)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST %q failed: %w", action, err)
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		return nil, checkAPIResp(resp, action, checkFunc)
	}

	return resp.Body, nil
}

// PostAPIWithResponse makes a post request to the Aviatrix API, decodes the response, checks for any errors
// and decodes the response into the return value v.
func (c *Client) PostAPIWithResponse(v interface{}, action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	return c.PostAPIContextWithResponse(context.Background(), v, action, d, checkFunc)
}

// PostAPIContextWithResponse makes a post request to the Aviatrix API, decodes the response, checks for any errors
// and decodes the response into the return value v.
func (c *Client) PostAPIContextWithResponse(ctx context.Context, v interface{}, action string, d interface{}, checkFunc CheckAPIResponseFunc) error {
	resp, err := c.PostContext(ctx, c.baseURL, d)
	if err != nil {
		return fmt.Errorf("HTTP POST %q failed: %w", action, err)
	}
	defer resp.Body.Close()
	return checkAndReturnAPIResp(resp, v, "POST", action, checkFunc)
}

// PostFileAPI will encode the files and parameters with multipart form encoding and POST to the API.
// The API response is decoded and checked with the provided checkFunc
func (c *Client) PostFileAPI(params map[string]string, files []File, checkFunc CheckAPIResponseFunc) error {
	if params["action"] == "" {
		return fmt.Errorf("cannot PostFileAPI without an 'action' in params map")
	}
	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return fmt.Errorf("HTTP POST %q failed: %w", params["action"], err)
	}
	defer resp.Body.Close()
	return checkAPIResp(resp, params["action"], checkFunc)
}

// PostFileAPIContext will encode the files and parameters with multipart form encoding and POST to the API.
// The API response is decoded and checked with the provided checkFunc
func (c *Client) PostFileAPIContext(ctx context.Context, params map[string]string, files []File, checkFunc CheckAPIResponseFunc) error {
	if params["action"] == "" {
		return fmt.Errorf("cannot PostFileAPIContext without an 'action' in params map")
	}
	resp, err := c.PostFileContext(ctx, c.baseURL, params, files)
	if err != nil {
		return fmt.Errorf("HTTP POST %q failed: %w", params["action"], err)
	}
	defer resp.Body.Close()
	return checkAPIResp(resp, params["action"], checkFunc)
}

// ResponseHook is called with the raw response (both initial and poll responses),
// allowing callers to extract custom fields like ha_gw_name.
type ResponseHook func(raw map[string]interface{})

// AsyncPollPayloadFunc returns the payload for polling task status.
type AsyncPollPayloadFunc func(requestID string) interface{}

// asyncCfg holds optional configuration for PostAsyncAPIContext.
type asyncCfg struct {
	onResponse  ResponseHook
	pollPayload AsyncPollPayloadFunc
}

// AsyncOption configures async API behavior.
type AsyncOption func(*asyncCfg)

// WithResponseHook sets a hook to be called with the raw response.
// The hook is called on both the initial response and each poll response,
// allowing callers to extract custom fields like ha_gw_name whenever they appear.
func WithResponseHook(h ResponseHook) AsyncOption {
	return func(c *asyncCfg) { c.onResponse = h }
}

// WithPollPayload sets a custom function to generate the poll payload.
func WithPollPayload(f AsyncPollPayloadFunc) AsyncOption {
	return func(c *asyncCfg) { c.pollPayload = f }
}

// PostAsyncAPI submits an async request and waits for completion.
func (c *Client) PostAsyncAPI(action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error {
	return c.PostAsyncAPIContext(context.Background(), action, i, checkFunc, opts...)
}

//nolint:cyclop,funlen
func (c *Client) PostAsyncAPIContext(ctx context.Context, action string, i interface{}, checkFunc CheckAPIResponseFunc, opts ...AsyncOption) error {
	// Build config with defaults
	cfg := asyncCfg{
		pollPayload: func(requestID string) interface{} {
			return map[string]string{
				"action":     "check_task_status",
				"CID":        c.CID,
				"request_id": requestID,
			}
		},
	}
	for _, o := range opts {
		o(&cfg)
	}

	log.Printf("[DEBUG] Post AsyncAPI %s: %v", action, i)
	resp, err := c.PostContext(ctx, c.baseURL, i)
	if err != nil {
		return fmt.Errorf("HTTP POST %s failed: %w", action, err)
	}

	buf := new(bytes.Buffer)
	if _, readErr := buf.ReadFrom(resp.Body); readErr != nil {
		_ = resp.Body.Close()
		return fmt.Errorf("read response body for %s failed: %w", action, readErr)
	}
	_ = resp.Body.Close()
	bodyString := buf.String()

	var data struct {
		Return bool   `json:"return"`
		Result string `json:"results"`
		Reason string `json:"reason"`
	}

	if err = json.NewDecoder(strings.NewReader(bodyString)).Decode(&data); err != nil {
		return fmt.Errorf("json decode %s failed: %w (body: %s)", action, err, bodyString)
	}
	if !data.Return {
		return fmt.Errorf("rest API %s POST failed to initiate async action: %s", action, data.Reason)
	}

	// Call the start response hook if provided
	if cfg.onResponse != nil {
		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(bodyString), &raw); err == nil {
			cfg.onResponse(raw)
		}
	}

	requestID := data.Result

	const maxPoll = 360
	sleepDuration := time.Second * 10
	var j int
	for ; j < maxPoll; j++ {
		form := cfg.pollPayload(requestID)
		resp, err = c.PostContext(ctx, c.baseURL, form)
		if err != nil {
			// Could be transient HTTP error, e.g. EOF error
			time.Sleep(sleepDuration)
			continue
		}

		buf = new(bytes.Buffer)
		if _, readErr := buf.ReadFrom(resp.Body); readErr != nil {
			_ = resp.Body.Close()
			return fmt.Errorf("read check_task_status body failed: %w", readErr)
		}
		_ = resp.Body.Close()

		pollBodyString := buf.String()
		err = json.Unmarshal([]byte(pollBodyString), &data)
		if err != nil {
			// Only check for status codes after trying to parse JSON because we may get an error with a valid JSON body
			// and that is a valid and actionable response...
			if resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable {
				time.Sleep(sleepDuration)
				continue
			}
			return fmt.Errorf("decode check_task_status failed: %w\n Body: %s", err, pollBodyString)
		}

		// Call the hook on each poll response to capture fields like ha_gw_name
		if cfg.onResponse != nil {
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(pollBodyString), &raw); err == nil {
				cfg.onResponse(raw)
			}
		}

		if !data.Return {
			if data.Reason != "" && data.Reason != "REQUEST_IN_PROGRESS" {
				return fmt.Errorf("rest API %s POST failed: %s", action, data.Reason)
			}

			// Not done yet
			time.Sleep(sleepDuration)
			continue
		}

		// Async API is done, return result of checkFunc
		return checkFunc(action, "Post", data.Result, data.Return)
	}

	// Waited for too long and async API never finished
	return fmt.Errorf("waited %s but upgrade never finished. Please manually verify the upgrade status", maxPoll*sleepDuration)
}

// checkAPIResp will decode the response and check for any errors with the provided checkFunc
func checkAPIResp(resp *http.Response, action string, checkFunc CheckAPIResponseFunc) error {
	var data APIResp
	var b bytes.Buffer
	_, err := b.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body %q failed: %w", action, err)
	}
	body := b.String()
	if err = json.Unmarshal([]byte(body), &data); err != nil {
		return fmt.Errorf("json Decode %q failed: %w\n Body: %s", action, err, body)
	}

	return checkFunc(action, "Post", data.Reason, data.Return)
}

// checkAndReturnAPIResp will decode the response and check for any errors with the provided checkFunc.
// If there are no errors, the response will be put into the return value v.
func checkAndReturnAPIResp(resp *http.Response, v interface{}, method, action string, checkFunc CheckAPIResponseFunc) error {
	var data APIResp
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body %q failed: %w", action, err)
	}
	bodyString := buf.String()

	if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&data); err != nil {
		return fmt.Errorf("json decode into standard format failed: %w (body: %s)", err, bodyString)
	}
	if err := checkFunc(action, method, data.Reason, data.Return); err != nil {
		return err
	}
	if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&v); err != nil {
		return fmt.Errorf("json Decode failed: %w\n Body: %s", err, bodyString)
	}
	return nil
}

// GetAPI makes a GET request to the Aviatrix API
// First, we decode into the generic APIResp struct, then check for errors
// If no errors, we will decode into the user defined structure that is passed in
func (c *Client) GetAPI(v interface{}, action string, d map[string]string, checkFunc CheckAPIResponseFunc) error {
	return c.GetAPIContext(context.Background(), v, action, d, checkFunc)
}

// GetAPIContext makes a GET request to the Aviatrix API
// If the GET request fails we will retry
// First, we decode into the generic APIResp struct, then check for errors
// If no errors, we will decode into the user defined structure that is passed in
func (c *Client) GetAPIContext(
	ctx context.Context,
	v interface{},
	action string,
	d map[string]string,
	checkFunc CheckAPIResponseFunc,
) error {
	Url, err := c.urlEncode(d)
	if err != nil {
		return fmt.Errorf("could not url encode values for action %q: %w", action, err)
	}

	try, maxTries := 0, 5
	backoff := 500 * time.Millisecond

	var resp *http.Response
	for {
		try++

		resp, err = c.GetContext(ctx, Url, nil)
		if err == nil {
			break
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		log.WithFields(log.Fields{
			"try":    try,
			"action": action,
			"err":    err.Error(),
		}).Warnf("HTTP GET request failed")

		if try == maxTries {
			return fmt.Errorf("HTTP Get %s failed after %d tries: %w", action, try, err)
		}

		select {
		case <-time.After(backoff):
			backoff *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if resp == nil || resp.Body == nil {
		return fmt.Errorf("HTTP Get %s returned nil response/body", action)
	}
	defer resp.Body.Close()

	const maxBody = 256 << 10 // 256 KiB
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if readErr != nil {
		return fmt.Errorf("read response body for %s failed: %w", action, readErr)
	}
	if len(body) > maxBody {
		body = body[:maxBody]
	}

	var data APIResp
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("json decode into standard format failed for %s: %w\nBody: %s", action, err, string(body))
	}

	if err := checkFunc(action, "Get", data.Reason, data.Return); err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("json decode into target failed for %s: %w\nBody: %s", action, err, string(body))
	}

	return nil
}

func (c *Client) urlEncode(d map[string]string) (string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("parsing url: %w", err)
	}
	v := url.Values{}
	for key, val := range d {
		v.Add(key, val)
	}
	Url.RawQuery = v.Encode()
	return Url.String(), nil
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
	Path           string
	ParamName      string
	UseFileContent bool   // set to true when using the file content instead of file path
	FileName       string // use when UseFileContent is true
	FileContent    string // use when UseFileContent is true
}

// PostFile will encode the files and parameters with multipart form encoding.
func (c *Client) PostFile(path string, params map[string]string, files []File) (*http.Response, error) {
	body, contentType, err := encodeMultipartFormData(params, files)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.HTTPClient.Do(req)
}

// PostFileContext will encode the files and parameters with multipart form encoding.
func (c *Client) PostFileContext(ctx context.Context, path string, params map[string]string, files []File) (*http.Response, error) {
	body, contentType, err := encodeMultipartFormData(params, files)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.HTTPClient.Do(req)
}

func encodeMultipartFormData(params map[string]string, files []File) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Encode the files
	for _, f := range files {
		if !f.UseFileContent {
			if f.Path == "" {
				continue
			}

			file, err := os.Open(f.Path)
			if err != nil {
				return nil, "", err
			}
			fileContents, err := io.ReadAll(file)
			if err != nil {
				return nil, "", err
			}
			fi, err := file.Stat()
			if err != nil {
				return nil, "", err
			}
			_ = file.Close()
			part, err := createFormFile(f.ParamName, fi.Name(), http.DetectContentType(fileContents), writer)
			if err != nil {
				return nil, "", err
			}
			_, _ = part.Write(fileContents)
		} else {
			fileContents := []byte(f.FileContent)

			part, err := createFormFile(f.ParamName, f.FileName, http.DetectContentType(fileContents), writer)
			if err != nil {
				return nil, "", err
			}
			_, _ = part.Write(fileContents)
		}
	}

	// Encode the other params
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
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

// Request makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) Request(verb string, path string, i interface{}) (*http.Response, error) {
	return c.RequestContext(context.Background(), verb, path, i)
}

// RequestContext makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) RequestContext(ctx context.Context, verb string, path string, i interface{}) (*http.Response, error) {
	log.Tracef("%s %s", verb, path)

	try, maxTries, backoff := 0, 2, 500*time.Millisecond
	var req *http.Request
	var err error
	var data *APIResp
	var resp *http.Response

	for {
		try++

		if i != nil {
			body, err := form.EncodeToValues(i)
			if err != nil {
				return nil, err
			}
			reader := strings.NewReader(body.Encode())
			req, err = http.NewRequestWithContext(ctx, verb, path, reader)
			if err == nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
		if _, readErr := buf.ReadFrom(resp.Body); readErr != nil {
			_ = resp.Body.Close()
			return resp, readErr
		}
		_ = resp.Body.Close()

		// Replace resp.Body with new ReadCloser so that other methods can read the buffer again
		resp.Body = io.NopCloser(buf)

		if strings.Contains(resp.Header.Get("Content-Type"), "json") {
			bodyString := buf.String()
			data = new(APIResp)
			if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(data); err != nil {
				return resp, fmt.Errorf("json Decode into standard format failed: %w\n Body: %s", err, bodyString)
			}

			// Any error not related to expired CID should return
			if !(strings.Contains(data.Reason, "CID is invalid") || strings.Contains(data.Reason, "Invalid session. Please login again.")) {
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
					return resp, fmt.Errorf("failed to parse url: %w", err)
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

func (c *Client) RequestContextLogin(ctx context.Context, verb string, path string, i interface{}, token string) (*http.Response, error) {
	log.Tracef("%s %s", verb, path)

	try, maxTries, backoff := 0, 2, 500*time.Millisecond
	var req *http.Request
	var err error
	var data *APIResp
	var resp *http.Response

	for {
		try++

		if i != nil {
			buf := new(bytes.Buffer)
			if err = form.NewEncoder(buf).Encode(i); err != nil {
				return nil, err
			}
			body := buf.String()
			log.Tracef("%s %s Body: %s", verb, path, body)

			reader := strings.NewReader(body)
			req, err = http.NewRequestWithContext(ctx, verb, path, reader)
			if err == nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Add("X-Access-Key", token)
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
		if _, readErr := buf.ReadFrom(resp.Body); readErr != nil {
			_ = resp.Body.Close()
			return resp, readErr
		}
		_ = resp.Body.Close()

		// Replace resp.Body with new ReadCloser so that other methods can read the buffer again
		resp.Body = io.NopCloser(buf)

		if strings.Contains(resp.Header.Get("Content-Type"), "json") {
			bodyString := buf.String()
			data = new(APIResp)
			if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(data); err != nil {
				return resp, fmt.Errorf("json Decode into standard format failed: %w\n Body: %s", err, bodyString)
			}

			// Any error not related to expired CID should return
			if !(strings.Contains(data.Reason, "CID is invalid") || strings.Contains(data.Reason, "Invalid session. Please login again.")) {
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
					return resp, fmt.Errorf("failed to parse url: %w", err)
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
