data "terraform_remote_state" "shared" {
  backend = "gcs"
  config = {
    bucket = "UPDATE_BACKEND_BUCKET"
    prefix = "02-shared/"
  }
}
