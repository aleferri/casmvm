package main

import (
	"testing"
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
	}
	if ret.Returns().Peek(1) != 60 {
		t.Error("Expected 60 as first return value")
	}
}

func TestCasmVMErr(t *testing.T) {
	vm, err := ParseLineByLine("../../tests/tryerr.csm", true)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	_, runErr := vm.Enter(0)
	if runErr == nil {
		t.Error(runErr.Error())
		t.Fail()
	}
}
