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
	for i := 0; i < dataCount; i++ {
		dataPipe = append(dataPipe, toolkit.RandInt(600)+1)
	}
	if len(dataPipe) != dataCount {
		t.Fatalf("Error: want %d data got %d", dataCount, len(dataPipe))
	}
	toolkit.Println("Data (20 samples): ", toolkit.JsonString(dataPipe[:20]))
}

func TestLoad(t *testing.T) {
	var outs []int
	ds := new(crowd.PipeSource).SetData(&dataPipe)
	pipe1 := new(crowd.Pipe).From(ds).SetOutput(&outs)
	pipe1.ParseAndExec(nil, false)
	if pipe1.ErrorTxt() != "" {
		t.Fatalf("Error load: " + pipe1.ErrorTxt())
	}
	if len(outs) != len(dataPipe) {
		t.Fatalf("Error: want %d data got %d", len(dataPipe), len(outs))
	}
	t.Logf("Data: " + toolkit.JsonString(outs[0:20]))
}

/*
func TestPipe(t *testing.T) {
	t.Skip()
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
i}
*/
