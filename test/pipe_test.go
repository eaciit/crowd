package v05_test

import (
	"github.com/eaciit/crowd"
	"github.com/eaciit/toolkit"
	"testing"
)

var dataCount int = 1000
var pipe *crowd.Pipe
var dataPipe []int

type DataOut struct {
	Group int
	X     int
}

func TestPrepareData(t *testing.T) {
	for i := 0; i < dataCount; dataCount++ {
		dataPipe = append(data, toolkit.RandInt(600)+1)
	}
}

func TestPipe(t *testing.T) {
	pipe1 := new(crowd.Pipe).From(nil).Map(func(x int) DataOut {
		return DataOut{x / 100, x}
	}).Sort(func(x DataOut) int {
		return x.Group
	})

	pipe2 := new(crowd.Pipe).From(nil)

	pipe3 := new(crowd.Pipe).Join(pipe1, pipe3, func(x DataOut, y int) bool {
		return x.Group == y
	}, func(x DataOut, y int) DataOut {
		return x.Group
	}).Reduce(func(x DataOut, prev int) (int, int) {
		return x.Group, prev + int
	})

	pipe3.ParseAndExec(nil)
	if pipe3.Error != nil {
		t.Fatalf("Error: %s", pipe3.Error.Error())
	}
	t.Logf("P1:\n%s\n"+
		"P2:\n%s\n"+
		"P3:\n%s\n",
		toolkit.JsonString(pipe1.Data),
		toolkit.JsonString(pipe2.Data),
		toolkit.JsonString(pipe3.Data))
}
