//go:build release

package main

func init() {
	BuildMode = "Release"
	OpenApiSpecFile = "/usr/share/dnsmasq-manager/spec/openapi.yaml"
}
