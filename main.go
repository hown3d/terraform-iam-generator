package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/hown3d/terraform-iam-generator/internal/core"
	"github.com/hown3d/terraform-iam-generator/internal/terraform"
)

var workDir, _ = os.Getwd()
var directory *string = flag.String("dir", workDir, "terraform directory to use")
var tfVarFiles = core.NewSliceFlag(func(s string) (terraform.VariableFile, error) {
	return terraform.VariableFile(s), nil
})
var tfVars = core.NewSliceFlag(func(s string) (terraform.Variable, error) {
	key, val, found := strings.Cut(s, "=")
	if !found {
		return terraform.Variable{}, errors.New("Invalid specification of Terraform variable: Pass variable in the following format \"KEY=VALUE\"")
	}
	return terraform.Variable{
		Key:   key,
		Value: val,
	}, nil
})

func main() {
	flag.Var(&tfVarFiles, "tf-var-file", "Path to a terraform variables file. Must be relative to the passed directory. Can be used multiple times")
	flag.Var(&tfVars, "tf-var", "	Terraform variables to use. Specify like this: \"KEY=VALUE\". Can be used multiple times")
	flag.Parse()
	core.Run(*directory, tfVars.GetValues(), tfVarFiles.GetValues())
}
