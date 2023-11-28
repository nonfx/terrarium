locals {
  tr_component_folder = {
    apac = {
      display_name = "test2"
      parent       = "folders/802470219928"
    }
    canada = {
      display_name = "test3"
      parent       = "folders/802470219928"
    }
    common = {
      display_name = "test1"
      parent       = "folders/802470219928"
    }
  }
  tr_component_project = {
    dev-project = {
      auto_create_network = false
      billing_account     = "017769-C68281-DA1844"
      id                  = "apac"
      labels              = {}
      name                = "d-prj-test"
      project-id          = "test-1213"
      project_id          = "d-prj-100"
      service             = []
      skip_delete         = false
    }
    prod-project = {
      auto_create_network = false
      billing_account     = "017769-C68281-DA1844"
      id                  = "apac"
      labels              = {}
      name                = "p-proj-test"
      project-id          = "test-1213"
      project_id          = "p-proj-700"
      service             = []
      skip_delete         = false
    }
  }
  tr_component_subfolder = {
    dev = {
      display_name = "test5"
      id           = "apac"
      parent       = "test1"
    }
    shared = {
      display_name = "test4"
      id           = "apac"
      parent       = "test1"
    }
  }
  tr_component_subfolder1 = {
    monitoring = {
      display_name = "monitoring-folder"
      id           = "shared"
      parent       = "test1"
    }
  }
}
