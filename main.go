package main

import (
	"os"

	"github.com/joshvanl/cert-managerctl/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
