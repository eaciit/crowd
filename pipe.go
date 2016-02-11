package crowd

type ApplyScope string

const (
	ScopeLocal   ApplyScope = "local"
	ScopeGlobal  ApplyScope = "global"
	ScopeCluster ApplyScope = "cluster"
)

type PipeItem struct {
}

type IPipeSource interface {
}

type PipeSource struct {
}

type Pipe struct {
	c *Crowd

	Items []*PipeItem

	parsed bool
	err    error
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

func (p *Pipe) Exec(inputs interface{}) interface{} {
	return nil
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

func (p *Pipe) Join(p1 *Pipe, p2 *Pipe, fnJoin interface{}) *Pipe {
	return p
}

func (p *Pipe) From(s IPipeSource) *Pipe {
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
