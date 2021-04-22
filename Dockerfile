FROM golang:1.16 AS builder
ARG VERSION=development
LABEL stage=builder
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./kubectl-opslevel -ldflags="-X 'github.com/opslevel/kubectl-opslevel/cmd.version=${VERSION}'"


FROM golang:1.16 AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]
WORKDIR /
RUN curl -o /usr/local/bin/jq http://stedolan.github.io/jq/download/linux64/jq && \
  chmod +x /usr/local/bin/jq
COPY --from=builder /workspace/kubectl-opslevel /usr/local/bin/

