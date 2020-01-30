# Flow k8 sql ![Build and Publish](https://github.com/flow-lab/flow-k8-sql/workflows/Build%20and%20Publish/badge.svg?branch=master)

This project was created by cookiecutter https://github.com/flow-lab/ms-template.

It is a helper container for running SQL migration scripts with _PostgreSQL_ on _Kubernetes_. It should be used as 
[InitContainer](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) in deployment definition.

Check [minikube.yml](./minikube.yml) for example configuration.

### Main functionality:

- read all sql files in given directory and execute based on timestamp (from oldest to newest)
- postgres db support with TLS
- golang implementation with modules
- kubernetes support
- docker size ~9MB

### Requirements

- file name should be **UNIXTIME**-**WHAT**.sql, eg: 1580247785-user-table.sql or just **INDEX**-**WHAT**.sql, eg:
  0-new-user-db.sql, 1-birthdate-column-in-user.sql etc  
- sql scripts must be idempotent
- set up all required env variables for db connection and db scripts, check [cmd/main.go](./cmd/main.go) for more details

Check Makefile for all important stuff.

## Project structure

- [./cmd/](./cmd/) - app implementation
- [./internal](./internal) - internal packages that should not be shared with other projects
- [./github](./.github/) - GitHub Actions workflows

## Requirements

- [golang](https://golang.org/doc/install) installation
- gui editor, eg [goland](https://www.jetbrains.com/go)

### GitHub Actions

Project is using GitHub Actions for deployment. Workflows are located in [./github/workflows](./github/workflows), where:

- google.yml - tests, builds and deploys to docker GCR and DockerHub repositories

### Secrets

Projects requires secrets for GitHub Actions. Secrets should be located in GitHub project secrets.

- GKE_PROJECT - Google Cloud project that cluster is located
- GKE_EMAIL - cluster email
- GKE_KEY - base64 encoded service account key that has access to deploy to Docker registry
- DOCKERHUB_TOKEN - DockerHub access token

## Logging

Project is using standard logger from `log` library. It is configured in `main.go` and should be used in all logging
statements. Log is in format like: 

**NAME** : (**VERSION**, **SHA**) : **DATE-TIME** **FILE**: **MSG**

where:

- **NAME**: microservice name
- **VERSION**: [semver](https://semver.org/) version taken from annotated tag, `dev` otherwise
- **SHA**: git SHA in short version
- **DATE-TIME**: date time with microseconds
- **FILE**: source file name and line information
- **MSG**: log message


## Running locally

You can test it locally with [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/). Start minikube
 and apply `kubectl apply -f minikube.yml`.