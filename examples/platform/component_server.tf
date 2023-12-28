# A server that hosts and serves static files.
# @title: Static Server
module "tr_component_server_static" {
  source = "./mock-modules/bucket-static-site"

  for_each = local.tr_component_server_static

  bucket_name = each.key
}

module "k8s_cluster" {
  source = "./mock-modules/k8s-cluster"
}

# A server that hosts web applications and handles HTTP requests.
# @title: Web Server
module "tr_component_server_web" {
  source = "./mock-modules/k8s-service"

  for_each = local.tr_component_server_web

  cluster_id = module.k8s_cluster.cluster_id
  is_public = true
  name = each.key
  port = each.value.port
}

# A server that is not exposed to the public internet.
# @title: Private Server
module "tr_component_server_private" {
  source = "./mock-modules/k8s-service"

  for_each = local.tr_component_server_private

  cluster_id = module.k8s_cluster.cluster_id
  is_public = false
  name = each.key
  port = each.value.port
}
