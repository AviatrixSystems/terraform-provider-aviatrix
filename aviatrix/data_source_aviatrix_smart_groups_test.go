package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccDataSourceAviatrixSmartGroups_basic(t *testing.T) {
	resourceName := "data.aviatrix_smart_groups.test"

	skipAcc := os.Getenv("SKIP_DATA_SMART_GROUPS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Smart Groups tests as SKIP_DATA_SMART_GROUPS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixSmartGroupsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "smart_groups.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "smart_groups.0.uuid"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixSmartGroupsConfigBasic() string {
	return `
resource "aviatrix_smart_group" "test" {
	name = "aaa-smart-group"
	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}
	}
}
data "aviatrix_smart_groups" "test"{
	depends_on = [
        aviatrix_smart_group.test
  ]
}
`
}

func TestAccDataSourceAviatrixSmartGroups_k8s(t *testing.T) {
	resourceName := "data.aviatrix_smart_groups.test"

	skipAcc := os.Getenv("SKIP_DATA_SMART_GROUPS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Smart Groups tests as SKIP_DATA_SMART_GROUPS is set")
	}

	clusterId1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	namespace1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	service1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	clusterId2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	namespace2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	pod2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)

	expect := func(resourceName string, keySuffix string, value string) func(state *terraform.State) error {
		return func(state *terraform.State) error {
			rm := state.RootModule()
			resource := rm.Resources[resourceName]
			attrs := resource.Primary.Attributes

			for k, v := range attrs {
				if value == v {
					if strings.HasSuffix(k, keySuffix) {
						return nil
					}
					return fmt.Errorf("invalid key %s for value %s", k, value)
				}
			}
			return fmt.Errorf("value %s not found", value)
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixSmartGroupsConfigK8s(clusterId1, namespace1, service1, clusterId2, namespace2, pod2),
				Check: resource.ComposeTestCheckFunc(
					expect(resourceName, "."+goaviatrix.K8sClusterIdKey, clusterId1),
					expect(resourceName, "."+goaviatrix.K8sNamespaceKey, namespace1),
					expect(resourceName, "."+goaviatrix.K8sServiceKey, service1),
					expect(resourceName, "."+goaviatrix.K8sClusterIdKey, clusterId2),
					expect(resourceName, "."+goaviatrix.K8sNamespaceKey, namespace2),
					expect(resourceName, "."+goaviatrix.K8sPodNameKey, pod2),
					expect(resourceName, "."+goaviatrix.TypeKey, "k8s"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixSmartGroupsConfigK8s(clusterId1, namespace1, service1, clusterId2, namespace2, pod2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_smart_group" "test" {
	name = "test-smart-group"
	selector {
		match_expressions {
            type           = "k8s"
			k8s_cluster_id = "%s"
		    k8s_namespace  = "%s"
		    k8s_service    = "%s"
		}

		match_expressions {
			type           = "k8s"
			k8s_cluster_id = "%s"
			k8s_namespace  = "%s"
			k8s_pod        = "%s"
		}
	}
}
data "aviatrix_smart_groups" "test"{
	depends_on = [
        aviatrix_smart_group.test
  ]
}
`, clusterId1, namespace1, service1, clusterId2, namespace2, pod2)
}
