locals {
  tr_component_postgres = {
    "default" : {
      "version" : "11",
      "db_name" : "default_db"
    }
  }

  tr_component_redis = {
    "default": {
      "version": "5.0.6"
    }
  }
}
