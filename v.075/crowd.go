package crowd075

import (
	"errors"

	"github.com/eaciit/toolkit"
)

type FnCrowd func(x interface{}) interface{}

var Self FnCrowd = func(x interface{}) interface{} {
	return x
}

func _fn(f FnCrowd) FnCrowd {
	if f == nil {
		return Self
	} else {
		return f
	}
}

type CrowdResult struct {
	Avg   interface{}
	Min   interface{}
	Max   interface{}
	Sum   interface{}
	Group map[interface{}][]interface{}
	Sort  interface{}
}
type Crowd struct {
	SliceBase
	Error  error
	Result CrowdResult
}
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

var cmd = Command{}

func GetCommand() Command {
	return cmd
}

func From(data interface{}) *Crowd {
	c := new(Crowd)
	c.SetData(data)
	cmd.isFrom = true
	return c
}

func (c *Crowd) Avg(fn FnCrowd) *Crowd {
	cmd.isAvg = true
	cmd.fnAvg = fn
	return c
}
func (c *Crowd) Min(fn FnCrowd) *Crowd {
	cmd.isMin = true
	cmd.fnMin = fn
	return c

}
func (c *Crowd) Max(fn FnCrowd) *Crowd {
	cmd.isMax = true
	cmd.fnMax = fn
	return c
}
func (c *Crowd) Sum(fn FnCrowd) *Crowd {
	cmd.isSum = true
	cmd.fnSum = fn
	return c
}
func (c *Crowd) Group(fnGroupKey FnCrowd, fnGroupChild FnCrowd) *Crowd {
	cmd.isGroup = true
	cmd.fnGroupKey = fnGroupKey
	cmd.fnGroupChild = fnGroupChild
	return c
}
func (c *Crowd) Sort(sortDirection SortDirection, fn FnCrowd) *Crowd {
	cmd.isSort = true
	cmd.fnSort = fn
	cmd.sortDir = sortDirection
	return c
}

func (c *Crowd) Exec() (*Crowd, error) {
	var e error
	if !cmd.isFrom {
		return c, errors.New("From data not defined.")
	}
	if cmd.isAvg {
		c.Result.Avg = 0
		l := c.Len()
		if l == 0 {
			return c, nil
		}
		ret, _ := toolkit.GetEmptySliceElement(c.data)
		//toolkit.Println("Value: ", ret, reflect.TypeOf(ret).String())
		if !toolkit.IsNumber(ret) {
			return c, nil
		}

		fn := _fn(cmd.fnAvg)
		sum := float64(0)
		for i := 0; i < l; i++ {
			item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
			sum += item
		}
		//e := toolkit.Serde(sum, &ret, "json")
		c.Result.Avg = sum / float64(l)
		//		return c, nil
	}
	if cmd.isMin {
		var min interface{}
		l := c.Len()

		//min, _ = toolkit.GetEmptySliceElement(c.data)
		fn := _fn(cmd.fnMin)
		for i := 0; i < l; i++ {
			item := fn(c.Item(i))
			if item == int(0) {
				toolkit.Println("Item ", i, "=0")
			}
			if i == 0 {
				min = item
			} else if toolkit.Compare(min, item, "$gt") {
				min = item
			}
		}
		c.Result.Min = min
		//		return c, nil
	}
	if cmd.isMax {
		var max interface{}
		l := c.Len()

		max, _ = toolkit.GetEmptySliceElement(c.data)
		fn := _fn(cmd.fnMax)
		for i := 0; i < l; i++ {
			item := fn(c.Item(i))
			if i == 0 {
				max = item
			} else if toolkit.Compare(max, item, "$lt") {
				max = item
			}
		}
		c.Result.Max = max
	}
	if cmd.isSum {
		l := c.Len()

		ret, _ := toolkit.GetEmptySliceElement(c.data)
		//toolkit.Println("Value: ", ret, reflect.TypeOf(ret).String())
		if !toolkit.IsNumber(ret) {
			c.Result.Sum = 0
		}

		fn := _fn(cmd.fnSum)
		sum := float64(0)
		for i := 0; i < l; i++ {
			item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
			sum += item
		}
		//e := toolkit.Serde(sum, &ret, "json")

		c.Result.Sum = sum
	}
	if cmd.isGroup {
		ret := map[interface{}][]interface{}{}
		l := c.Len()
		fnKey := _fn(cmd.fnGroupKey)
		fnChild := _fn(cmd.fnGroupChild)
		for i := 0; i < l; i++ {
			item := c.Item(i)
			k := fnKey(item)
			v := fnChild(item)
			_, has := ret[k]
			if !has {
				ret[k] = []interface{}{}
			}
			ret[k] = append(ret[k], v)
		}
		c.Result.Group = ret
	}
	if cmd.isSort {
		l := c.Len()
		if l == 0 {
			c.Result.Sort = 0
		}

		type sk struct {
			Index   int
			SortKey interface{}
		}
		c.Result.Sort = c.data
		fn := _fn(cmd.fnSort)
		keysorter, esort := NewSorter(c.Result.Sort, fn)
		if esort != nil {
			e = errors.New("crowd.Sort: " + esort.Error())
		}
		keysorter.Sort(cmd.sortDir)
		//		c.Result.Sort = c.data
	}
	return c, e
}
