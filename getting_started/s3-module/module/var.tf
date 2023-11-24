variable "acl" {
    default = "private" # for public just change the value to "public-read" 
}

variable "control_object_ownership" {
    default = true
}

variable "object_ownership" {
    default = "ObjectWriter"
}


variable "versioning" {
    default = {
        "enabled" = true
    }
}
