package main

import (
	"flag"

	"github.com/hown3d/terraform-iam-generator/internal/core"
)

var directory *string = flag.String("dir", "", "terraform directory to use")

func main() {
	flag.Parse()
	core.Run(*directory)
}
