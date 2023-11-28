output "tr_component_folder_name" {
  value = { for k, v in module.tr_component_folder : k => v.name }

}

output "tr_component_project_number" {
  value = { for k, v in module.tr_component_project : k => v.number }
}

output "tr_component_subfolder_name" {
  value = { for k, v in module.tr_component_subfolder : k => v.name }

}
