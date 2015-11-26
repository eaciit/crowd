package crowd

import (
	//"fmt"
	"reflect"
	"sync"
	"time"
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

func (c *Crowd) Max(fn FnTo) interface{} {
	if fn == nil {
		fn = self
	}

	var (
		maxValue   interface{}
		maxInt     int
		maxFloat64 float64
		maxString  string
		maxUint    uint
		maxUint8   uint8
		maxUint16  uint16
		maxUint32  uint32
		maxUint64  uint64
		maxDate    string
	)

	for _, val := range c.Data {
		fnResult := fn(val)
		v := reflect.ValueOf(fnResult).Kind()
		b := IsDate(val)
		if b == true {
			dateTime := fnResult.(time.Time).Format("2006-1-2")
			switch {
			case dateTime > maxDate:
				maxDate = dateTime
			}
			maxValue = maxDate
		}

		if v == reflect.String {
			// fmt.Printf("%s\n", fnResult)
			switch {
			case fnResult.(string) > maxString:
				maxString = fnResult.(string)
			}
			maxValue = maxString
		} else if v == reflect.Int || v == reflect.Int8 {
			switch {
			case fnResult.(int) > maxInt:
				maxInt = fnResult.(int)
			}
			maxValue = maxInt
		} else if v == reflect.Float32 || v == reflect.Float64 {
			switch {
			case fnResult.(float64) > maxFloat64:
				maxFloat64 = fnResult.(float64)
			}
			maxValue = maxFloat64
		} else if v == reflect.Uint {
			switch {
			case fnResult.(uint) > maxUint:
				maxUint = fnResult.(uint)
			}
			maxValue = maxUint
		} else if v == reflect.Uint8 {
			switch {
			case fnResult.(uint8) > maxUint8:
				maxUint8 = fnResult.(uint8)
			}
			maxValue = maxUint8
		} else if v == reflect.Uint16 {
			switch {
			case fnResult.(uint16) > maxUint16:
				maxUint16 = fnResult.(uint16)
			}
			maxValue = maxUint16
		} else if v == reflect.Uint32 {
			switch {
			case fnResult.(uint32) > maxUint32:
				maxUint32 = fnResult.(uint32)
			}
			maxValue = maxUint32
		} else if v == reflect.Uint64 {
			switch {
			case fnResult.(uint64) < maxUint64:
				maxUint64 = fnResult.(uint64)
			}
			maxValue = maxUint64
		}
	}

	return maxValue
}

func (c *Crowd) Min(fn FnTo) interface{} {
	if fn == nil {
		fn = self
	}

	var (
		minValue   interface{}
		minInt     int
		minFloat64 float64
		minString  string
		minUint    uint
		minUint8   uint8
		minUint16  uint16
		minUint32  uint32
		minUint64  uint64
		minDate    []string
		b          bool
	)

	for key, val := range c.Data {
		fnResult := fn(val)
		v := reflect.ValueOf(fnResult).Kind()

		b = IsDate(val)
		if b == true {
			dateTime := fnResult.(time.Time).Format("2006-1-2")
			minDate = append(minDate, dateTime)
		}

		if v == reflect.String {
			switch {
			case key == 0:
				minString = fnResult.(string)
			case fnResult.(string) < minString:
				minString = fnResult.(string)
			}
			minValue = minString
		} else if v == reflect.Int || v == reflect.Int8 {
			switch {
			case key == 0:
				minInt = fnResult.(int)
			case fnResult.(int) < minInt:
				minInt = fnResult.(int)
			}
			minValue = minInt
		} else if v == reflect.Float32 || v == reflect.Float64 {
			switch {
			case key == 0:
				minFloat64 = fnResult.(float64)
			case val.(float64) < minFloat64:
				minFloat64 = val.(float64)
			}
			minValue = minFloat64
		} else if v == reflect.Uint {
			switch {
			case key == 0:
				minUint = fnResult.(uint)
			case val.(uint) < minUint:
				minUint = val.(uint)
			}
			minValue = minUint
		} else if v == reflect.Uint8 {
			switch {
			case key == 0:
				minUint8 = fnResult.(uint8)
			case val.(uint8) < minUint8:
				minUint8 = val.(uint8)
			}
			minValue = minUint8
		} else if v == reflect.Uint16 {
			switch {
			case key == 0:
				minUint16 = fnResult.(uint16)
			case val.(uint16) < minUint16:
				minUint16 = val.(uint16)
			}
			minValue = minUint16
		} else if v == reflect.Uint32 {
			switch {
			case key == 0:
				minUint32 = fnResult.(uint32)
			case val.(uint32) < minUint32:
				minUint32 = val.(uint32)
			}
			minValue = minUint32
		} else if v == reflect.Uint64 {
			switch {
			case key == 0:
				minUint64 = fnResult.(uint64)
			case val.(uint64) < minUint64:
				minUint64 = val.(uint64)
			}
			minValue = minUint64
		}
	}

	if b == true {
		var getMinDate = func(n []string) string {
			var min string
			for i, e := range n {
				switch {
				case i == 0:
					min = e
				case e < min:
					min = e
				}
			}
			return min
		}

		var min = getMinDate(minDate)
		minValue = min
	}

	return minValue
}

func (c *Crowd) FindOne(fn func(interface{}) bool) interface{} {
	var v interface{}
	for _, val := range c.Data {
		if fn(val) == true {
			return val
		}
	}
	return v
}

func (c *Crowd) Find(fn func(interface{}) bool) *Crowd {
	var v []interface{}
	for _, val := range c.Data {
		if fn(val) == true {
			v = append(v, val)
		}

	}
	return From(v)
}

func (c *Crowd) Median(fn FnTo) interface{} {
	var v []interface{}
	var result float64

	for _, value := range c.Data {
		v = append(v, value)

	}

	devided := len(v) / 2
	result = toF64(v[devided])
	if len(v)%2 == 0 {
		result = (result + toF64(v[devided-1])) / 2
	}
	return result
}

func (c *Crowd) Mean(fn FnTo) interface{} {
	var v []interface{}
	var TotalSum float64
	var result float64
	// v := make([]int, 0, c.Len())
	for _, value := range c.Data {
		v = append(v, value)
	}

	for _, each := range v {
		TotalSum += toF64(each)
	}
	result = TotalSum / toF64(c.Len())
	return result
}

func IsDate(o interface{}) bool {
	t := reflect.TypeOf(o)
	name := t.PkgPath() + "." + t.Name()
	if name == "time.Time" {
		return true
	}
	return false
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
