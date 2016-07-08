package crowd_test

import (
	//"log"
	"math"
	"testing"
	//"sync"
	"github.com/eaciit/crowd"
	"github.com/eaciit/toolkit"
)

var c *crowd.Crowd

type Obj struct {
	F float64
	I int
}

var (
	objs          []Obj
	sampleCount   int = 100
	sum, min, max int
	avg           float64
)

func skipIfNil(t *testing.T) {
	if c == nil || c.GetData() == nil {
		t.Skip()
	}
}

func check(t *testing.T, e error, section string) {
	if e != nil {
		if section == "" {
			t.Fatalf("%s", e.Error())
		} else {
			t.Fatalf("%s: %s", section, e.Error())
		}
	}
}

func TestFrom(t *testing.T) {
	for i := 0; i < sampleCount; i++ {
		obj := Obj{}
		obj.F = toolkit.ToFloat64(toolkit.RandInt(1000), 2, toolkit.RoundingAuto)
		obj.I = toolkit.RandInt(1000) + 5
		objs = append(objs, obj)

		if i == 0 {
			min = obj.I
			max = obj.I
			sum = obj.I
			avg = toolkit.ToFloat64(obj.I, 4, toolkit.RoundingAuto)
		} else {
			sum += obj.I
			avg = toolkit.ToFloat64(sum, 4, toolkit.RoundingAuto) /
				toolkit.ToFloat64(i+1, 4, toolkit.RoundingAuto)
			if min > obj.I {
				min = obj.I
			}
			if max < obj.I {
				max = obj.I
			}
		}
	}

	c = crowd.From(&objs)
	check(t, c.Error, "")
	toolkit.Printf("Data len: %d, max: %d, min: %d, sum: %d, avg: %5.4f\n",
		c.Len(), max, min, sum, avg)
}

func fn(x interface{}) interface{} {
	return x.(Obj).I
}

func TestAggr(t *testing.T) {
	skipIfNil(t)
	c2 := *c
	c2.Min(fn).Max(fn).Sum(fn).Avg(fn).Exec()
	check(t, c2.Error, "Aggr")
	if toolkit.ToInt(c2.Result.Min, toolkit.RoundingAuto) != min ||
		toolkit.ToInt(c2.Result.Max, toolkit.RoundingAuto) != max ||
		c2.Result.Sum != toolkit.ToFloat64(sum, 4, toolkit.RoundingAuto) ||
		c2.Result.Avg != avg {
		t.Fatalf("Error aggr. Got %v\n", toolkit.JsonString(c2.Result))
	}
	toolkit.Println("Value: ", toolkit.JsonString(c2.Result))
}

func TestWhereSelectGroup(t *testing.T) {
	skipIfNil(t)

	c1 := *c
	cwhere := c1.Where(func(x interface{}) interface{} {
		if x.(Obj).I < 200 {
			return true
		}
		return false
	}).Exec()
	check(t, cwhere.Error, "")
	toolkit.Println("First 20 data: ", toolkit.JsonString(cwhere.Result.Data().([]Obj)))

	cselect := cwhere.Apply(func(x interface{}) interface{} {
		return x.(Obj).F
	}).Exec()
	check(t, cselect.Error, "")
	toolkit.Println("Select : First 20 data: ", toolkit.JsonString(cselect.Result.Data().([]float64)[:20]))

	cgroup := cselect.Group(func(x interface{}) interface{} {
		return (x.(float64) - math.Mod(x.(float64), float64(100))) / float64(100)
	}, nil).Exec()
	check(t, cgroup.Error, "")
	datas := cgroup.Result.Data().([]crowd.KV)
	for _, v := range datas {
		toolkit.Printf("Group %2.0f: %d data: %v\n",
			v.Key,
			len(v.Value.([]float64)),
			v.Value.([]float64))
	}

	cgroupaggr := cselect.Apply(func(x interface{}) interface{} {
		kv := x.(crowd.KV)
		vs := kv.Value.([]float64)
		sum := crowd.From(&vs).Sum(nil).Exec().Result.Sum
		return crowd.KV{kv.Key, sum}
	}).Exec()
	check(t, cgroupaggr.Error, "")
	toolkit.Println("GroupAggr: First 20 data: ", toolkit.JsonString(cgroupaggr.Result.Data().([]crowd.KV)))

	cgroupaggrmax := cgroupaggr.Max(func(x interface{}) interface{} {
		return x.(crowd.KV).Value
	}).Exec()
	check(t, cgroupaggrmax.Error, "")
	toolkit.Println("GroupAggrMax: ", cgroupaggrmax.Result.Max)

	//c = crowd.From(&objs)
	toolkit.Printfn("Data len: %d", c.Len())
	oneshot := c.Apply(func(x interface{}) interface{} {
		return float64(x.(Obj).I)
	}).Group(func(x interface{}) interface{} {
		return (x.(float64) - math.Mod(x.(float64), float64(100))) / float64(100)
	}, nil).Apply(func(x interface{}) interface{} {
		kv := x.(crowd.KV)
		vs := kv.Value.([]float64)
		sum := crowd.From(&vs).Sum(nil).Exec().Result.Sum
		return crowd.KV{kv.Key, sum}
	}).Sum(func(x interface{}) interface{} {
		return x.(crowd.KV).Value
	}).Exec()
	check(t, oneshot.Error, "")
	toolkit.Println("Oneshot: ", oneshot.Result.Sum)
}
