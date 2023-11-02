## Getting started with developer perspective

### Lets take an example of a Go application(banking app)

#### Steps

Step 1: Clone the repository
```
 https://github.com/cldcvr/t8-demo-workflow-go-psql
```
Step 2: Checkout if the application directory contains terrarium.yaml file if not then we have to create one

Step 3: Now just run terrarium command to generate the terraform code for the particular application

```
terrarium generate -p /Users/himanshuuniyal/test-terrarium/terrarium/examples/platform-demo/platform -a . -o /Users/himanshuuniyal/test-terrarium/output
```

#### Command flags

|  flags  |  description |
| ------  |  ----------- |
|  -p     |  specifies the application dir that is being used to generate application infra |
|  -a     |  specifies the terrarium.yaml file path            |
|  -o     |   specifies the terraform code path for the application infra           |
