package helper

import (
	"log"
	"strings"

	"context"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/hiera/api"
	"github.com/lyraproj/hiera/hiera"
	"github.com/lyraproj/hiera/provider"
	sdk "github.com/lyraproj/hierasdk/hiera"

	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// Lookup is a wrapper for lyraproj's hiera/hiera.LookupAndRender
func Lookup(config string, strategy string, key string, valueType string, vars map[string]interface{}) ([]byte, error) {
	var args []string
	var out []byte
	var explain []byte
	var err error
	var b bytes.Buffer

	var cmdOpts hiera.CommandOptions

	cfgOpts := vf.MutableMap()
	cfgOpts.Put(
		provider.LookupKeyFunctions, []sdk.LookupKey{provider.ConfigLookupKey, provider.Environment})

	dir, err := os.Getwd()
	log.Printf("[DEBUG] PWD is %s", dir)
	log.Printf("[DEBUG] Config file is %s", config)
	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Printf("[DEBUG] ERROR '%s' reading config %s", err.Error(), config)
		return out, err
	}
	cmdOpts.Merge = strategy
	log.Printf("[DEBUG] Lookup strategy is %s", strategy)
	cfgOpts.Put(api.HieraConfig, config)

	//TODO: Implement type
	log.Printf("[DEBUG] Lookup value type is %s", valueType)
	//if valueType != "" {
	//	//cmdOpts.Type = valueType
	//	cmdOpts.Type = valueType
	//}

	cfgOpts.Put(api.HieraDialect, "pcore")
	for key, value := range vars {
		//log.Printf("[DEBUG] Var: %s=%s", key, value)
		cmdOpts.Variables = append(cmdOpts.Variables, strings.Join([]string{key, value.(string)}, "="))
	}
	log.Printf("[DEBUG] Lookup variables are %s", cmdOpts.Variables)
	cmdOpts.RenderAs = "json"

	args = append(args, key)
	log.Printf("[DEBUG] Lookup key is %s", key)
	hiera.DoWithParent(context.TODO(), provider.MuxLookupKey, cfgOpts, func(c api.Session) {
		hiera.LookupAndRender(c, &cmdOpts, args, &b)
	})
	if out, err = ioutil.ReadAll(io.Reader(&b)); err != nil {
		log.Printf("[DEBUG] ERROR %s", err.Error())
		return out, err
	}
	log.Printf("[DEBUG] out is %s", string(out))

	cmdOpts.RenderAs = "yaml"
	cmdOpts.ExplainOptions = true
	cmdOpts.ExplainData = true
	hiera.DoWithParent(context.TODO(), provider.MuxLookupKey, cfgOpts, func(c api.Session) {
		hiera.LookupAndRender(c, &cmdOpts, args, &b)
	})

	if explain, err = ioutil.ReadAll(io.Reader(&b)); err != nil {
		log.Printf("[DEBUG] ERROR %s", err.Error())
		return out, err
	}
	log.Printf("[DEBUG] explain is %s", string(explain))

	return out, nil
}
