package crowd

import (
	"fmt"
	"github.com/eaciit/toolkit"
	"testing"
)

var (
	data []int = []int{1, 2, 3, 4}
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

func TestMax(t *testing.T) {
	i := From(data).Max(nil)
	if i != 4 {
		t.Errorf("Expect %d got %2.0f", 4, i)
		return
	}
}

func TestMin(t *testing.T) {
	i := From(data).Min(nil)
	if i != 1 {
		t.Errorf("Expect %d got %2.0f", 1, i)
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
	var i interface{}
	i = From(data).FindOne(func(x interface{}) bool {
		return x == 2
	})
	var val interface{}
	m := i.([]int)
	for _, each := range m {
		if each == 2 {
			val = each
		}
	}

	if val == nil || val.(int) != 2 {
		t.Errorf("Expect %d got %v", 2, val)
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
