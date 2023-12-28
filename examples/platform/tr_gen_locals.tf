locals {

  # Inputs for 'postgres' component instances.
  tr_component_postgres = {
    "default": {
      # Version of the PostgreSQL engine to use
      # @enum: 11, 12, 13
      "version": "11",

      # The name provided here may get prefix and suffix based
      # @title: Database Name
      "db_name": "default_db"
    }
  }

  # Inputs for 'redis' component instances.
  tr_component_redis = {
    "default": {
      # Version of the Redis engine to use
      "version": "5.0.6"
    }
  }

  tr_component_server_static = {
    "default": {
      # The port number on which the server should listen.
      port: 80
    }
  }

  tr_component_server_web = {
    "default": {
      # The port number on which the server should listen.
      port: 80
    }
  }

  tr_component_server_private = {
    "default": {
      # The port number on which the server should listen.
      port: 80
    }
  }

  tr_component_job_queue = {
    "default": {}
  }

  tr_component_job_scheduled = {
    "default": {
      # The schedule on which the job should run, in cron format.
      "schedule": "0 0 12 * * ?"
    }
  }
}
