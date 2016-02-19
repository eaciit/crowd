package crowd

import (
	"errors"
	"github.com/eaciit/toolkit"
	"reflect"
)

type ApplyScope string

const (
	ScopeLocal   ApplyScope = "local"
	ScopeGlobal  ApplyScope = "global"
	ScopeCluster ApplyScope = "cluster"
)

type Pipe struct {
	c *Crowd

	Items []*PipeItem

	source IPipeSource
	parsed bool
	err    error
	output interface{}
}

func (p *Pipe) SetError(s string) {
	p.err = errors.New(s)
}

func (p *Pipe) Error() error {
	return p.err
}

func (p *Pipe) ErrorTxt() string {
	if p.err == nil {
		return toolkit.Sprintf("")
	}
	return p.err.Error()
}

func (p *Pipe) Parsed() bool {
	return p.parsed
}

func (p *Pipe) Parse() error {
	p.err = nil
	p.parsed = true
	return p.err
}

func (p *Pipe) Exec(parm toolkit.M) error {
	if p.source == nil {
		return errors.New("Pipe.Exec: Source is invalid")
	}

	if len(p.Items) == 0 {
		if p.output != nil {
			e := toolkit.Serde(p.source.Data(), p.output, "json")
			if e != nil {
				return errors.New("Pipe.Exec: unable to serde the result " + e.Error())
			}
		}
		return nil
	}

	if parm == nil {
		parm = toolkit.M{}
	}

	sLen := p.source.Len()
	for sIndex := 0; sIndex < sLen; sIndex++ {
		parm = parm.Set("dataindex", sIndex)
		p.Items[0].Set("parm", parm)
		p.Items[0].Set("in", p.source.Seek(sIndex, SeekFromStart))
		erun := p.Items[0].Run()
		if erun != nil {
			return errors.New("Pipe.Exec: " + erun.Error())
		} else {
			//toolkit.Println("Executed")
		}
	}

	return nil
}

/*
func (p *Pipe) ParseAndExec(inputs interface{}, reparse bool) {
	if reparse || p.parsed == false {
		p.Parse()
	}
	if p.Error() != nil {
		return
	}
	p.Exec(inputs)
}
*/

func (p *Pipe) Partition(i int) *Pipe {
	//p.Set("partition", i)
	pi := new(PipeItem)
	pi.Set("op", "partition")
	pi.Set("partition", i)
	eadd := p.addItem(pi)
	if eadd != nil {
		p.SetError(eadd.Error())
		return p
	}
	return p
}

func (p *Pipe) SetOutput(o interface{}) *Pipe {
	pi := new(PipeItem)
	pi.Set("op", "setoutput")
	pi.Set("fn", func(x interface{}) {
		if toolkit.IsSlice(o) {
			toolkit.AppendSlice(o, x)
		} else {
			reflect.ValueOf(o).Elem().Set(reflect.ValueOf(x))
		}
	})
	eadd := p.addItem(pi)
	if eadd != nil {
		p.SetError(eadd.Error())
		return p
	}
	p.output = o
	return p
}

func (p *Pipe) Join(p1 *Pipe, p2 *Pipe, fnJoin interface{}) *Pipe {
	return p
}

func (p *Pipe) From(s IPipeSource) *Pipe {
	p.source = s
	return p
}

func (p *Pipe) Where(fn interface{}) *Pipe {
	pi := new(PipeItem)
	pi.Set("op", "where")
	pi.Set("fn", fn)
	p.addItem(pi)
	return p
}

func (p *Pipe) Map(fn interface{}) *Pipe {
	pi := new(PipeItem)
	pi.Set("op", "map")
	pi.Set("fn", fn)
	p.addItem(pi)
	return p
}

func (p *Pipe) Sort(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Reduce(fn interface{}) *Pipe {
	pi := new(PipeItem)
	pi.Set("op", "mapreduce")
	pi.Set("fn", fn)
	_ = p.addItem(pi)
	/*
		if eadd != nil {
			toolkit.Println("AddReduce:", eadd.Error())
		}
	*/
	return p
}

func (p *Pipe) addItem(pi *PipeItem) error {
	if p.ErrorTxt() != "" {
		return errors.New("Pipe.addPipeItem: " + p.ErrorTxt())
	}

	if pi == nil {
		return errors.New("Pipe.addPipeItem: PipeItem is nil")
	}

	if len(p.Items) > 0 {
		lastpi := p.Items[len(p.Items)-1]
		if lastpi.Get("op", "") == "setoutput" {
			return errors.New("Pipe.addPipeItem: Last PipeItem is SetOutput. No more PipeItem can't be inserted after SetOutput")
		}
		lastpi.nextItem = pi
	}

	pi.Set("index", len(p.Items))
	p.Items = append(p.Items, pi)

	return nil
}
