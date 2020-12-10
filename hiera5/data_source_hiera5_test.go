package hiera5

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceHiera5_Basic(t *testing.T) {
	key := "aws_instance_size"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHiera5Config(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5Check(key),
				),
			},
			{
				Config: testAccDataSourceHiera5Config(keyUnavailable),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5Check(keyUnavailable),
				),
				ExpectError: regexp.MustCompile("key '" + keyUnavailable + "' not found"),
			},
			{
				Config: testAccDataSourceHiera5Config(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5DefaultValueCheck("default"),
				),
			},
		},
	})
}

func testAccDataSourceHiera5Check(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5.%s", key)

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

		if attr["value"] != "t2.large" {
			return fmt.Errorf(
				"value is %s; want %s",
				attr["value"],
				"t2.large",
			)
		}

		return nil
	}
}

func testAccDataSourceHiera5DefaultValueCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5.%s", key)

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

		if attr["value"] != "default_value" {
			return fmt.Errorf(
				"value is %s; want %s",
				attr["value"],
				"default_value",
			)
		}

		return nil
	}
}

func testAccDataSourceHiera5Config(key string) string {
	return fmt.Sprintf(`
		provider "hiera5" {
			config = "test-fixtures/hiera.yaml"
			scope = {
				environment = "live"
				service     = "api"
			}
		        merge = "deep"
		}

		data "hiera5" "%s" {
		  key = "%s"
		}

		data "hiera5" "default" {
			key = "default"
			default = "default_value"
		}
		`, key, key)
}
