package aviatrix

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixFirewallInstanceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixFirewallInstanceImagesRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID.",
			},
			"firewall_images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of firewall instances associated with fireNet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"firewall_image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the firewall image.",
						},
						"firewall_image_version": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Versions of the firewall image.",
						},
						"firewall_size": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Instance sizes of the firewall image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixFirewallInstanceImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Get("vpc_id").(string)

	firewallInstanceImages, err := client.GetFirewallInstanceImages(vpcId)
	if err != nil {
		return fmt.Errorf("couldn't get firewall instance images: %s", err)
	}

	var images []map[string]interface{}
	for _, image := range *firewallInstanceImages {
		fI := make(map[string]interface{})
		fI["firewall_image"] = image.Image
		listVersion := image.Version
		sort.Sort(sort.Reverse(byVersion(listVersion)))
		fI["firewall_image_version"] = listVersion
		listSize := image.Size
		sort.Sort(sort.Reverse(bySize(listSize)))
		fI["firewall_size"] = listSize
		images = append(images, fI)
	}

	if err = d.Set("firewall_images", images); err != nil {
		return fmt.Errorf("couldn't set firewall_images: %s", err)
	}

	d.SetId(vpcId)
	return nil
}

type byVersion []string

func (p byVersion) Len() int {
	return len(p)
}
func (p byVersion) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p byVersion) Less(i, j int) bool {
	if checkFirstCharacter(p[i]) == "R" {
		if strings.Contains(p[i], "_") {
			return versionFormat1(p[i], p[j], "_")
		} else if strings.Contains(p[i], "-") {
			return versionFormat1(p[i], p[j], "-")
		} else {
			log.Printf("need to add a new method sort this version format")
		}
	} else if checkFirstCharacter(p[i]) == "P" {
		return versionFormat2(p[i], p[j], "-")
	} else {
		return compareVersion(p[i], p[j])
	}
	return false
}

//Rxx.xx-xxx.xxx
//Rxx.xx_revx.x
func versionFormat1(version1, version2, flag string) bool {
	//[Rxx.xx, xxx.xxx]
	//[Rxx.xx, revx.x]
	versionArray1 := strings.Split(version1, flag)
	versionArray2 := strings.Split(version2, flag)
	reg, _ := regexp.Compile("[^0-9.-]+")
	if reg.ReplaceAllString(versionArray1[0], "") == reg.ReplaceAllString(versionArray2[0], "") {
		return compareVersion(reg.ReplaceAllString(versionArray1[1], ""), reg.ReplaceAllString(versionArray2[1], ""))
	}
	return compareVersion(reg.ReplaceAllString(versionArray1[0], ""), reg.ReplaceAllString(versionArray2[0], ""))
}

//PA-VM-xx.xxx.xx
func versionFormat2(version1, version2, flag string) bool {
	versionArray1 := strings.Split(version1, flag)
	versionArray2 := strings.Split(version2, flag)
	return compareVersion(versionArray1[2], versionArray2[2])
}

func compareVersion(version1, version2 string) bool {
	v1, _ := version.NewVersion(version1)
	v2, _ := version.NewVersion(version2)
	return v1.LessThan(v2)
}

func checkFirstCharacter(input string) string {
	firstCharacter := input[0:1]
	return firstCharacter
}

type bySize []string

func (s bySize) Len() int {
	return len(s)
}
func (s bySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bySize) Less(i, j int) bool {
	if strings.Contains(s[i], "-") {
		return compareImageSizeTest(s[i], s[j], "-", 2)
	} else if strings.Contains(s[i], ".") {
		return compareImageSizeTest(s[i], s[j], ".", 1)
	} else if strings.Contains(s[i], "_") {
		return compareImageSizeTest(s[i], s[j], "_", 1)
	}
	return false
}

func compareImageSizeTest(imageSize1, imageSize2, flag string, indexFlag int) bool {
	imageSizeArray1 := strings.Split(imageSize1, flag)
	imageSizeArray2 := strings.Split(imageSize2, flag)
	for index := range imageSizeArray1 {
		if index >= indexFlag {
			reg, err := regexp.Compile("[^0-9]+")
			if err != nil {
				log.Printf("[WARN] Failed to remove character value %s: %v", imageSizeArray1[index], err)
			}
			imageSizeIndex1 := reg.ReplaceAllString(imageSizeArray1[index], "")
			imageSizeIndex2 := reg.ReplaceAllString(imageSizeArray2[index], "")
			int1, err := strconv.Atoi(imageSizeIndex1)
			if err != nil {
				log.Printf("[WARN] Failed to convert string to float %s: %v", imageSizeIndex1, err)
			}
			int2, err := strconv.Atoi(imageSizeIndex2)
			if err != nil {
				log.Printf("[WARN] Failed to convert string to float %s: %v", imageSizeIndex2, err)
			}
			if int1 > int2 {
				return true
			}
			if int1 < int2 {
				return false
			}
		}
		if imageSizeArray1[index] > imageSizeArray2[index] {
			return true
		}
		if imageSizeArray1[index] < imageSizeArray2[index] {
			return false
		}
	}
	return true
}
