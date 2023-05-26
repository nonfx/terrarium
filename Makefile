######################################################
# Following targets need docker installed on the system
######################################################

# Define variables for PostgreSQL container
POSTGRES_CONTAINER := postgres
POSTGRES_DB := cc_terrarium
POSTGRES_USER := postgres

# Define variables for pg_dump command
DUMP_DIR := ./data

# Define phony targets (targets that don't correspond to files)
.PHONY: db-dump docker-run docker-stop docker-stop-clean seed_resources seed_mappings seed_modules seed test

db-dump:  ## Target for dumping PostgreSQL database to a file
	docker compose exec -T $(POSTGRES_CONTAINER) pg_dump --username=$(POSTGRES_USER) --create --file=/docker-entrypoint-initdb.d/$(POSTGRES_DB).sql $(POSTGRES_DB)

docker-run:  ## Starts app in docker containers
	docker compose up -d

docker-stop:  ## Stops and removes docker containers
	docker compose down

docker-stop-clean:  ## Stops and removes containers as well as volumes to cleanup database
	docker compose down -v

######################################################
# Following targets need terraform installed on the system
######################################################

TERRAFORM_DIR := ./terraform
TF_FILES := $(shell find $(TERRAFORM_DIR) -name '*.tf')

$(TERRAFORM_DIR)/.terraform: $(TF_FILES)
	@rm -rf terraform/.terraform
	@cd $(TERRAFORM_DIR) && terraform init

# generate tf_resources.json file for set terraform providers
cache_data/tf_resources.json: $(TERRAFORM_DIR)/.terraform
	@echo "generating ./cache_data/tf_resources.json"
	@cd terraform && terraform providers schema -json > ../cache_data/tf_resources.json

# run terraform init to have terraform modules downloaded
terraform/.terraform/modules/modules.json: terraform/modules.tf
	@echo "running terraform init"
	@cd terraform && terraform init

######################################################
# Following targets need Go installed on the system
######################################################

GOPATH = $(shell go env GOPATH|cut -d ":" -f 1)

test:  ## Run go unit tests
	@$(GOPATH)/bin/godotenv go test `go list github.com/cldcvr/terrarium/... | grep -v /pkg/terraform-config-inspect/`

seed: seed_resources seed_mappings seed_modules

seed_resources: docker-run cache_data/tf_resources.json  ## Seed tf-provider resources into db from terraform/provider.tf
	@echo "Running resource seed..."
	@$(GOPATH)/bin/godotenv go run ./api/cmd/seed_resources

seed_mappings: docker-run  ## Load .env file and run seed_mappings
	@echo "Running mapping seed..."
	@$(GOPATH)/bin/godotenv go run ./api/cmd/seed_mappings

seed_modules: docker-run $(TERRAFORM_DIR)/.terraform  ## Seed tf-modules into db from terraform/modules.tf
	@echo "Running module seed..."
	@$(GOPATH)/bin/godotenv go run ./api/cmd/seed_modules

seed_mappings: docker-run  ## Load .env file and run seed_mappings
	@echo "Running mapping seed..."
	@$(GOPATH)/bin/godotenv go run ./api/cmd/seed_mappings

include scripts/mocks.mak
include scripts/protoc.mak
