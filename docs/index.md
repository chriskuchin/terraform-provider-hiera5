---
page_title: "hiera5 Provider"
subcategory: ""
description: |-
  
---

# hiera5 Provider



## Example Usage

```terraform
provider "hiera5" {
  # Optional
  config = "~/hiera.yaml"
  # Optional
  scope = {
    environment = "live"
    service     = "api"
    # Complex variables are supported using pdialect
    facts = "{timezone=>'CET'}"
  }
  # Optional
  merge = "deep"
}
```

## Schema

### Optional

- **config** (String, Optional)
- **merge** (String, Optional)
- **scope** (Map of String, Optional)
