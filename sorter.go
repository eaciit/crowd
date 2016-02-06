package crowd

import (
	"errors"
	"github.com/eaciit/toolkit"
	"sort"
)

type SortDirection string

const (
	SortAscending  SortDirection = "asc"
	SortDescending SortDirection = "desc"
)

type Sorter struct {
	SliceBase
	FnSort    FnCrowd
	Direction SortDirection
}

func NewSorter(data interface{}, fnsort FnCrowd) (s *Sorter, e error) {
	if !toolkit.IsPointer(data) || !toolkit.IsSlice(data) {
		e = errors.New("crowd.NewSorter: data is not pointer of slice")
		return
	}

	s = new(Sorter)
	s.data = data
	s.FnSort = _fn(fnsort)
	return
}

func (s *Sorter) Swap(i, j int) {
	//s.data[i], s.data[j] = s.data[j], s.data[i
	//toolkit.Printf("Swapping: %d and %d \n", i, j)
	si := s.Item(i)
	sj := s.Item(j)

	s.Set(i, sj)
	s.Set(j, si)
}

func (s *Sorter) Less(i, j int) bool {
	v0 := s.FnSort(s.Item(i))
	v1 := s.FnSort(s.Item(j))
	return toolkit.Compare(v0, v1, "$lt")
}

func (s *Sorter) Sort(direction SortDirection) interface{} {
	sort.Sort(s)
	return s.data
}
