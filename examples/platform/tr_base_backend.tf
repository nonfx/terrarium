# All tf files starting with `tr_base_` gets copied to the generated
# code as it is. regardless of the requirement.

terraform {
  backend "gcs" {
    bucket = "bkt-tfstate-001"
    prefix = "01-common/"
  }
}
