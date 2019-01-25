package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVersionUpgrade,
		Read:   resourceAviatrixVersionRead,
		Update: resourceAviatrixVersionUpgrade,
		Delete: resourceAviatrixVersionDelete,

		Schema: map[string]*schema.Schema{
			"target_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAviatrixVersionUpgrade(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	version := &goaviatrix.Version{
		Version: d.Get("target_version").(string),
	}
	log.Printf("[INFO] Upgrading Aviatrix controller")

	current, verF, err := client.GetCurrentVersion()
	if err != nil {
		if err.Error() == "valid action required" {
			// assume pre 3.2
			verF = &goaviatrix.AviatrixVersion{
				Major: 3,
				Minor: 1,
				Build: 0,
			}
		} else {
			return fmt.Errorf("unable to get current Controller version: %s (%s)", err, current)
		}
	}

	var err1 error
	if verF.Major <= 3 && verF.Minor <= 2 {
		err1 = client.Pre32Upgrade()
	} else {
		err1 = client.Upgrade(version)
	}

	if err1 != nil {
		return fmt.Errorf("failed to upgrade Aviatrix Controller: %s", err1)
	}
	newCurrent, _, _ := client.GetCurrentVersion()
	log.Printf("Upgrade complete (now %s)", newCurrent)
	d.SetId(newCurrent)

	return resourceAviatrixVersionRead(d, meta)
}

func resourceAviatrixVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	current, _, err := client.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("unable to read current Controller version: %s (%s)", err, current)
	}
	log.Printf("Version read completes (now %s)", current)
	d.Set("version", current)

	return nil
}

func resourceAviatrixVersionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
