package v05_test

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/toolkit"
	"testing"
)

func TestSort(t *testing.T) {
	data := []int{20, 30, 21, 24, 55, 102, 120, 180, 2, 95, 67, 1000, 210}
	toolkit.Println("Before sorting: ", toolkit.JsonString(data))
	sorter, _ := crowd.NewSorter(&data, nil)

	sorter.Sort(crowd.SortAscending)
	toolkit.Println("After sorting Ascending: ", toolkit.JsonString(data))
	min := 0
	for _, v := range data {
		if v < min {
			t.Fatalf("Error: %d and %d", min, v)
		} else {
			min = v
		}
	}

	sorter.Sort(crowd.SortDescending)
	toolkit.Println("After sorting Descending: ", toolkit.JsonString(data))
	max := 0
	for i, v := range data {
		if i == 0 {
			max = v
		}
		if v > max {
			t.Fatalf("Error: %d and %d", max, v)
		} else {
			max = v
		}
	}
}
