FROM golang:1.16 AS builder
ARG VERSION=development
LABEL stage=builder
WORKDIR /workspace
COPY ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download
COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./kubectl-opslevel -ldflags="-X 'github.com/opslevel/kubectl-opslevel/cmd.version=${VERSION}'"


FROM ubuntu:impish AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]
WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get purge && apt-get clean && apt-get autoclean && \
    curl -L -o /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
    chmod +x /usr/local/bin/jq
COPY --from=builder /workspace/kubectl-opslevel /usr/local/bin/

