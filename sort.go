package crowd

import (
	"sort"
)

type SortItem struct {
	Key   interface{}
	Value interface{}
}

type FnSortGet func(SortItem) interface{}
type FnSortCompare func(interface{}, interface{}) bool

type Sorter struct {
	Items    []SortItem
	Getter   FnSortGet
	Comparer FnSortCompare
}

func (s *Sorter) Len() int {
	return len(s.Items)
}

func (s *Sorter) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

func (s *Sorter) Less(i, j int) bool {
	//_ = "breakpoint"
	if s.Getter == nil {
		s.Getter = func(so SortItem) interface{} {
			return so.Value
		}
	}

	if s.Comparer == nil {
		return true
	}

	fi := s.Getter(s.Items[i])
	fj := s.Getter(s.Items[j])
	return s.Comparer(fi, fj)
}

func NewSorter(mis []SortItem, fn FnSortGet, fnc FnSortCompare) *Sorter {
	so := new(Sorter)
	so.Items = mis
	so.Getter = fn
	so.Comparer = fnc
	return so
}

func NewSortSlice(is []interface{}, fn FnSortGet, fnc FnSortCompare) *Sorter {
	mis := []SortItem{}
	for i, v := range is {
		mis = append(mis, SortItem{i, v})
	}
	return NewSorter(mis, fn, fnc)
}

func NewSortMap(is E, fn FnSortGet, fnc FnSortCompare) *Sorter {
	mis := []SortItem{}
	for i, v := range is {
		mis = append(mis, SortItem{i, v})
	}
	return NewSorter(mis, fn, fnc)
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
