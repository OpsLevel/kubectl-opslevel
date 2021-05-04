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
      - '.metadata.annotations."opslevel.com/aliases" | fromjson?'
      - .metadata.annotations."app.kubernetes.io/instance"
      - '"\(.metadata.name)-\(.metadata.namespace)"'
      tags:
      - {"imported": "kubectl-opslevel"}
      - '.metadata.annotations."opslevel.com/tags" | fromjson?'
      - .metadata.labels
      - .spec.template.metadata.labels
      tools:
      - '.metadata.annotations."opslevel.com/tools" | fromjson?'
`
