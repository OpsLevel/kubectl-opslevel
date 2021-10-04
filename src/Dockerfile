FROM ubuntu:focal AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/opslevel"]
WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get purge && apt-get clean && apt-get autoclean && \
    curl -o /usr/local/bin/jq http://stedolan.github.io/jq/download/linux64/jq && \
    chmod +x /usr/local/bin/jq
COPY kubectl-opslevel /usr/local/bin