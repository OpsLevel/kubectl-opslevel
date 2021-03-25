FROM golang:1.16 AS builder
LABEL stage=builder
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./kubectl-opslevel


FROM golang:1.16 AS release
ENV USER_UID=1001 USER_NAME=opslevel
COPY --from=builder /workspace/kubectl-opslevel /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]