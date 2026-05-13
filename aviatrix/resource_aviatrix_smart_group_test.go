package aviatrix

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixSmartGroup_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.cidr", "11.0.0.0/16"),

					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.type", "vm"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.account_name", "devops"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.region", "us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.tags.k3", "v3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupBasic() string {
	return `
resource "aviatrix_smart_group" "test" {
	name = "test-smart-group"

	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}

		match_expressions {
			type         = "vm"
			account_name = "devops"
			region       = "us-west-2"
			tags         = {
				k3 = "v3"
			}
		}
	}
}
`
}

func TestAccAviatrixSmartGroup_k8s(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.k8s"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupK8s(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "k8s-test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.type", "k8s"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.k8s_cluster_id", "test-cluster-id"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.k8s_namespace", "test-namespace"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.k8s_service", "test-service"),

					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.type", "k8s"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.k8s_cluster_id", "test-cluster-id"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.k8s_namespace", "test-namespace"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.k8s_pod", "test-pod"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupK8s() string {
	return `
resource "aviatrix_smart_group" "k8s" {
	name = "k8s-test-smart-group"

	selector {
		match_expressions {
			type           = "k8s"
			k8s_cluster_id = "test-cluster-id"
			k8s_namespace  = "test-namespace"
			k8s_service    = "test-service"
		}

		match_expressions {
			type           = "k8s"
			k8s_cluster_id = "test-cluster-id"
			k8s_namespace  = "test-namespace"
			k8s_pod        = "test-pod"
		}
	}
}
`
}

func TestAccAviatrixSmartGroup_reject_bad_k8s_combinations(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "aviatrix_smart_group" "bad-k8s" {
					name = "bad-k8s-test-smart-group"

					selector {
						match_expressions {
							type = "k8s"
							zone = "us-east-2a"
						}
					}
				}
				`,
				ExpectError: regexp.MustCompile("invalid selector combination for k8s resource type"),
			},
		},
	})
}

func TestAccAviatrixSmartGroup_reject_bad_k8s_names(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "aviatrix_smart_group" "bad-k8s" {
					name = "bad-k8s-test-smart-group"

					selector {
						match_expressions {
							type = "k8s"
							k8s_namespace = "-invalid name-"
						}
					}
				}
				`,
				ExpectError: regexp.MustCompile("must be a valid Kubernetes Namespace name"),
			},
		},
	})
}

func TestAccAviatrixSmartGroup_k8s_node(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.k8s_node"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupK8sNode(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "k8s-test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.type", "k8s_node"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.k8s_cluster_id", "test-cluster-id"),

					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.type", "k8s_node"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.account_name", "aws"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupK8sNode() string {
	return `
resource "aviatrix_smart_group" "k8s_node" {
	name = "k8s-test-smart-group"

	selector {
		match_expressions {
			type           = "k8s_node"
			k8s_cluster_id = "test-cluster-id"
		}

		match_expressions {
			type           = "k8s_node"
			account_name   = "aws"
		}
	}
}
`
}

func TestAccAviatrixSmartGroup_reject_bad_k8s_node_combinations(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "aviatrix_smart_group" "bad-k8s" {
					name = "bad-k8s-test-smart-group"

					selector {
						match_expressions {
							type          = "k8s_node"
							k8s_namespace = "test-namespace"
						}
					}
				}
				`,
				ExpectError: regexp.MustCompile(`invalid selector keys for k8s_node resource type. invalid keys \['k8s_namespace']`),
			},
		},
	})
}

func TestAccAviatrixSmartGroup_s2c(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.s2c"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupS2C(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "s2c-test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.s2c", "test-mapped-s2c"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupS2C() string {
	return `
resource "aviatrix_smart_group" "s2c" {
	name = "s2c-test-smart-group"

	selector {
		match_expressions {
			s2c           = "test-mapped-s2c"
		}
	}
}
`
}

func TestAccAviatrixSmartGroup_external_threat(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.external_threat"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupExternalThreat(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "threat-test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.external", "threatiq"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.ext_args.threat_type", "tor"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.ext_args.region", "eastus"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupExternalThreat() string {
	return `
resource "aviatrix_smart_group" "external_threat" {
	name = "threat-test-smart-group"

	selector {
		match_expressions {
			external           = "threatiq"
			ext_args = {
				threat_type   = "tor"
				confidence    = "high"
				region = "eastus"
			}
		}
	}
}
`
}

func TestAccAviatrixSmartGroup_external_geo(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.external_geo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupExternalGeo(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "geo-test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.external", "geo"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.ext_args.country_iso_code", "FR"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.external", "geo"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.ext_args.country_iso_code", "RU"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupExternalGeo() string {
	return `
resource "aviatrix_smart_group" "external_geo" {
	name = "geo-test-smart-group"

	selector {
		match_expressions {
			external           = "geo"
			ext_args = {
				country_iso_code   = "FR"
			}
		}
		match_expressions {
			external           = "geo"
			ext_args = {
				country_iso_code   = "RU"
			}
		}
	}
}
`
}

func TestAccAviatrixSmartGroup_serverless(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.serverless"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupServerless(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "serverless-test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.type", "serverless"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.name", "test-function"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSmartGroupServerless() string {
	return `
resource "aviatrix_smart_group" "serverless" {
	name = "serverless-test-smart-group"

	selector {
		match_expressions {
			type = "serverless"
			name = "test-function"
		}
	}
}
`
}

func testAccCheckSmartGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("no Smart Group resource found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Smart Group ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		smartGroup, err := client.GetSmartGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get Smart Group status: %w", err)
		}

		if smartGroup.UUID != rs.Primary.ID {
			return fmt.Errorf("smart Group ID not found")
		}

		return nil
	}
}

func testAccSmartGroupDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetSmartGroup(context.Background(), rs.Primary.ID)
		expectedError := "App domain not found"
		if err == nil || !strings.Contains(err.Error(), expectedError) {
			return fmt.Errorf("smart group configured when it should be destroyed, want %s, got: %w", expectedError, err)
		}
	}

	return nil
}
