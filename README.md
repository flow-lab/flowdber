# FlowDBer [![Build and Publish on Release Tag](https://github.com/flow-lab/flowdber/actions/workflows/docker-release.yml/badge.svg)](https://github.com/flow-lab/flowdber/actions/workflows/docker-release.yml)

FlowDBer project includes a helper container that makes it easy to run SQL migration scripts using PostgreSQL on
Kubernetes. To use it, simply include it as an InitContainer in your deployment definition, following the instructions
in the Kubernetes documentation.

Check [minikube.yml](./minikube.yml) for example configuration.

### Requirements

- file name should be **UNIXTIME**-**WHAT**.sql, eg: 1580247785-user-table.sql or just **INDEX**-**WHAT**.sql, eg:
  0-new-user-db.sql, 1-birthdate-column-in-user.sql etc
- sql scripts must be idempotent, which means that they can be run multiple times without causing any problems :)

## Project structure

- [./cmd](./cmd) - app implementation
- [./internal](./internal) - internal packages that should not be shared with other projects
- [./github](./.github) - GitHub Actions workflows
- [./certs](./certs) - certificates for TLS when running locally with docker compose
- [./db-scripts](./db-scripts) - database scripts that will be executed on startup when using docker compose
- [./docker-compose.yaml](./docker-compose.yaml) - docker compose file for running locally
- [./Makefile](./Makefile) - makefile for running commands during development
- [LICENSE](./LICENSE) - license file

## Requirements

- [golang](https://golang.org/doc/install) installation
- gui editor, e.g. [goland](https://www.jetbrains.com/go)

## DockerHub

[https://hub.docker.com/r/flowlab/flowdber](https://hub.docker.com/r/flowlab/flowdber)

## Logging

App is using [dlog](https://github.com/flow-lab/dlog) for logging. It is configured to log to stdout and stderr.

## Running locally

You can test it locally with [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/). Start minikube
and apply `kubectl apply -f minikube.yml`.

## Credits

- This project was created by cookiecutter https://github.com/flow-lab/ms-template
