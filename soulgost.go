package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"strings"
	"io/ioutil"
	"path/filepath"
	"go/token"
	"go/parser"

	"github.com/s0ulw1sh/soulgost/db/gen"
)

type appSoulgost struct {
	dbGen    bool
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
		self.argPath = "."
	} else {
		self.argPath = args[0]
	}

	list := strings.Split(*modes, ",")

	for _, m := range list {
		if strings.ToLower(m) == "db" {
			self.dbGen = true
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

func (self *appSoulgost) readFile(file_path string) ([]byte, error) {
	return os.ReadFile(file_path)
}

func (self *appSoulgost) processFile(file_path string) {
	data, err := self.readFile(file_path)

	if err != nil {
		log.Fatal(err)
	}

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
		gen.Generate(astf, fdir, fname)
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