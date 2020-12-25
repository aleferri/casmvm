package casmvm

import "testing"

func TestCasmVM(t *testing.T) {
	vm, err := LineByLine("../../tests/try.csm", true)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	runErr := vm.Run()
	if runErr != nil {
		t.Error(runErr.Error())
		t.Fail()
	}
}
