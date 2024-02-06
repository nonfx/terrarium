variable "name" {
  type = string
  description = "Name of the scheduler"
  default = "default"
}

variable "event_receiver_url" {
  type = string
  description = "URL to push the events to"
  default = ""
}

resource "random_id" "random" {
  byte_length = 8
}

output "scheduler_id" {
  description = "Id of the scheduler"
  value = random_id.random.id
}
