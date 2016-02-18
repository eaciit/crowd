package crowd

import (
	"github.com/eaciit/toolkit"
	"reflect"
	"strings"
)

type PipeItem struct {
	attributes toolkit.M
	nextItem   *PipeItem
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

func (p *PipeItem) SetError(err string) error {
	return nil
}

func (p *PipeItem) Run() error {
	parm := p.Get("parm", toolkit.M{}).(toolkit.M)
	verbose := parm.Get("verbose", false).(bool)
	op := strings.ToLower(p.Get("op", "").(string))
	if op == "" {
		//p.Set("error", "OP is mandatory")
		return p.SetError("OP is mandatory")
	}

	//fn := p.Get("fn_"+op, nil)
	fn := p.Get("fn", nil)
	if fn == nil {
		return p.SetError(toolkit.Sprintf("Function %s is not available", op))
	}

	vfn := reflect.Indirect(reflect.ValueOf(fn))
	if vfn.Kind() != reflect.Func {
		return p.SetError(toolkit.Sprintf("Function %s is not a function", op))
	}

	var ins []reflect.Value
	var outs []reflect.Value

	pIn := p.Get("in", nil)
	if !toolkit.IsSlice(pIn) {
		ins = append(ins, reflect.ValueOf(pIn))
	} else {
		pLen := toolkit.SliceLen(pIn)
		for pIndex := 0; pIndex < pLen; pIndex++ {
			ins = append(ins, reflect.ValueOf(toolkit.SliceItem(pIn, pIndex)))
		}
	}

	//toolkit.Println(toolkit.JsonString(ins))
	outs = vfn.Call(ins)

	var iouts []interface{}
	for _, out := range outs {
		iouts = append(iouts, out.Interface())
	}

	if verbose {
		toolkit.Printf("Data %d Pipe %d %s: %s => %s\n",
			p.Get("parm", nil).(toolkit.M).Get("dataindex", 0).(int),
			p.Get("index", 0).(int), op,
			toolkit.JsonString(pIn),
			toolkit.JsonString(iouts))
	}

	if op == "where" && iouts[0] == false {
		return nil
	} else if op == "where" && iouts[0] == true {
		iouts = []interface{}{}
		for _, in := range ins {
			iouts = append(iouts, in.Interface())
		}
	}

	//p.Set("output", outs)
	if p.nextItem != nil {
		p.nextItem.Set("parm", parm)
		p.nextItem.Set("in", iouts)
		return p.nextItem.Run()
	} else {
		p.Set("output", iouts)
	}

	return nil
}
