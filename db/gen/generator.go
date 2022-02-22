package gen

import (
	"os"
	"strings"
	"strconv"
	"go/ast"
	"go/token"

	"github.com/s0ulw1sh/soulgost/utils"
)

type dbGenerator struct {
	structs []dbstruct
}

type field struct {
	pk     bool
	uk     bool
	fk     bool
	ai     bool
	nn     bool
	pg     bool
	xx     bool
	xw     bool
	name   string
	goname string
	gotype string
}

type dbstruct struct {
	name     string
	table    string
	stype    *ast.StructType
	fcount   int
	pkcount  int
	fields   []field

	dbcount   bool
	dbload    bool
	dbsave    bool
	dbremove  bool
}

func (self *dbstruct) prepare() bool {
	var (
		tag string
		tags []string
		ok bool
	)

	self.fcount = 0

	for _, f := range self.stype.Fields.List {
		if f.Tag == nil { continue }

		tagraw := f.Tag.Value

		if tagraw[0] == '`' {
			tagraw = tagraw[1:len(tagraw)-1]
		}

		if tag, ok = utils.TagLookup(tagraw, "sg"); !ok { continue }

		tags = strings.Split(tag, ",")

		if len(tags) == 0 { return false }

		newfld := field{
			name:   tags[0],
			goname: f.Names[0].Name,
		}

		for _, tag := range tags[1:] {
			switch tag {
			case "ai": newfld.ai = true
			case "nn": newfld.nn = true
			case "uk": newfld.uk = true
			case "pk":
				newfld.pk = true
				self.pkcount += 1
			case "fk": newfld.fk = true
			case "xx": newfld.xx = true
			case "xw": newfld.xw = true
			case "pg": newfld.pg = true
			}
		}

		switch t := f.Type.(type) {
		case *ast.SelectorExpr:
			if tx, ok := t.X.(*ast.Ident); ok {
				newfld.gotype = tx.Name + "." + t.Sel.Name
			} else {
				return false
			}
		case *ast.Ident:
			newfld.gotype = t.Name
		default:
			return false
		}

		self.fields = append(self.fields, newfld)
		self.fcount += 1
	}

	return self.fcount != 0
}

func (self *dbGenerator) getTableName(root *ast.File, typename string) string {
	var (
		fd *ast.FuncDecl
		se *ast.StarExpr
		id *ast.Ident
		rs *ast.ReturnStmt
		bl *ast.BasicLit
		ok bool
	)

	for _, d := range root.Decls {
		if fd, ok = d.(*ast.FuncDecl); !ok || fd.Recv == nil || len(fd.Recv.List) != 1 { continue }
		if se, ok = fd.Recv.List[0].Type.(*ast.StarExpr); !ok { continue }
		if id, ok = se.X.(*ast.Ident); !ok || id.Name != typename { continue }
		if fd.Name.Name != "TableName" { continue }

		for _, stmt := range fd.Body.List {
			if rs, ok = stmt.(*ast.ReturnStmt); !ok && len(rs.Results) != 1 { continue }
			if bl, ok = rs.Results[0].(*ast.BasicLit); !ok && bl.Kind != token.STRING { continue }

			value, err := strconv.Unquote(bl.Value)
			if err == nil {
				return value
			} 
		}
	}

	return strings.ToLower(typename)
}

func (self *dbGenerator) hasFuncs(root *ast.File, s *dbstruct) bool {
	var (
		fd *ast.FuncDecl
		se *ast.StarExpr
		id *ast.Ident
		ok bool
		cnt int = 0
	)
	for _, d := range root.Decls {
		if fd, ok = d.(*ast.FuncDecl); !ok || fd.Recv == nil || len(fd.Recv.List) != 1 { continue }
		if se, ok = fd.Recv.List[0].Type.(*ast.StarExpr); !ok { continue }
		if id, ok = se.X.(*ast.Ident); !ok || id.Name != s.name { continue }

		switch fd.Name.Name {
		case "DbLoad":    s.dbload    = true
		case "DbSave":    s.dbsave    = true
		case "DbRemove":  s.dbremove  = true
		default:
			continue
		}

		cnt += 1
	}

	return cnt == 7
}

func (self *dbGenerator) checkType(root *ast.File, decl *ast.GenDecl) bool {
	var (
		ts  *ast.TypeSpec
		st  *ast.StructType
		ok  bool
		ret bool
	)

	for _, spec := range decl.Specs {
		if ts, ok = spec.(*ast.TypeSpec); !ok { continue }
		if st, ok = ts.Type.(*ast.StructType); !ok { continue }

		for _, f := range st.Fields.List {
			if f.Tag == nil {
				continue
			}

			tag := f.Tag.Value

			if tag[0] == '`' {
				tag = tag[1:len(tag)-1]
			}

			if _, ok := utils.TagLookup(tag, "sg"); ok {

				self.structs = append(self.structs, dbstruct{
					name:  ts.Name.Name,
					table: self.getTableName(root, ts.Name.Name),
					stype: st,
				})

				ret = true
				break
			}
		}
	}

	return ret
}

func (self *dbGenerator) checkDecls(root *ast.File) bool {
	ret := false

	for _, d := range root.Decls {
		if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			if self.checkType(root, gd) {
				ret = true
			}
		}
	}

	return ret
}

// ===================================================
// GENERATIONS
// ===================================================

func (self *dbGenerator) genPaginations(fw *os.File, s *dbstruct) {
	var (
		pgarr []string
		pgcnt string
		pgnmn string
	)

	pgcnt = strings.ToLower(s.name)+"_pagi_max"
	pgnmn = s.name + "DbPagi"

	fw.WriteString("const " + pgcnt + " = 30\n\n")

	for _, f := range s.fields {
		if f.pg {
			pgarr = append(pgarr, f.goname + " " + f.gotype + " `json:\""+f.name+",omitempty\"`")
		}
	}

	fw.WriteString("type " + pgnmn + " struct {\n\t")
	fw.WriteString(strings.Join(pgarr, "\n\t") + "\n\t")
	fw.WriteString("Count int `json:\"count,omitempty\"`\n\t")
	fw.WriteString("Max   int `json:\"max,omitempty\"`\n\t")
	fw.WriteString("Page  int `json:\"page,omitempty\"`\n\t")
	fw.WriteString("Pages int `json:\"pages,omitempty\"`\n")
	fw.WriteString("}\n\n")
	
	fw.WriteString(`func (self *`+pgnmn+`) Limit() int {
	if self.Max == 0 {
		return `+pgcnt+`
	} else {
		return self.Max
	}
}` + "\n\n")

fw.WriteString(`func (self *`+pgnmn+`) Offset() int {
	p := self.Page
	if p > 0 { p = p - 1 }
	return self.Limit() * p
}` + "\n\n")

fw.WriteString(`func (self *`+pgnmn+`) MarshalJSON() ([]byte, error) {
	if self.Max == 0 { self.Max = `+pgcnt+` }
	self.Pages  = self.Count / self.Max
	if self.Pages == 0 { self.Pages = 1 }
	if self.Page  == 0 { self.Page  = 1 }
	return json.Marshal(*self)
}` + "\n\n")

}

func (self *dbGenerator) genSubStructs(fw *os.File, s *dbstruct) {
	fw.WriteString("type " + s.name + "DbList struct {\n\t")
	fw.WriteString("List []" + s.name + " `json:\"list\"`\n\t")
	fw.WriteString("Pagi "+s.name+"DbPagi `json:\"pagi\"`\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbCount(fw *os.File, s *dbstruct) {
	fw.WriteString("func "+s.name+"DbCount(dbx *sql.DB) (c int, err error) {\n\t")
	if s.pkcount == 0 {
		fw.WriteString("err = dbx.QueryRow(\"SELECT COUNT(*) FROM `" + s.table + "`\").Scan(&c)\n")
	} else {
		for _, f := range s.fields {
			if f.pk {
				fw.WriteString("err = dbx.QueryRow(\"SELECT COUNT(`"+f.name+"`) FROM `" + s.table + "`\").Scan(&c)\n")
				break
			}
		}
	}
	fw.WriteString("\treturn\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbCountPk(fw *os.File, s *dbstruct) {

	var (
		idsarr []string
		whrarr []string
		whsarr []string
		pkname string
	)

	for _, f := range s.fields {
		if f.pk {
			if len(pkname) == 0 {
				pkname = "`" + f.name + "`"
			}
			idsarr = append(idsarr, f.name + "_v " + f.gotype)
			whrarr = append(whrarr, "`" + f.name + "`=?")
			whsarr = append(whsarr, f.name + "_v")
		}
	}

	if len(idsarr) > 0 {
		fw.WriteString("func "+s.name+"DbCountByPk(dbx *sql.DB, "+strings.Join(idsarr, ",")+") (c int, err error) {\n\t")
		fw.WriteString("err = dbx.QueryRow(\"SELECT COUNT("+pkname+") FROM `" + s.table + "` WHERE "+strings.Join(whrarr, " AND ")+"\", "+strings.Join(whsarr, ",")+").Scan(&c)\n")
		fw.WriteString("\treturn\n")
		fw.WriteString("}\n\n")
	}
}

func (self *dbGenerator) genFnSelfDbLoad(fw *os.File, s *dbstruct) {
	var (
		idsarr []string
		selarr []string
		scnarr []string
		whrarr []string
		whsarr []string
	)

	for _, f := range s.fields {
		if f.xx { continue }

		if f.pk {
			idsarr = append(idsarr, f.name + "_v " + f.gotype)
			whrarr = append(whrarr, "`" + f.name + "`=?")
			whsarr = append(whsarr, f.name + "_v")
		}
		selarr = append(selarr, "`" + f.name + "`")
		scnarr = append(scnarr, "&self." + f.goname)
	}

	if len(idsarr) == 0 {
		return
	}

	fw.WriteString("func (self *" + s.name + ") DbLoad(dbx *sql.DB, "+strings.Join(idsarr, ", ")+") error {\n\t")
	fw.WriteString("return dbx.QueryRow(\"SELECT "+strings.Join(selarr, ",")+" FROM `" + s.table + "` WHERE "+strings.Join(whrarr, " AND ")+"\", "+strings.Join(whsarr, ", ")+").Scan("+strings.Join(scnarr, ", ")+")\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbInsert(fw *os.File, s *dbstruct) {

	var (
		fldarr []string
		qsarr  []string
		vsarr  []string
	)

	for _, f := range s.fields {
		if f.xw { continue }

		if !f.ai {
			fldarr = append(fldarr, "`" + f.name + "`")
			qsarr  = append(qsarr, "?")
			vsarr  = append(vsarr, "item." + f.goname)
		}
	}

	fw.WriteString("func "+s.name+"DbInsert(dbx *sql.DB, item *"+s.name+") (int64, error) {\n\t")
	fw.WriteString("res, err := dbx.Exec(\"INSERT INTO `"+s.table+"` ("+strings.Join(fldarr, ",")+") VALUES ("+strings.Join(qsarr, ",")+")\", "+strings.Join(vsarr, ",")+")\n\t")
	fw.WriteString("if err != nil {\n\t\t")
	fw.WriteString("return 0, err\n\t")
	fw.WriteString("}\n\t")
	fw.WriteString("return res.LastInsertId()\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnSelfDbSave(fw *os.File, s *dbstruct) {
	var (
		fldarr []string
		vsarr  []string
		pkarr  []string
		wharr  []string
	)

	for _, f := range s.fields {
		if f.xx || f.xw { continue }

		if !f.ai && !f.pk {
			fldarr = append(fldarr, "`" + f.name + "`=?")
			vsarr  = append(vsarr, "self." + f.goname)
		}
		if f.pk {
			pkarr  = append(pkarr, "`" + f.name + "`=?")
			wharr  = append(wharr, "self." + f.goname)
		}
	}

	if len(pkarr) == 0 {
		return
	}

	fw.WriteString("func (self *" + s.name + ") DbSave(dbx *sql.DB) (int64, error) {\n\t")
	fw.WriteString("res, err := dbx.Exec(\"UPDATE `"+s.table+"` SET "+strings.Join(fldarr, ",")+" WHERE "+strings.Join(pkarr, " AND ")+"\", "+strings.Join(vsarr, ",")+", "+strings.Join(wharr, ",")+")\n\t")
	fw.WriteString("if err != nil {\n\t\t")
	fw.WriteString("return 0, err\n\t")
	fw.WriteString("}\n\t")
	fw.WriteString("return res.RowsAffected()\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnSelfDbRemove(fw *os.File, s *dbstruct) {
	var (
		pkarr  []string
		wharr  []string
	)

	for _, f := range s.fields {
		if f.pk {
			pkarr  = append(pkarr, "`" + f.name + "`=?")
			wharr  = append(wharr, "self." + f.goname)
		}
	}

	if len(pkarr) == 0 {
		return
	}

	fw.WriteString("func (self *" + s.name + ") DbRemove(dbx *sql.DB) (int64, error) {\n\t")
	fw.WriteString("res, err := dbx.Exec(\"DELETE FROM `"+s.table+"` WHERE "+strings.Join(pkarr, " AND ")+"\", "+strings.Join(wharr, ",")+")\n\t")
	fw.WriteString("if err != nil {\n\t\t")
	fw.WriteString("return 0, err\n\t")
	fw.WriteString("}\n\t")
	fw.WriteString("return res.RowsAffected()\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbSelect(fw *os.File, s *dbstruct) {
	var (
		selarr []string
		scnarr []string
		whrarr []string
		whsarr []string
		pksarr string = "*"
	)

	for _, f := range s.fields {
		if f.xx { continue }

		if f.pk && len(pksarr) == 1 {
			pksarr = "`" + f.name + "`"
		}

		if f.pg {
			whrarr = append(whrarr, "`" + f.name + "`=?")
			whsarr = append(whsarr, "pagi_v." + f.goname)
		}
		

		selarr = append(selarr, "`" + f.name + "`")
		scnarr = append(scnarr, "&s." + f.goname)
	}

	fw.WriteString("func "+s.name+"DbSelect(dbx *sql.DB, pagi_v *"+s.name+"DbPagi, list_v *"+s.name+"DbList) error {\n\t")

	fw.WriteString("var s " + s.name + "\n\t")
	fw.WriteString("list_v.Pagi = *pagi_v\n\t")
	if len(whrarr) > 0 {
		fw.WriteString("err := dbx.QueryRow(\"SELECT COUNT("+pksarr+") FROM `" + s.table + "` WHERE "+strings.Join(whrarr, " AND ")+"\", "+strings.Join(whsarr, ", ")+").Scan(&list_v.Pagi.Count)\n\t")
		fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
		fw.WriteString("rows, err := dbx.Query(\"SELECT "+strings.Join(selarr, ",")+" FROM `" + s.table + "` WHERE "+strings.Join(whrarr, " AND ")+" LIMIT ? OFFSET ? \", "+strings.Join(whsarr, ", ")+", pagi_v.Limit(), pagi_v.Offset())\n\t")
	} else {
		fw.WriteString("err := dbx.QueryRow(\"SELECT COUNT("+pksarr+") FROM `" + s.table + "`\").Scan(&list_v.Pagi.Count)\n\t")
		fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
		fw.WriteString("rows, err := dbx.Query(\"SELECT "+strings.Join(selarr, ",")+" FROM `" + s.table + "` LIMIT ? OFFSET ?\", pagi_v.Limit(), pagi_v.Offset())\n\t")
	}

	fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
	fw.WriteString("defer rows.Close()\n\t")
	fw.WriteString("list_v.List = make([]"+s.name+", 0)\n\t")
	fw.WriteString("for rows.Next() {\n\t\t")
	fw.WriteString("if err = rows.Scan("+strings.Join(scnarr, ",")+"); err != nil {\n\t\t\treturn err\n\t\t}\n\t\t")
	fw.WriteString("list_v.List = append(list_v.List, s)\n\t")
	fw.WriteString("}\n\t")
	fw.WriteString("return nil\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbAll(fw *os.File, s *dbstruct) {
	var (
		selarr []string
		scnarr []string
		whrarr []string
		whsarr []string
		pksarr string = "*"
	)

	for _, f := range s.fields {
		if f.xx { continue }

		if f.pk && len(pksarr) == 1 {
			pksarr = "`" + f.name + "`"
		}

		if f.pg {
			whrarr = append(whrarr, "`" + f.name + "`=?")
			whsarr = append(whsarr, "pagi_v." + f.goname)
		}

		selarr = append(selarr, "`" + f.name + "`")
		scnarr = append(scnarr, "&s." + f.goname)
	}

	fw.WriteString("func "+s.name+"DbAll(dbx *sql.DB, list_v *"+s.name+"DbList) error {\n\t")

	fw.WriteString("var s " + s.name + "\n\t")

	fw.WriteString("rows, err := dbx.Query(\"SELECT "+strings.Join(selarr, ",")+" FROM `" + s.table + "`\")\n\t")

	fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
	fw.WriteString("list_v.List = make([]"+s.name+", 0)\n\t")
	fw.WriteString("defer rows.Close()\n\t")
	fw.WriteString("for rows.Next() {\n\t\t")
	fw.WriteString("if err = rows.Scan("+strings.Join(scnarr, ",")+"); err != nil {\n\t\t\treturn err\n\t\t}\n\t\t")
	fw.WriteString("list_v.List = append(list_v.List, s)\n\t")
	fw.WriteString("}\n\t")
	fw.WriteString("return nil\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbLoadById(fw *os.File, s *dbstruct) {
	var (
		idsarr []string
		selarr []string
		scnarr []string
		whrarr []string
		whsarr []string
	)

	for _, f := range s.fields {
		if f.xx { continue }

		if f.pk {
			idsarr = append(idsarr, f.name + "_v " + f.gotype)
			whrarr = append(whrarr, "`" + f.name + "`=?")
			whsarr = append(whsarr, f.name + "_v")
		}
		selarr = append(selarr, "`" + f.name + "`")
		scnarr = append(scnarr, "&item." + f.goname)
	}

	if len(idsarr) == 0 {
		return
	}

	fw.WriteString("func "+s.name+"DbLoadById(dbx *sql.DB, item *"+s.name+", "+strings.Join(idsarr, ", ")+") error {\n\t")
	fw.WriteString("return dbx.QueryRow(\"SELECT "+strings.Join(selarr, ",")+" FROM `" + s.table + "` WHERE "+strings.Join(whrarr, " AND ")+"\", "+strings.Join(whsarr, ", ")+").Scan("+strings.Join(scnarr, ", ")+")\n")
	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genFnDbUk(fw *os.File, s *dbstruct, fuk *field) {
	var (
		selarr []string
		scnarr []string
		whrarr []string
		whsarr []string
	)

	whrarr = append(whrarr, "`" + fuk.name + "`=?")
	whsarr = append(whsarr, fuk.name + "_v")

	for _, f := range s.fields {
		if f.xx { continue }

		selarr = append(selarr, "`" + f.name + "`")
		scnarr = append(scnarr, "&s." + f.goname)
	}

	fw.WriteString("func "+s.name+"DbAll"+fuk.goname+"(dbx *sql.DB, list_v *"+s.name+"DbList, "+fuk.name + "_v " + fuk.gotype+") error {\n\t")
	
	fw.WriteString("var s " + s.name + "\n\t")
	fw.WriteString("rows, err := dbx.Query(\"SELECT "+strings.Join(selarr, ",")+" FROM `" + s.table + "` WHERE "+strings.Join(whrarr, " AND ")+"\", "+strings.Join(whsarr, ",")+")\n\t")
	fw.WriteString("if err != nil {\n\t\treturn err\n\t}\n\t")
	fw.WriteString("list_v.List = make([]"+s.name+", 0)\n\t")
	fw.WriteString("defer rows.Close()\n\t")
	fw.WriteString("for rows.Next() {\n\t\t")
	fw.WriteString("if err = rows.Scan("+strings.Join(scnarr, ",")+"); err != nil {\n\t\t\treturn err\n\t\t}\n\t\t")
	fw.WriteString("list_v.List = append(list_v.List, s)\n\t")
	fw.WriteString("}\n\t")
	fw.WriteString("return nil\n")

	fw.WriteString("}\n\n")
}

func (self *dbGenerator) genStruct(f *os.File, s *dbstruct) {
	
	if !s.prepare() { return }

	self.genPaginations(f, s)
	self.genSubStructs(f, s)

	self.genFnDbCount(f, s)
	self.genFnDbCountPk(f, s)

	if !s.dbload {
		self.genFnSelfDbLoad(f, s)
	}

	if !s.dbsave {
		self.genFnSelfDbSave(f, s)
	}

	if !s.dbremove {
		self.genFnSelfDbRemove(f, s)
	}

	self.genFnDbInsert(f, s)
	self.genFnDbSelect(f, s)
	self.genFnDbAll(f, s)
	self.genFnDbLoadById(f, s)

	for _, fuk := range s.fields {
		if fuk.uk {
			self.genFnDbUk(f, s, &fuk)
		}
	}
}

func Generate(root *ast.File, f *os.File) bool {
	gen := dbGenerator{}

	if !gen.checkDecls(root) {
		return false
	}

	f.WriteString("package " + root.Name.Name + "\n\n")
	f.WriteString("//-soulgost\n\n")
	f.WriteString("// WARNING!!! \n")
	f.WriteString("// Code generated by \"soulgost -modes=db\"; DO NOT EDIT!\n")
	f.WriteString("// URL - https://github.com/s0ulw1sh/soulgost\n")
	f.WriteString("// by Pavel Rid aka s0ulw1sh\n\n")
	f.WriteString("import (\n")
	f.WriteString("\t\"strings\"\n")
	f.WriteString("\t\"database/sql\"\n")
	f.WriteString("\t\"encoding/json\"\n")
	f.WriteString("\t\"github.com/s0ulw1sh/soulgost/db\"\n")
	f.WriteString(")\n\n")

	for _, s := range gen.structs {
		if gen.hasFuncs(root, &s) { continue }
		gen.genStruct(f, &s)
	}

	return true
}