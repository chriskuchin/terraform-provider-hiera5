package hiera5

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceHiera5Array_Basic(t *testing.T) {
	key := "java_opts"
	keyUnavailable := "doesnt_exists"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHiera5ArrayConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5ArrayCheck(key),
				),
			},
			{
				Config: testAccDataSourceHiera5ArrayConfig(keyUnavailable),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5ArrayCheck(keyUnavailable),
				),
				ExpectError: regexp.MustCompile("unexpected end of JSON input"),
			},
		},
	})
}

func testAccDataSourceHiera5ArrayCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5_array.%s", key)

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
		if attr["value.0"] != "-Xms512m" {
			return fmt.Errorf(
				"value.0 is %s; want %s",
				attr["value.0"],
				"-Xms512m",
			)
		}
		if attr["value.1"] != "-Xmx2g" {
			return fmt.Errorf(
				"value.1 is %s; want %s",
				attr["value.1"],
				"-Xmx2g",
			)
		}
		if attr["value.2"] != "-Dspring.profiles.active=live" {
			return fmt.Errorf(
				"value.2 is %s; want %s",
				attr["value.2"],
				"-Dspring.profiles.active=live",
			)
		}
		return nil
	}
}

func testAccDataSourceHiera5ArrayConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera5" {
			config = "test-fixtures/hiera.yaml"
			scope = {
				environment = "live"
				service     = "api"
			}
		        merge = "deep"
		}
		
		data "hiera5_array" "%s" {
		  key = "%s"
		}
		`, key, key)
}
