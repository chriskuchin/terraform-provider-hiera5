# Terraform Hiera 5 Provider

[![pipeline status](https://gitlab.com/sbitio/terraform-provider-hiera5/badges/master/pipeline.svg)](https://gitlab.com/sbitio/terraform-provider-hiera5/-/commits/master) [![coverage report](https://gitlab.com/sbitio/terraform-provider-hiera5/badges/master/coverage.svg)](https://gitlab.com/sbitio/terraform-provider-hiera5/-/commits/master) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/sbitio/terraform-provider-hiera5)](https://goreportcard.com/report/sbitio/terraform-provider-hiera5)

This provider implements data sources that can be used to perform hierachical data lookups with Hiera.

This is useful for providing configuration values in an environment with a high level of dimensionality or for making values from an existing Puppet deployment available in Terraform.

It's based on [Terraform hiera provider](https://github.com/ribbybibby/terraform-provider-hiera) and [SmilingNavern's fork](https://github.com/SmilingNavern/terraform-provider-gohiera)

## Goals
* Clean implementation based on [Terraform Plugin SDK](https://www.terraform.io/docs/extend/plugin-sdk.html)
* Clean API implementatation based on [Lyra](https://lyraproj.github.io/)'s [Hiera in golang](https://github.com/lyraproj/hiera)
* Painless migration from [Terraform hiera provider](https://github.com/ribbybibby/terraform-provider-hiera), keeping around some naming and data sources

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.12.x

## Usage

### Configuration
To configure the provider:
```hcl
provider "hiera5" {
  # Optional
  config = "~/hiera.yaml"
  # Optional
  scope = {
    environment = "live"
    service     = "api"
    # Complex variables are supported using pdialect
    facts       = "{timezone=>'CET'}"
  }
  # Optional
  merge  = "deep"
}
```

### Data Sources
This provider only implements data sources.

#### Hash
To retrieve a hash:
```hcl
data "hiera5_hash" "aws_tags" {
    key = "aws_tags"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the hash, represented as a map

Terraform doesn't support nested maps or other more complex data structures. Any keys containing nested elements won't be returned.

#### Array
To retrieve an array:
```hcl
data "hiera5_array" "java_opts" {
    key = "java_opts"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the array (list)

#### Value
To retrieve any other flat value:
```hcl
data "hiera5" "aws_cloudwatch_enable" {
    key = "aws_cloudwatch_enable"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the value

All values are returned as strings because Terraform doesn't implement other types like int, float or bool. The values will be implicitly converted into the appropriate type depending on usage.

#### Json
To retrieve anything JSON encoded:
```hcl
data "hiera5_json" "aws_tags" {
    key = "aws_tags"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the returned value, JSON encoded

As Terraform doesn't support nested maps or other more complex data structures this data source makes perfect fit dealing with complex values.

## Example

Take a look at [test-fixtures](./hiera5/test-fixtures)

## Thanks to
* Julien Andrieux for writting [Go tools and GitLab: How to do continuous integration like a boss](https://about.gitlab.com/blog/2017/11/27/go-tools-and-gitlab-how-to-do-continuous-integration-like-a-boss/), a really good starting point.

## Develpment

### Requirements

* [Go](https://golang.org/doc/install) 1.12

### Notes

[This repository is vendored as recomended on Terraform's docs](https://www.terraform.io/docs/extend/terraform-0.12-compatibility.html#upgrading-to-the-latest-terraform-sdk)

### Whishlist
* [ ] Support overriding merge strategy in Data Sources
* [ ] Support overriding scope variables in Data Sources
