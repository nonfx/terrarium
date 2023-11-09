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

  tr_component_logs_bucket = {
    default = {
      retention_period = 365
    }
  }
}

module "tr_component_environment" {
  source = "../../modules/noop"

  depends_on = [ module.vpc ]
}

module "tr_component_postgres" {
  source = "../../modules/noop"
}

module "tr_component_bucket" {
  source = "../../modules/noop"
}

module "tr_component_logs_bucket" {
  source = "../../modules/noop"

  depends_on = [ module.log_sync ]
}
