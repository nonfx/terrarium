module "ecs" {
  source = "../modules/ecs"

  extract_resource_name = local.extract_resource_name
  environment           = var.environment

  tr_component_ecs_services = merge(local.app_vars_service_web, local.app_vars_service_private)
  ecs_config                = var.ecs_config
  vpc_id                    = module.vpc.vpc_id
  public_subnet_ids         = module.vpc.public_subnet_ids
  private_subnet_ids        = module.vpc.private_subnet_ids
  vpc_cidr_block            = module.vpc.vpc_cidr_blocks
  lb_config                 = var.lb_config
  certificate_arn           = module.acm.acm_certificate_arn
  domain_name               = var.domain_name
  zone_id                   = var.zone_id

  depends_on = [module.acm]

}

locals {
  app_vars_service_web     = { for k, v in local.tr_component_service_web : k => merge(v,
    {
      public = true
      name = k
      image = "public.ecr.aws/nginx/nginx:alpine-slim"
      site = format("%s.%s", k, var.domain_name)
    }
  )}
  app_vars_service_private = { for k, v in local.tr_component_service_private : k => merge(v,
    {
      public = false
      name = k
      image = "public.ecr.aws/nginx/nginx:alpine-slim"
      site = format("%s.%s", k, var.domain_name)
    }
  )}
}

module "tr_component_service_web" {
  source     = "cloudposse/label/null"
  depends_on = [module.ecs]
}

module "tr_component_service_private" {
  source     = "cloudposse/label/null"
  depends_on = [module.ecs]
}

