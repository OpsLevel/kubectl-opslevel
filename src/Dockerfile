FROM flant/jq:b6be13d5-glibc as libjq

FROM golang:1.22 AS builder
ENV CGO_ENABLED=1 CGO_CFLAGS="-I/libjq/include" CGO_LDFLAGS="-L/libjq/lib"
WORKDIR /src

COPY --from=libjq /libjq /libjq

# Cache-friendly download of go dependencies.
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . /src

ARG COMMIT
ARG VERSION
RUN GOOS=linux \
  go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o kubectl-opslevel main.go

FROM golang:1.22 as release
COPY --from=builder /src/kubectl-opslevel /usr/local/bin
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]
