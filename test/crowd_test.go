package crowd_test

import (
    //"log"
	//"math"
	"testing"
	//"sync"
    "github.com/eaciit/toolkit"
    "github.com/eaciit/crowd.dev"
)

var c *crowd.Crowd
type Obj struct{
    F float64
    I int
}

var (
    objs []Obj 
    sampleCount int = 10000    
    sum, min, max int
    avg float64
)

func skipIfNil(t *testing.T){
    if c==nil || c.GetData()==nil {
        t.Skip()
    }
}

func check(t *testing.T, e error, section string){
    if e!=nil {
        if section=="" {
            t.Fatalf("%s", e.Error())
        } else {
            t.Fatalf("%s: %s", section, e.Error())
        }
    }
}

func TestFrom(t *testing.T){
    for i:=0;i<sampleCount;i++{
        obj := Obj{}
        obj.F = toolkit.ToFloat64(toolkit.RandInt(1000),2,toolkit.RoundingAuto)
        obj.I = toolkit.RandInt(1000)+5
        objs = append(objs, obj)   
        
        if i==0{
            min=obj.I
            max=obj.I
            sum=obj.I
            avg=toolkit.ToFloat64(obj.I, 4, toolkit.RoundingAuto)
        } else {
            sum +=obj.I
            avg=toolkit.ToFloat64(sum,4,toolkit.RoundingAuto) / 
                toolkit.ToFloat64(i+1,4,toolkit.RoundingAuto)
            if min > obj.I {
                min=obj.I
            }
            if max < obj.I {
                max = obj.I
            }
        }
    }
    
    c = crowd.From(&objs)
    check(t, c.Error, "")
    toolkit.Printf("Data len: %d, max: %d, min: %d, sum: %d, avg: %5.4f\n", c.Len(), max, min, sum, avg)
}

func fn(x interface{})interface{}{
    return x.(Obj).I
}

func TestAggr(t *testing.T){
    skipIfNil(t)
    c2 := c.Min(fn).Max(fn).Sum(fn).Avg(fn).Exec()
    check(t, c2.Error, "Aggr")
    if c2.Result.Min!=toolkit.ToFloat64(min,4,toolkit.RoundingAuto) || 
        c2.Result.Max!=toolkit.ToFloat64(max,4,toolkit.RoundingAuto) || 
        c2.Result.Sum!=toolkit.ToFloat64(sum,4,toolkit.RoundingAuto) || 
        c2.Result.Avg!=avg {
        t.Fatalf("Error aggr. Got %v\n", toolkit.JsonString(c2.Result))
    }
}


/*
//func TestFrom(t *testing.T) {
//	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	c = crowd.From(&data)
//	log.Printf("c.data => %v", c.SliceBase.GetData())
//}
func TestChainLittleData(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 99, 100, 34, 23, 55, 60, 23, 100, 183, 24, 23, 34}
	fn := func(x interface{}) interface{} {
		return float64(x.(int))
	}
	//	fnGroupKey := func(x interface{}) interface{} {
	//		return x.(int)
	//	}
	//	fnGroupChild := func(x interface{}) interface{} {
	//		return struct {
	//			X   int
	//			Mod int
	//		}{
	//			x.(int),
	//			int(math.Mod(float64(x.(int)), float64(100))),
	//		}
	//	}
	//	c, err := crowd.From(&data).Avg(fn).Min(fn).Max(fn).Sum(fn).Group(fnGroupKey, fnGroupChild).Sort(crowd.SortAscending, nil).Exec()
	c, err := crowd.From(&data).Avg(fn).Min(fn).Max(fn).Sum(fn).Sort(crowd.SortAscending, nil).Exec()
	log.Printf("c.Command => %#v", crowd.GetCommand())
	log.Printf("c.result => %#v ; error => %#v", c.Result, err)
	log.Printf("After sorting Ascending:  %#v ", toolkit.JsonString(c.Result.Sort))
	log.Printf("c.data => %v", c.SliceBase.GetData())
}

func TestChainBigData(t *testing.T) {
	data := []int{2, 1000}
	l := 100000
	mtx := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func() {
			defer wg.Done()
			rnd := 5 + toolkit.RandInt(900)
			mtx.Lock()
			data = append(data, rnd)
			mtx.Unlock()
		}()
	}
	wg.Wait()

	fn := func(x interface{}) interface{} {
		return float64(x.(int))
	}
	fnGroupKey := func(x interface{}) interface{} {
		return x.(int)
	}
	fnGroupChild := func(x interface{}) interface{} {
		return struct {
			X   int
			Mod int
		}{
			x.(int),
			int(math.Mod(float64(x.(int)), float64(100))),
		}
	}
	c, err := crowd.From(&data).Avg(fn).Min(fn).Max(fn).Sum(fn).Group(fnGroupKey, fnGroupChild).Sort(crowd.SortAscending, nil).Exec()
	//	c, err := crowd.From(&data).Group(fnGroupKey, fnGroupChild).Exec()
	//	c, err := v.From(&data).Avg(fn).Min(fn).Max(fn).Sum(fn).Sort(crowd.SortAscending, nil).Exec()
	log.Printf("c.Command => %#v", crowd.GetCommand())
	//	log.Printf("c.result => %#v ; error => %#v", c.Result, err)
	log.Printf("c.Result.Avg => %#v", c.Result.Avg)
	log.Printf("c.Result.Min => %#v", c.Result.Min)
	log.Printf("c.Result.Max => %#v", c.Result.Max)
	log.Printf("c.Result.Sum => %#v", c.Result.Sum)
	log.Printf("Data: %#v ", toolkit.JsonString(data[:100]))
	log.Printf("c.error => %#v", err)
	//	log.Printf("After sorting Ascending: %v", toolkit.JsonString(c.Result.Sort))

	//	log.Printf("c.data => %v", c.SliceBase.GetData())
}

func TestChainLittleString(t *testing.T) {
	data := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	//	fn := func(x interface{}) interface{} {
	//		return float64(x.(int))
	//	}
	//	fnGroupKey := func(x interface{}) interface{} {
	//		return x.(int)
	//	}
	//	fnGroupChild := func(x interface{}) interface{} {
	//		return struct {
	//			X   int
	//			Mod int
	//		}{
	//			x.(int),
	//			int(math.Mod(float64(x.(int)), float64(100))),
	//		}
	//	}
	//	c, err := crowd.From(&data).Avg(fn).Min(fn).Max(fn).Sum(fn).Group(fnGroupKey, fnGroupChild).Sort(crowd.SortAscending, nil).Exec()
	c, err := crowd.From(&data).Avg(nil).Min(nil).Max(nil).Sum(nil).Sort(crowd.SortAscending, nil).Exec()
	log.Printf("c.Command => %#v", crowd.GetCommand())
	log.Printf("c.result => %#v ; error => %#v", c.Result, err)
	toolkit.Println("After sorting Ascending: ", toolkit.JsonString(c.Result.Sort))
	log.Printf("c.data => %v", c.SliceBase.GetData())
}
*/