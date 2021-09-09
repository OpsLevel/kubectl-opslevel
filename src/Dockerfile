FROM ubuntu:impish AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/kubectl-opslevel"]
WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get purge && apt-get clean && apt-get autoclean
RUN curl -L -o /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
    chmod +x /usr/local/bin/jq
COPY kubectl-opslevel /
