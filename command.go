package crowd 

import (
    "reflect"
    "errors"
    "github.com/eaciit/toolkit"
)

type FnJoinKey func(interface{},interface{})bool
type FnJoinSelect func(interface{},interface{})interface{}
type CommandType string
const (
    CommandMin  CommandType = "min"
    CommandMax = "max"
    CommandSum = "sum"
    CommandAvg = "avg"
    CommandSort = "sort"
    CommandGroup = "group"
    CommandGroupAggr = "groupaggr"
    CommandWhere = "where"
    CommandApply = "apply"
    CommandJoin = "join"
)
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

func (c *Crowd) GroupAggr(fnGroupKey FnCrowd, fnGroupChild FnCrowd) *Crowd {
	fnGroupKey = _fn(fnGroupKey)
    fnGroupChild = _fn(fnGroupChild)
	c.commands = append(c.commands, newCommand(CommandGroupAggr, fnGroupKey, fnGroupChild))
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
		fn := cmd.Fns[0]
		sum := float64(0)
		for i := 0; i < l; i++ {
            el := fn(c.Item(i))
            if !toolkit.IsNumber(el){
                c.Result.Sum=0
                return nil
            }
			item := toolkit.ToFloat64(el, 4, toolkit.RoundingAuto)
			sum += item
		}
		c.Result.Sum = sum
    } else if cmd.CommandType==CommandMin{
        fn := cmd.Fns[0]
		var ret interface{}
        for i := 0; i < l; i++ {
			item := fn(c.Item(i))
			if i==0 {
                ret=item
            } else if toolkit.Compare(ret, item, "gt") {
                ret=item
            }
		}
		c.Result.Min=ret
    } else if cmd.CommandType==CommandMax{
        fn := cmd.Fns[0]
		var ret interface{}
        for i := 0; i < l; i++ {
			item := fn(c.Item(i))
			if i==0 {
                ret=item
            } else if toolkit.Compare(ret, item, "lt") {
                ret=item
            }
		}
		c.Result.Max=ret
    } else if cmd.CommandType==CommandAvg{
        fn := cmd.Fns[0]
		ret := float64(0)
		for i := 0; i < l; i++ {
			el := fn(c.Item(i))
            if !toolkit.IsNumber(el){
                c.Result.Sum=0
                return nil
            }
            item := toolkit.ToFloat64(el, 4, toolkit.RoundingAuto)
			ret += item
		}
		c.Result.Avg = ret / toolkit.ToFloat64(l,0,toolkit.RoundingAuto)
    } else if cmd.CommandType==CommandWhere{
        fn := cmd.Fns[0]
        el, _ := toolkit.GetEmptySliceElement(c.data)
        tel := reflect.TypeOf(el)
		array := reflect.MakeSlice(reflect.SliceOf(tel),0,0)
        for i := 0; i < l; i++ {
			item := c.Item(i)
            if fn(item).(bool) {
                array = reflect.Append(array, reflect.ValueOf(item))
            }
        }
        c.Result.data = array.Interface()
    } else if cmd.CommandType==CommandApply{
        fn := cmd.Fns[0]
        var array reflect.Value
        for i := 0; i < l; i++ {
			item := fn(c.Item(i))
            if i==0{
                //toolkit.Println(reflect.ValueOf(item).Type().String())
                array = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(item)),0,0)
            } else {
                array = reflect.Append(array, reflect.ValueOf(item))
            }
        }
        c.Result.data = array.Interface()
        c.data = c.Result.data
    } else if cmd.CommandType==CommandGroup{
        fng := cmd.Fns[0]
        fnc := cmd.Fns[1]
        mvs := map[interface{}]reflect.Value{}
        mvo := map[interface{}]interface{}{}
        for i := 0; i < l; i++ {
		    item := c.Item(i)
            g := fng(item)
            gi := fnc(item)
            array, exist := mvs[g]
            if !exist{
                //array = []interface{}{}
                array = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(gi)),0,0)
            }
            array = reflect.Append(array,reflect.ValueOf(gi))
            //toolkit.Println("Data:",g,array)
            mvs[g]=array
        }
        for k, v := range mvs{
            mvo[k]=v.Interface()
        }
        c.Result.data = mvo
        c.data = mvo
    } else {
        return errors.New(string(cmd.CommandType) + ": not yet applicable")
    }
    return nil
}