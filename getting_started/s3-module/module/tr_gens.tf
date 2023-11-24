locals{
    tr_component_s3_bucket_for_logs = {
        "defaults" : {
         "bucket" : "terrarium-test-bucket"
         "acl"    : "log-delivery-write"
         "force_destroy" : "true"

        "control_object_ownership" : "true"
        "object_ownership"         : "ObjectWriter"

        "attach_elb_log_delivery_policy" : "true"  
        "attach_lb_log_delivery_policy"  : "true"
        }
    }
}