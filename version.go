package main

import "fmt"

var (
	Version     string
	BuildNumber string
)

func getVersion() string {
	return fmt.Sprintf("Build version: %s, Build number: %s", Version, BuildNumber)
}
