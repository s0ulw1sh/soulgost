package gen

import (
	"strings"
	"go/ast"
	"go/parser"
	"go/token"
)

const (
	phase_check_ignore = iota
	parse_check_calls
)

const (
	db_mode_none = iota
	db_mode_first
	db_mode_select
	db_mode_update
	db_mode_insert
	db_mode_delete
)

type tableItem struct {
	name  string
	alias string
}

type Generator struct {
	data string
	output string
	items []genQuery
	counter int
	file *ast.File
	node ast.Node
	fset *token.FileSet
	phaseprev int
	phase int
}

type genQuery struct {
	mode int
	pos token.Pos
	end token.Pos
	dbobj string
	tables []tableItem
	args int
	fields []qfiled
	conds  []qcondition
	limit  string
	offset string
}

const (
	field_mode_normal = iota
	field_mode_expr
	field_mode_scan
)

type qfiled struct {
	name  string
	value string
	ftype int
}

type qcondition struct {
	name  string
	cnd   string
	value interface{}
}

func Generate(filename string, filedata string) (string, error) {
	var (
		pos token.Pos
	)

	fset := token.NewFileSet()

	pfile, err := parser.ParseFile(fset, filename, filedata, parser.ParseComments)

	if err != nil {
		return "", err
	}

	existimport := false

	g := Generator{
		data:   filedata,
		output: "",
		file:   pfile,
		node:   pfile,
		fset:   fset,
	}

	l := len(g.data)

	for _, s := range pfile.Imports {

		if s.Path.Value == `"github.com/s0ulw1sh/soulgost/db"` {

			g.output = g.data[0:s.Pos()-1]
			pos      = s.End()-1

			existimport = true

			break
		}
	}

	if !existimport {
		return g.data, nil
	}

	ast.Walk(&g, g.node)

	for _, item := range g.items {
		g.output += g.data[pos:item.pos]

		g.output += g.generateQuery(&item)

		pos = item.end
	}

	g.output += g.data[pos:l-1]

	return g.output, nil
}

func (v *Generator) Visit(node ast.Node) (w ast.Visitor) {
	if v.inspectExpr(node) {
		return v
	} 

	return nil
}

func (v *Generator) inspectExpr(node ast.Node) bool {
	res := true

	switch n := node.(type) {

	case *ast.CallExpr:

		for _, a := range n.Args {
			v.inspectExpr(a)
		}

		res = v.visitExpr(n)

		if !res {
			v.items[v.counter].pos = n.Pos()-1
			v.items[v.counter].end = n.End()-1

			v.counter += 1
		}

	case *ast.ExprStmt:
		return v.inspectExpr(n.X)

	case *ast.AssignStmt:
		if len(n.Rhs) == 1 {
			return v.inspectExpr(n.Rhs[0])
		}

	case *ast.IfStmt:

		v.inspectExpr(n.Init)

		if b, ok := n.Cond.(*ast.BinaryExpr); ok {
			v.inspectExpr(b.X)
		}

		return false

	case *ast.ReturnStmt:
		for _, a := range n.Results {
			res = v.inspectExpr(a)
		}

	case *ast.BlockStmt:
		for _, l := range n.List {
			res = v.inspectExpr(l)
		}
	}

	return res
}

func (v *Generator) visitExpr(expr ast.Expr) bool {

	switch e := expr.(type) {
	case *ast.CallExpr:

		s, ok := e.Fun.(*ast.SelectorExpr)

		if ok {

			switch x := s.X.(type) {
			case *ast.Ident:
				if x.Name == "db" && s.Sel.Name == "Q" && len(e.Args) == 1 {

					v.items = append(v.items, genQuery{})

					dblnk := v.data[e.Args[0].Pos()-1:e.Args[0].End()-1]

					if strings.HasPrefix(dblnk, "&") {
						dblnk = dblnk[1:]
					}

					v.items[v.counter].dbobj = dblnk

					return false
				}
			default:
				r := v.visitExpr(s.X)

				if !r {
					switch s.Sel.Name {

					case "Table":
						var (
							tname  string
							talias string
						)
						
						for i, a := range e.Args {
							if l, ok := a.(*ast.BasicLit); ok && l.Kind == token.STRING {
								if i == 0 { tname  = l.Value[1:len(l.Value)-1] }
								if i == 1 { talias = l.Value[1:len(l.Value)-1] }
							}
						}

						v.items[v.counter].tables = append(v.items[v.counter].tables, tableItem{tname, talias})

					case "First":
						v.items[v.counter].mode  = db_mode_first
						fallthrough

					case "Select":

						if v.items[v.counter].mode == db_mode_none {
							v.items[v.counter].mode = db_mode_select
						}
						
						for i, a := range e.Args {
							if i == 0 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseFields(l)
								}
							}
							if i == 1 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseConds(l)
								}
							}
						}

					case "Update":
						v.items[v.counter].mode = db_mode_update

						for i, a := range e.Args {
							if i == 0 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseFields(l)
								}
							}
							if i == 1 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseConds(l)
								}
							}
						}

					case "Insert":
						v.items[v.counter].mode = db_mode_insert

						for i, a := range e.Args {
							if i == 0 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseFields(l)
								}
							}
						}

					case "Delete":
						v.items[v.counter].mode = db_mode_delete

						for i, a := range e.Args {
							if i == 0 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseConds(l)
								}
							}
						}

					case "Limit":
						for i, a := range e.Args {
							if i == 0 {
								v.items[v.counter].limit = v.data[a.Pos()-1:a.End()-1]
							}
						}

					case "Offset":
						for i, a := range e.Args {
							if i == 0 {
								v.items[v.counter].offset = v.data[a.Pos()-1:a.End()-1]
							}
						}

					case "Count":
						if len(e.Args) == 2 {
							v.items[v.counter].mode  = db_mode_first
							fldname := ""
							for i, a := range e.Args {
								if i == 0 {
									if l, ok := a.(*ast.BasicLit); ok && l.Kind == token.STRING {
										fldname = l.Value[1:len(l.Value)-1]
										v.items[v.counter].fields = append(v.items[v.counter].fields, qfiled{fldname, "COUNT(`"+fldname+"`)", field_mode_expr})
									}
								}

								if i == 1 {
									valname := v.data[a.Pos()-1:a.End()-1]
									v.items[v.counter].fields = append(v.items[v.counter].fields, qfiled{fldname, valname, field_mode_scan})
								}
							}
						}

					case "Where":
						for i, a := range e.Args {
							if i == 0 {
								if l, ok := a.(*ast.CompositeLit); ok {
									v.parseConds(l)
								}
							}
						}

					}
				}

				return r
			}

		}
	}

	return true
}

func (v *Generator) isField(cl ast.Expr) bool {
	if s, oks := cl.(*ast.SelectorExpr); oks {
		if x, okx := s.X.(*ast.Ident); okx {
			return x.Name == "db" && s.Sel.Name == "Field"
		}
	}

	return false
}

func (v *Generator) isExpr(cl ast.Expr) bool {
	if s, oks := cl.(*ast.SelectorExpr); oks {
		if x, okx := s.X.(*ast.Ident); okx {
			return x.Name == "db" && s.Sel.Name == "Expr"
		}
	}

	return false
}

func (v *Generator) isConds(cl ast.Expr) bool {
	if s, oks := cl.(*ast.SelectorExpr); oks {
		if x, okx := s.X.(*ast.Ident); okx {
			return x.Name == "db" && s.Sel.Name == "Condition"
		}
	}

	return false
}


func (v *Generator) parseFields(cl *ast.CompositeLit) bool {

	switch t := cl.Type.(type) {
	case *ast.ArrayType:
		if v.isField(t.Elt) {
			for _, e := range cl.Elts {
				if l, ok := e.(*ast.CompositeLit); ok && v.isField(l.Type) {
					var (
						fldname string
						valname string
					)

					for z, m := range l.Elts {
						if z == 0 {
							if bls, okbls := m.(*ast.BasicLit); okbls {
								if bls.Kind == token.STRING {
									fldname = bls.Value[1:len(bls.Value)-1]
								}
							}
						}

						if z == 1 {
							if cls, clsok := m.(*ast.CompositeLit); clsok && v.isExpr(cls.Type) && len(cls.Elts) == 1 {
								valname = v.data[cls.Elts[0].Pos():cls.Elts[0].End()-2]
								v.items[v.counter].fields = append(v.items[v.counter].fields, qfiled{fldname, valname, field_mode_expr})
							} else {
								valname = v.data[m.Pos()-1:m.End()-1]
								v.items[v.counter].fields = append(v.items[v.counter].fields, qfiled{fldname, valname, field_mode_normal})
							}

						}
					}

				}
			}
		}
	}

	return true
}

func (v *Generator) parseConds(cl *ast.CompositeLit) bool {

	switch t := cl.Type.(type) {
	case *ast.ArrayType:
		if v.isConds(t.Elt) {
			for _, e := range cl.Elts {
				if l, ok := e.(*ast.CompositeLit); ok && v.isConds(l.Type) {
					var (
						fldname string
						condexp string
						valname string
					)

					for z, m := range l.Elts {
						if z == 0 {
							if bls, okbls := m.(*ast.BasicLit); okbls {
								if bls.Kind == token.STRING {
									fldname = bls.Value[1:len(bls.Value)-1]
								}
							}
						}

						if z == 1 {
							if bls, okbls := m.(*ast.BasicLit); okbls {
								if bls.Kind == token.STRING {
									condexp = bls.Value[1:len(bls.Value)-1]
								}
							}
						}

						if z == 2 {
							valname = v.data[m.Pos()-1:m.End()-1]
							v.items[v.counter].conds = append(v.items[v.counter].conds, qcondition{fldname, condexp, valname})
						}
					}
				}
			}
		}
	}

	return true
}

func (v *Generator) generateQuery(item *genQuery) string {
	out := item.dbobj

	switch item.mode {
	case db_mode_first:  out += `.QueryRow("SELECT ` + v.generateQuerySelect(item) + ` FROM ` + v.generateQueryTables(item)
	case db_mode_select: out += `.Query("SELECT ` + v.generateQuerySelect(item) + ` FROM ` + v.generateQueryTables(item)
	case db_mode_update: out += `.Exec("UPDATE ` + v.generateQueryTables(item) + ` SET ` + v.generateUpdateSet(item)
	case db_mode_insert: out += `.Exec("INSERT ` + v.generateQueryTables(item) + v.generateInsertSet(item)
	case db_mode_delete: out += `.Exec("DELETE ` + v.generateQueryTables(item)
	}

	if len(item.conds) > 0 {
		out += ` WHERE ` + v.generateQueryConds(item)
	}

	if item.limit != "" {
		out += ` LIMIT ?`
	}

	if item.offset != "" {
		out += ` OFFSET ?`
	}

	qvars := v.generateQueryVars(item)

	if len(qvars) > 0 {
		out += `", ` + qvars + `)`
	} else {
		out += `")`
	}

	switch item.mode {
	case db_mode_first: out += `.Scan(` + v.generateQueryScans(item) + `)`
	}

	return out
}

func (v *Generator) generateQuerySelect(item *genQuery) string {
	var out []string

	for _, field := range item.fields {
		if field.ftype == field_mode_expr {
			out = append(out, field.value + " AS `" + field.name + "`")
		} else if field.ftype == field_mode_normal {
			out = append(out, "`" + field.name + "`")
		}
	}

	return strings.Join(out, ", ")
}

func (v *Generator) generateUpdateSet(item *genQuery) string {
	var out []string

	for _, field := range item.fields {
		if field.ftype == field_mode_expr {
			out = append(out, "`" + field.name + "`=" + field.value)
		} else {
			out = append(out, "`" + field.name + "`=?")
		}
	}

	return strings.Join(out, ", ")
}

func (v *Generator) generateInsertSet(item *genQuery) string {
	var (
		flds []string
		vars []string
	)

	for _, field := range item.fields {
		flds = append(flds, "`" + field.name + "`")
		if field.ftype == field_mode_expr {
			vars = append(vars, field.value)
		} else {
			vars = append(vars, "?")
		}
	}

	return `(` + strings.Join(flds, ", ") + `) VALUES (` + strings.Join(vars, ", ") + `)`
}

func (v *Generator) generateQueryTables(item *genQuery) string {
	var out []string

	for _, table := range item.tables {

		if item.mode == db_mode_first || item.mode == db_mode_select {
			if table.alias != "" {
				out = append(out, "`" + table.name + "` " + table.alias)
			} else {
				out = append(out, "`" + table.name + "`")
			}
		}

		if item.mode == db_mode_update {
			if table.alias != "" {
				out = append(out, "`" + table.name + "` AS " + table.alias)
			} else {
				out = append(out, "`" + table.name + "`")
			}
		}

		if item.mode == db_mode_insert || item.mode == db_mode_delete {
			out = append(out, "`" + table.name + "`")
		}
	}

	return strings.Join(out, ", ")
}

func (v *Generator) generateQueryConds(item *genQuery) string {
	var out []string

	for _, cond := range item.conds {
		cnd := "`" + cond.name + "`" + cond.cnd + "?"
		out = append(out, cnd)
	}

	return strings.Join(out, " AND ")
}

func (v *Generator) generateQueryVars(item *genQuery) string {
	var out []string

	if item.mode == db_mode_update || item.mode == db_mode_insert {
		for _, field := range item.fields {
			if field.ftype != field_mode_expr {
				out = append(out, field.value)
			}
		}
	}

	for _, cond := range item.conds {
		out = append(out, cond.value.(string))
	}

	if item.limit != "" {
		out = append(out, item.limit)
	}

	if item.offset != "" {
		out = append(out, item.offset)
	}

	return strings.Join(out, ", ")
}

func (v *Generator) generateQueryScans(item *genQuery) string {
	var out []string

	for _, field := range item.fields {
		if field.ftype != field_mode_expr {
			out = append(out, field.value)
		}
	}

	return strings.Join(out, ", ")
}