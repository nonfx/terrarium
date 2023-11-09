module "tr_component_postgres" {
  source = "../../modules/cloudsql"

  for_each = local.tr_component_postgres

  family = "postgres"
  instance_type = var.instance_type
  version = each.value.version

  # use the appropriate subnet based on the environment
  subnet = data.terraform_remote_state.shared.subnet[var.env_type]
}
