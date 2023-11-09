variable "name" {
  description = "name of the network"
}

output "vpc_name" {
  value = var.name
}
