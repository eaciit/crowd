package crowd

import (
	"fmt"
	"github.com/eaciit/toolkit"
	"testing"
	"time"
)

var t = time.Now()
var (
	data       []int       = []int{1, 2, 3, 4}
	dataString []string    = []string{"1", "2", "3", "4"}
	dataFloat  []float64   = []float64{2.5, 3.2, 6.7, 5.5}
	dataDate   []time.Time = []time.Time{
		time.Date(2015, time.November, 26, 0, 0, 0, 0, time.UTC),
		time.Date(2015, time.November, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2015, time.November, 28, 0, 0, 0, 0, time.UTC),
		time.Date(2015, time.November, 29, 0, 0, 0, 0, time.UTC),
	}
)

var randoms []int

func prepareData() []int {
	dataNo := 10000
	if randoms == nil {
		randoms = []int{}
		for i := 0; i < dataNo; i++ {
			randoms = append(randoms, toolkit.RandInt(dataNo))
		}
	}
	return randoms
}

func TestLen(t *testing.T) {
	i := From(data).Len()
	if i != 4 {
		t.Errorf("Expect %d got %d", 4, i)
		return
	}
}

func TestSum(t *testing.T) {
	i := From(data).Sum(nil)
	if i != 10 {
		t.Errorf("Expect %d got %2.0f", 10, i)
		return
	}
}

func TestAvg(t *testing.T) {
	i := From(data).Avg(nil)
	if i != 2.5 {
		t.Errorf("Expect %d got %2.0f", 2.5, i)
		return
	}
}

func TestMaxInt(t *testing.T) {
	i := From(data).Max(func(x interface{}) interface{} {
		return x
	})

	if i != 4 {
		t.Errorf("Expect %d got %2.0f", 4, i)
		return
	}
}

func TestMaxString(t *testing.T) {
	i := From(dataString).Max(func(x interface{}) interface{} {
		return x
	})

	if i.(string) != "4" {
		t.Errorf("Expect %d got %2.0f", 4, i)
		return
	}
}

func TestMaxFloat(t *testing.T) {
	i := From(dataFloat).Max(func(x interface{}) interface{} {
		return x
	})

	if i.(float64) != 6.7 {
		t.Errorf("Expect %d got %2.0f", 6.7, i)
		return
	}
}

func TestMaxDate(t *testing.T) {
	i := From(dataDate).Max(nil)

	timeDate := time.Date(2015, time.November, 29, 0, 0, 0, 0, time.UTC)
	if i.(time.Time) != timeDate {
		t.Errorf("Expect %s got %s", "2015-11-29 07:00:00 +0700 ICT", i)
		return
	}
}

func TestMinInt(t *testing.T) {
	i := From(data).Min(nil)
	if i != 1 {
		t.Errorf("Expect %d got %2.0f", 1, i)
		return
	}
}

func TestMinString(t *testing.T) {
	i := From(dataString).Min(func(x interface{}) interface{} {
		return x
	})

	if i.(string) != "1" {
		t.Errorf("Expect %d got %2.0f", "1", i)
		return
	}
}

func TestMinFloat(t *testing.T) {
	i := From(dataFloat).Min(func(x interface{}) interface{} {
		return x
	})

	if i.(float64) != 2.5 {
		t.Errorf("Expect %d got %2.0f", 2.5, i)
		return
	}
}

func TestMinDate(t *testing.T) {
	i := From(dataDate).Min(func(x interface{}) interface{} {
		return x
	})

	timeDate := time.Date(2015, time.November, 26, 0, 0, 0, 0, time.UTC)
	if i.(time.Time) != timeDate {
		t.Errorf("Expect %s got %s", "2015-11-26 07:00:00 +0700 ICT", i)
		return
	}
}

func TestMean(t *testing.T) {
	i := From(data).Mean(nil)
	if i == nil {
		t.Errorf("Got None")
		return
	}
}

func TestMedian(t *testing.T) {
	i := From(data).Median(nil)
	if i == nil {
		t.Errorf("Got None")
		return
	}
}

func FindOne(t *testing.T) {
	i := From(data).FindOne(func(x interface{}) bool {
		return x == 2
	})
	if i == nil || i.(int) != 2 {
		t.Errorf("Expect %d got %v", 2, i)
		return
	}
}

func Find(t *testing.T) {
	es := From(data).Find(func(x interface{}) bool {
		return x.(int) <= 2
	}).Data
	if From(es).Len() == 0 {
		t.Errorf("Got none")
		return
	}
}

func fg(x interface{}) interface{} {
	i := x.(int)
	return i / 1000
	return i
}

func TestGroupSubset(t *testing.T) {
	//t.Skip()
	g := From(prepareData()).Group(fg, nil).Subset(5, 0).Data
	for k, v := range g {
		fmt.Printf("k:%v, v:%s\n", k, toolkit.JsonString(v.([]interface{})[0:2]))
	}
}

func TestSliceSort(t *testing.T) {
	g := From(prepareData()).Group(fg, nil).Slice()
	x := []interface{}{}
	for _, v := range g {
		ints := v.([]interface{})
		for _, i := range ints {
			x = append(x, i)
		}
	}

	sorted := NewSortSlice(x, fsort, fcompare).
		Sort().
		Slice()[0:100]
	fmt.Printf("Sample Results:\n%v\nSorted:\n%v\n", x[0:100], sorted)
}

func fsort(so SortItem) interface{} {
	return so.Value
}

func fcompare(a, b interface{}) bool {
	return a.(int) < b.(int)
}
