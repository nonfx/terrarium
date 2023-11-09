variable "retention_period" {
  default = 365
}

output "self_link" {
  value = base64encode(var.retention_period)
}
