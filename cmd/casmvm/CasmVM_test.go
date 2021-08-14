package main

import (
	"testing"
)

func TestCasmVM(t *testing.T) {
	vm, err := ParseLineByLine("../../tests/try.csm", true)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	_, runErr := vm.Enter(0)
	if runErr != nil {
		t.Error(runErr.Error())
		t.Fail()
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
