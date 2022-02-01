CONTAINER ?= goapp
DOCKER_FILE = docker/docker-compose.yml
DOCKER_COMPOSE = docker-compose -f $(DOCKER_FILE)

# DOCKER TASKS

# Build the container
build: ## Build the release and develoment container. The development
	$(DOCKER_COMPOSE) build

exec: ## Execute the container
	docker exec -it $(CONTAINER) $(COMMAND)

# Execute goapp test inside the container
test: exec 
 	COMMAND=go test -v ./test

# Build and run the container
up: ## Spin up the project
	$(DOCKER_COMPOSE) up

down: ## Stop running containers
	$(DOCKER_COMPOSE) down

rm: stop ## Stop and remove running containers
	docker rm $(APP_NAME)
