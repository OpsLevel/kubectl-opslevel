FROM ubuntu:focal AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/usr/local/bin/kubectl-opslevel"]
WORKDIR /app
RUN apt-get update && \
    apt-get install -y jq && \
    apt-get purge && apt-get clean && apt-get autoclean
COPY kubectl-opslevel /usr/local/bin
