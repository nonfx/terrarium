locals {
  tr_component_folder = {
    "default" = {
      "display_name" : "test"
      "parent" : "test1",
    }
  }
  tr_component_project = {
    "default" = {
      "name" : "my-project",
      "project-id" : "test-1213"
      #"folder_id" : "folders/80123",
      "billing_account" : "dev-123",
      "auto_create_network" : false,
      "labels" : {},
      "skip_delete" : false,
      "service" : []
      "id" : ""
    }
  }
  tr_component_subfolder = {
    "default" = {
      "display_name" : "test",
      "parent" : "test1",
      "id" : ""
    }
  }
  tr_component_subfolder1 = {
    "default" = {
      "display_name" : "test",
      "parent" : "test1",
      "id" : ""
    }
  }
}
