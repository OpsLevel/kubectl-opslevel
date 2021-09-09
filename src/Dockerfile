FROM stedolan/jq AS release
ENV USER_UID=1001 USER_NAME=opslevel
ENTRYPOINT ["/kubectl-opslevel"]
WORKDIR /app
COPY kubectl-opslevel /
