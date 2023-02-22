SHELL := /bin/bash

SRV_NAME := flowdber
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

docker-compose-up:
	docker compose up --force-recreate --build

# minikube
minikube-start:
	minikube start --cpus 4 --memory 4096 --disk-size 10g --kubernetes-version=v1.26.1

minikube-stop:
	minikube stop

minikube-build:
	eval $$(minikube docker-env) && make docker-build

minikube-deploy:
	kubectl apply -f minikube.yml

minikube-get-pod:
	kubectl get pod # or just use k9s FTW

