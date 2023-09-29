locals {
  tr_component_ecs_combined = {
    for cluster, config in var.ecs_config : cluster => {
      name         = config.name
      compute_type = config.engine.type
      engine = {
        type           = config.engine.type
        default_weight = config.engine.default_weight
        spot_weight    = config.engine.spot_weight
      }
      services = {
        for service, details in var.tr_component_ecs_services : service => {
          name     = try(details.name, null)
          cpu      = try(details.cpu, 512)
          memory   = try(details.memory, 1024)
          image    = try(details.image, null)
          port     = try(details.port, null)
          protocol = try(details.protocol, null)
        }
      }
    }
  }
}


module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 5.2"

  for_each = local.tr_component_ecs_combined

  cluster_name = each.value.name
  cluster_configuration = {
    execute_command_configuration = {
      logging = "OVERRIDE"
      log_configuration = {
        cloud_watch_log_group_name = "/aws/ecs/aws-ec2"
      }
    }
  }

  fargate_capacity_providers = each.value.compute_type == "FARGATE" ? {
    FARGATE = {
      default_capacity_provider_strategy = {
        weight = each.value.engine.default_weight
      }
    }
    FARGATE_SPOT = {
      default_capacity_provider_strategy = {
        weight = each.value.engine.spot_weight
      }
    }
  } : null

  services = {
    for service_key, service_value in var.tr_component_ecs_services : service_key => {
      cpu                    = service_value.cpu
      memory                 = service_value.memory
      environment            = try(service_value.environment, null)
      create_task_definition = try(service_key, null) != null ? true : false

      container_definitions = {
        # Container definition(s)
        service_key = {
          cpu                      = service_value.cpu
          memory                   = service_value.memory
          essential                = true
          image                    = service_value.image
          readonly_root_filesystem = false
          name                     = service_value.name

          port_mappings = try(service_value.port, null) != null ? [
            {
              name          = "${service_key}-port"
              containerPort = service_value.port
              protocol      = service_value.protocol
            }
          ] : []
        }
      }
      # To add service discovery
      # service_connect_configuration = service_value.port != null ? {
      #   namespace = var.environment
      #   service = {
      #     client_alias = {
      #       port     = service_value.port
      #       dns_name = "${service_key}-svc"
      #     }
      #     port_name      = "${service_key}-port"
      #     discovery_name = "${service_key}-svc"
      #   }
      # } : { }

      service_connect_configuration = { for k, v in {
        namespace = var.environment
        service = {
          client_alias = {
            port     = try(service_value.port, null)
            dns_name = "${service_key}-svc"
          }
          port_name      = "${service_key}-port"
          discovery_name = "${service_key}-svc"
        }
      } : k => v if try(service_value.connect, null) != null }


      load_balancer = { for k, v in {
        service = {
          target_group_arn = try(service_value.port, null) != null ? (length([for arn in module.alb["default"].target_group_arns : arn if length(regexall(substr(service_key, 0, 5), arn)) > 0]) > 0 ? ([for arn in module.alb["default"].target_group_arns : arn if length(regexall(substr(service_key, 0, 5), arn)) > 0][0]) : null) : null
          container_name   = try(service_value.port, null) != null ? service_value.name : null
          container_port   = try(service_value.port, null)
        }
      } : k => v if try(service_value.port, null) != null }

      enable_cloudwatch_logging = var.environment == "production" || var.environment == "prod" ? true : false
      subnet_ids                = var.private_subnet_ids["default"]
      create_security_group     = try(service_value.port, null) == null ? false : true
      security_group_rules = try(service_value.port, null) != null ? {
        "alb_ingress_${service_key}" = {
          type                     = "ingress"
          from_port                = service_value.port
          to_port                  = service_value.port
          protocol                 = service_value.protocol
          description              = "${service_key}-port"
          source_security_group_id = module.sg[service_key].security_group_id
        }
        egress_all = {
          type        = "egress"
          from_port   = 0
          to_port     = 0
          protocol    = "-1"
          cidr_blocks = ["0.0.0.0/0"]
        }
        } : {
        "alb_ingress_${service_key}" = {
          type                     = null
          from_port                = null
          to_port                  = null
          protocol                 = null
          description              = null
          source_security_group_id = null
        }
        egress_all = {
          type        = "egress"
          from_port   = 0
          to_port     = 0
          protocol    = "-1"
          cidr_blocks = ["0.0.0.0/0"]
        }
      }

    }
  }

  depends_on = [time_sleep.wait_60_seconds]

  tags = {
    Environment = var.environment
    Project     = "Example"
  }
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [module.alb]

  create_duration = "60s"
}

resource "aws_lb_listener_rule" "lb-listener-rule" {
  for_each     = { for k, v in var.tr_component_ecs_services : k => v if lookup(v, "port", null) != null }
  listener_arn = module.alb["default"].https_listener_ids[0]
  priority     = 100 + index(keys(var.tr_component_ecs_services), each.key)

  action {
    type             = "forward"
    target_group_arn = try(each.value.port, null) != null ? (length([for arn in module.alb["default"].target_group_arns : arn if length(regexall(substr(each.key, 0, 5), arn)) > 0]) > 0 ? ([for arn in module.alb["default"].target_group_arns : arn if length(regexall(substr(each.key, 0, 5), arn)) > 0][0]) : null) : null
  }

  condition {
    path_pattern {
      values = [each.value.path]
    }
  }
  condition {
    host_header {
      values = [each.value.site]
    }
  }
}

module "sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.1.0"

  for_each = var.tr_component_ecs_services

  name        = var.extract_resource_name
  description = "Security group for ${each.key}-service with custom ports open within VPC"
  vpc_id      = var.vpc_id["default"]

  ingress_cidr_blocks = [var.vpc_cidr_block["default"]]
  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      description = "Service name"
      cidr_blocks = "0.0.0.0/0"
    },
  ]
  ingress_with_cidr_blocks = try(each.value.port, []) != [] ? [
    {
      from_port   = each.value.port
      to_port     = each.value.port
      protocol    = each.value.protocol
      description = "${each.key}-service ports"
      cidr_blocks = var.vpc_cidr_block["default"]
    }
  ] : []
}


module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 8.0"

  for_each = var.lb_config

  name = "${var.extract_resource_name}-alb"

  create_lb = each.value.create_alb == true ? true : false

  load_balancer_type    = each.value.load_balancer_type
  create_security_group = true

  vpc_id          = var.vpc_id["default"]
  subnets         = var.public_subnet_ids["default"]
  security_groups = [for sg in values(module.sg) : sg.security_group_id]

  access_logs = {
    bucket = module.s3_bucket[each.key].s3_bucket_id
  }

  security_group_rules = [
    {
      type        = "ingress"
      from_port   = 80
      to_port     = 80
      protocol    = "TCP"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {
      type        = "ingress"
      from_port   = 443
      to_port     = 443
      protocol    = "TCP"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {
      type        = "egress"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]


  target_groups = [
    for service_key, service_value in var.tr_component_ecs_services : {
      name_prefix      = substr("${service_key}", 0, 6)
      backend_protocol = "HTTP"
      backend_port     = try(service_value.port, null)
      target_type      = "ip"
    }
    if try(service_value.port, null) != null
  ]

  https_listeners = [
    {
      port               = 443
      protocol           = "HTTPS"
      certificate_arn    = var.certificate_arn
      target_group_index = 0
    }
  ]

  http_tcp_listeners = [
    {
      port        = 80
      protocol    = "HTTP"
      action_type = "redirect"
      redirect = {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }
  ]

  depends_on = [module.sg]

  tags = {
    Environment = "Test"
  }
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
  keepers = {
    bucket_base_name = var.extract_resource_name
  }
}

module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.14.0"

  for_each = var.lb_config

  bucket = "${var.extract_resource_name}-alb-logs-${random_id.bucket_suffix.hex}"
  acl    = "log-delivery-write"

  # Allow deletion of non-empty bucket
  force_destroy = var.environment == "production" || var.environment == "prod" ? false : true

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  attach_elb_log_delivery_policy = true # Required for ALB logs
  attach_lb_log_delivery_policy  = true # Required for ALB/NLB logs
}


resource "aws_route53_record" "alb_route53_record" {
  for_each = module.alb

  zone_id = var.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = each.value.lb_dns_name
    zone_id                = each.value.lb_zone_id
    evaluate_target_health = false
  }
}
