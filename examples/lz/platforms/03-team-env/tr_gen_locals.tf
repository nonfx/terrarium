# this file is meant for keeping the component inputs,
# anything else kept here will not be added to the generated code

locals {
  tr_component_bucket = {
    default = {
      name = ""
    }
  }

  tr_component_postgres = {
    default = {
      version = "11"
    }
  }

  tr_component_environment = {
    default = {}
  }
}

module "tr_component_environment" {
  source = "../../modules/noop"

  depends_on = [ module.subnet_dev ]
}
