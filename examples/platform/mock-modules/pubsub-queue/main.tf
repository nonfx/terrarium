variable "name" {
  type = string
  description = "Name of the queue"
  default = "default"
}

variable "event_receiver_url" {
  type = string
  description = "URL for push based model where the events are sent to the given url"
  default = ""
}

resource "random_id" "random" {
  byte_length = 8
}

output "queue_id" {
  description = "Id of the queue"
  value = random_id.random.id
}

output "queue_url" {
  description = "URL of the queue to send events at and query for messages using sdk"
  value = "http://${random_id.random.id}.mockqueue.com"
}
