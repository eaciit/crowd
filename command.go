package crowd 

import (
    "github.com/eaciit/toolkit"
)

type FnJoinKey func(interface{},interface{})bool
type FnJoinSelect func(interface{},interface{})interface{}

type Command struct{
    CommandType CommandType
    Parms *toolkit.M
    Fns []FnCrowd
    FnJoinKey FnJoinKey
    FnJoinSelect FnJoinSelect
}

func newCommand(commandType CommandType, functions ...FnCrowd) *Command{
    c := new(Command)
    c.CommandType = commandType
    c.Fns = functions
    c.Parms = &toolkit.M{}
    return c
}

func From(data interface{}) *Crowd {
	c := new(Crowd)
	c.SetData(data)
    //cmd.isFrom = true
	return c
}

func (c *Crowd) Avg(fn FnCrowd) *Crowd {
	c.commands = append(c.commands, newCommand(CommandAvg, _fn(fn)))
    return c
}
func (c *Crowd) Min(fn FnCrowd) *Crowd {
    fn = _fn(fn)
	c.commands = append(c.commands, newCommand(CommandMin, fn))
    return c
}

func (c *Crowd) Max(fn FnCrowd) *Crowd {
	fn = _fn(fn)
	c.commands = append(c.commands, newCommand(CommandMax, fn))
    return c
}

func (c *Crowd) Sum(fn FnCrowd) *Crowd {
	fn = _fn(fn)
	c.commands = append(c.commands, newCommand(CommandSum, fn))
    return c
}

func (c *Crowd) Group(fnGroupKey FnCrowd, fnGroupChild FnCrowd) *Crowd {
	fnGroupKey = _fn(fnGroupKey)
    fnGroupChild = _fn(fnGroupChild)
	c.commands = append(c.commands, newCommand(CommandGroup, fnGroupKey, fnGroupChild))
    return c
}

func (c *Crowd) Where(fn FnCrowd) *Crowd {
    fn = _fn(fn)
    cmd := newCommand(CommandWhere, fn)
  	c.commands = append(c.commands, cmd)
    return c
}

func (c *Crowd) Apply(fn FnCrowd) *Crowd {
    fn = _fn(fn)
	cmd := newCommand(CommandApply, fn)
  	c.commands = append(c.commands, cmd)
    return c
}

func (c *Crowd) Join(data interface{}, fnKey FnJoinKey, fnSelect FnJoinSelect) *Crowd{
    cmd := newCommand(CommandApply)
    cmd.FnJoinKey = fnKey
    cmd.FnJoinSelect = fnSelect
  	c.commands = append(c.commands, cmd)
    return c
}

func (cmd *Command) Exec(c *Crowd)error{
    return nil
}