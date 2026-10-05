package aviatrix

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

//go:embed terraform_provider_version.txt
var Version string

// Config contains the configuration for the Aviatrix provider
type Config struct {
	// Username is the username for accessing the Aviatrix Controller.
	Username string
	// Password is the password for accessing the Aviatrix Controller.
	Password string
	// ControllerIP Is the IP address of the Aviatrix Controller.
	ControllerIP string
	// VerifyCert signals whether to verify the server's certificate chain and
	// hostname.
	VerifyCert bool
	// PathToCACert represents the path to the CA Certificate to use when
	// communicating with the Aviatrix Controller.
	PathToCACert string
	// IgnoreTags represents keys or key prefixes that should be ignored
	// across all resources handled by this provider for situations where
	// external systems are managing certain tags.
	IgnoreTags *goaviatrix.IgnoreTagsConfig
}

// wrapTransport represents an HTTP transport used for setting the user-agent
// for all requests.
type wrapTransport struct {
	transport http.RoundTripper
	userAgent string
}

// RoundTrip implements the HTTP transport interface sending user-agent for all
// requests.
func (wtr *wrapTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("User-Agent", wtr.userAgent)
	return wtr.transport.RoundTrip(req)
}

// defaultTransport returns the default HTTP transport to use when accessing the
// Aviatrix Controller.
func defaultTransport(caCertPath string, verifyCert bool) (*http.Transport, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !verifyCert,
	}

	if verifyCert && caCertPath != "" {
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}
		tlsConfig.RootCAs = caCertPool
	}

	return &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}, nil
}

// getUserAgent returns a string representing the user-agent used by the terraform client.
func getUserAgent() string {
	return fmt.Sprintf("terraform-provider-aviatrix/%s (%s; %s; %s)", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}

// Client returns a client for accessing the Aviatrix Controller
func (c *Config) Client() (*goaviatrix.Client, error) {
	tr, err := defaultTransport(c.PathToCACert, c.VerifyCert)
	if err != nil {
		return nil, err
	}

	// Wrap the transport so we always send the user-agent on all requests.
	wtr := &wrapTransport{
		userAgent: getUserAgent(),
		transport: tr,
	}
	client, err := goaviatrix.NewClient(c.Username, c.Password, c.ControllerIP, &http.Client{Transport: wtr}, c.IgnoreTags)

	log.Printf("[INFO] Aviatrix Client configured for use")

	if client == nil || err != nil {
		log.Printf("[ERROR] unable to create client: %s", err)
	}
	return client, err
}

// mustClient asserts that the meta interface is a valid *goaviatrix.Client.
// This is a helper to satisfy forcetypeassert lints while ensuring the
// provider has its required API client.
func mustClient(meta interface{}) *goaviatrix.Client {
	if client, ok := meta.(*goaviatrix.Client); ok && client != nil {
		return client
	}
	panic("internal error: provider meta is not a valid *goaviatrix.Client; check provider configuration")
}
