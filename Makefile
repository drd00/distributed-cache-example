# Variables
APP_NAME=distributed-cache-example
DOCKER_IMAGE=distributed-cache-example
DOCKER_TAG=latest
K8S_NAMESPACE=default

.PHONY: help build run test docker-build docker-push k8s-deploy k8s-delete clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go binary
	go build -o bin/cache ./cmd/cache

run: ## Run locally
	go run ./cmd/cache

test: ## Run tests
	go test -v ./...

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push: ## Push Docker image (update with your registry)
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) your-registry/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push your-registry/$(DOCKER_IMAGE):$(DOCKER_TAG)

k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f k8s/

k8s-delete: ## Delete from Kubernetes
	kubectl delete -f k8s/

k8s-logs: ## Tail logs from all cache pods
	kubectl logs -f -l app=cache --all-containers=true

k8s-port-forward: ## Port forward to first pod
	kubectl port-forward cache-0 8080:8080

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

.DEFAULT_GOAL := help
