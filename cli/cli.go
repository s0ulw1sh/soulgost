package cli

import (
	"fmt"
)

var (
	clis = make(map[string]CliType)
)

type CliType interface {
	CliRun([]string) error
}

func Register(name string, cli CliType) {
	clis[name] = cli
}

func Run(args_s []string) {
	if len(args_s) < 2 {
		return
	}

	if cli, ok := clis[args_s[0]]; ok {
		cli.CliRun(args_s[1:])
		return
	}

	if len(clis) > 0 {
		fmt.Println("command not found, available commands:")

		for name, _ := range clis {
			fmt.Println("\t" + name)
		}
	}

}