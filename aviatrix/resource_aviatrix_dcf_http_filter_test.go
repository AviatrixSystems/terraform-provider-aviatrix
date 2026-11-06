package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestMarshalDCFHTTPFilterInput(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAviatrixDCFHTTPFilter().Schema, map[string]any{
		"name":        "write-methods",
		"description": "POST/PUT to web",
		"filters": []any{
			map[string]any{"all": []any{
				map[string]any{"header": ":method", "values": []any{"POST", "PUT"}},
				map[string]any{"header": "app", "values": []any{"web"}},
			}},
			map[string]any{"all": []any{
				map[string]any{"header": "x-present"},
			}},
		},
	})

	filter := marshalDCFHTTPFilterInput(d)
	require.NotNil(t, filter.Desc)
	assert.Equal(t, "POST/PUT to web", *filter.Desc)
	assert.Equal(t, "write-methods", filter.Name)
	assert.Equal(t, []goaviatrix.DCFHTTPFilterAndClause{
		{All: []goaviatrix.DCFHTTPFilterElement{
			{Hdr: ":method", Values: []string{"POST", "PUT"}},
			{Hdr: "app", Values: []string{"web"}},
		}},
		{All: []goaviatrix.DCFHTTPFilterElement{
			{Hdr: "x-present", Values: []string{}},
		}},
	}, filter.Filters)
}

func TestMarshalDCFHTTPFilterInput_OmitsEmptyDescription(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAviatrixDCFHTTPFilter().Schema, map[string]any{
		"name": "f",
		"filters": []any{
			map[string]any{"all": []any{map[string]any{"header": "app"}}},
		},
	})

	assert.Nil(t, marshalDCFHTTPFilterInput(d).Desc)
}

func TestFlattenDCFHTTPFilterClauses_RoundTrip(t *testing.T) {
	clauses := []goaviatrix.DCFHTTPFilterAndClause{
		{All: []goaviatrix.DCFHTTPFilterElement{
			{Hdr: ":method", Values: []string{"GET"}},
			{Hdr: "x-present"},
		}},
	}

	d := schema.TestResourceDataRaw(t, resourceAviatrixDCFHTTPFilter().Schema, map[string]any{})
	mustSet(d, "filters", flattenDCFHTTPFilterClauses(clauses))

	// Values absent from the API response read back as an empty list, the same
	// as an omitted values argument, so no diff is produced.
	assert.Equal(t, []goaviatrix.DCFHTTPFilterAndClause{
		{All: []goaviatrix.DCFHTTPFilterElement{
			{Hdr: ":method", Values: []string{"GET"}},
			{Hdr: "x-present", Values: []string{}},
		}},
	}, expandDCFHTTPFilterClauses(getList(d, "filters")))
}

func TestDCFHTTPFilterHeaderValidation(t *testing.T) {
	filtersElem, ok := resourceAviatrixDCFHTTPFilter().Schema["filters"].Elem.(*schema.Resource)
	require.True(t, ok)
	allElem, ok := filtersElem.Schema["all"].Elem.(*schema.Resource)
	require.True(t, ok)
	validate := allElem.Schema["header"].ValidateFunc

	for _, valid := range []string{":method", ":scheme", ":authority", ":path", "app", "x-forwarded-for", "user-agent", "x_custom.1~"} {
		_, errs := validate(valid, "header")
		assert.Empty(t, errs, "expected %q to be valid", valid)
	}
	for _, invalid := range []string{"", "Method", "X-Forwarded-For", ":host", ":status", "bad header", "a:b", "hdr/1"} {
		_, errs := validate(invalid, "header")
		assert.NotEmpty(t, errs, "expected %q to be invalid", invalid)
	}
}

func TestAccAviatrixDcfHTTPFilter_basic(t *testing.T) {
	if os.Getenv("SKIP_DCF_HTTP_FILTER") == "yes" {
		t.Skip("Skipping DCF HTTP Filter test as SKIP_DCF_HTTP_FILTER is set")
	}
	resourceName := "aviatrix_dcf_http_filter.test"

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
				Config: testAccDcfHTTPFilterBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfHTTPFilterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-http-filter"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "filters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.0.header", ":method"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.0.values.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.0.values.0", "POST"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.0.values.1", "PUT"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				Config: testAccDcfHTTPFilterUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfHTTPFilterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-http-filter-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated"),
					resource.TestCheckResourceAttr(resourceName, "filters.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.all.1.header", "app"),
					resource.TestCheckResourceAttr(resourceName, "filters.1.all.0.header", "x-present"),
					resource.TestCheckResourceAttr(resourceName, "filters.1.all.0.values.#", "0"),
				),
			},
			{
				// Removing the description must clear it on the controller.
				Config: testAccDcfHTTPFilterBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfHTTPFilterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-http-filter"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "filters.#", "1"),
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

func testAccDcfHTTPFilterBasic() string {
	return `
resource "aviatrix_dcf_http_filter" "test" {
	name = "test-dcf-http-filter"

	filters {
		all {
			header = ":method"
			values = ["POST", "PUT"]
		}
	}
}
`
}

func testAccDcfHTTPFilterUpdated() string {
	return `
resource "aviatrix_dcf_http_filter" "test" {
	name        = "test-dcf-http-filter-updated"
	description = "updated"

	filters {
		all {
			header = ":method"
			values = ["POST", "PUT"]
		}
		all {
			header = "app"
			values = ["web"]
		}
	}

	filters {
		all {
			header = "x-present"
		}
	}
}
`
}

func testAccCheckDcfHTTPFilterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF HTTP Filter resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF HTTP Filter ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())
		if _, err := client.GetDCFHTTPFilter(context.Background(), rs.Primary.ID); err != nil {
			return fmt.Errorf("failed to get DCF HTTP Filter: %w", err)
		}
		return nil
	}
}

func testAccCheckDcfHTTPFilterDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_http_filter" {
			continue
		}
		_, err := client.GetDCFHTTPFilter(ctx, rs.Primary.ID)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("DCF HTTP Filter %s still exists or could not be checked: %w", rs.Primary.ID, err)
		}
	}
	return nil
}
