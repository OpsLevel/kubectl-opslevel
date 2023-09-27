FROM flant/jq:b6be13d5-glibc as libjq

FROM golang:1.21 AS builder
WORKDIR /src

# Cache-friendly download of go dependencies.
ADD go.mod go.sum /src/
RUN go mod download

COPY --from=libjq /libjq /libjq
COPY . /src

FROM builder as tester

RUN CGO_ENABLED=1 \
    CGO_CFLAGS="-I/libjq/include" \
    CGO_LDFLAGS="-L/libjq/lib" \
    go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

FROM builder as complier

RUN CGO_ENABLED=1 \
    CGO_CFLAGS="-I/libjq/include" \
    CGO_LDFLAGS="-L/libjq/lib" \
    GOOS=linux \
    go build -ldflags="-s -w" -o kubectl-opslevel main.go

FROM golang:1.21 as release
COPY --from=complier /src/kubectl-opslevel /usr/local/bin
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]
