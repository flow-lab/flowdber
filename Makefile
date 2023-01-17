SHELL := /bin/bash

SRV_NAME := flow-k8-sql
PROJECT := diatom-ai
HOSTNAME := eu.gcr.io
DOCKER_IMG := flowlab/${SRV_NAME}

deps:
	go get -u -t ./...

deps-reset:
	git checkout -- go.mod

tidy:
	go mod tidy

verify:
	go mod verify

test:
	go test -mod=readonly -covermode=atomic -v ./...

docker-build:
	docker build -t ${DOCKER_IMG}:latest .

docker-tag:
	docker tag ${DOCKER_IMG} ${HOSTNAME}/${PROJECT}/${SRV_NAME}

docker-push:
	gcloud docker -- push ${HOSTNAME}/${PROJECT}/${SRV_NAME}

docker-clean:
	docker system prune -f

# minikube
minikube-init:
	eval $(minikube docker-env)

minikube-deploy:
	kubectl apply -f minikube.yml

minikube-get-pod:
	kubectl get pod

