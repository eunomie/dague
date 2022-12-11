package config

var (
	// BuildImage is the base Golang docker image used to run all the tools.
	BuildImage = "golang:1.19.4-alpine3.17"
	// AppDir is the default local to the container folder where to copy/mount sources.
	AppDir = "/go/src"
)
