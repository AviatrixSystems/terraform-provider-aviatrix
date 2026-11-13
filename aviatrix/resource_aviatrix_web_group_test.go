package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixWebGroup_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_WEB_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Web Group test as SKIP_WEB_GROUP is set")
	}
	resourceName := "aviatrix_web_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: testAccWebGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-web-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.urlfilter", "aviatrix.com/test"),
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

func testAccWebGroupBasic() string {
	return `
resource "aviatrix_web_group" "test" {
	name = "test-web-group"

	selector {
		match_expressions {
			urlfilter = "aviatrix.com/test"
		}
	}
}
`
}

func TestAccAviatrixWebGroup_external(t *testing.T) {
	skipAcc := os.Getenv("SKIP_WEB_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Web Group test as SKIP_WEB_GROUP is set")
	}
	resourceName := "aviatrix_web_group.external"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: testAccWebGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebGroupExternal("external_sni = \"malicious-domain-feed\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.external_sni", "malicious-domain-feed"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.ext_args.%", "0"),
				),
			},
			{
				Config: testAccWebGroupExternal(`external_sni = "malicious-domain-feed"
			ext_args = {
				filterA = "valueA"
			}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.external_sni", "malicious-domain-feed"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.ext_args.filterA", "valueA"),
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

func testAccWebGroupExternal(matchExpression string) string {
	return fmt.Sprintf(`
resource "aviatrix_web_group" "external" {
	name = "test-external-web-group"

	selector {
		match_expressions {
			%s
		}
	}
}
`, matchExpression)
}

func TestMarshalWebGroupInput_external(t *testing.T) {
	data := schema.TestResourceDataRaw(t, resourceAviatrixWebGroup().Schema, map[string]any{
		"name": "external",
		"selector": []any{
			map[string]any{
				"match_expressions": []any{
					map[string]any{
						goaviatrix.ExternalSNIKey: "malicious-domain-feed",
						goaviatrix.ExtArgsPrefix:  map[string]any{"filterA": "valueA"},
					},
					map[string]any{
						goaviatrix.ExternalSNIKey: "malicious-domain-feed",
					},
				},
			},
		},
	})

	webGroup, err := marshalWebGroupInput(data)
	require.NoError(t, err)
	assert.Equal(t, []*goaviatrix.WebGroupMatchExpression{
		{ExternalSNI: "malicious-domain-feed", ExtArgs: map[string]string{"filterA": "valueA"}},
		{ExternalSNI: "malicious-domain-feed"},
	}, webGroup.Selector.Expressions)
}

func TestMarshalWebGroupInput_rejectsInvalidMatchExpressions(t *testing.T) {
	tests := []struct {
		name            string
		matchExpression map[string]any
		wantErr         string
	}{
		{
			name:            "snifilter and urlfilter",
			matchExpression: map[string]any{"snifilter": "aviatrix.com", "urlfilter": "aviatrix.com/test"},
			wantErr:         "only one of snifilter, urlfilter, external_url and external_sni",
		},
		{
			name:            "external_sni and snifilter",
			matchExpression: map[string]any{"external_sni": "malicious-domain-feed", "snifilter": "aviatrix.com"},
			wantErr:         "only one of snifilter, urlfilter, external_url and external_sni",
		},
		{
			name:            "ext_args without a feed",
			matchExpression: map[string]any{"snifilter": "aviatrix.com", "ext_args": map[string]any{"filterA": "valueA"}},
			wantErr:         "ext_args can only be set together with external_url or external_sni",
		},
		{
			name:            "ext_args with a filter key",
			matchExpression: map[string]any{"external_sni": "malicious-domain-feed", "ext_args": map[string]any{"snifilter": "aviatrix.com"}},
			wantErr:         `ext_args cannot contain the key "snifilter"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, resourceAviatrixWebGroup().Schema, map[string]any{
				"name": "invalid",
				"selector": []any{
					map[string]any{"match_expressions": []any{tt.matchExpression}},
				},
			})

			_, err := marshalWebGroupInput(data)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

// Setting the flattened selector must fit the schema and clear fields the API no longer returns.
func TestFlattenWebGroupSelector(t *testing.T) {
	data := schema.TestResourceDataRaw(t, resourceAviatrixWebGroup().Schema, map[string]any{
		"name": "external",
		"selector": []any{
			map[string]any{
				"match_expressions": []any{
					map[string]any{"snifilter": "aviatrix.com"},
				},
			},
		},
	})

	webGroup := &goaviatrix.WebGroup{
		Selector: goaviatrix.WebGroupSelector{
			Expressions: []*goaviatrix.WebGroupMatchExpression{
				{ExternalSNI: "phishing-domain-feed", ExtArgs: map[string]string{"filterA": "valueA"}},
			},
		},
	}
	require.NoError(t, data.Set("selector", flattenWebGroupSelector(webGroup)))

	assert.Equal(t, "", data.Get("selector.0.match_expressions.0.snifilter"))
	assert.Equal(t, "phishing-domain-feed", data.Get("selector.0.match_expressions.0.external_sni"))
	assert.Equal(t, map[string]any{"filterA": "valueA"}, data.Get("selector.0.match_expressions.0.ext_args"))
}

func testAccCheckWebGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Web Group resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Web Group ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		webGroup, err := client.GetWebGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get Web Group status: %w", err)
		}

		if webGroup.UUID != rs.Primary.ID {
			return fmt.Errorf("web Group ID not found")
		}

		return nil
	}
}

func testAccWebGroupDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_web_group" {
			continue
		}

		_, err := client.GetWebGroup(context.Background(), rs.Primary.ID)
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("web group configured when it should be destroyed")
		}
	}

	return nil
}
