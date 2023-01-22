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
	config = "test-fixtures/hiera.yml"
	scope = {
		"service" = "api"
		"environment" = "live"
		"facts" = "{timezone=>'CET'}"
	}
	merge = "deep"
}`
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

// func getDefaultSUTProvider() provider.Provider {
// 	&Hiera5ProviderModel{
// 		Config: "test-fixtures/hiera.yaml",
// 		Scope:  map[string]interface{}{"service": "api", "environment": "live", "facts": "{timezone=>'CET'}"},
// 		Merge:  "deep",
// 	}

// 	return nil
// }

// var testAccProviders map[string]*schema.Provider
// var testAccProvider *schema.Provider

// func init() {
// 	testAccProvider = Provider()
// 	testAccProviders = map[string]*schema.Provider{
// 		"hiera5": testAccProvider,
// 	}
// }

// func TestProvider(t *testing.T) {
// 	if err := Provider().InternalValidate(); err != nil {
// 		t.Fatalf("err: %v", err)
// 	}
// }

// func TestProviderConfigure(t *testing.T) {
// 	rp := Provider()

// 	raw := map[string]interface{}{
// 		"config": "test-fixtures/hiera.yaml",
// 		"scope":  map[string]interface{}{"service": "api", "environment": "live", "facts": "{timezone=>'CET'}"},
// 		"merge":  "deep",
// 	}

// 	err := rp.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
// 	if err != nil {
// 		t.Fatalf("err: %v", err)
// 	}
// }

// func TestProviderImpl(t *testing.T) {
// 	var _ *schema.Provider = Provider()
// }

// func testAccPreCheck(t *testing.T) {
// }
