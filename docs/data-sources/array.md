---
page_title: "hiera5_array Data Source - terraform-provider-hiera5"
subcategory: ""
description: |-
  
---

# Data Source `hiera5_array`



## Example Usage

```terraform
data "hiera5_array" "java_opts" {
  key = "java_opts"
}
```

## Schema

### Required

- **key** (String, Required)

### Optional

- **default** (List of String, Optional)
- **id** (String, Optional) The ID of this resource.

### Read-only

- **value** (List of String, Read-only)


