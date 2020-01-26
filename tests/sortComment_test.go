package tests

import "testing"

import "sort"

type Comm struct {
	Id     int64
	Parent int64
	Thread int64
	Name   string
	Child  []Comm
}

func TestSorting(t *testing.T) {

	list := []Comm{
		Comm{
			Id:     1,
			Parent: 0,
			Thread: 0,
			Name:   "1-0",
		}, Comm{
			Id:     2,
			Parent: 0,
			Thread: 0,
			Name:   "2-0",
		}, Comm{
			Id:     3,
			Parent: 1,
			Thread: 0,
			Name:   "3-1",
		}, Comm{
			Id:     4,
			Parent: 1,
			Thread: 0,
			Name:   "4-1",
		}, Comm{
			Id:     5,
			Parent: 3,
			Thread: 0,
			Name:   "5-3",
		},
		Comm{
			Id:     6,
			Parent: 5,
			Thread: 0,
			Name:   "6-5",
		},
	}

	m := make(map[int64]Comm)
	for _, v := range list {
		m[v.Id] = v
	}

	for key, val := range m {
		if v, ok := m[val.Parent]; ok {
			if v.Child == nil {
				v.Child = []Comm{val}
			} else {
				v.Child = append(v.Child, val)
				sort.Slice(v.Child, func(i, j int) bool { return v.Child[i].Id < v.Child[j].Id })
			}
			m[val.Parent] = v
			delete(m, key)
		}
	}

	t.Errorf("%v", m)
}

/*
[{1 0 1-0 [{3 1 3-1 []} {4 1 4-1 []}]}
{2 0 2-0 []}
{3 1 3-1 [{5 3 5-3 []}]}
{4 1 4-1 []}
{5 3 5-3 [{6 5 6-5 []}]}
{6 5 6-5 []}]

*/
