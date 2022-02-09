package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"flag"
	"strings"
	"io/ioutil"
	"path/filepath"
	"go/token"
	"go/parser"

	db_gen  "github.com/s0ulw1sh/soulgost/db/gen"
	cli_gen "github.com/s0ulw1sh/soulgost/cli/gen"
	api_gen "github.com/s0ulw1sh/soulgost/api/gen"
)

type appSoulgost struct {
	dbGen    bool
	cliGen   bool
	apiGen   bool
	argPath  string
	fset     *token.FileSet
}

func (self *appSoulgost) init() {
	log.SetFlags(0)
	log.SetPrefix("soulgost: ")

	self.fset = token.NewFileSet()
}

func (self *appSoulgost) flagUsage() {
	fmt.Fprintf(os.Stderr, "Usage of soulgost:\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func (self *appSoulgost) flagParse() {
	modes     := flag.String("modes", "", "operating modes; db")
	
	flag.Usage = self.flagUsage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		wd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		self.argPath = wd
	} else {
		self.argPath = filepath.Clean(args[0])

		abs, err := filepath.Abs(self.argPath)

		if err != nil {
			log.Fatal(err)
		}

		self.argPath = abs
	}

	list := strings.Split(*modes, ",")

	for _, m := range list {
		switch strings.ToLower(m) {
		case "db":  self.dbGen  = true
		case "cli": self.cliGen = true
		case "api": self.apiGen = true
		}
	}
}

func (self *appSoulgost) isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func (self *appSoulgost) readFile(file_path string) []byte {
	data, err := os.ReadFile(file_path)
	
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (self *appSoulgost) moveFile(src, dest string) error {
	i, err := os.Open(src)
    if err != nil {
        return err
    }
    o, err := os.Create(dest)
    if err != nil {
        i.Close()
        return err
    }
    defer o.Close()
    _, err = io.Copy(o, i)
    i.Close()
    return err
}

func (self *appSoulgost) processFile(file_path string) {
	data := self.readFile(file_path)

	astf, err := parser.ParseFile(self.fset, file_path, data, parser.ParseComments)

	if err != nil {
		log.Fatal(err)
	}

	for _, cmng := range astf.Comments {
		for _, l := range cmng.List {
			if strings.HasPrefix(l.Text, "//-soulgost") ||
			   strings.HasPrefix(l.Text, "//+build exclude") {
				return
			}
		}
	}

	fdir  := filepath.Dir(file_path)
	fname := strings.TrimSuffix(filepath.Base(file_path), ".go")

	if self.dbGen {
		f, err := os.CreateTemp("", "soulgost-db")
		defer os.Remove(f.Name())

		if err != nil {
			log.Fatal(err)
		}

		if db_gen.Generate(astf, f) {
			err = self.moveFile(f.Name(), filepath.Join(fdir, fname + "_sgdb.go"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if self.cliGen {
		f, err := os.CreateTemp("", "soulgost-cli")
		defer os.Remove(f.Name())

		if err != nil {
			log.Fatal(err)
		}

		if cli_gen.Generate(astf, f) {
			err = self.moveFile(f.Name(), filepath.Join(fdir, fname + "_sgcli.go"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if self.apiGen {
		f, err := os.CreateTemp("", "soulgost-api")
		defer os.Remove(f.Name())

		if err != nil {
			log.Fatal(err)
		}

		if api_gen.Generate(astf, f) {
			err = self.moveFile(f.Name(), filepath.Join(fdir, fname + "_sgapi.go"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (self *appSoulgost) run() {
	if !self.isDir(self.argPath) {
		self.processFile(self.argPath)
		return
	}

	list, err := ioutil.ReadDir(self.argPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range list {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}

		fp := filepath.Join(self.argPath, f.Name())
		self.processFile(fp)
	}
}

func main() {
	app := appSoulgost{}

	app.init()

	app.flagParse()

	app.run()
}