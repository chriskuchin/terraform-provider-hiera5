---
page_title: "hiera5_hash Data Source - terraform-provider-hiera5"
subcategory: ""
description: |-
  
---

# Data Source `hiera5_hash`



## Example Usage

```terraform
data "hiera5_hash" "aws_tags" {
  key = "aws_tags"
}
```

## Schema

### Required

- **key** (String, Required)

### Optional

- **default** (Map of String, Optional)
- **id** (String, Optional) The ID of this resource.

### Read-only

- **value** (Map of String, Read-only)


