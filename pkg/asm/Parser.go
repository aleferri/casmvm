package asm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/operators"
)

//MakeAssert create either a SigWarn or a SigErr
type MakeAssert func(msg string, ref uint16) opcodes.Opcode

//ParseAssertOpcode parse a SigWarn or a SigErr
func ParseAssertOpcode(makeAssert MakeAssert, line string) (opcodes.Opcode, error) {
	indexString := strings.IndexByte(line, '"')
	if indexString < 3 {
		return nil, errors.New("Expected <reference> \"Message\", but found " + line)
	}
	left := line[1 : indexString-1]
	message := line[indexString+1:]
	ref, err := strconv.ParseUint(strings.TrimSpace(left), 10, 16)
	if err != nil {
		return nil, errors.New("Expected a valid reference, but found " + left)
	}
	opcode := makeAssert(strings.TrimRight(message, "\""), uint16(ref))
	return opcode, nil
}

//ParseValueOpcode parse opcodes that produces values
func ParseValueOpcode(indexName string, words []string, verbose bool) (opcodes.Opcode, error) {
	if len(words) < 3 { // opname shape operand
		return nil, fmt.Errorf("Expected <opname> <shape> <operand> [<operand>], but %v found", words)
	}
	opName := words[0]
	if verbose {
		fmt.Printf("Found opcode %s\n", opName)
	}
	local, _ := strconv.ParseUint(strings.TrimSpace(indexName), 10, 16)
	refARaw := strings.TrimPrefix(words[2], "%")
	refA, _ := strconv.ParseUint(refARaw, 10, 16)
	switch opName {
	case "local":
		{
			return opcodes.MakeAssignment(uint16(local), uint16(refA), opcodes.IntShape), nil
		}
	case "const":
		{
			fullInt, _ := strconv.ParseInt(refARaw, 10, 64)
			return opcodes.MakeIConst(uint16(local), fullInt), nil
		}
	default:
		{
			for str, fn := range operators.BinaryOperatorsNames {
				if opName == str {
					if len(words) < 4 {
						return nil, errors.New("Missing operand")
					}
					refBRaw := strings.TrimPrefix(words[3], "%")
					refB, _ := strconv.ParseUint(refBRaw, 10, 16)
					if verbose {
						fmt.Println("References: ", refA, refB)
					}
					return opcodes.MakeBinaryOp(uint16(local), opName, opcodes.IntShape, uint16(refA), uint16(refB), fn), nil
				}
			}

			for str, fn := range operators.UnaryOperatorsNames {
				if opName == str {
					return opcodes.MakeUnaryOp(uint16(local), opName, opcodes.IntShape, uint16(refA), fn), nil
				}
			}

			return nil, errors.New("Missing opcode " + opName)
		}
	}
}

//ParseOpcode for the CasmVM
func ParseOpcode(line string, verbose bool) (opcodes.Opcode, error) {
	if verbose {
		fmt.Println("Processing line", line)
	}
	if strings.Index(line, "sigerr") == 0 {
		return ParseAssertOpcode(opcodes.MakeSigError, line[7:])
	}
	if strings.Index(line, "sigwarn") == 0 {
		return ParseAssertOpcode(opcodes.MakeSigWarning, line[7:])
	}
	indexComment := strings.IndexByte(line, ';')
	if indexComment > 0 {
		line = line[0:indexComment]
	}
	if line[0] == '%' {
		indexEq := strings.IndexByte(line, '=')
		if indexEq > 2 {
			leftPart := line[1:indexEq]
			if verbose {
				fmt.Printf("Assign to reference %s\n", leftPart)
			}
			words := strings.Fields(line[indexEq+1:])
			return ParseValueOpcode(leftPart, words, verbose)
		} else {
			return nil, fmt.Errorf("Expected <ref> '=' <rest-of-opcode>, but found %s instead\n", line)
		}
	}
	words := strings.Fields(line)
	opcodeName := words[0]
	if verbose {
		fmt.Printf("Found opcode '%s'\n", opcodeName)
	}
	switch opcodeName {
	case "branch":
		{
			ref, _ := strconv.ParseUint(words[1], 10, 16)
			val, _ := strconv.ParseInt(words[2], 10, 64)
			offset, _ := strconv.ParseInt(words[3], 10, 64)
			return opcodes.MakeBranch(val, uint16(ref), int32(offset)), nil
		}
	case "goto":
		{
			offset, _ := strconv.ParseInt(words[3], 10, 64)
			return opcodes.MakeGoto(int32(offset)), nil
		}
	case "leave":
		{
			refs := []uint16{}
			for _, r := range words[1:] {
				ref, _ := strconv.ParseUint(r, 10, 16)
				refs = append(refs, uint16(ref))
			}
			return opcodes.MakeLeave(refs...), nil
		}
	case "enter":
		{
			start, _ := strconv.ParseUint(words[1], 10, 16)
			frame, _ := strconv.ParseUint(words[2], 10, 32)
			refs := []uint16{}
			for _, r := range words[1:] {
				ref, _ := strconv.ParseUint(r, 10, 16)
				refs = append(refs, uint16(ref))
			}
			return opcodes.MakeEnter(uint16(start), uint32(frame), refs), nil
		}
	}
	return nil, errors.New("Opcode is invalid")
}
