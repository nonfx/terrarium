output "tr_component_sqs_arn" {
    value = {for k, v in module.tr_component_sqs : k => v.queue_arn }
}


output "tr_component_sqs_name" {
    value = {for k, v in module.tr_component_sqs : k => v.queue_name }
}

output "tr_component_sqs_id" {
    value = {for k, v in module.tr_component_sqs : k => v.queue_id }
}

output "tr_component_sqs_url" {
    value = {for k, v in module.tr_component_sqs : k => v.queue_url }
}
