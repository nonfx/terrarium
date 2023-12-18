variable "subnet" {
  description = "subnet to put the db in"
  type = string
}

variable "engine_version" {
  description = "version of the db"
  type = string
}

variable "family" {
  description = "database family like postgres, mysql, etc."
  type = string
}

variable "instance_type" {
  default = "small"
}

resource "random_string" "host" {
  length = 8
}

resource "random_password" "password" {
  length = 8
}

resource "random_integer" "port" {
  min = 10000
  max = 100000
}

output "host" {
  value = random_string.host.result
}

output "port" {
  value = random_integer.port.result
}

output "name" {
  value = base64encode(format("%s_%s", var.subnet, var.engine_version))
}

output "password" {
  value = random_password.password.result
  sensitive = true
}
