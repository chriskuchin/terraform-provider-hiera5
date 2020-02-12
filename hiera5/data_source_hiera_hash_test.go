package hiera5

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceHiera5Hash_Basic(t *testing.T) {
	key := "aws_tags"
	keyUnavailable := "doesnt_exists"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHiera5HashConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5HashCheck(key),
				),
			},
			{
				Config: testAccDataSourceHiera5HashConfig(keyUnavailable),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHiera5HashCheck(keyUnavailable),
				),
				ExpectError: regexp.MustCompile("Key '" + keyUnavailable + "' not found"),
			},
		},
	})
}

func testAccDataSourceHiera5HashCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera5_hash.%s", key)

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
		if attr["value.tier"] != "1" {
			return fmt.Errorf(
				"value.tier is %s; want %s",
				attr["value.tier"],
				"1",
			)
		}
		if attr["value.team"] != "A" {
			return fmt.Errorf(
				"value.team is %s; want %s",
				attr["value.team"],
				"A",
			)
		}
		return nil
	}
}

func testAccDataSourceHiera5HashConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera5" {
			config = "test-fixtures/hiera.yaml"
			scope = {
				environment = "live"
				service     = "api"
			}
		        merge = "deep"
		}
		
		data "hiera5_hash" "%s" {
		  key = "%s"
		}
		`, key, key)
}
