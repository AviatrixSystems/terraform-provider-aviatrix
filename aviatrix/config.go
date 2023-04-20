package aviatrix

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

// Config contains the configuration for the Aviatrix provider
// (Username, Password, and Controller IP)
type Config struct {
	Username     string
	Password     string
	ControllerIP string
	VerifyCert   bool
	PathToCACert string
	IgnoreTags   *goaviatrix.IgnoreTagsConfig
}

// Client gets the Aviatrix client to access the Controller
// Arguments:
//
//	None
//
// Returns:
//
//	the aviatrix client (from goaviatrix)
//	error (if any)
func (c *Config) Client() (*goaviatrix.Client, error) {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !c.VerifyCert,
		},
	}

	if c.VerifyCert && c.PathToCACert != "" {
		caCert, err := ioutil.ReadFile(c.PathToCACert)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tr.TLSClientConfig.RootCAs = caCertPool
	}

	client, err := goaviatrix.NewClient(c.Username, c.Password, c.ControllerIP, &http.Client{Transport: tr}, c.IgnoreTags)

	log.Printf("[INFO] Aviatrix Client configured for use")

	if client == nil || err != nil {
		log.Printf("[ERROR] unable to create client: %s", err)
	}
	return client, err
}
