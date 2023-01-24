package hiera5

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHiera5Bool_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_bool" "sut" {
						key = "enable_spot_instances"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_bool.sut", "value", "true"),
					resource.TestCheckResourceAttrSet("data.hiera5_bool.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Bool_Default_Found(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_bool" "sut" {
						key = "enable_spot_instances"
						default = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_bool.sut", "value", "true"),
					resource.TestCheckResourceAttrSet("data.hiera5_bool.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Bool_Default_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_bool" "sut" {
						key = "disable_spot_instances"
						default = false
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_bool.sut", "value", "false"),
					resource.TestCheckResourceAttrSet("data.hiera5_bool.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Bool_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_bool" "sut" {
						key = "disable_spot_instances"
					}`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestAccDataSourceHiera5Bool_ScopeOverride(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_bool" "sut" {
						key = "enable_spot_instances"
						scope = {
							"service" = "worker"
							"environment" = "stage"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_bool.sut", "value", "false"),
					resource.TestCheckResourceAttrSet("data.hiera5_bool.sut", "id"),
				),
			},
		},
	})
}
