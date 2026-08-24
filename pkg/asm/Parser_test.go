package asm

import (
	"strings"
	"testing"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

func TestParseSigOpcodes(t *testing.T) {
	warn, err := ParseOpcode(`sigwarn %3 "value looks off: "`, false)
	if err != nil {
		t.Fatalf("sigwarn should parse, got error '%s'", err.Error())
	}
	if warn.String() != `sigwarn %3 value looks off: ` {
		t.Errorf("Unexpected sigwarn: '%s'", warn.String())
	}

	sigerr, err := ParseOpcode(`sigerr %1 "should be less than 15, but is " ; trailing comment`, false)
	if err != nil {
		t.Fatalf("sigerr should parse, got error '%s'", err.Error())
	}
	if !strings.Contains(sigerr.String(), "should be less than 15, but is ") {
		t.Errorf("The message should survive intact, got '%s'", sigerr.String())
	}
	if strings.Contains(sigerr.String(), "trailing comment") {
		t.Errorf("The comment should not leak into the message: '%s'", sigerr.String())
	}
}

func TestParseHonorsShapes(t *testing.T) {
	op, err := ParseOpcode(`%1 = add u16 %0 %0`, false)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if !strings.Contains(op.String(), "u16") {
		t.Errorf("The declared shape should be honored, got '%s'", op.String())
	}

	iconst, err := ParseOpcode(`%0 = const i8 200`, false)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if iconst.(*opcodes.IConst).Constant() != -56 {
		t.Errorf("const i8 200 should be reshaped to -56, got %d", iconst.(*opcodes.IConst).Constant())
	}
}

func TestParseToleratesCommas(t *testing.T) {
	op, err := ParseOpcode(`%5 = sub i16 %4, %3`, false)
	if err != nil {
		t.Fatalf("Commas are separators, got error '%s'", err.Error())
	}
	if !strings.Contains(op.String(), "%4 %3") {
		t.Errorf("Expected operands %%4 and %%3, got '%s'", op.String())
	}
}

func TestParseEnterWithoutArguments(t *testing.T) {
	op, err := ParseOpcode(`[ %0 ] = enter 1`, false)
	if err != nil {
		t.Fatalf("enter without arguments is valid, got error '%s'", err.Error())
	}
	if len(op.References()) != 0 || len(op.Locals()) != 1 {
		t.Errorf("Unexpected enter: '%s'", op.String())
	}
}

func TestParseControlFlow(t *testing.T) {
	branch, err := ParseOpcode(`branch %2 0 -3`, false)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if branch.(*opcodes.Branch).Offset() != -3 {
		t.Errorf("Expected offset -3, got %d", branch.(*opcodes.Branch).Offset())
	}

	jump, err := ParseOpcode(`goto 4`, false)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if jump.(*opcodes.Goto).Offset() != 4 {
		t.Errorf("Expected offset 4, got %d", jump.(*opcodes.Goto).Offset())
	}

	leave, err := ParseOpcode(`leave %1 %2`, false)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(leave.References()) != 2 {
		t.Errorf("Expected 2 references, got '%s'", leave.String())
	}
}

func TestParseErrors(t *testing.T) {
	malformed := []string{
		`%5 = sub i16 %4x %3`,  //riferimento non numerico
		`%5 = sub w16 %4 %3`,   //shape sconosciuta
		`%5 = frobnicate i16 %4 %3`, //opcode inesistente
		`%x = const i16 5`,     //indice locale non numerico
		`goto abc`,             //offset non numerico
		`branch %1 0`,          //branch incompleto
		`%5 = add i16 %4`,      //operando mancante
		`[ %0 ] = enter abc`,   //indice di callable non numerico
		`; only a comment`,     //vuota dopo lo strip del commento
	}
	for _, line := range malformed {
		if _, err := ParseOpcode(line, false); err == nil {
			t.Errorf("Expected an error for '%s'", line)
		}
	}
}
