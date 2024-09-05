FROM flant/jq:b6be13d5-glibc AS libjq

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
RUN GOOS=linux go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o kubectl-opslevel main.go &&\
    chmod +x /src/kubectl-opslevel

FROM gcr.io/distroless/base-debian12 AS release
WORKDIR /
COPY --from=builder /src/kubectl-opslevel /kubectl-opslevel
ENTRYPOINT ["/kubectl-opslevel"]
