locals {
  tr_component_postgres = {
    "default" : {
      "engine_version" : "11.20",
      "family" : "postgres"
    }
  }
  tr_component_redis = {
    "default" : {
      "version" : "5.0.6"
    },
  }
  tr_component_service_web = {
    "default" : {
      port : 80
    }
  }

  tr_component_service_enabled  = length(local.tr_component_service_web) > 0
  tr_component_redis_enabled    = length(local.tr_component_redis) > 0
  tr_component_postgres_enabled = length(local.tr_component_postgres) > 0
  tr_taxon_sql_enabled          = anytrue([local.tr_component_postgres_enabled])
  tr_taxon_database_enabled     = anytrue([local.tr_taxon_sql_enabled])
}
