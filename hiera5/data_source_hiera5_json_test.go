package hiera5

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHiera5JSON_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_json" "sut" {
						key = "aws_tags"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_json.sut", "value", `{"team":"A","tier":1}`),
					resource.TestCheckResourceAttrSet("data.hiera5_json.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5JSON_Default_Found(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_json" "sut" {
						key = "aws_tags"
						default = "{\"team\":\"B\",\"tier\":\"10\"}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_json.sut", "value", `{"team":"A","tier":1}`),
					resource.TestCheckResourceAttrSet("data.hiera5_json.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5JSON_Default_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_json" "sut" {
						key = "gcp_tags"
						default = "{\"team\":\"B\",\"tier\":\"10\"}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_json.sut", "value", `{"team":"B","tier":"10"}`),
					resource.TestCheckResourceAttrSet("data.hiera5_json.sut", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceHiera5JSON_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_json" "sut" {
						key = "gcp_tags"
					}`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestAccDataSourceHiera5JSON_ScopeOverride(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hiera5_json" "sut" {
						key = "aws_tags"
						scope = {
							"service" = "api"
							"environment" = "stage"
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hiera5_json.sut", "value", `{"team":"A"}`),
					resource.TestCheckResourceAttrSet("data.hiera5_json.sut", "id"),
				),
			},
		},
	})
}
