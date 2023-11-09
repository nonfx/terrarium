module "log_sync" {
  source = "../../modules/log_sync"
}

output "log_sync_self_link" {
  value = module.log_sync.self_link
}
