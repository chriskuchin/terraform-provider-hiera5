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
