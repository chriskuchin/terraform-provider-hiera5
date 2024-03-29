package hiera5

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cast"

	"github.com/chriskuchin/terraform-provider-hiera5/hiera5/helper"
)

type override func(h *hiera5) *hiera5

type hiera5 struct {
	Config string
	Scope  map[string]interface{}
	Merge  string
}

func WithScopeOverride(scope map[string]interface{}) override {
	return func(h *hiera5) *hiera5 {
		if scope == nil {
			return h
		}

		return &hiera5{
			Config: h.Config,
			Scope:  scope,
			Merge:  h.Merge,
		}
	}
}

func handleOverrides(h *hiera5, opts ...override) *hiera5 {
	override := h
	for _, opt := range opts {
		override = opt(override)
	}

	return override
}

func newHiera5(config string, scope map[string]interface{}, merge string) hiera5 {
	return hiera5{
		Config: config,
		Scope:  scope,
		Merge:  merge,
	}
}

func (h *hiera5) lookup(ctx context.Context, key string, valueType string) ([]byte, error) {
	out, err := helper.Lookup(ctx, h.Config, h.Merge, key, valueType, h.Scope)
	if err == nil && string(out) == "" {
		return out, fmt.Errorf("key '%s' not found", key)
	}

	if !json.Valid(out) {
		return out, fmt.Errorf("key '%s''s lookup returned invalid JSON: '%s'", key, out)
	}

	return out, err
}

func (h *hiera5) array(ctx context.Context, key string, opts ...override) ([]interface{}, error) {
	var (
		f interface{}
		e []interface{}
	)

	out, err := handleOverrides(h, opts...).lookup(ctx, key, "Array")
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(out, &f)

	if _, ok := f.([]interface{}); ok {
		for _, v := range f.([]interface{}) {
			e = append(e, cast.ToString(v))
		}
	} else {
		return nil, fmt.Errorf("key '%s' does not return a valid array", key)
	}

	return e, nil
}

func (h *hiera5) hash(ctx context.Context, key string, opts ...override) (map[string]interface{}, error) {
	var f interface{}

	e := make(map[string]interface{})

	out, err := handleOverrides(h, opts...).lookup(ctx, key, "Hash")
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(out, &f)

	if _, ok := f.(map[string]interface{}); ok {
		for k, v := range f.(map[string]interface{}) {
			e[k] = cast.ToString(v)
		}
	} else {
		return nil, fmt.Errorf("key '%s' does not return a valid hash", key)
	}

	return e, nil
}

func (h *hiera5) value(ctx context.Context, key string, opts ...override) (string, error) {
	var f interface{}

	out, err := handleOverrides(h, opts...).lookup(ctx, key, "")
	if err != nil {
		return "", err
	}

	_ = json.Unmarshal(out, &f)

	return cast.ToString(f), nil
}

func (h *hiera5) bool(ctx context.Context, key string, opts ...override) (bool, error) {
	var f interface{}

	out, err := handleOverrides(h, opts...).lookup(ctx, key, "")
	if err != nil {
		return false, err
	}

	_ = json.Unmarshal(out, &f)

	return cast.ToBool(f), nil
}

func (h *hiera5) json(ctx context.Context, key string, opts ...override) (string, error) {
	var b bytes.Buffer

	out, err := handleOverrides(h, opts...).lookup(ctx, key, "")
	if err != nil {
		return "", err
	}

	_ = json.Compact(&b, out)
	out, _ = io.ReadAll(io.Reader(&b))

	return string(out), nil
}
