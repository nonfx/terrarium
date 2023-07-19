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
	@touch ${HOME}/.netrc

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

.PHONY: test mod-tidy seed seed_resources seed_modules seed_mappings

-include .env
export

SEED_SRCS := $(shell find ./src/pkg ./src/cli \( -name '*.go' -o -name 'go.mod' \))

mod-clean:  # delete go*.sum files
	@echo "deleting .sum files..."
	@rm -f ./src/api/go.sum ./src/cli/go.sum ./src/pkg/go.sum ./go.work.sum

mod-tidy:  # run go mod tidy on each workspace entity, and then sync workspace
	@echo "running tidy on api go module..."
	@cd src/api && go mod tidy -e
	@echo "running tidy on cli go module..."
	@cd src/cli && go mod tidy -e
	@echo "running tidy on pkg go module..."
	@cd src/pkg && go mod tidy -e
	@echo "running sync on go workspace..."
	@go mod download && go work sync

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

######################################################
# Following targets are used for integration tests
######################################################

## These targets are temporary until we have ability to pass the terraform and cache directory in terrarium resource, module and mappings commands

# generate tf_resources.json file for set terraform providers
TERRAFORM_VERSION := 1.5.3
GO_VERSION := 1.20
INSTALL_DIR := /usr/local/bin
GO_INSTALL_DIR := /usr/local

install_terraform:
	@echo "Installing Terraform $(TERRAFORM_VERSION)..."
	@curl -sSL -o terraform.zip https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_linux_amd64.zip
	@unzip -o terraform.zip -d $(INSTALL_DIR)
	@rm -f terraform.zip
	@echo "Terraform $(TERRAFORM_VERSION) is installed at $(INSTALL_DIR)/terraform"
	@terraform --version

install_go:
	# @echo "Installing Go $(GO_VERSION)..."
	# @curl -sSL -o go.tar.gz https://dl.google.com/go/go$(GO_VERSION).linux-amd64.tar.gz
	# @tar -C $(GO_INSTALL_DIR) -xzf go.tar.gz
	# @rm -f go.tar.gz
	# @echo "Go $(GO_VERSION) is installed at $(GO_INSTALL_DIR)/go"
	# @export PATH=$PATH:/usr/local/go/bin:$GO_INSTALL_DIR/go/bin
	@which go
	@ls -la /usr/local/go/bin/
	@echo $PATH
	@go version

test_cache_data/tf_resources.json: $(TERRAFORM_DIR)/.terraform
	@echo "generating ./cache_data/tf_resources.json"
	@mkdir -p ./src/cli/int_test/test/cache_data
	@cd terraform && terraform version && terraform providers schema -json > ../src/cli/int_test/test/cache_data/tf_resources.json

test-int: clean_tf install_terraform tf_init test_cache_data/tf_resources.json install_go
	 	@go test ./src/cli/int_test/test -v -testify.m="${TEST_REGEX}" --tags="${TEST_TAG}" -timeout 600000s
