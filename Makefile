GOPATH := $(shell go env GOPATH|cut -d ":" -f 1)

# Define variables for PostgreSQL container
POSTGRES_CONTAINER := postgres
POSTGRES_DB := cc_terrarium
POSTGRES_USER := postgres

# Define variables for pg_dump command
DUMP_DIR := ./data

# Define phony targets (targets that don't correspond to files)
.PHONY: db-dump docker-run docker-stop docker-stop-clean seed test

db-dump:  ## Target for dumping PostgreSQL database to a file
	docker compose exec -T $(POSTGRES_CONTAINER) pg_dump -U $(POSTGRES_USER) -C $(POSTGRES_DB) | dos2unix > $(DUMP_DIR)/$(POSTGRES_DB).sql

docker-run:  ## Starts app in docker containers
	docker compose up -d

docker-stop:  ## Stops and removes docker containers
	docker compose down

docker-stop-clean:  ## Stops and removes containers as well as volumes to cleanup database
	docker compose down -v

test:  ## Run go unit tests
	@$(GOPATH)/bin/godotenv go test `go list github.com/cldcvr/terrarium/... | grep -v /pkg/terraform-config-inspect/`

# generate tf_resources.json file for set terraform providers
cache_data/tf_resources.json: terraform/providers.tf
	@echo "generating ./cache_data/tf_resources.json"
	@mkdir -p cache_data
	@cd terraform && terraform init && terraform providers schema -json > ../cache_data/tf_resources.json

seed_resources: cache_data/tf_resources.json docker-run  ## Load .env file and run seed_resources
	@echo "Running resource seed..."
	@$(GOPATH)/bin/godotenv go run ./api/cmd/seed_resources

seed_modules: docker-run  ## Load .env file and run seed_modules
	@echo "Running module seed..."
	@$(GOPATH)/bin/godotenv go run ./api/cmd/seed_modules

seed: seed_resources seed_modules
