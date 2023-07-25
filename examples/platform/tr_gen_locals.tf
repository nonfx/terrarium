locals {
  tr_component_postgres = {
    "default" : {
      "version" : "11.12",
      "db_name" : "default_db"
    }
  }


  tr_component_postgres_enabled = length(local.tr_component_postgres) > 0
  tr_taxon_sql_enabled          = anytrue([local.tr_component_postgres_enabled])
  tr_taxon_database_enabled     = anytrue([local.tr_taxon_sql_enabled])
}
