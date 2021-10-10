package main

import (
	"github.com/faradey/madock/src/cli/commands"
	"github.com/faradey/madock/src/cli/helper"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		command := strings.ToLower(os.Args[1])
		flag := ""
		if len(os.Args) > 2 {
			flag = strings.ToLower(os.Args[2])
		}

		switch command {
		case "setup":
			commands.Setup()
		case "start":
			commands.Start()
		case "stop":
			commands.Stop(flag)
		case "restart":
			commands.Restart()
		case "rebuild":
			commands.Rebuild()
		case "magento":
			flag = strings.Join(os.Args[2:], " ")
			commands.Magento(flag)
		case "composer":
			flag = strings.Join(os.Args[2:], " ")
			commands.Composer(flag)
		case "db":
			option := ""
			if len(os.Args) > 3 {
				option = strings.ToLower(os.Args[3])
			}
			commands.DB(flag, option)
		case "cron":
			commands.Cron(flag)
		case "debug":
			commands.Debug(flag)
		case "bash":
			flag2 := ""
			if len(os.Args) > 3 {
				flag2 = strings.ToLower(os.Args[3])
			}
			commands.Bash(flag, flag2)
		case "help":
			helper.Help()
		case "logs":
			helper.Help()
		case "add":
			flags := ""
			if len(os.Args) > 3 {
				flags = strings.ToLower(os.Args[3])
			}
			commands.Add(flag, flags)
		default:
			commands.IsNotDefine()
		}
	} else {
		helper.Help()
	}
}
