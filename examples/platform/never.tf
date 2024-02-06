# this resource is not used in any component's dependency tree,
# so it'll always be excluded in the generated code
resource "random_string" "random_never" {
  length  = 8
  special = false
  upper   = false
  lower   = true
}
