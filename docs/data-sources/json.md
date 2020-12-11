---
page_title: "hiera5_json Data Source - terraform-provider-hiera5"
subcategory: ""
description: |-
  
---

# Data Source `hiera5_json`



## Example Usage

```terraform
data "hiera5_json" "aws_tags" {
  key = "aws_tags"
}

locals {
  aws_tags = jsondecode(data.hiera5_json.aws_tags.value)
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


