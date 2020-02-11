package hiera5

import (
	"encoding/json"
	"fmt"
	//"log"

	"github.com/spf13/cast"

	"gitlab.com/sbitio/terraform-provider-hiera5/hiera5/helper"
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
	return helper.Lookup(h.Config, h.Merge, key, valueType, h.Scope)
}

func (h *hiera5) array(key string) ([]interface{}, error) {
	var f interface{}
	var e []interface{}

	out, err := h.lookup(key, "Array")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		return nil, err
	}
	if _, ok := f.([]interface{}); ok {
		for _, v := range f.([]interface{}) {
			e = append(e, cast.ToString(v))
		}
	} else {
		return nil, fmt.Errorf("Key '%s' does not return a valid array", key)
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

	err = json.Unmarshal(out, &f)
	if err != nil {
		return nil, err
	}

	if _, ok := f.(map[string]interface{}); ok {
		for k, v := range f.(map[string]interface{}) {
			e[k] = cast.ToString(v)
		}
	} else {
		return nil, fmt.Errorf("Key '%s' does not return a valid hash", key)
	}
	return e, nil
}

func (h *hiera5) value(key string) (string, error) {
	var f interface{}

	out, err := h.lookup(key, "")
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		return "", err
	}

	return cast.ToString(f), nil
}
