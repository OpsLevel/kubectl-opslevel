# opslevel-jq-parser
A jq wrapper which aids in converting data to opslevel-go input structures

This library leverages https://github.com/flant/libjq-go which are CGO bindings to the JQ library which provide C native speed

#  Installation

```bash
go get github.com/opslevel/opslevel-jq-parser/v2024
```

Then wherever you compile or test that project you'll need to add

```bash
docker run --name "libjq" -d flant/jq:b6be13d5-glibc
docker cp libjq:/libjq ./libjq 
docker rm libjq
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/libjq/include"
export CGO_LDFLAGS="-L$(pwd)/libjq/lib"
```

Here is a nice stanza you can put into your github actions workflow files

> NOTE: the version is important - please see https://github.com/flant/libjq-go#notes

```yaml
      - name: Setup LibJQ
        run: |-
          docker run --name "libjq" -d flant/jq:b6be13d5-glibc
          docker cp libjq:/libjq ./libjq 
          docker rm libjq
          echo CGO_ENABLED=1 >> $GITHUB_ENV
          echo CGO_CFLAGS="-I$(pwd)/libjq/include" >> $GITHUB_ENV
          echo CGO_LDFLAGS="-L$(pwd)/libjq/lib" >> $GITHUB_ENV
```