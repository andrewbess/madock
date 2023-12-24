package config

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/faradey/madock/src/helper/cli/attr"
	configs2 "github.com/faradey/madock/src/helper/configs"
	"github.com/faradey/madock/src/helper/paths"
	"log"
	"os"
	"strings"
)

type ArgsStruct struct {
	attr.Arguments
	Name  string `arg:"-n,--name" help:"Parameter name"`
	Value string `arg:"-v,--value" help:"Parameter value"`
}

func ShowEnv() {
	lines := configs2.GetProjectConfig(configs2.GetProjectName())
	for key, line := range lines {
		fmt.Println(key + " " + line)
	}
}

func SetEnvOption() {
	args := getArgs()
	name := strings.ToUpper(args.Name)
	val := args.Value
	if len(name) > 0 && configs2.IsOption(name) {
		configPath := paths.GetExecDirPath() + "/projects/" + configs2.GetProjectName() + "/config.xml"
		configs2.SetParam(configPath, name, val)
	}
}

func getArgs() *ArgsStruct {
	args := new(ArgsStruct)
	if attr.IsParseArgs && len(os.Args) > 2 {
		argsOrigin := os.Args[2:]
		p, err := arg.NewParser(arg.Config{
			IgnoreEnv: true,
		}, args)

		if err != nil {
			log.Fatal(err)
		}

		err = p.Parse(argsOrigin)

		if err != nil {
			log.Fatal(err)
		}
	}

	attr.IsParseArgs = false
	return args
}
