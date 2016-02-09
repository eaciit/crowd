package crowd

type ApplyScope string

const (
	ScopeLocal   ApplyScope = "local"
	ScopeGlobal  ApplyScope = "global"
	ScopeCluster ApplyScope = "cluster"
)

type PipeItem struct {
}

type Pipe struct {
	c     *Crowd
	Items []*PipeItem
}

func (p *Pipe) Exec() interface{} {
	return nil
}

func (p *Pipe) From(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Map(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Sort(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Apply(scope ApplyScope, fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Reduce(fn interface{}) *Pipe {
	return p
}
