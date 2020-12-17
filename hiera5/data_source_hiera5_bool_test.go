package hiera5

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceHiera5Bool_Basic(t *testing.T) {
	key := "enable_spot_instances"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHiera5BoolConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5BoolCheck(key),
				),
			},
			{
				Config: testAccDataSourceHiera5BoolConfig(keyUnavailable),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5BoolCheck(keyUnavailable),
				),
				ExpectError: regexp.MustCompile("key '" + keyUnavailable + "' not found"),
			},
			{
				Config: testAccDataSourceHiera5BoolConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5BoolDefaultValueCheck("default"),
				),
			},
		},
	})
}

func testAccDataSourceHiera5BoolCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5_bool.%s", key)

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

		if attr["value"] != "true" {
			return fmt.Errorf(
				"value is %s; want %s",
				attr["value"],
				"true",
			)
		}

		return nil
	}
}

func testAccDataSourceHiera5BoolDefaultValueCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5_bool.%s", key)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes
		if attr["id"] != fmt.Sprintf("%s", key) {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				fmt.Sprintf("%s", key),
			)
		}

		if attr["value"] != "false" {
			return fmt.Errorf(
				"value is %s; want %s",
				attr["value"],
				"false",
			)
		}

		return nil
	}
}

func testAccDataSourceHiera5BoolConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera5" {
			config = "test-fixtures/hiera.yaml"
			scope = {
				environment = "live"
				service     = "api"
			}
		        merge = "deep"
		}

		data "hiera5_bool" "%s" {
		  key = "%s"
		}

		data "hiera5_bool" "default" {
			key = "default"
			default = false
		}
		`, key, key)
}
