.PHONY: build run test stop clean

build:
	@go build -o bin/api

run: build
	@echo "Checking Docker network..."
	@docker network create mongo-net || true
	@echo "Starting MongoDB container..."
	@docker run -d --name local-mongo --network mongo-net --network-alias mongo -p 27017:27017 mongo:latest 2>/dev/null || docker start local-mongo
	@echo "Starting Mongo Express Web Console..."
	@docker run -d --name mongo-express --network mongo-net -p 8081:8081 -e ME_CONFIG_MONGODB_URL=mongodb://mongo:27017 mongo-express:latest 2>/dev/null || docker start mongo-express
	@echo "Database UI ready at http://localhost:8081"
	@echo "Starting Go API..."
	@./bin/api

test:
	@go test -v ./...

stop:
	@echo "Stopping database and UI containers..."
	@docker stop local-mongo mongo-express || true

clean: stop
	@echo "Cleaning up binaries and containers..."
	@rm -rf bin/
	@docker rm local-mongo mongo-express || true
	@docker network rm mongo-net || true
