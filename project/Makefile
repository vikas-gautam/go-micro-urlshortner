

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"




# ## start: starts the front end
# start:
# 	@echo "Starting front end"
# 	cd ../front-end && ./${FRONT_END_BINARY} &

# ## stop: stop the front end
# stop:
# 	@echo "Stopping front end..."
# 	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
# 	@echo "Stopped front end!"