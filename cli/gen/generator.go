package gen

import (
	"os"
	"fmt"
	"strings"
	"go/ast"

	"github.com/s0ulw1sh/soulgost/hash"
)

type cliGenerator struct {
	ctypes map[string]*clitype
}

type clitype struct {
	typename string
	funcs []clifuncs
}

func (self *clitype) addfn(fn clifuncs) {
	self.funcs = append(self.funcs, fn)
}

func (self *clitype) call_var() string {
	return self.typename + "CliStruct"
}

func (self *clitype) call_help() string {
	return self.typename + "CliHelp"
}

type clifuncs struct {
	name     string
	command  string
	typename string
	cmdhash  uint32
	ft       *ast.FuncType
	params    []cliparams
}

func (self *clifuncs) call_name() string {
	return self.typename + "CliCmd_" + self.command
}

type cliparams struct {
	name   string
	gotype string
}

func (self *cliGenerator) checkDecls(root *ast.File) bool {
	var (
		fd  *ast.FuncDecl
		se  *ast.StarExpr
		id  *ast.Ident
		rt  *ast.Ident
		ok  bool
		hw  bool
		cmd string
		tn  string
	)

	for _, d := range root.Decls {
		if fd, ok = d.(*ast.FuncDecl); !ok || fd.Recv == nil || len(fd.Recv.List) != 1 { continue }
		if se, ok = fd.Recv.List[0].Type.(*ast.StarExpr); !ok { continue }
		if id, ok = se.X.(*ast.Ident); !ok { continue }
		if len(fd.Type.Results.List) != 1 { continue }
		if rt, ok = fd.Type.Results.List[0].Type.(*ast.Ident); !ok || rt.Name != "error" { continue }
		if !strings.HasPrefix(fd.Name.Name, "Cli") { continue }

		tn = id.Name

		if _, ok = self.ctypes[tn]; !ok {
			self.ctypes[tn] = &clitype{
				typename: tn,
				funcs:    make([]clifuncs, 0),
			}
		}

		hw = true

		cmd = strings.ToLower(strings.TrimPrefix(fd.Name.Name, "Cli"))

		clf := clifuncs{
			name:     fd.Name.Name,
			command:  cmd,
			typename: tn,
			cmdhash:  hash.MurMur2([]byte(cmd)),
			ft:       fd.Type,
		}

		if fd.Type.Params != nil {
			for _, l := range fd.Type.Params.List {

				if id, ok = l.Type.(*ast.Ident); !ok { continue }

				for _, n := range l.Names {
					clf.params = append(clf.params, cliparams{
						name:   n.Name,
						gotype: id.Name,
					})
				}
			}
		}

		if entry, ok := self.ctypes[tn]; ok {
			entry.addfn(clf)
		}
	}

	return hw
}


func (self *cliGenerator) genCliTop(fw *os.File, t *clitype) {
	fw.WriteString("type "+t.typename+"CliCmdFunc func([]string) error\n\n")
	fw.WriteString("var "+t.call_var()+" "+t.typename+"\n\n")
}

func (self *cliGenerator) genCliHelp(fw *os.File, t *clitype) {
	fw.WriteString("func " + t.call_help() + "() {\n\t")

	fw.WriteString("fmt.Println(\"List of commands:\")\n\t")
	for _, fn := range t.funcs {
		fw.WriteString("fmt.Println(\"\\t"+fn.command)
		for _, p := range fn.params {
			fw.WriteString(" "+p.name+":<"+p.gotype+">")
		}
		fw.WriteString("\")\n\t")
	}

	fw.WriteString("fmt.Println(\"\")\n")

	fw.WriteString("}\n\n")
}

func (self *cliGenerator) genCliCallerFunc(fw *os.File, t *clitype, fn *clifuncs) {
	var (
		parr []string
		iserr bool
	)
	fw.WriteString("func " + fn.call_name() + "(s_args []string) error {\n\t")
	
	for i, p := range fn.params {
		if p.gotype == "string" {
			parr = append(parr, fmt.Sprintf("s_args[%d]", i))
		} else {
			iserr = true
			switch p.gotype {
			case "int8", "int16", "int32", "int64":
				fw.WriteString(fmt.Sprintf("var p%d int64\n\t", i+1))
			case "uint8", "uint16", "uint32", "uint64":
				fw.WriteString(fmt.Sprintf("var p%d uint64\n\t", i+1))
			case "float32", "float64":
				fw.WriteString(fmt.Sprintf("var p%d float64\n\t", i+1))
			case "bool":
				fw.WriteString(fmt.Sprintf("var p%d bool\n\t", i+1))
			}
		}
	}

	if iserr {
		fw.WriteString("var err error\n\t")
	}

	fw.WriteString(fmt.Sprintf("if len(s_args) != %d {\n\t\t", len(fn.params)))
	fw.WriteString("fmt.Println(`Error: incorrect number of parameters`)\n\t\t")
	fw.WriteString("fmt.Println(`Example: "+fn.command)
	for _, p := range fn.params {
		fw.WriteString(" "+p.name+":<"+p.gotype+">")
	}
	fw.WriteString("`)\n\t}\n\t")

	for i, p := range fn.params {
		switch p.gotype {
		case "int8", "int16", "int32", "int64":
			parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
			fw.WriteString(fmt.Sprintf("p%d, err = strconv.ParseInt(s_args[%d], 10, %s)\n\t", i+1, i, strings.TrimPrefix(p.gotype, "int")))
			fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
		case "uint8", "uint16", "uint32", "uint64":
			parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
			fw.WriteString(fmt.Sprintf("p%d, err = strconv.ParseUint(s_args[%d], 10, %s)\n\t", i+1, i, strings.TrimPrefix(p.gotype, "uint")))
			fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
		case "float32", "float64":
			parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
			fw.WriteString(fmt.Sprintf("p%d, err = strconv.ParseFloat(s_args[%d], %s)\n\t", i+1, i, strings.TrimPrefix(p.gotype, "float")))
			fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
		case "bool":
			parr = append(parr, fmt.Sprintf("p%d", i+1))
			fw.WriteString(fmt.Sprintf("p%d, err = strconv.ParseBool(s_args[%d])\n\t", i+1, i))
			fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
		}
	}

	fw.WriteString("return " + t.call_var() + "." + fn.name + "("+strings.Join(parr, ", ")+")\n")

	fw.WriteString("}\n\n")
}

func (self *cliGenerator) genCliCaller(fw *os.File, t *clitype) {
	var (
		mapvars []string
	)

	for _, fn := range t.funcs {
		mapvars = append(mapvars, fmt.Sprintf("%d: %s", fn.cmdhash, fn.call_name()))
		self.genCliCallerFunc(fw, t, &fn)
	}

	fw.WriteString("func " + t.typename + "CliRun(s_args []string) error {\n\t")
	fw.WriteString("var cmd_hash uint32\n\t")
	fw.WriteString("var cmd_map = map[uint32]TAppBrendCliCmdFunc{\n\t\t"+strings.Join(mapvars, ",\n\t\t")+",\n\t}\n\t")
	fw.WriteString("if len(s_args) == 0 || s_args[0] == \"?\" {\n\t\t"+t.call_help()+"()\n\t\treturn nil\n\t}\n\t")
	fw.WriteString("cmd_hash = hash.MurMur2([]byte(s_args[0]))\n\t")
	fw.WriteString("if fn, ok := cmd_map[cmd_hash]; ok { return fn(s_args[1:]) }\n\t")
	fw.WriteString("fmt.Println(`Error: command \"`+s_args[0]+`\" not found`)\n\t")
	fw.WriteString(t.call_help() + "()\n\t")
	fw.WriteString("return nil\n")
	fw.WriteString("}")
}

func Generate(root *ast.File, f *os.File) bool {
	gen := cliGenerator{}

	gen.ctypes = make(map[string]*clitype)

	if !gen.checkDecls(root) {
		return false
	}

	f.WriteString("package " + root.Name.Name + "\n\n")
	f.WriteString("//-soulgost\n\n")
	f.WriteString("// WARNING!!! \n")
	f.WriteString("// Code generated by \"soulgost -modes=cli\"; DO NOT EDIT!\n")
	f.WriteString("// URL - https://github.com/s0ulw1sh/soulgost\n")
	f.WriteString("// by Pavel Rid aka s0ulw1sh\n\n")
	
	f.WriteString("import (\n")
	f.WriteString("\t\"fmt\"\n")
	f.WriteString("\t\"strconv\"\n")
	f.WriteString("\t\"github.com/s0ulw1sh/soulgost/hash\"\n")
	f.WriteString(")\n\n")

	for _, t := range gen.ctypes {
		gen.genCliTop(f, t)
		gen.genCliHelp(f, t)
		gen.genCliCaller(f, t)
	}

	return true
}