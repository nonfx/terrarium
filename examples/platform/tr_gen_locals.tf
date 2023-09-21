locals {

  # Inputs for 'postgres' component instances.
  tr_component_postgres = {
    "default" : {
      # Version of the PostgreSQL engine to use
      # @enum: 11, 12, 13
      "version" : "11",
      # The name provided here may get prefix and suffix based
      # @title: Database Name
      "db_name" : "default_db"
    }
  }

  # Inputs for 'redis' component instances.
  tr_component_redis = {
    "default" : {
      # Version of the Redis engine to use
      "version" : "5.0.6"
    }
  }
}
