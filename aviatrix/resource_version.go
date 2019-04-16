package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVersionCreate,
		Read:   resourceAviatrixVersionRead,
		Update: resourceAviatrixVersionUpdate,
		Delete: resourceAviatrixVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"target_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The release version number to which the controller will be upgraded to.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current version of the controller.",
			},
		},
	}
}

func resourceAviatrixVersionCreate(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("version", newCurrent)
	d.SetId(newCurrent)

	return nil
}

func resourceAviatrixVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	curVersion := d.Get("version").(string)
	if curVersion == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no version received. Import Id is %s", id)
		d.Set("version", id)
		d.Set("target_version", id)
		d.SetId(id)
		return nil
	}

	current, _, err := client.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("unable to read current Controller version: %s (%s)", err, current)
	}
	log.Printf("Version read completes (now %s)", current)
	d.SetId(current)

	latestVersion, _ := client.GetLatestVersion()
	if latestVersion != "" && current != latestVersion {
		d.Set("target_version", current)
	}
	d.Set("version", current)

	return nil
}

func resourceAviatrixVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	curVersion := d.Get("version").(string)
	cur := strings.Split(curVersion, ".")

	latestVersion, _ := client.GetLatestVersion()
	latest := strings.Split(latestVersion, ".")

	targetVersion := d.Get("target_version").(string)

	if targetVersion == "latest" {
		if latestVersion != "" {
			for i := range cur {
				if cur[i] != latest[i] {
					return resourceAviatrixVersionCreate(d, meta)
				}
			}
		}
		return nil
	}
	return resourceAviatrixVersionCreate(d, meta)
}

func resourceAviatrixVersionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
