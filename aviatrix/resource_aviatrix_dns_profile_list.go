package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDNSProfileList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDNSProfileListCreate,
		ReadWithoutTimeout:   resourceAviatrixDNSProfileListRead,
		UpdateWithoutTimeout: resourceAviatrixDNSProfileListUpdate,
		DeleteWithoutTimeout: resourceAviatrixDNSProfileListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"profiles": {
				Type:             schema.TypeList,
				Required:         true,
				Description:      "DNS profiles",
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncDNSProfileList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"global": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"lan": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"local_domain_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"wan": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func marshalDNSProfileListInput(d *schema.ResourceData) map[string]interface{} {
	data := make(map[string]interface{})
	var templateName []string
	template := make(map[string]interface{})

	profiles := d.Get("profiles").([]interface{})
	for _, p0 := range profiles {
		p1 := p0.(map[string]interface{})

		templateName = append(templateName, p1["name"].(string))

		for k, v := range p1 {
			if k != "name" {
				var sl []string
				for _, v1 := range v.([]interface{}) {
					sl = append(sl, v1.(string))
				}
				template[k] = sl
			}
		}

		data[p1["name"].(string)] = template
		template = make(map[string]interface{})
	}

	data["template_names"] = templateName

	return data
}

func resourceAviatrixDNSProfileListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	data := marshalDNSProfileListInput(d)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixDNSProfileListReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateDNSProfileList(ctx, data); err != nil {
		return diag.Errorf("could not create DNS profiles: %v", err)
	}

	return resourceAviatrixDNSProfileListReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixDNSProfileListReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixDNSProfileListRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixDNSProfileListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	profileList, err := client.GetDNSProfileList(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DNS profiles: %s", err)
	}

	data := marshalDNSProfileListInput(d)
	var profiles []map[string]interface{}
	for k, v := range profileList {
		if k != "template_names" {
			if stringInSlice(k, data["template_names"].([]string)) {
				profile := make(map[string]interface{})
				profile["name"] = k
				for k1, v1 := range v.(map[string]interface{}) {
					profile[k1] = v1
				}
				profiles = append(profiles, profile)
			}
		}
	}

	if err = d.Set("profiles", profiles); err != nil {
		return diag.Errorf("failed to set DNS profiles: %s\n", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDNSProfileListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)

	if d.HasChange("profiles") {
		pOld, pNew := d.GetChange("profiles")

		var oldProfileNames []string
		var newProfileNames []string

		dataCreate := make(map[string]interface{})
		var nameCreate []string
		templateCreate := make(map[string]interface{})
		dataUpdate := make(map[string]interface{})
		var nameUpdate []string
		templateUpdate := make(map[string]interface{})
		dataDelete := make(map[string]interface{})
		var nameDelete []string

		for _, p0 := range pOld.([]interface{}) {
			p1 := p0.(map[string]interface{})

			oldProfileNames = append(oldProfileNames, p1["name"].(string))
		}

		for _, p0 := range pNew.([]interface{}) {
			p1 := p0.(map[string]interface{})

			newProfileNames = append(newProfileNames, p1["name"].(string))

			if !stringInSlice(p1["name"].(string), oldProfileNames) {
				nameCreate = append(nameCreate, p1["name"].(string))

				for k, v := range p1 {
					if k != "name" {
						var sl []string
						for _, v1 := range v.([]interface{}) {
							sl = append(sl, v1.(string))
						}
						templateCreate[k] = sl
					}
				}

				dataCreate[p1["name"].(string)] = templateCreate
				templateCreate = make(map[string]interface{})
			} else {
				nameUpdate = append(nameUpdate, p1["name"].(string))

				for k, v := range p1 {
					if k != "name" {
						var sl []string
						for _, v1 := range v.([]interface{}) {
							sl = append(sl, v1.(string))
						}
						templateUpdate[k] = sl
					}
				}

				dataUpdate[p1["name"].(string)] = templateUpdate
				templateUpdate = make(map[string]interface{})
			}

		}

		dataCreate["template_names"] = nameCreate
		dataUpdate["template_names"] = nameUpdate

		if len(nameCreate) != 0 {
			err := client.CreateDNSProfileList(ctx, dataCreate)
			if err != nil {
				return diag.Errorf("failed to create new DNS profiles during update: %s", err)
			}
		}

		if len(nameUpdate) != 0 {
			err := client.UpdateDNSProfileList(ctx, dataUpdate)
			if err != nil {
				return diag.Errorf("failed to update DNS profiles: %s", err)
			}
		}

		for _, n := range oldProfileNames {
			if !stringInSlice(n, newProfileNames) {
				nameDelete = append(nameDelete, n)
			}
		}

		dataDelete["template_names"] = nameDelete

		if len(nameDelete) != 0 {
			err := client.DeleteDNSProfileList(ctx, dataDelete)
			if err != nil {
				return diag.Errorf("failed to delete DNS profiles during update: %s", err)
			}
		}
	}

	d.Partial(false)

	return resourceAviatrixDNSProfileListRead(ctx, d, meta)
}

func resourceAviatrixDNSProfileListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	data := marshalDNSProfileListInput(d)

	err := client.DeleteDNSProfileList(ctx, data)
	if err != nil {
		return diag.Errorf("could not delete DNS profiles: %v", err)
	}

	return nil
}
