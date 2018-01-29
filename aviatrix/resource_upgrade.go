package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixUpgrade() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVersionUpgrade,
		Read:   resourceAviatrixVersionRead,
		Update: resourceAviatrixVersionUpgrade,
		Delete: resourceAviatrixVersionDelete,

		Schema: map[string]*schema.Schema{
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixVersionUpgrade(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	version := &goaviatrix.Version{
		Version: d.Get("version").(string),
	}

	log.Printf("[INFO] Upgrading Aviatrix controller")

	err := client.Upgrade(version)
	if err != nil {
		return fmt.Errorf("Failed to upgrade Aviatrix Controller: %s", err)
	}
	return nil
}

func resourceAviatrixVersionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAviatrixVersionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
