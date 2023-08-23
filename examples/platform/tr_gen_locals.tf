locals {

  # Inputs for 'postgres' component instances.
  # version: Version of the PostgreSQL engine to use
  # db_name[Database Name]: The name provided here may get prefix and suffix based
  tr_component_postgres = {
    "default" : {
      "version" : "11",
      "db_name" : "default_db"
    }
  }

  # Inputs for 'redis' component instances.
  # version: Version of the Redis engine to use
  tr_component_redis = {
    "default" : {
      "version" : "5.0.6"
    }
  }
}
