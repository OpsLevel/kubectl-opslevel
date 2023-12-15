FROM flant/jq:b6be13d5-glibc as libjq

FROM golang:1.21 AS builder
ENV CGO_ENABLED=1 CGO_CFLAGS="-I/libjq/include" CGO_LDFLAGS="-L/libjq/lib"
WORKDIR /app

COPY --from=libjq /libjq /libjq

# Cache-friendly download of go dependencies.
COPY ./src/go.mod ./src/go.sum /app/src/
RUN cd ./src && go mod download

COPY . /app

FROM builder as tester
WORKDIR /app/src

RUN go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

FROM builder as complier
WORKDIR /app/src

RUN go install github.com/goreleaser/goreleaser@v1.22.1

ARG RELEASE_VERSION
ARG GITHUB_TOKEN
ARG ORG_GITHUB_TOKEN
ARG GPG_FINGERPRINT
ENV GITHUB_TOKEN=${GITHUB_TOKEN} ORG_GITHUB_TOKEN=${ORG_GITHUB_TOKEN} GPG_FINGERPRINT=${GPG_FINGERPRINT}

RUN goreleaser release --clean --release-notes=../.changes/${RELEASE_VERSION}.md
