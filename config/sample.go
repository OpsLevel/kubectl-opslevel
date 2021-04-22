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
      language: .metadata.annotations."opslevel.com/languague"
      framework: .metadata.annotations."opslevel.com/framework"
      aliases:
      - '"\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - .metadata.labels
`