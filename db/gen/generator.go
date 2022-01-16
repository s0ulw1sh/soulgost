package gen

import (
	"strconv"
	"go/ast"
	"go/token"
)

type dbGenerator struct {
	structs []*ast.StructType
}


// from src/reflect/type.go
func (self *dbGenerator) tagLookup(tag, key string) (value string, ok bool) {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}

func (self *dbGenerator) checkType(decl *ast.GenDecl) bool {
	var (
		ts  *ast.TypeSpec
		st  *ast.StructType
		ok  bool
		ret bool
	)

	for _, spec := range decl.Specs {
		if ts, ok = spec.(*ast.TypeSpec); !ok {
			continue
		}

		if st, ok = ts.Type.(*ast.StructType); !ok {
			continue
		}

		for _, f := range st.Fields.List {
			if f.Tag == nil {
				continue
			}

			_, ok := self.tagLookup(f.Tag.Value, "sg")

			if ok {
				self.structs = append(self.structs, st)
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
			if self.checkType(gd) {
				ret = true
			}
		}
	}

	return ret
}

func Generate(root *ast.File, fdir string, fname string) {
	gen := dbGenerator{}

	if !gen.checkDecls(root) {
		return
	}

}