package hiera5

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cast"

	"github.com/chriskuchin/terraform-provider-hiera5/hiera5/helper"
)

type hiera5 struct {
	Config string
	Scope  map[string]interface{}
	Merge  string
}

func newHiera5(config string, scope map[string]interface{}, merge string) hiera5 {
	return hiera5{
		Config: config,
		Scope:  scope,
		Merge:  merge,
	}
}

func (h *hiera5) lookup(key string, valueType string) ([]byte, error) {
	out, err := helper.Lookup(h.Config, h.Merge, key, valueType, h.Scope)
	if err == nil && string(out) == "" {
		return out, fmt.Errorf("key '%s' not found", key)
	}

	if !json.Valid(out) {
		return out, fmt.Errorf("key '%s''s lookup returned invalid JSON: '%s'", key, out)
	}

	return out, err
}

func (h *hiera5) array(key string) ([]interface{}, error) {
	var (
		f interface{}
		e []interface{}
	)

	out, err := h.lookup(key, "Array")
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

func (h *hiera5) hash(key string) (map[string]interface{}, error) {
	var f interface{}

	e := make(map[string]interface{})

	out, err := h.lookup(key, "Hash")
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

func (h *hiera5) value(key string) (string, error) {
	var f interface{}

	out, err := h.lookup(key, "")
	if err != nil {
		return "", err
	}

	_ = json.Unmarshal(out, &f)

	return cast.ToString(f), nil
}

func (h hiera5) bool(key string) (bool, error) {
	var f interface{}

	out, err := h.lookup(key, "")
	if err != nil {
		return false, err
	}

	_ = json.Unmarshal(out, &f)

	return cast.ToBool(f), nil
}

func (h *hiera5) json(key string) (string, error) {
	var b bytes.Buffer

	out, err := h.lookup(key, "")
	if err != nil {
		return "", err
	}

	_ = json.Compact(&b, out)
	out, _ = io.ReadAll(io.Reader(&b))

	return string(out), nil
}
