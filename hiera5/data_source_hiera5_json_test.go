package hiera5

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceHiera5Json_Basic(t *testing.T) {
	key := "aws_tags"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHiera5JsonConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5JsonCheck(key),
				),
			},
			{
				Config: testAccDataSourceHiera5JsonConfig(keyUnavailable),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5JsonCheck(keyUnavailable),
				),
				ExpectError: regexp.MustCompile("key '" + keyUnavailable + "' not found"),
			},
			{
				Config: testAccDataSourceHiera5JsonConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5DefaultJSONCheck("default"),
				),
			},
		},
	})
}

func testAccDataSourceHiera5JsonCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5_json.%s", key)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes
		if attr["id"] != key {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				key,
			)
		}

		if attr["value"] != `{"team":"A","tier":1}` {
			return fmt.Errorf(
				"value is %s; want %s",
				attr["value"],
				`{"team":"A","tier":1}`,
			)
		}

		return nil
	}
}

func testAccDataSourceHiera5DefaultJSONCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5_json.%s", key)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes
		if attr["id"] != key {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				key,
			)
		}

		if attr["value"] != `{"team":"B","tier":"10"}` {
			return fmt.Errorf(
				"value is %s; want %s",
				attr["value"],
				`{"team":"B","tier":"10"}`,
			)
		}

		return nil
	}
}

func testAccDataSourceHiera5JsonConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera5" {
			config = "test-fixtures/hiera.yaml"
			scope = {
				environment = "live"
				service     = "api"
			}
		        merge = "deep"
		}

		data "hiera5_json" "%s" {
		  key = "%s"
		}

		data "hiera5_json" "default" {
			key = "default"
			default = "{\"team\":\"B\",\"tier\":\"10\"}"
		}
		`, key, key)
}
