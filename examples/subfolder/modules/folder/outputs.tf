output "name" {
  value       = google_folder.folder.name
  description = "The resource name of the Folder. Its format is folders/{folder_id}."
}
output "lifecycle_state" {
  value       = google_folder.folder.lifecycle_state
  description = "The lifecycle state of the folder such as ACTIVE or DELETE_REQUESTED."
}
output "create_time" {
  value       = google_folder.folder.create_time
  description = "Timestamp when the Folder was created."
}