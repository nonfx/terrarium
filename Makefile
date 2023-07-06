######################################################
# Following targets need docker installed on the system
######################################################

# Define phony targets (targets that don't correspond to files)
.PHONY: docker-init db-dump docker-build docker-run start-db docker-stop docker-stop-clean docker-tools-build docker-seed docker-api-test

# Define variables for PostgreSQL container
POSTGRES_CONTAINER := postgres
POSTGRES_DB := cc_terrarium
POSTGRES_USER := postgres

# Define variables for pg_dump command
DUMP_DIR := ./data

docker-init:  ## Initialize the environment before running docker commands
	touch ${HOME}/.netrc

db-dump:  ## Target for dumping PostgreSQL database to a file
	docker compose exec -T $(POSTGRES_CONTAINER) pg_dump -U $(POSTGRES_USER) $(POSTGRES_DB) | dos2unix > data/$(POSTGRES_DB).sql

docker-build:  ## Build container image
	docker compose build

docker-run:  ## Starts app in docker containers
	docker compose up -d

start-db:  ## Starts database in docker containers
	docker compose up -d postgres

docker-stop:  ## Stops and removes docker containers
	docker compose down

docker-stop-clean:  ## Stops and removes containers as well as volumes to cleanup database
	docker compose down -v

docker-tools-build:
	docker compose --profile tooling build

docker-seed: docker-tools-build start-db
	docker compose run --rm seeder

docker-api-test: docker-tools-build
	docker compose run --rm test

######################################################
# Following targets need terraform installed on the system
######################################################

.PHONY: clean_tf tf_init

TERRAFORM_DIR := ./terraform
TF_FILES := $(shell find $(TERRAFORM_DIR) -name '*.tf' -not -path '$(TERRAFORM_DIR)/.terraform/*')

clean_tf:
	rm -rf $(TERRAFORM_DIR)/.terraform
	rm -f $(TERRAFORM_DIR)/.terraform.lock.hcl

tf_init: $(TERRAFORM_DIR)/.terraform

# generate tf_resources.json file for set terraform providers
cache_data/tf_resources.json: $(TERRAFORM_DIR)/.terraform
	@echo "generating ./cache_data/tf_resources.json"
	@mkdir -p cache_data
	@cd terraform && terraform version && terraform providers schema -json > ../cache_data/tf_resources.json

$(TERRAFORM_DIR)/.terraform: $(TF_FILES)
	@cd $(TERRAFORM_DIR) && terraform version && terraform init || (terraform providers && exit 1)
	@touch $(TERRAFORM_DIR)/.terraform

######################################################
# Following targets need Go installed on the system
######################################################

.PHONY: test seed seed_resources seed_modules seed_mappings

-include .env
export

SEED_SRCS := $(shell find ./src/pkg ./src/cli -name '*.go')

test:  ## Run go unit tests
	go test `go list github.com/cldcvr/terrarium/...`

seed: seed_resources seed_modules seed_mappings

seed_resources: .bin/cli cache_data/tf_resources.json  ## Seed tf-provider resources into db from terraform/provider.tf
	@echo "Running resource seed..."
	@./.bin/cli farm resources

seed_modules: .bin/cli $(TERRAFORM_DIR)/.terraform  ## Seed tf-modules into db from terraform/modules.tf
	@echo "Running module seed..."
	@./.bin/cli farm modules

seed_mappings: .bin/cli  ## Load .env file and run seed_mappings
	@echo "Running mapping seed..."
	@./.bin/cli farm mappings

.bin/cli: $(SEED_SRCS)
	@echo "Building cli..."
	@mkdir -p ./.bin
	@go build -o "./.bin/cli" ./src/cli
