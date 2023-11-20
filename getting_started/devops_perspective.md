## Getting started with Devops perspective

> For Devops we have terraform modules already created. You can check them in below path

```
examples/platform-demo/platform
```

- As a Devops we know that we are eager to know how terrarium is getting synced with terraform 

- So here are the details 

    - Source code repo for terrarium

        ```
        ```

    - Code is written on Golang




### Writing a basic module for s3 bucket integrated with terrarium

```
module "tr_component_s3" {
  source  = "../"
  bucket_name = "Ollion"
  name = "Ollion_Bucket"
  environment ="Ollion_dev"
}
```

### Generate terrarium.yaml file from terraform code

```
terrarium platform lint
```

### Generated terrarium yaml file 

```
components:
    - id: s3
      title: S3
      outputs:
        type: object
graph:
    - id: module.tr_component_s3
      requirements: []
```