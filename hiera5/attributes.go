package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	idAttribute = schema.StringAttribute{
		Computed: true,
	}

	keyAttribute = schema.StringAttribute{
		Required:    true,
		Description: "The key to lookup within the hiera data. Data Source will error if the key is not found and no default is provided",
	}

	scopeAttribute = schema.MapAttribute{
		ElementType: types.StringType,
		Description: "Map object defining the various hiera variables to determin how hiera merges files.",
		Optional:    true,
	}

	scopeOverrideAttribute = schema.MapAttribute{
		ElementType: types.StringType,
		Description: "Map object defining the various hiera variables to determin how hiera merges files. If present will override the provider scope setting for this datasource only.",
		Optional:    true,
	}

	valueDescription   = "The result of the lookup in the hiera data, or the default value if the key is not found."
	defaultDescription = "Default value to return if the value isn't found in the hiera data."
)

func processScopeOverrideAttribute(ctx context.Context, rawScope types.Map) (map[string]interface{}, []diag.Diagnostic) {
	var scopeOverride map[string]interface{}
	var diag []diag.Diagnostic
	if !rawScope.IsNull() {
		var scopeOverrideString map[string]string
		diag = rawScope.ElementsAs(ctx, &scopeOverrideString, true)

		scopeOverride = map[string]interface{}{}
		for k, v := range scopeOverrideString {
			scopeOverride[k] = v
		}
	}

	return scopeOverride, diag
}
