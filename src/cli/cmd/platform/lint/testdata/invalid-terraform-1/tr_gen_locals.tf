locals {
  # Inputs for 'postgres' component instances.
  tr_component_postgres = {
    "default" : {
      # Version of the PostgreSQL engine to use
      # @enum: 11.11,12.3, 13.9
      "version" : "15", # <------ ERROR: default value not in enum
      # The name provided here may get prefix and suffix based
      # @title: Database Name
      "db_name" : "default_db"
    }
  }


  tr_component_postgres_enabled = length(local.tr_component_postgres) > 0
  tr_taxon_sql_enabled          = anytrue([local.tr_component_postgres_enabled])
  tr_taxon_database_enabled     = anytrue([local.tr_taxon_sql_enabled])
}
