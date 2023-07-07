//go:build release

package main

func init() {
	BuildMode = "Release"
	OpenApiSpecFile = "/usr/share/pi-hole-manager/spec/openapi.yaml"
}
