package hiera5

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHiera5_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5" "sut" {
						key = "aws_instance_size"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5.sut", "value", "t2.large"),
					resource.TestCheckResourceAttrSet("data.hiera5.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5_Default_Found(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5" "sut" {
						key = "aws_instance_size"
						default = "t3.large"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5.sut", "value", "t2.large"),
					resource.TestCheckResourceAttrSet("data.hiera5.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5_Default_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5" "sut" {
						key = "gcp_instance_size"
						default = "t3.large"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5.sut", "value", "t3.large"),
					resource.TestCheckResourceAttrSet("data.hiera5.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5" "sut" {
						key = "gcp_instance_size"
					}`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestAccDataSourceHiera5_ScopeOverride(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5" "sut" {
						key = "aws_instance_size"
						scope = {
							"service" = "worker"
							"environment" = "live"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5.sut", "value", "t2.micro"),
					resource.TestCheckResourceAttrSet("data.hiera5.sut", "id"),
				),
			},
		},
	})
}
