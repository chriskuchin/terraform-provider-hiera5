package hiera5

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/spf13/cast"
)

const keyUnavailable = "doesnt_exists"

func TestHiera5Lookup(t *testing.T) {
	var f interface{}

	hiera := testHiera5Config()

	out, err := hiera.lookup(context.TODO(), "aws_cloudwatch_enable", "")
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

func TestHiera5Array(t *testing.T) {
	hiera := testHiera5Config()

	v, err := hiera.array(context.TODO(), "java_opts")
	if err != nil {
		t.Errorf("Error running hiera.Array: %s", err)
	}

	if v[0] != "-Xms512m" {
		t.Errorf(
			"v[0] is %s; want %s",
			v[0],
			"-Xms512m",
		)
	}

	if v[1] != "-Xmx2g" {
		t.Errorf(
			"v[1] is %s; want %s",
			v[1],
			"-Xmx2g",
		)
	}

	if v[2] != "-Dspring.profiles.active=live" {
		t.Errorf(
			"v[2] is %s; want %s",
			v[2],
			"-Dspring.profiles.active=live",
		)
	}

	v2, err2 := hiera.array(context.TODO(), keyUnavailable)
	if err2 == nil || v2 != nil {
		t.Errorf("Error running hiera.Array: %s", v2)
	}

	v3, err3 := hiera.array(context.TODO(), "aws_tags")
	if err3 == nil || v3 != nil {
		t.Errorf("Error running hiera.Array: %s", v3)
	}

	hieraBad := testHiera5ConfigBad()

	v4, err4 := hieraBad.array(context.TODO(), "java_opts")
	if err4 == nil || v4 != nil {
		t.Errorf("Error running hiera.Array: %s", v4)
	}
}

func TestHiera5Hash(t *testing.T) {
	hiera := testHiera5Config()

	v, err := hiera.hash(context.TODO(), "aws_tags")
	if err != nil {
		t.Errorf("Error running hiera.Hash: %s", err)
	}

	if v["team"] != "A" {
		t.Errorf("aws_tags.team is %s; want %s", v, "A")
	}

	if v["tier"] != "1" {
		t.Errorf("aws_tags.tier is %s; want %s", v, "1")
	}

	v2, err2 := hiera.hash(context.TODO(), keyUnavailable)
	if err2 == nil || v2 != nil {
		t.Errorf("Error running hiera.Hash: %s", v2)
	}

	v3, err3 := hiera.hash(context.TODO(), "java_opts")
	if err3 == nil || v3 != nil {
		t.Errorf("Error running hiera.Hash: %s", v3)
	}

	hieraBad := testHiera5ConfigBad()

	v4, err4 := hieraBad.hash(context.TODO(), "aws_tags")
	if err4 == nil || v4 != nil {
		t.Errorf("Error running hiera.Hash: %s", v4)
	}
}

func TestHiera5Value(t *testing.T) {
	hiera := testHiera5Config()

	v, err := hiera.value(context.TODO(), "aws_cloudwatch_enable")
	if err != nil {
		t.Errorf("Error running hiera.Value: %s", err)
	}

	if v != "true" {
		t.Errorf("aws_cloudwatch_enable is %s; want %s", v, "true")
	}

	v2, err2 := hiera.value(context.TODO(), keyUnavailable)
	if err2 == nil || v2 != "" {
		t.Errorf("Error running hiera.value: %s", v2)
	}

	hieraBad := testHiera5ConfigBad()

	v4, err4 := hieraBad.value(context.TODO(), "aws_cloudwatch_enable")
	if err4 == nil || v4 != "" {
		t.Errorf("Error running hiera.value: %s", v4)
	}
}

func TestHiera5Json(t *testing.T) {
	hiera := testHiera5Config()

	v, err := hiera.json(context.TODO(), "aws_tags")
	if err != nil {
		t.Errorf("Error running hiera.json: %s", err)
	}

	if v != `{"team":"A","tier":1}` {
		t.Errorf("aws_tags is %s; want %s", v, `{"team":"A","tier":1}`)
	}

	v2, err2 := hiera.json(context.TODO(), keyUnavailable)
	if err2 == nil || v2 != "" {
		t.Errorf("Error running hiera.json: %s", v2)
	}

	hieraBad := testHiera5ConfigBad()

	v4, err4 := hieraBad.json(context.TODO(), "aws_cloudwatch_enable")
	if err4 == nil || v4 != "" {
		t.Errorf("Error running hiera.json: %s", v4)
	}
}

func testHiera5Config() hiera5 {
	return newHiera5(
		"test-fixtures/hiera.yaml",
		map[string]interface{}{"service": "api", "environment": "live", "facts": "{timezone=>'CET'}"},
		"deep",
	)
}
func testHiera5ConfigBad() hiera5 {
	return newHiera5(
		"doesnt_exists/hiera.yaml",
		map[string]interface{}{"service": "api", "environment": "live", "facts": "{timezone=>'CET'}"},
		"deep",
	)
}
