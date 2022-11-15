<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# config

```go
import "github.com/eunomie/dague/config"
```

## Index

- [Variables](<#variables>)


## Variables

```go
var (
    // BuildImage is the base Golang docker image used to run all the tools.
    BuildImage = "golang:1.19.3-alpine3.16"
    // AppDir is the default local to the container folder where to copy/mount sources.
    AppDir = "/go/src"
)
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->