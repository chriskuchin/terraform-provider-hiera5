package hiera5

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHiera5Hash_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_hash" "sut" {
						key = "aws_tags"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.%", "2"),
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.tier", "1"),
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.team", "A"),
					resource.TestCheckResourceAttrSet("data.hiera5_hash.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Hash_Default_Found(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_hash" "sut" {
						key = "aws_tags"
						default = {
							"service" = "unknown"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.%", "2"),
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.tier", "1"),
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.team", "A"),
					resource.TestCheckResourceAttrSet("data.hiera5_hash.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Hash_Default_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_hash" "sut" {
						key = "gcp_tags"
						default = {
							"service" = "unknown"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.%", "1"),
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.service", "unknown"),
					resource.TestCheckResourceAttrSet("data.hiera5_hash.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5Hash_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_hash" "sut" {
						key = "gcp_tags"
					}`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestAccDataSourceHiera5Hash_ScopeOverride(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_hash" "sut" {
						key = "aws_tags"
						scope = {
							"service" = "api"
							"environment" = "stage"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.%", "1"),
					resource.TestCheckResourceAttr("data.hiera5_hash.sut", "value.team", "A"),
					resource.TestCheckResourceAttrSet("data.hiera5_hash.sut", "id"),
				),
			},
		},
	})
}
