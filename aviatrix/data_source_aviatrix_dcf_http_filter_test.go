package aviatrix

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataSourceAviatrixDcfHTTPFilter_basic(t *testing.T) {
	if os.Getenv("SKIP_DATA_DCF_HTTP_FILTER") == "yes" {
		t.Skip("Skipping Data Source DCF HTTP Filter test as SKIP_DATA_DCF_HTTP_FILTER is set")
	}
	dataSourceName := "data.aviatrix_dcf_http_filter.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: testAccCheckDcfHTTPFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcfHTTPFilterBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", "aviatrix_dcf_http_filter.test", "uuid"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "data source test"),
					resource.TestCheckResourceAttr(dataSourceName, "filters.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "filters.0.all.0.header", "app"),
					resource.TestCheckResourceAttr(dataSourceName, "filters.0.all.0.values.0", "web"),
				),
			},
		},
	})
}

func TestAccDataSourceAviatrixDcfHTTPFilter_notFound(t *testing.T) {
	if os.Getenv("SKIP_DATA_DCF_HTTP_FILTER") == "yes" {
		t.Skip("Skipping Data Source DCF HTTP Filter test as SKIP_DATA_DCF_HTTP_FILTER is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		Steps: []resource.TestStep{
			{
				Config:      `data "aviatrix_dcf_http_filter" "test" { name = "non-existent-http-filter-12345" }`,
				ExpectError: regexp.MustCompile(`DCF HTTP filter with name .* not found`),
			},
		},
	})
}

func testAccDataSourceDcfHTTPFilterBasic() string {
	return `
resource "aviatrix_dcf_http_filter" "test" {
	name        = "test-dcf-http-filter-data-source"
	description = "data source test"

	filters {
		all {
			header = "app"
			values = ["web"]
		}
	}
}

data "aviatrix_dcf_http_filter" "test" {
	name = aviatrix_dcf_http_filter.test.name
}
`
}
