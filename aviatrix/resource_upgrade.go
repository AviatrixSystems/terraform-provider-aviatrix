package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
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

	current, verf, err := client.GetCurrentVersion()
	if err != nil {
		if err.Error() == "valid action required" {
			// assume pre 3.2
			verf = &goaviatrix.AviatrixVersion{
				Major: 3,
				Minor: 1,
				Build: 0,
			}
		} else {
			return fmt.Errorf("Unable to get current Controller version: %s (%s)", err, current)
		}
	}

	var errv error
	if verf.Major <= 3 && verf.Minor <= 2 {
		errv = client.Pre32Upgrade()
	} else {
		errv = client.Upgrade(version)
	}
	if errv != nil {
		return fmt.Errorf("Failed to upgrade Aviatrix Controller: %s", errv)
	}
	newCurrent, _, _ := client.GetCurrentVersion()
	log.Printf("Upgrade complete (now %s)", newCurrent)
	d.SetId(newCurrent)

	return nil
}

func resourceAviatrixVersionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAviatrixVersionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
