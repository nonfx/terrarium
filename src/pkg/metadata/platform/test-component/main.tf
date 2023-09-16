locals {

  # Inputs for 'postgres' component instances.
  tr_component_postgres = {
    "default" : {
      # Version of the PostgreSQL engine to use
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
      # @description: Redis description override
      "version" : "5.0.6"
    }
  }
}

# A relational database management system using SQL.
# @title: PostgreSQL Database
# Dolores fugiat dolor illo omnis optio ipsam.
module "tr_component_postgres" {
  source = "terraform-aws-modules/rds/aws"
}
