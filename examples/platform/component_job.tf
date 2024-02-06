module "job_queue_runner" {
  source = "./mock-modules/k8s-service"

  for_each = local.tr_component_job_queue

  name = "queue_watcher_${each.key}"
  cluster_id = module.k8s_cluster.cluster_id
}

module "job_scheduled_runner" {
  source = "./mock-modules/k8s-service"

  for_each = local.tr_component_job_scheduled

  name = "scheduler_${each.key}"
  cluster_id = module.k8s_cluster.cluster_id
}

# A job that performs tasks in the queue.
# @title: Queue Job
module "tr_component_job_queue" {
  source = "./mock-modules/pubsub-queue"

  for_each = local.tr_component_job_queue

  event_receiver_url = module.job_queue_runner[each.key].host
  name = each.key
}

# A job that is run at scheduled intervals.
# @title: Scheduled Job
module "tr_component_job_scheduled" {
  source = "./mock-modules/cloud-scheduler"

  for_each = local.tr_component_job_scheduled

  event_receiver_url = module.job_scheduled_runner[each.key].host
  name = each.key
}
