data "hiera5_bool" "enable_spot_instances" {
  key     = "enable_spot_instances"
  default = false
}