package main

import "testing"

func TestCasmVM(t *testing.T) {
	vm, err := ParseLineByLine("../../tests/try.csm", true)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	runErr := vm.Run(true)
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

	runErr := vm.Run(true)
	if runErr == nil {
		t.Error(runErr.Error())
		t.Fail()
	}
}
