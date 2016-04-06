package crowd 

import (
    "errors"
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
    if c.data==nil {
        return errors.New("Exec: Data is empty")
    }
    l := c.Len()
    if cmd.CommandType==CommandSum{
		el, _ := toolkit.GetEmptySliceElement(c.data)
		if !toolkit.IsNumber(el) {
			c.Result.Sum = float64(0)
            return nil
		}
		fn := cmd.Fns[0]
		sum := float64(0)
		for i := 0; i < l; i++ {
			item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
			sum += item
		}
		c.Result.Sum = sum
    } else if cmd.CommandType==CommandMin{
        el, _ := toolkit.GetEmptySliceElement(c.data)
		if !toolkit.IsNumber(el) {
			c.Result.Min = float64(0)
            return nil
		}
        fn := cmd.Fns[0]
		ret := float64(0)
		for i := 0; i < l; i++ {
			item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
			if i==0 {
                ret=item
            } else if ret > item {
                ret=item
            }
		}
		c.Result.Min=ret
    } else if cmd.CommandType==CommandMax{
        el, _ := toolkit.GetEmptySliceElement(c.data)
		if !toolkit.IsNumber(el) {
			c.Result.Max = float64(0)
            return nil
		}
        fn := cmd.Fns[0]
		ret := float64(0)
		for i := 0; i < l; i++ {
			item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
			if i==0 {
                ret=item
            } else if ret < item {
                ret=item
            }
		}
		c.Result.Max=ret
    } else if cmd.CommandType==CommandAvg{
        el, _ := toolkit.GetEmptySliceElement(c.data)
		if !toolkit.IsNumber(el) {
			c.Result.Avg = float64(0)
            return nil
		}
        fn := cmd.Fns[0]
		ret := float64(0)
		for i := 0; i < l; i++ {
			item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
			ret += item
		}
		c.Result.Avg = ret / toolkit.ToFloat64(l,0,toolkit.RoundingAuto)
    }
    return nil
}