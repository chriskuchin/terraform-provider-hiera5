package hiera5

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"hiera5": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"config": "test-fixtures/hiera.yaml",
		"scope":  map[string]interface{}{"service": "api", "environment": "live", "facts": "{timezone=>'CET'}"},
		"merge":  "deep",
	}

	err := rp.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
}
