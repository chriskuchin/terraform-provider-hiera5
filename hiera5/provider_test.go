package hiera5

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	providerConfig = `
provider "hiera5" {
	config = "test-fixtures/hiera.yaml"
	scope = {
		"service" = "api"
		"environment" = "live"
		"facts" = "{'timezone'=>'CET'}"
	}
	merge = "deep"
}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hiera5": providerserver.NewProtocol6WithError(New()),
}

func TestAccProvider_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig,
			},
		},
	})
}
