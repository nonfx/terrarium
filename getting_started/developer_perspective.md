## Getting started with developer perspective

### Lets take an example of a Go application(banking app)

#### Steps

Step 1: Clone the repository
```
 https://github.com/cldcvr/t8-demo-workflow-go-psql
```
Step 2: Checkout if the application directory contains terrarium.yaml file if not then we have to create one. Below is the sample example for this particular application

```
id: banking_app
name: Banking App
env_prefix: BA

compute:
  use: service_web
  inputs:
    port: 3000

dependencies:
  - id: ledgerdb
    use: postgres
    inputs:
      version: "11.20"
  - id: cachedb
    use: redis
```

Step 3: Now just run terrarium command to generate the terraform code for the particular application

```
terrarium generate -p /Users/himanshuuniyal/test-terrarium/terrarium/examples/platform-demo/platform -a . -o /Users/himanshuuniyal/test-terrarium/output
```
> The above command will give you the following output

![Alt text](getting_started/images/one.png?raw=true "generate")
![Alt text](getting_started/images/two.png?raw=true "output")

#### Command flags

|  flags  |  description |
| ------  |  ----------- |
|  -p     |  specifies the application dir that is being used to generate application infra |
|  -a     |  specifies the terrarium.yaml file path            |
|  -o     |   specifies the terraform code path for the application infra           |
