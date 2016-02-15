package crowd

import (
	"errors"
	"github.com/eaciit/toolkit"
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
		return ""
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

func (p *Pipe) Exec(inputs interface{}) {
	if p.output != nil && p.source != nil {
		e := toolkit.Serde(p.source.Data(), p.output, "json")
		if e != nil {
			p.SetError("Exec: unable to serde the result " + e.Error())
		}
	}
	return
}

func (p *Pipe) ParseAndExec(inputs interface{}, reparse bool) {
	if reparse || p.parsed == false {
		p.Parse()
	}
	if p.Error() != nil {
		return
	}
	p.Exec(inputs)
}

func (p *Pipe) SetOutput(o interface{}) *Pipe {
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
	return p
}

func (p *Pipe) Map(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Sort(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Reduce(fn interface{}) *Pipe {
	return p
}
