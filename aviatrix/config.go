package aviatrix

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

// Config contains the configuration for the Aviatrix provider
// (Username, Password, and Controller IP)
type Config struct {
	Username     string
	Password     string
	ControllerIP string
}

// Client gets the Aviatrix client to access the Controller
// Arguments:
//    None
// Returns:
//    the aviatrix client (from goaviatrix)
//    error (if any)
func (c *Config) Client() (*goaviatrix.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client, err := goaviatrix.NewClient(c.Username, c.Password, c.ControllerIP, &http.Client{Transport: tr})

	log.Printf("[INFO] Aviatrix Client configured for use")

	if client == nil || err != nil {
		log.Printf("[ERROR] unable to create client: %s", err)
	}
	return client, err
}
