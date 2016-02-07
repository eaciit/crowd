package v05_test

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/toolkit"
	"math"
	"testing"
)

var c *crowd.Crowd
var data []int

func stopInvalidC(t *testing.T) {
	if c == nil {
		t.Fatalf("C is NIL")
	} else if c.Error != nil {
		t.Fatalf(c.Error.Error())
	}
}

func TestFromSlice(t *testing.T) {
	data = []int{20, 30, 21, 24, 55, 102, 120, 180, 2, 95, 67, 1000, 210}
	c = crowd.From(&data)
	if c.Error != nil {
		t.Fatalf("Error: " + c.Error.Error())
	}
}

func TestMin(t *testing.T) {
	stopInvalidC(t)
	m := c.Min(func(x interface{}) interface{} {
		return x.(int) / 2
	}).(int)

	if m != int(1) {
		t.Fatalf("Want 1 got %v", m)
	}
}

func TestMax(t *testing.T) {
	stopInvalidC(t)
	m := c.Max(func(x interface{}) interface{} {
		return x.(int) / 2
	}).(int)

	if m != int(500) {
		t.Fatalf("Want 1 got %v", m)
	}
}

func TestSum(t *testing.T) {
	stopInvalidC(t)
	var m int
	m = int(c.Sum(func(x interface{}) interface{} {
		return x.(int) / 2
	}))

	sum := int(0)
	for _, d := range data {
		sum += int(d / 2)
	}

	if m != sum {
		t.Fatalf("Want %d got %d", sum, m)
	}
}

func TestAvg(t *testing.T) {
	stopInvalidC(t)
	m := c.Avg(func(x interface{}) interface{} {
		return float64(x.(int)) / 2.0
	})

	sum := float64(0)
	for _, d := range data {
		sum += float64(float64(d) / 2.0)
	}
	avg := sum / float64(len(data))

	if m != avg {
		t.Fatalf("Want %d got %d", avg, m)
	}
}

func TestGroup(t *testing.T) {
	stopInvalidC(t)
	groups := c.Group(func(x interface{}) interface{} {
		return x.(int) / 100
	},
		func(x interface{}) interface{} {
			return struct {
				X   int
				Mod int
			}{
				x.(int),
				int(math.Mod(float64(x.(int)), float64(100))),
			}
		})

	x := 0
	for _, childs := range groups {
		x += len(childs)
	}
	toolkit.Println("Data: ", toolkit.JsonString(data))
	toolkit.Println("Groups: ", groups)
	l := len(data)
	if x != l {
		t.Fatalf("Expect %d got %d", l, x)
	} else {
		t.Logf("Expect %d got %d", l, x)
	}
}
