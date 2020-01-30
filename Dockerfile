FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev make

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
ENV GO111MODULE=on

ARG GITHUB_SHA
ARG VERSION

RUN make test
RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -ldflags="-w -s -X 'main.version=${VERSION}' -X 'main.commit=${GITHUB_SHA}' -X 'main.date=$(date)'" -o /go/bin/app ./cmd/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]