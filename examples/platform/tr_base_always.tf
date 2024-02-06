# All tf files starting with `tr_base_` gets copied to the generated
# code as it is. regardless of the requirement.

resource "random_string" "random_always" {
  length  = 8
  special = false
  upper   = false
  lower   = true
}
