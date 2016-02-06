package crowd

type ApplyScope string

const (
	ScopeLocal   ApplyScope = "local"
	ScopeGlobal  ApplyScope = "global"
	ScopeCluster ApplyScope = "cluster"
)

type Pipe struct {
	c *Crowd
}

func (p *Pipe) Exec() interface{} {
	return nil
}

func (p *Pipe) From(fn interface{}) *Pipe {
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
