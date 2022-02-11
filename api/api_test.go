package api

import "testing"
import "context"

type TestApi struct {}
type TestApiStruct struct {
	Z int
}

func (self *TestApi) ApiHello(req Request, in *int, out *int) error {
	*out += 123 + *in
	return nil
}

func (self *TestApi) ApiStruct(req Request, in *TestApiStruct, out *TestApiStruct) error {
	out.Z += 123 + in.Z
	return nil
}

func (self *TestApi) CallApi(method string, req Request, res Response) (err error) {
	switch method {
	case "hello":
		var par_in int
		var par_out int
		if err = req.GetParams(&par_in); err != nil { return res.WriteResult(nil, err) }
		if err = self.ApiHello(req, &par_in, &par_out); err != nil { return res.WriteResult(nil, err) }
		return res.WriteResult(&par_out, nil)
	case "struct":
		var par_in TestApiStruct
		var par_out TestApiStruct
		if err = req.GetParams(&par_in); err != nil { return res.WriteResult(nil, err) }
		if err = self.ApiStruct(req, &par_in, &par_out); err != nil { return res.WriteResult(nil, err) }
		return res.WriteResult(&par_out, nil)
	default:
		return res.WriteResult(nil, ErrMethodNotFound)
	}
	return nil
}

func TestRpcApiFnCaller(t *testing.T) {
	var a TestApi
	var r Rpc
	var c context.Context = context.Background()

	var in int  = 5
	var out int = 0

	Register("Test", &a)

	err := r.Call(c, "Test", "Hello", &in, &out)

	if out != 128 || err != nil {
		t.Error("Not equal params", in, out, err)
	}
}

func TestRpcApiFnCallerStruct(t *testing.T) {
	var a TestApi
	var r Rpc
	var c context.Context = context.Background()

	var in TestApiStruct  = TestApiStruct{100}
	var out TestApiStruct = TestApiStruct{0}

	Register("Test", &a)

	err := r.Call(c, "Test", "Struct", &in, &out)

	if out.Z != 223 || err != nil {
		t.Error("Not equal params", in, out, err)
	}
}