// Package hiera provides the lookup function types and the ProviderContext
package hiera

import (
	"net/url"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/vf"
)

type (
	// ProviderContext provides utility functions to a provider function
	ProviderContext interface {
		// Option returns the given option or nil if no such option exists
		Option(option string) dgo.Value

		// StringOption returns the option for the given name as a string and true provided that the option is present
		// and is a string. If its missing, or if its found to be something other than a string, this
		// method returns the empty string, false
		StringOption(option string) (string, bool)

		// BoolOption returns the option for the given name as a bool and true provided that the option is present
		// and is a bool. If its missing, or if its found to be something other than a bool, this
		// method returns false, false
		BoolOption(option string) (bool, bool)

		// IntOption returns the option for the given name as an int and true provided that the option is present
		// and is an int. If its missing, or if its found to be something other than an int, this method returns 0, false
		IntOption(option string) (int, bool)

		// FloatOption returns the option for the given name as a float64 and true provided that the option is present
		// and is an float64. If its missing, or if its found to be something other than a float64, this method
		// returns 0.0, false
		FloatOption(option string) (float64, bool)

		// OptionMap returns all options as an immutable Map
		OptionsMap() dgo.Map

		// ToData converts the given value into Data
		ToData(value interface{}) dgo.Value
	}

	providerContext struct {
		options dgo.Map
	}
)

// NewProviderContext creates a context containing the values of the the "options" key in the given url.Values.
func NewProviderContext(q url.Values) ProviderContext {
	var opts dgo.Map
	if jo := q.Get(`options`); jo != `` {
		v := streamer.UnmarshalJSON([]byte(jo), nil)
		if om, ok := v.(dgo.Map); ok {
			opts = om
		}
	}
	return &providerContext{options: opts}
}

// ProviderContextFromMap returns a ProviderContext that contains a frozen version of the given map
func ProviderContextFromMap(m dgo.Map) ProviderContext {
	if m == nil {
		m = vf.Map()
	} else {
		m = m.FrozenCopy().(dgo.Map)
	}
	return &providerContext{options: m}
}

func (c *providerContext) Option(name string) (d dgo.Value) {
	if c.options != nil {
		d = c.options.Get(name)
	}
	return
}

func (c *providerContext) StringOption(name string) (s string, ok bool) {
	var o dgo.String
	if o, ok = c.Option(name).(dgo.String); ok {
		s = o.GoString()
	}
	return
}

func (c *providerContext) IntOption(name string) (i int, ok bool) {
	var o dgo.Integer
	if o, ok = c.Option(name).(dgo.Integer); ok {
		i = int(o.GoInt())
	}
	return
}

func (c *providerContext) FloatOption(name string) (f float64, ok bool) {
	var o dgo.Float
	if o, ok = c.Option(name).(dgo.Float); ok {
		f = o.GoFloat()
	}
	return
}

func (c *providerContext) BoolOption(name string) (b bool, ok bool) {
	var o dgo.Boolean
	if o, ok = c.Option(name).(dgo.Boolean); ok {
		b = o.GoBool()
	}
	return
}

func (c *providerContext) OptionsMap() dgo.Map {
	return c.options
}

func (c *providerContext) ToData(value interface{}) dgo.Value {
	return vf.Value(value)
}
