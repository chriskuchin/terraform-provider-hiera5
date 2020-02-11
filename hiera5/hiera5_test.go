package hiera5

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cast"
)

func TestHiera5Lookup(t *testing.T) {
	var f interface{}

	hiera := testHiera5Config()

	out, err := hiera.lookup("aws_cloudwatch_enable", "")
	if err != nil {
		t.Errorf("Error running hiera: %s", err)
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %s", err)

	}

	v := cast.ToString(f)

	if v != "true" {
		t.Errorf("aws_cloudwatch_enable is %s; want %s", v, "true")
	}
}

func TestHiera5Hash(t *testing.T) {
	hiera := testHiera5Config()

	v, err := hiera.hash("aws_tags")
	if err != nil {
		t.Errorf("Error running hiera.Hash: %s", err)
	}

	if v["team"] != "A" {
		t.Errorf("aws_tags.team is %s; want %s", v, "A")
	}

	if v["tier"] != "1" {
		t.Errorf("aws_tags.tier is %s; want %s", v, "1")
	}
}

func TestHiera5Value(t *testing.T) {
	hiera := testHiera5Config()

	v, err := hiera.value("aws_cloudwatch_enable")
	if err != nil {
		t.Errorf("Error running hiera.Value: %s", err)
	}

	if v != "true" {
		t.Errorf("aws_cloudwatch_enable is %s; want %s", v, "true")
	}
}

func testHiera5Config() hiera5 {
	return newHiera5(
		"test-fixtures/hiera.yaml",
		map[string]interface{}{"service": "api", "environment": "live", "facts": "{timezone=>'CET'}"},
		"deep",
	)
}
