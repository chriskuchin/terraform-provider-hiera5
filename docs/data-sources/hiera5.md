---
page_title: "hiera5 Data Source - terraform-provider-hiera5"
subcategory: ""
description: |-
  
---

# Data Source `hiera5`



## Example Usage

```terraform
data "hiera5" "aws_cloudwatch_enable" {
  key = "aws_cloudwatch_enable"
}
```

## Schema

### Required

- **key** (String, Required)

### Optional

- **default** (String, Optional)
- **id** (String, Optional) The ID of this resource.

### Read-only

- **value** (String, Read-only)


