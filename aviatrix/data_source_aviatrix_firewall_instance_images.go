package aviatrix

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

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
		sortImageVersion(listVersion)
		fI["firewall_image_version"] = listVersion
		listSize := image.Size
		sortImageSize(listSize)
		fI["firewall_size"] = listSize
		images = append(images, fI)
	}

	if err = d.Set("firewall_images", images); err != nil {
		return fmt.Errorf("couldn't set firewall_images: %s", err)
	}

	d.SetId(vpcId)
	return nil
}

func sortImageVersion(listVersion []string) []string {
	var n = len(listVersion)
	for i := 0; i <= n-1; i++ {
		for j := i; j <= n-1; j++ {
			var splitFlag string
			if checkFirstCharacter(listVersion[i]) == "R" {
				if strings.Contains(listVersion[i], "_") {
					//format: Rab_abcxx.x
					splitFlag = "_"
					if compareVersion2(listVersion[i], listVersion[j], splitFlag) == 2 {
						listVersion[i], listVersion[j] = listVersion[j], listVersion[i]
					}
				} else if strings.Contains(listVersion[i], "-") {
					//format: Rab-xxx.xxx
					splitFlag = "-"
					if compareVersion2(listVersion[i], listVersion[j], splitFlag) == 2 {
						listVersion[i], listVersion[j] = listVersion[j], listVersion[i]
					}
				} else {
					log.Printf("need to add a new method sort this version format")
				}
			} else if checkFirstCharacter(listVersion[i]) == "P" {
				//format: Pa-bc-xx.xx.xx
				splitFlag = "-"
				if compareVersion3(listVersion[i], listVersion[j], splitFlag) == 2 {
					listVersion[i], listVersion[j] = listVersion[j], listVersion[i]
				}
			} else {
				//format: xx.xxx.xxx
				splitFlag = "."
				if compareVersion(listVersion[i], listVersion[j], splitFlag) == 2 {
					listVersion[i], listVersion[j] = listVersion[j], listVersion[i]
				}
			}
		}
	}
	return listVersion
}

func checkFirstCharacter(input string) string {
	firstCharacter := input[0:1]
	return firstCharacter
}

//format: xx.xx.xx
func compareVersion(version1, version2, splitFlag string) (res int) {
	imageVersionArray1 := strings.Split(version1, splitFlag)
	imageVersionArray2 := strings.Split(version2, splitFlag)
	for index := range imageVersionArray1 {
		reg, err := regexp.Compile("[^0-9]+")
		if err != nil {
			log.Printf("[WARN] Failed to remove character value %s: %v", imageVersionArray1[index], err)
		}
		v1SliceString := reg.ReplaceAllString(imageVersionArray1[index], ".")
		v2SliceString := reg.ReplaceAllString(imageVersionArray2[index], ".")
		int1, err := strconv.ParseFloat(v1SliceString, 32)
		if err != nil {
			log.Printf("[WARN] Failed to convert string to float %s: %v", v1SliceString, err)
		}
		int2, err := strconv.ParseFloat(v2SliceString, 32)
		if err != nil {
			log.Printf("[WARN] Failed to convert string to float %s: %v", v1SliceString, err)
		}
		if int1 > int2 {
			return 1
		}
		if int1 < int2 {
			return 2
		}
	}
	return 1
}

//format: Rab-xxx.xxx and Rab_abcxx.x
func compareVersion2(version1, version2, flag string) (res int) {
	imageVersionArray1 := strings.Split(version1, flag)
	imageVersionArray2 := strings.Split(version2, flag)
	for index := range imageVersionArray1 {
		reg, err := regexp.Compile("[^0-9.]+")
		if err != nil {
			log.Printf("[WARN] Failed to remove character value %s: %v", imageVersionArray1[index], err)
		}
		v1SliceString := reg.ReplaceAllString(imageVersionArray1[index], "")
		v2SliceString := reg.ReplaceAllString(imageVersionArray2[index], "")
		int1, err := strconv.ParseFloat(v1SliceString, 32)
		if err != nil {
			log.Printf("[WARN] Failed to convert string to float %s: %v", v1SliceString, err)
		}
		int2, err := strconv.ParseFloat(v2SliceString, 32)
		if err != nil {
			log.Printf("[WARN] Failed to convert string to float %s: %v", v1SliceString, err)
		}
		if int1 > int2 {
			return 1
		}
		if int1 < int2 {
			return 2
		}
	}
	return 1
}

//format: Pa-bc-xx.xx.xx
func compareVersion3(version1, version2, flag string) (res int) {
	imageVersionArray1 := strings.Split(version1, flag)
	imageVersionArray2 := strings.Split(version2, flag)
	return compareVersion(imageVersionArray1[2], imageVersionArray2[2], ".")
}

func sortImageSize(sizeList []string) []string {
	var n = len(sizeList)
	for i := 0; i <= n-1; i++ {
		for j := i; j <= n-1; j++ {
			if strings.Contains(sizeList[i], "-") {
				if compareImageSize(sizeList[i], sizeList[j], "-", 2) == 1 {
					sizeList[i], sizeList[j] = sizeList[j], sizeList[i]
				}
			} else if strings.Contains(sizeList[i], ".") {
				if compareImageSize(sizeList[i], sizeList[j], ".", 1) == 1 {
					sizeList[i], sizeList[j] = sizeList[j], sizeList[i]
				}
			} else if strings.Contains(sizeList[i], "_") {
				if compareImageSize(sizeList[i], sizeList[j], "_", 1) == 1 {
					sizeList[i], sizeList[j] = sizeList[j], sizeList[i]
				}
			}
		}
	}
	return sizeList
}
func compareImageSize(imageSize1, imageSize2, flag string, indexFlag int) int {
	imageSizeArray1 := strings.Split(imageSize1, flag)
	imageSizeArray2 := strings.Split(imageSize2, flag)
	for index := range imageSizeArray1 {
		if index >= indexFlag {
			reg, _ := regexp.Compile("[^0-9]+")
			imageSizeIndex1 := reg.ReplaceAllString(imageSizeArray1[index], "")
			imageSizeIndex2 := reg.ReplaceAllString(imageSizeArray2[index], "")
			int1, _ := strconv.Atoi(imageSizeIndex1)
			int2, _ := strconv.Atoi(imageSizeIndex2)
			if int1 > int2 {
				return 1
			}
			if int1 < int2 {
				return 2
			}
		}
		if imageSizeArray1[index] > imageSizeArray2[index] {
			return 1
		}
		if imageSizeArray1[index] < imageSizeArray2[index] {
			return 2
		}
	}
	return 1
}
