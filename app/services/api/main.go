package main

import (
	"os"

	"github.com/nhaancs/bhms/app/services/api/v1/cmd"
	"github.com/nhaancs/bhms/app/services/api/v1/cmd/all"
)

var build = "develop"

func main() {
	if err := cmd.Main(build, all.Routes()); err != nil {
		os.Exit(1)
	}
}
