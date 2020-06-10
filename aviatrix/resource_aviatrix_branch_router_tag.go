package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixBranchRouterTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixBranchRouterTagCreate,
		Read:   resourceAviatrixBranchRouterTagRead,
		Update: resourceAviatrixBranchRouterTagUpdate,
		Delete: resourceAviatrixBranchRouterTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the tag.",
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						return true
					}
					return false
				},
				Description: "Config to apply to branches that are attached to the tag.",
			},
			"branches": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of branch names to attach to this tag.",
			},
			"commit": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to commit the config to the branches.",
			},
		},
	}
}

func marshalBranchRouterTagInput(d *schema.ResourceData) *goaviatrix.BranchRouterTag {
	brt := &goaviatrix.BranchRouterTag{
		Name:   d.Get("name").(string),
		Config: d.Get("config").(string),
		Commit: d.Get("commit").(bool),
	}

	var brs []string
	for _, s := range d.Get("branches").([]interface{}) {
		brs = append(brs, s.(string))
	}
	brt.Branches = brs

	return brt
}

func resourceAviatrixBranchRouterTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	brt := marshalBranchRouterTagInput(d)

	if err := client.CreateBranchRouterTag(brt); err != nil {
		return err
	}

	d.SetId(brt.Name)
	return nil
}

func resourceAviatrixBranchRouterTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)
	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no branch_router_tag name received. Import Id is %s", id)
		d.SetId(id)
		name = id
	}

	brt := &goaviatrix.BranchRouterTag{
		Name: name,
	}

	brt, err := client.GetBranchRouterTag(brt)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find branch_router_tag %s: %v", name, err)
	}

	d.Set("name", brt.Name)
	d.Set("config", brt.Config)
	if err := d.Set("branches", brt.Branches); err != nil {
		return err
	}

	d.SetId(brt.Name)
	return nil
}

func resourceAviatrixBranchRouterTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	brt := marshalBranchRouterTagInput(d)

	if d.HasChange("config") {
		if err := client.UpdateBranchRouterTagConfig(brt); err != nil {
			return err
		}
	}

	if d.HasChange("branches") {
		if err := client.UpdateBranchRouterTagBranches(brt); err != nil {
			return err
		}
	}

	// User had changed 'commit' to true, commit the tag
	if d.HasChange("commit") && d.Get("commit").(bool) {
		if err := client.CommitBranchRouterTag(brt); err != nil {
			return err
		}
	}

	return nil
}

func resourceAviatrixBranchRouterTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	brt := marshalBranchRouterTagInput(d)

	if err := client.DeleteBranchRouterTag(brt); err != nil {
		return err
	}

	return nil
}
