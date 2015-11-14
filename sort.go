package crowd

import (
	"sort"
)

type SortItem struct {
	Key   interface{}
	Value interface{}
}

type FnSort func(SortItem) float64
type Sorter struct {
	Items  []SortItem
	FnSort FnSort
}

func (s *Sorter) Len() int {
	return len(s.Items)
}

func (s *Sorter) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

func (s *Sorter) Less(i, j int) bool {
	//_ = "breakpoint"
	if s.FnSort == nil {
		return false
	}
	fi := s.FnSort(s.Items[i])
	fj := s.FnSort(s.Items[j])
	return fi < fj
}

func NewSorter(mis []SortItem, fn FnSort) *Sorter {
	so := new(Sorter)
	so.Items = mis
	so.FnSort = fn
	return so
}

func NewSortSlice(is []interface{}, fn FnSort) *Sorter {
	mis := []SortItem{}
	for i, v := range is {
		mis = append(mis, SortItem{i, v})
	}
	return NewSorter(mis, fn)
}

func NewSortMap(is E, fn FnSort) *Sorter {
	mis := []SortItem{}
	for i, v := range is {
		mis = append(mis, SortItem{i, v})
	}
	return NewSorter(mis, fn)
}

func (s *Sorter) Sort() *Sorter {
	sort.Sort(s)
	return s
}

func (s *Sorter) Slice() []interface{} {
	rets := []interface{}{}
	for _, v := range s.Items {
		rets = append(rets, v.Value)
	}
	return rets
}
