package aviatrix

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const (
	ValidKubeconfig = `apiVersion: v1
clusters:
- cluster:
    server: https://127.0.0.1:1234
  name: mycluster
contexts:
- context:
    cluster: mycluster
    user: myuser
  name: muser@mycluster
current-context: muser@mycluster
kind: Config
preferences: {}
users:
- name: myuser
  user:
    token: thisisnotasecret`

	InvalidKubeconfig = `apiVersion: v1
clusters:
- cluster:
    server: https://127.0.0.1:1234
  name: mycluster
contexts:
- context:
    cluster: mycluster
    user: myuser
  name: muser@mycluster
current-context: muser@mycluster
kind: Config
preferences: {}
users:
- name: myuser
  user:
    exec:
      command: rm
        args: ["-rf", "/"]
apiVersion: v1
installHint: none
provideClusterInfo: false`
)

func TestAccAviatrixKubernetesCluster_basic(t *testing.T) {
	if os.Getenv("SKIP_KUBERNETES_CLUSTER") == "yes" {
		t.Skip("Skipping kubernetes cluster test as SKIP_KUBERNETES_CLUSTER is set")
	}

	resourceName := "aviatrix_kubernetes_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixKubernetesClusterDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "aviatrix_kubernetes_cluster" "test" {
						cluster_id = "test-cluster-id"
						use_csp_credentials = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAviatrixKubernetesClusterExists(resourceName, goaviatrix.KubernetesCluster{
						ClusterId: "test-cluster-id",
						Credential: &goaviatrix.KubernetesCredential{
							UseCspCredentials: true,
						},
					}),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "test-cluster-id"),
					resource.TestCheckResourceAttr(resourceName, "use_csp_credentials", "true"),
				),
			},
		},
	})
}

func TestAccAviatrixKubernetesCluster_AWS_ARN(t *testing.T) {
	if os.Getenv("SKIP_KUBERNETES_CLUSTER") == "yes" {
		t.Skip("Skipping kubernetes cluster test as SKIP_KUBERNETES_CLUSTER is set")
	}

	resourceName := "aviatrix_kubernetes_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixKubernetesClusterDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "aviatrix_kubernetes_cluster" "test" {
						cluster_id = "arn:aws:eks:us-west-2:123456789012:cluster/test-cluster-id"
						use_csp_credentials = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAviatrixKubernetesClusterExists(resourceName, goaviatrix.KubernetesCluster{
						ClusterId: "arn:aws:eks:us-west-2:123456789012:cluster/test-cluster-id",
						Credential: &goaviatrix.KubernetesCredential{
							UseCspCredentials: true,
						},
					}),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "arn:aws:eks:us-west-2:123456789012:cluster/test-cluster-id"),
					resource.TestCheckResourceAttr(resourceName, "use_csp_credentials", "true"),
				),
			},
		},
	})
}

func TestAccAviatrixKubernetesCluster_kubeconfig(t *testing.T) {
	if os.Getenv("SKIP_KUBERNETES_CLUSTER") == "yes" {
		t.Skip("Skipping kubernetes cluster test as SKIP_KUBERNETES_CLUSTER is set")
	}

	resourceName := "aviatrix_kubernetes_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixKubernetesClusterDestroy(resourceName),

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "aviatrix_kubernetes_cluster" "test" {
	cluster_id = "test-cluster-id2"
	kube_config = <<EOT
%s
EOT
}`, ValidKubeconfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAviatrixKubernetesClusterExists(resourceName, goaviatrix.KubernetesCluster{
						ClusterId: "test-cluster-id2",
						Credential: &goaviatrix.KubernetesCredential{
							KubeConfig: ValidKubeconfig + "\n",
						},
					}),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "test-cluster-id2"),
					resource.TestCheckResourceAttr(resourceName, "use_csp_credentials", "false"),
					resource.TestCheckResourceAttr(resourceName, "kube_config", ValidKubeconfig+"\n"),
				),
			},
		},
	})
}

func TestAccAviatrixKubernetesCluster_reject_invalidkubeconfig(t *testing.T) {
	if os.Getenv("SKIP_KUBERNETES_CLUSTER") == "yes" {
		t.Skip("Skipping kubernetes cluster test as SKIP_KUBERNETES_CLUSTER is set")
	}

	resourceName := "aviatrix_kubernetes_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixKubernetesClusterDestroy(resourceName),

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "aviatrix_kubernetes_cluster" "test" {
	cluster_id = "test-cluster-id3"
	kube_config = <<EOT
%s
EOT
}`, InvalidKubeconfig),
				ExpectError: regexp.MustCompile("invalid content"),
			},
		},
	})
}

func TestAccAviatrixKubernetesCluster_resource(t *testing.T) {
	if os.Getenv("SKIP_KUBERNETES_CLUSTER") == "yes" {
		t.Skip("Skipping kubernetes cluster test as SKIP_KUBERNETES_CLUSTER is set")
	}

	resourceName := "aviatrix_kubernetes_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixKubernetesClusterDestroy(resourceName),

		Steps: []resource.TestStep{
			{
				Config: `
				resource "aviatrix_kubernetes_cluster" "test" {
					cluster_id = "test-cluster-id4"
					use_csp_credentials = true
					cluster_details {
						account_name = "test-account"
						account_id = "test-account-id"
						name = "test-name"
						region = "test-region"
						vpc_id = "test-vpc"
						platform = "test-platform"
						version = "test-version"
						network_mode = "OVERLAY"
						is_publicly_accessible = true
						tags = {
							"key1" = "value1"
						}
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAviatrixKubernetesClusterExists(resourceName, goaviatrix.KubernetesCluster{
						ClusterId: "test-cluster-id4",
						Credential: &goaviatrix.KubernetesCredential{
							UseCspCredentials: true,
						},
						Resource: &goaviatrix.ClusterResource{
							Name:        "test-name",
							Region:      "test-region",
							VpcId:       "test-vpc",
							AccountId:   "test-account-id",
							AccountName: "test-account",
							Platform:    "test-platform",
							Version:     "test-version",
							NetworkMode: "OVERLAY",
							Public:      true,
							Tags: []goaviatrix.Tag{
								{
									Key:   "key1",
									Value: "value1",
								},
							},
						},
					}),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "test-cluster-id4"),
					resource.TestCheckResourceAttr(resourceName, "use_csp_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "cluster_details.0.account_name", "test-account"),
				),
			},
		},
	})
}

func TestAccAviatrixKubernetesCluster_update(t *testing.T) {
	if os.Getenv("SKIP_KUBERNETES_CLUSTER") == "yes" {
		t.Skip("Skipping kubernetes cluster test as SKIP_KUBERNETES_CLUSTER is set")
	}

	resourceName := "aviatrix_kubernetes_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAviatrixKubernetesClusterDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "aviatrix_kubernetes_cluster" "test" {
						cluster_id = "test-cluster-id6"
						use_csp_credentials = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAviatrixKubernetesClusterExists(resourceName, goaviatrix.KubernetesCluster{
						ClusterId: "test-cluster-id6",
						Credential: &goaviatrix.KubernetesCredential{
							UseCspCredentials: true,
						},
					}),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "test-cluster-id6"),
					resource.TestCheckResourceAttr(resourceName, "use_csp_credentials", "true"),
				),
			},
			{
				Config: `
					resource "aviatrix_kubernetes_cluster" "test" {
						cluster_id = "test-cluster-id6"
						use_csp_credentials = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAviatrixKubernetesClusterExists(resourceName, goaviatrix.KubernetesCluster{
						ClusterId: "test-cluster-id6",
						Credential: &goaviatrix.KubernetesCredential{
							UseCspCredentials: false,
						},
					}),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "test-cluster-id6"),
					resource.TestCheckResourceAttr(resourceName, "use_csp_credentials", "false"),
				),
			},
		},
	})
}

func testAccCheckAviatrixKubernetesClusterDestroy(name string) func(state *terraform.State) error {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if resource.Primary == nil || len(resource.Primary.ID) == 0 {
			return fmt.Errorf("No ID is set")
		}
		client := testAccProvider.Meta().(*goaviatrix.Client)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cluster, err := client.GetKubernetesCluster(ctx, resource.Primary.ID)
		if err == nil {
			return fmt.Errorf("Expected an error")
		}
		if !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("The object is still found after being destroyed: %v", cluster)
		}
		return nil
	}
}

func testAccCheckAviatrixKubernetesClusterExists(name string, expected goaviatrix.KubernetesCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if resource.Primary == nil || len(resource.Primary.ID) == 0 {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		actual, err := client.GetKubernetesCluster(ctx, resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("could not find cluster with resource id %s: %w", resource.Primary.ID, err)
		}

		if diff := cmp.Diff(expected, *actual,
			// Ignore the ID field
			cmp.FilterPath(func(path cmp.Path) bool {
				return path.String() == "Id"
			}, cmp.Ignore())); len(diff) > 0 {
			return fmt.Errorf("difference found between expected and actual cluster: %s", diff)
		}
		return nil
	}
}
