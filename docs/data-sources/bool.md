---
page_title: "hiera5_bool Data Source - terraform-provider-hiera5"
subcategory: ""
description: |-
  
---

# Data Source `hiera5_bool`



## Example Usage

```terraform
data "hiera5_bool" "enable_spot_instances" {
  key     = "enable_spot_instances"
  default = false
}
```

## Schema

### Required

- **key** (String, Required)

### Optional

- **default** (Boolean, Optional)
- **id** (String, Optional) The ID of this resource.

### Read-only

- **value** (Boolean, Read-only)


