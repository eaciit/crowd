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

/*
func TestMax(t *testing.T) {
	i := From(data).Max(nil)
	if i != 4 {
		t.Errorf("Expect %d got %2.0f", 4, i)
		return
	}
}

func TestMin(t *testing.T) {
	i := From(data).Max(nil)
	if i != 1 {
		t.Errorf("Expect %d got %2.0f", 1, i)
		return
	}
}


func TestMean(t *testing.T){
}

func TestMedian(t *testing.T){
}

func FindOne(t *testing.T){
	i := From(data).FindOne(func(x interface{})bool{
		return x==2
	});
	if i==nil || i.(int) != 2 {
		t.Errorf("Expect %d got %v", 2, i)
		return
	}
}

func Find(t *testing.T){
	es := From(data).FindOne(func(x interface{})bool{
		return x<=2
	}).Data;
	if es.Len()==0 {
		t.Errorf("Got none")
		return
	}
}
*/

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
