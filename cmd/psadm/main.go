// psadm is a tool for EC2 System Manager Parameter Store.
package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Options struct{}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	os.Exit(realmain())
}

func realmain() int {
	if _, err := parser.Parse(); err != nil {
		return 1
	}
	return 0
}
