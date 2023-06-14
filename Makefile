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
.PHONY: docker-init db-dump docker-build docker-run start-db docker-stop docker-stop-clean

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
	@touch $(HOME)/.gitconfig && touch $(HOME)/.netrc
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

$(TERRAFORM_DIR)/.terraform: $(TF_FILES)
	@cd $(TERRAFORM_DIR) && terraform version && terraform init && terraform providers
	@touch $(TERRAFORM_DIR)/.terraform

clean_tf:
	rm -rf $(TERRAFORM_DIR)/.terraform
	rm -f $(TERRAFORM_DIR)/.terraform.lock.hcl

tf_init: $(TERRAFORM_DIR)/.terraform

# generate tf_resources.json file for set terraform providers
cache_data/tf_resources.json: $(TERRAFORM_DIR)/.terraform
	@echo "generating ./cache_data/tf_resources.json"
	@mkdir -p cache_data
	@cd terraform && terraform version && terraform providers schema -json > ../cache_data/tf_resources.json

######################################################
# Following targets need Go installed on the system
######################################################

.PHONY: test seed seed_resources seed_mappings seed_modules

-include .env
export

test:  ## Run go unit tests
	go test `go list github.com/cldcvr/terrarium/... | grep -v /pkg/terraform-config-inspect/`

seed: seed_resources seed_modules seed_mappings

SOURCES := $(shell find ./api/ -name '*.go')

.bin/seed_resources: $(SOURCES)
	@echo "building seed_resources"
	@mkdir -p ./.bin
	@cd ./api/cmd/seed_resources && go build -o "../../../.bin"

.bin/seed_modules: $(SOURCES)
	@echo "building seed_modules"
	@mkdir -p ./.bin
	@cd ./api/cmd/seed_modules && go build -o "../../../.bin"

.bin/seed_mappings: $(SOURCES)
	@echo "building seed_mappings"
	@mkdir -p ./.bin
	@cd ./api/cmd/seed_mappings && go build -o "../../../.bin"

seed_resources: .bin/seed_resources cache_data/tf_resources.json  ## Seed tf-provider resources into db from terraform/provider.tf
	@echo "Running resource seed..."
	@./.bin/seed_resources

seed_modules: .bin/seed_modules $(TERRAFORM_DIR)/.terraform  ## Seed tf-modules into db from terraform/modules.tf
	@echo "Running module seed..."
	@./.bin/seed_modules

seed_mappings: .bin/seed_mappings  ## Load .env file and run seed_mappings
	@echo "Running mapping seed..."
	@./.bin/seed_mappings

-include scripts/mocks.mak
-include scripts/protoc.mak
