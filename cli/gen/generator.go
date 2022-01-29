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

func (self *clitype) start_func() string {
	return "func (self *"+self.typename+") "
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

func (self *clifuncs) example_cmd() string {
	example := self.command

	for _, p := range self.params {
		example += " "+p.name+":<"+p.gotype+">"
	}

	return example
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
		if fd.Type.Results == nil || len(fd.Type.Results.List) != 1 { continue }
		if rt, ok = fd.Type.Results.List[0].Type.(*ast.Ident); !ok || rt.Name != "error" { continue }
		if !strings.HasPrefix(fd.Name.Name, "CliCmd") { continue }

		tn = id.Name

		if _, ok = self.ctypes[tn]; !ok {
			self.ctypes[tn] = &clitype{
				typename: tn,
				funcs:    make([]clifuncs, 0),
			}
		}

		hw = true

		cmd = strings.ToLower(strings.TrimPrefix(fd.Name.Name, "CliCmd"))

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

	fw.WriteString("fmt.Println(\"\")\n\t")
	fw.WriteString("strconv.FormatBool(true)\n")

	fw.WriteString("}\n\n")
}

func (self *cliGenerator) genCliCaller(fw *os.File, t *clitype) {

	fw.WriteString(t.start_func() + "CliRun(s_args []string) error {\n\t")
	fw.WriteString("if len(s_args) == 0 || s_args[0] == \"?\" {\n\t\t"+t.call_help()+"()\n\t\treturn nil\n\t}\n\t")

	fw.WriteString("s_args_cmd := strings.ToLower(strings.Trim(s_args[0]))\n\t")
	fw.WriteString("s_args = s_args[1:]\n\t")
	fw.WriteString("switch hash.MurMur2([]byte(s_args_cmd)) {\n\t")
	for _, fn := range t.funcs {
		var parr []string
		fw.WriteString(fmt.Sprintf("case %d:\n\t\t", fn.cmdhash))

		fw.WriteString(fmt.Sprintf("if len(s_args) != %d {\n\t\t\t", len(fn.params)))
		fw.WriteString("fmt.Println(\"Error: incorrect number of parameters\")\n\t\t\t")
		fw.WriteString("fmt.Println(\"Example: "+fn.example_cmd() + "\")\n\t\t}\n\t\t")

		for i, p := range fn.params {
			var iserr = false
			if p.gotype == "string" {
				parr = append(parr, fmt.Sprintf("s_args[%d]", i))
			} else {
				iserr = true
				switch p.gotype {
				case "int8", "int16", "int32", "int64", "int":
					fw.WriteString(fmt.Sprintf("p%d, err := strconv.ParseInt(s_args[%d], 10, %s)\n\t\t", i+1, i, strings.TrimPrefix(p.gotype, "int")))
					parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
				case "uint8", "uint16", "uint32", "uint64", "uint":
					fw.WriteString(fmt.Sprintf("p%d, err := strconv.ParseUint(s_args[%d], 10, %s)\n\t\t", i+1, i, strings.TrimPrefix(p.gotype, "uint")))
					parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
				case "float32", "float64":
					fw.WriteString(fmt.Sprintf("p%d, err := strconv.ParseFloat(s_args[%d], %s)\n\t\t", i+1, i, strings.TrimPrefix(p.gotype, "float")))
					parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
				case "bool":
					fw.WriteString(fmt.Sprintf("p%d, err := strconv.ParseBool(s_args[%d])\n\t\t", i+1, i))
					parr = append(parr, fmt.Sprintf(p.gotype+"(p%d)", i+1))
				}
			}

			if iserr {
				fw.WriteString("if err != nil {\n\t\t\t")
				fw.WriteString("fmt.Println(\"CLI error `"+p.name+"` parameter must be "+p.gotype+"\")\n\t\t\t")
				fw.WriteString("return err\n\t\t")
				fw.WriteString("}\n\t\t")
			}
		}

		fw.WriteString("return " + "self." + fn.name + "("+strings.Join(parr, ", ")+")\n\t")
	}

	fw.WriteString("default:\n\t\t")
	fw.WriteString("fmt.Println(`Error: command \"`+s_args_cmd+`\" not found`)\n\t\t")
	fw.WriteString(t.call_help() + "()\n\t")
	fw.WriteString("}\n\t")
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
	f.WriteString("\t\"strings\"\n")
	f.WriteString("\t\"strconv\"\n")
	f.WriteString("\t\"github.com/s0ulw1sh/soulgost/hash\"\n")
	f.WriteString(")\n\n")

	for _, t := range gen.ctypes {
		gen.genCliHelp(f, t)
		gen.genCliCaller(f, t)
	}

	return true
}