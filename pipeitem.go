package crowd

import (
	"github.com/eaciit/toolkit"
	"reflect"
	"strings"
)

type PipeItem struct {
	attributes toolkit.M
}

func (p *PipeItem) initAttributes() {
	if p.attributes == nil {
		p.attributes = toolkit.M{}
	}
}

func (p *PipeItem) Set(k string, v interface{}) {
	p.initAttributes()
	p.attributes.Set(k, v)
}

func (p *PipeItem) Get(k string, def interface{}) interface{} {
	p.initAttributes()
	return p.attributes.Get(k, def)
}

func (p *PipeItem) Run() {
	op := strings.ToLower(p.Get("op", "").(string))
	if op == "" {
		p.Set("error", "OP is mandatory")
		return
	}
	fn := p.Get("fn_"+op, nil)
	if fn == nil {
		p.Set("error", toolkit.Sprintf("Function %s is not available", op))
		return
	}

	vfn := reflect.Indirect(reflect.ValueOf(fn))
	if vfn.Kind() != reflect.Func {
		p.Set("error", toolkit.Sprintf("Function %s is not a function", op))
	}

	var ins []reflect.Value
	var outs []reflect.Value
	outs = vfn.Call(ins)

	p.Set("output", outs)
}
