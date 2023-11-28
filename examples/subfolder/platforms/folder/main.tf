module "tr_component_folder" {
  source       = "../../modules/folder"
  for_each     = local.tr_component_folder
  display_name = each.value.display_name
  parent       = each.value.parent
}

module "tr_component_project" {
  source              = "../../modules/project"
  for_each            = local.tr_component_project
  name                = each.value.name
  project_id          = each.value.project_id
  folder_id           = module.tr_component_folder["${each.value.id}"].name
  billing_account     = each.value.billing_account
  auto_create_network = each.value.auto_create_network
  labels              = each.value.labels
  skip_delete         = each.value.skip_delete
}

module "tr_component_subfolder" {
  source       = "../../modules/folder"
  for_each     = local.tr_component_subfolder
  display_name = each.value.display_name
  parent       = module.tr_component_folder["${each.value.id}"].name
}

module "tr_component_subfolder1" {
  source       = "../../modules/folder"
  for_each     = local.tr_component_subfolder1
  display_name = each.value.display_name
  parent       = module.tr_component_subfolder["${each.value.id}"].name
}
