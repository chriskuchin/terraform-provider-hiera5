package hiera5

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHiera5Array_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_array" "sut" {
						key = "java_opts"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.#", "3"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.0", "-Xms512m"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.1", "-Xmx2g"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.2", "-Dspring.profiles.active=live"),
					resource.TestCheckResourceAttrSet("data.hiera5_array.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Array_Default_Found(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_array" "sut" {
						key = "java_opts"
						default = ["value1", "value2"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.#", "3"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.0", "-Xms512m"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.1", "-Xmx2g"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.2", "-Dspring.profiles.active=live"),
					resource.TestCheckResourceAttrSet("data.hiera5_array.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Array_Default_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_array" "sut" {
						key = "missing_key"
						default = ["value1", "value2"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.#", "2"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.0", "value1"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.1", "value2"),
					resource.TestCheckResourceAttrSet("data.hiera5_array.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Array_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_array" "sut" {
						key = "missing_key"
					}`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestAccDataSourceHiera5Array_ScopeOverride(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_array" "sut" {
						key = "java_opts"
						scope = {
							"service" = "api"
							"environment" = "stage"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.#", "2"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.0", "-Xms512m"),
					resource.TestCheckResourceAttr("data.hiera5_array.sut", "value.1", "-Xmx2g"),
					resource.TestCheckResourceAttrSet("data.hiera5_array.sut", "id"),
				),
			},
		},
	})
}
