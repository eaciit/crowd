package crowd

import (
	//"fmt"
	"reflect"
	"sync"
)

type E map[interface{}]interface{}
type FnTo func(interface{}) interface{}

var (
	self FnTo = func(x interface{}) interface{} {
		return x
	}
)

type Crowd struct {
	Data E
	Keys []interface{}
}

func From(d interface{}) *Crowd {
	c := new(Crowd)
	c.Data = E{}
	c.Keys = []interface{}{}
	tof := reflect.TypeOf(d).Kind()
	if tof == reflect.Map {
		es := d.(E)
		for k, v := range es {
			c.Data[k] = v
			c.Keys = append(c.Keys, k)
		}
	} else if tof == reflect.Slice {
		var slice reflect.Value
		slice = reflect.ValueOf(d)
		lenslice := slice.Len()
		for i := 0; i < lenslice; i++ {
			c.Data[i] = slice.Index(i).Interface()
			c.Keys = append(c.Keys, i)
		}
	}
	return c
}

func (c *Crowd) Slice() []interface{} {
	ret := []interface{}{}
	for _, v := range c.Data {
		ret = append(ret, v)
	}
	return ret
}

func (c *Crowd) Len() int {
	if c.Data == nil {
		return 0
	} else {
		return len(c.Data)
	}
}

func toF64(i interface{}) float64 {
	if f64, ok := i.(float64); ok {
		return f64
	} else if f32, ok := i.(float32); ok {
		return float64(f32)
	} else if fi, ok := i.(int); ok {
		return float64(fi)
	} else {
		return 0
	}
}

func (c *Crowd) Sum(fn FnTo) float64 {
	var ret float64 = 0
	if fn == nil {
		fn = self
	}

	for _, v := range c.Data {
		f64 := toF64(fn(v))
		ret += f64
	}
	return ret
}

func (c *Crowd) Avg(fn FnTo) float64 {
	var ret float64 = c.Sum(fn) / float64(c.Len())
	return ret
}

func (c *Crowd) Group(fnKey, fnValue FnTo) *Crowd {
	GroupData := E{}
	if fnKey == nil {
		fnKey = self
	}
	if fnValue == nil {
		fnValue = self
	}

	//_ = "breakpoint"
	wg := new(sync.WaitGroup)
	mtx := new(sync.Mutex)
	for _, v := range c.Data {
		wg.Add(1)
		go func(v interface{}, GroupData *E,
			wg *sync.WaitGroup, mtx *sync.Mutex) {
			gd := *GroupData
			groupId := fnKey(v)
			value := fnValue(v)
			var datas []interface{}
			//data, exist := GroupData[groupId]
			mtx.Lock()
			data, exist := gd[groupId]
			if !exist {
				datas = []interface{}{value}
			} else {
				datas = append(data.([]interface{}), value)
			}
			gd[groupId] = datas
			mtx.Unlock()
			wg.Done()
		}(v, &GroupData, wg, mtx)
	}
	wg.Wait()
	//_ = "breakpoint"

	return From(GroupData)
}

func (c *Crowd) Subset(take, skip int) *Crowd {
	idx := 0
	inloop := true
	skipped := 0
	takeFlag := false
	taken := 0
	dataLength := c.Len()
	ret := E{}

	for inloop {
		if skipped == skip && taken < take {
			takeFlag = true
		} else {
			skipped++
		}

		if takeFlag {
			ret[c.Keys[taken]] = c.Data[c.Keys[taken]]
			taken++
			if taken >= take {
				inloop = false
			}
		}

		idx++
		if idx >= dataLength {
			inloop = false
		}
	}

	return From(ret)
}

func (c *Crowd) Max(fn FnTo) int {
	var maximum interface{}
	maximum = c.Data[0]

	for _, val := range c.Data {
		if val.(int) > maximum.(int) {
			maximum = val
		}
	}
	return maximum.(int)
}

func (c *Crowd) Min(fn FnTo) int {
	var maximum interface{}
	maximum = c.Data[0]

	for _, val := range c.Data {
		if val.(int) < maximum.(int) {
			maximum = val
		}
	}
	return maximum.(int)
}

func (c *Crowd) FindOne(fn interface{}) interface{} {
	v := make([]int, 0, len(c.Data))
	for _, value := range c.Data {
		v = append(v, value.(int))
	}
	return v
}

func (c *Crowd) Find(fn interface{}) *Crowd {
	dataLength := c.Len()
	ret := E{}
	for i := 0; i < dataLength; i++ {
		ret[c.Keys[i]] = c.Data[c.Keys[i]]
	}
	return From(ret)
}

func (c *Crowd) Median(fn FnTo) interface{} {
	var v []float64
	var result float64

	for _, value := range c.Data {
		v = append(v, toF64(value.(int)))
	}

	devied := len(v) / 2
	result = v[devied]
	if len(v)%2 == 0 {
		result = (result + v[devied-1]) / 2
	}
	return result
}

func (c *Crowd) Mean(fn FnTo) interface{} {
	var v []float64
	var result float64
	// v := make([]int, 0, c.Len())
	for _, value := range c.Data {
		v = append(v, toF64(value.(int)))
	}

	for _, each := range v {
		result += each
	}
	return result
}

/*
func (c *Crowd) Sort(fn FnTo) *Crowd {
	type sortObj Crowd

	func(s *sortObj) Len()int{
		return s.Len()
	}

	func (s *sortObj) Swap(i, j int){
		s.Key[i], s.Keys[j] = s.Keys[j], s.Keys[i]
	}

	func (s *sortObj) Less(i, j int) bool{
		fi := fn(s.Data[s.Keys[i]])
		fj := fn(s.Data[s.Keys[j]])
		return fi < fj
	}

	so := c
	sorting.Sort(so)

	return c
}
*/
