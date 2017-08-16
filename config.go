package aviatrix

import (
	"crypto/tls"
	"log"
	"net/http"
	"github.com/go-aviatrix/goaviatrix"
)

type Config struct {
	Username string
	Password string
	ControllerIP      string
}

type AviatrixClient struct {
	AccessKey  string
	SecretKey  string
	Host       string
	HttpClient *http.Client
}

func (c *Config) Client() (*goaviatrix.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client,_ := goaviatrix.NewClient(c.Username, c.Password, c.ControllerIP, &http.Client{Transport: tr})

	log.Printf("[INFO] Aviatrix Client configured for use")

	return client, nil
}
