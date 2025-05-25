.PHONY: build run test clean docker-build docker-run k8s-deploy k8s-clean

# Build the application
build:
	go build -o bin/main .

# Run the application locally
run:
	go run .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t collaborative-notes:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env collaborative-notes:latest

# Deploy to Kubernetes
k8s-deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/postgres.yaml
	kubectl apply -f k8s/app.yaml

# Clean Kubernetes deployment
k8s-clean:
	kubectl delete -f k8s/app.yaml
	kubectl delete -f k8s/postgres.yaml
	kubectl delete -f k8s/namespace.yaml

# Development setup
dev-setup:
	go mod tidy
	go mod download

# Database migration (if needed)
migrate:
	go run . --migrate
