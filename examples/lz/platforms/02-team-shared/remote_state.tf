data "terraform_remote_state" "common" {
  backend = "gcs"
  config = {
    bucket = "UPDATE_BACKEND_BUCKET"
    prefix = "01-common/"
  }
}
