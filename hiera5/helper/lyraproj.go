package helper

import (
	"fmt"
	"strings"

	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/hiera/api"
	"github.com/lyraproj/hiera/hiera"
	"github.com/lyraproj/hiera/provider"
	sdk "github.com/lyraproj/hierasdk/hiera"

	"bytes"
	"io"
	"os"
)

// Lookup is a wrapper for lyraproj's hiera/hiera.LookupAndRender
// it returns either an empty string when key is not found or JSON encoded key's value
func Lookup(ctx context.Context, config string, strategy string, key string, valueType string, vars map[string]interface{}) ([]byte, error) {
	var (
		args    []string
		out     []byte
		b       bytes.Buffer
		cmdOpts hiera.CommandOptions
	)

	cfgOpts := vf.MutableMap()
	cfgOpts.Put(
		provider.LookupKeyFunctions, []sdk.LookupKey{provider.ConfigLookupKey, provider.Environment})

	tflog.Debug(ctx, fmt.Sprintf("Config file is %s", config))

	if _, err := os.Stat(config); os.IsNotExist(err) {
		tflog.Debug(ctx, fmt.Sprintf("[DEBUG] ERROR '%s' reading config %s", err.Error(), config))
		return out, err
	}

	cmdOpts.Merge = strategy
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Lookup strategy is %s", strategy))
	cfgOpts.Put(api.HieraConfig, config)

	//TODO: Implement type
	//if valueType != "" {
	//	cmdOpts.Type = valueType
	//}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Lookup value type is %s", valueType))

	cfgOpts.Put(api.HieraDialect, "pcore")

	for key, value := range vars {
		//tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Var: %s=%s", key, value)
		cmdOpts.Variables = append(cmdOpts.Variables, strings.Join([]string{key, value.(string)}, "="))
	}

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Lookup variables are %s", cmdOpts.Variables))

	cmdOpts.RenderAs = "json"

	args = append(args, key)
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Lookup key is %s", key))
	hiera.DoWithParent(context.TODO(), provider.MuxLookupKey, cfgOpts, func(c api.Session) {
		hiera.LookupAndRender(c, &cmdOpts, args, &b)
	})

	out, _ = io.ReadAll(io.Reader(&b))

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] out is %s", string(out)))

	//cmdOpts.RenderAs = "yaml"
	//cmdOpts.ExplainOptions = true
	//cmdOpts.ExplainData = true
	//hiera.DoWithParent(context.TODO(), provider.MuxLookupKey, cfgOpts, func(c api.Session) {
	//	hiera.LookupAndRender(c, &cmdOpts, args, &b)
	//})
	//explain, _ = ioutil.ReadAll(io.Reader(&b))
	//tflog.Debug(ctx, fmt.Sprintf("[DEBUG] explain is %s", string(explain))

	return out, nil
}
