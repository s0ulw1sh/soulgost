package cli

import (
	"os"
	"log"
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
	if len(clis) > 0 {
		if cli, ok := clis[args_s[0]]; ok && len(args_s) >= 2 {
			err := cli.CliRun(args_s[1:])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("command not found, available commands:")
			for name, _ := range clis {
				fmt.Println("\t" + name)
			}
		}
	} else {
		fmt.Println("invalid command line params")
	}

	os.Exit(0)
}