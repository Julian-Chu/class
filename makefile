SHELL := /bin/bash

build:
	docker build \
		-f zarf/docker/dockerfile \
		-t class-sales-api-amd64:1.0 \
		.

run:
	go run app/sales-api/main.go

admin:
	go run app/admin/main.go

kind-up:
	kind create cluster --name class-cluster --config zarf/k8s/dev/kind-config.yaml

kind-load:
	kind load docker-image class-sales-api-amd64:1.0 --name class-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -

kind-status:
	kubectl get nodes
	kubectl get pods --watch

kind-status-full:
	kubectl describe pod -lapp=sales-api

kind-logs:
	kubectl logs -lapp=sales-api --all-containers=true -f

kind-update: build
	kind load docker-image class-sales-api-amd64:1.0 --name class-cluster
	kubectl delete pods -lapp=sales-api