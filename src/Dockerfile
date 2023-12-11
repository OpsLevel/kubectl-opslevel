FROM flant/jq:b6be13d5-glibc as libjq

FROM golang:1.21 AS builder
ENV CGO_ENABLED=1 CGO_CFLAGS="-I/libjq/include" CGO_LDFLAGS="-L/libjq/lib"
WORKDIR /src

COPY --from=libjq /libjq /libjq

# Cache-friendly download of go dependencies.
COPY go.mod go.sum ./
RUN go mod download

COPY . /src

FROM builder as tester

RUN go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

FROM builder as complier

RUN GOOS=linux \
    go build -ldflags="-s -w" -o kubectl-opslevel main.go

FROM golang:1.21 as release
COPY --from=complier /src/kubectl-opslevel /usr/local/bin
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]
