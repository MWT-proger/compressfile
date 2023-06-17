COMPOSE_FILE=docker/docker-compose.yml
ENV_FILE=.env

build:
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE)  build -d

up:
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

build_up:
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d --build

dev_up:
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up

