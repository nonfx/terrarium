variable "cidr" {
  default = "10.0.1.0/24"
}

variable "network" {
  description = "name of the network to create the subnet in"
}

variable "env_type" {
  default = "dev"
}

output "name" {
  value = base64encode(concat(var.cidr, var.network))
}
