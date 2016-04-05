package crowd 

import (
    "github.com/eaciit/toolkit"
)

type Command struct{
    CommandType CommandType
    Parms *toolkit.M
    Fns []FnCrowd
}

/*
type Command struct {
	isFrom       bool
	isAvg        bool
	fnAvg        FnCrowd
	isMin        bool
	fnMin        FnCrowd
	isMax        bool
	fnMax        FnCrowd
	isSum        bool
	fnSum        FnCrowd
	isGroup      bool
	fnGroupKey   FnCrowd
	fnGroupChild FnCrowd
	isSort       bool
	fnSort       FnCrowd
	sortDir      SortDirection
}
*/


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
	c.commands = append(c.commands, newCommand(CommandAvg, fn))
    return c
}
func (c *Crowd) Min(fn FnCrowd) *Crowd {
	c.commands = append(c.commands, newCommand(CommandMin, fn))
    return c
}

func (c *Crowd) Max(fn FnCrowd) *Crowd {
	c.commands = append(c.commands, newCommand(CommandMax, fn))
    return c
}

func (c *Crowd) Sum(fn FnCrowd) *Crowd {
	c.commands = append(c.commands, newCommand(CommandSum, fn))
    return c
}

func (c *Crowd) Group(fnGroupKey FnCrowd, fnGroupChild FnCrowd) *Crowd {
	c.commands = append(c.commands, newCommand(CommandGroup, fnGroupKey, fnGroupChild))
    return c
}

func (c *Crowd) Where(fn FnCrowd) *Crowd {
    cmd := newCommand(CommandWhere, fn)
  	c.commands = append(c.commands, cmd)
    return c
}

func (c *Crowd) Apply(fn FnCrowd) *Crowd {
    cmd := newCommand(CommandApply, fn)
  	c.commands = append(c.commands, cmd)
    return c
}


func (c *Crowd) Sort(sortDirection SortDirection, fn FnCrowd) *Crowd {
    cmdSort := newCommand(CommandSort, fn)
    cmdSort.Parms.Set("direction", sortDirection)
	c.commands = append(c.commands, cmdSort)
    return c
}

func (cmd *Command) Exec(c *Crowd)error{
    return nil
}