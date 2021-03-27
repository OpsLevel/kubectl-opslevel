package config

var ConfigSample = `#Sample Opslevel CLI Config
service:
  import:
    - kind: deployment
      ol_service_name: "{$.metadata.name}"
`