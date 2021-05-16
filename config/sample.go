package config

var ConfigSample = `#Sample Opslevel CLI Config
service:
  import:
  - selector:
      kind: deployment
      namespace: ""
      labels: {}
    opslevel:
      name: .metadata.name
      description: .metadata.annotations."opslevel.com/description"
      owner: .metadata.annotations."opslevel.com/owner"
      lifecycle: .metadata.annotations."opslevel.com/lifecycle"
      tier: .metadata.annotations."opslevel.com/tier"
      product: .metadata.annotations."opslevel.com/product"
      language: .metadata.annotations."opslevel.com/language"
      framework: .metadata.annotations."opslevel.com/framework"
      aliases:
      - '"k8s:\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - '{"imported": "kubectl-opslevel"}'
      - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tags"))) | map({(.key | split(".")[2]): .value})'
      - .metadata.labels
      - .spec.template.metadata.labels
      tools:
      - '{"category": "other", "displayName": "my-cool-tool", "url": .metadata.annotations."example.com/my-cool-tool"} | if .url then . else empty end'
      - '.metadata.annotations | to_entries |  map(select(.key | startswith("opslevel.com/tools"))) | map({"category": .key | split(".")[2], "displayName": .key | split(".")[3], "url": .value})'
`
