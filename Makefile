# Define variables for PostgreSQL container
POSTGRES_CONTAINER := postgres
POSTGRES_DB := cc_terrarium
POSTGRES_USER := postgres

# Define variables for pg_dump command
DUMP_DIR := ./data

# Define phony targets (targets that don't correspond to files)
.PHONY: db-dump docker-run docker-stop docker-stop-clean

db-dump:  ## Target for dumping PostgreSQL database to a file
	docker compose exec -T $(POSTGRES_CONTAINER) pg_dump -U $(POSTGRES_USER) -C $(POSTGRES_DB) | dos2unix > $(DUMP_DIR)/$(POSTGRES_DB).sql

docker-run:  ## Starts app in docker containers
	docker compose up -d

docker-stop:  ## Stops and removes docker containers
	docker compose down

docker-stop-clean:  ## Stops and removes containers as well as volumes to cleanup database
	docker compose down -v
