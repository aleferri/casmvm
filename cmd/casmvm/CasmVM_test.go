package main

import (
	"testing"

	"github.com/aleferri/casmvm/pkg/vmex"
)

func TestCasmVM(t *testing.T) {
	vm, err := ParseLineByLine("../../tests/try.csm", true)
	if err != nil {
		t.Error(err.Error())
		return
	}

	ret, runErr := vm.Enter(1)
	if runErr != nil {
		t.Error(runErr.Error())
		return
	}
	if ret.Returns().Peek(0) != 0 {
		t.Error("Expected 0 as first return value")
		t.Fail()
	}
	if ret.Returns().Peek(1) != 60 {
		t.Error("Expected 60 as second return value")
		t.Fail()
	}
}

func TestCasmVMInvoke(t *testing.T) {
	vm, err := ParseLineByLine("../../tests/try.csm", true)
	if err != nil {
		t.Error(err.Error())
		return
	}

	ret := vmex.MakeVMFrame()

	runErr := vm.Start(1, &ret)
	if runErr != nil {
		t.Error(runErr.Error())
		return
	}
	if ret.Returns().Peek(0) != 0 {
		t.Error("Expected 0 as first return value")
		t.Fail()
	}
	if ret.Returns().Peek(1) != 60 {
		t.Error("Expected 60 as second return value")
		t.Fail()
	}
}

func TestCasmVMErr(t *testing.T) {
	vm, parseErr := ParseLineByLine("../../tests/tryerr.csm", true)
	if parseErr != nil {
		t.Error(parseErr.Error())
		t.Fail()
		return
	}

	_, runErr := vm.Enter(0)
	if runErr == nil {
		t.Error("Was expecting an error")
		t.Fail()
	}
}
