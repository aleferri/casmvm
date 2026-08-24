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

//parseReference parse a local reference like %12, the '%' prefix is optional
func parseReference(token string) (uint16, error) {
	raw := strings.TrimPrefix(token, "%")
	ref, err := strconv.ParseUint(raw, 10, 16)
	if err != nil {
		return 0, errors.New("Expected a valid reference, but found " + token)
	}
	return uint16(ref), nil
}

//ParseAssertOpcode parse a SigWarn or a SigErr
func ParseAssertOpcode(makeAssert MakeAssert, line string) (opcodes.Opcode, error) {
	line = strings.TrimSpace(line)
	indexString := strings.IndexByte(line, '"')
	if len(line) == 0 || line[0] != '%' || indexString < 0 {
		return nil, errors.New("Expected %<reference> \"Message\", but found " + line)
	}
	ref, err := parseReference(strings.TrimSpace(line[:indexString]))
	if err != nil {
		return nil, err
	}
	message := line[indexString+1:]
	if end := strings.IndexByte(message, '"'); end >= 0 {
		message = message[:end]
	}
	return makeAssert(message, ref), nil
}

//parseEnterOpcode parse the enter opcode, indexName is the return list without the leading '['
func parseEnterOpcode(indexName string, words []string) (opcodes.Opcode, error) {
	refsRet := []uint16{}
	for _, ret := range strings.Fields(strings.Trim(indexName, "[] ")) {
		val, err := parseReference(ret)
		if err != nil {
			return nil, err
		}
		refsRet = append(refsRet, val)
	}
	frame, err := strconv.ParseUint(words[1], 10, 32)
	if err != nil {
		return nil, errors.New("Expected a valid callable index, but found " + words[1])
	}
	refs := []uint16{}
	for _, r := range words[2:] {
		ref, err := parseReference(r)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return opcodes.MakeEnter(refsRet, uint32(frame), refs), nil
}

//ParseValueOpcode parse opcodes that produces values
func ParseValueOpcode(indexName string, words []string, verbose bool) (opcodes.Opcode, error) {
	if len(words) < 2 {
		return nil, fmt.Errorf("expected <opname> <shape> <operand> [<operand>], but %v found", words)
	}
	opName := words[0]
	if verbose {
		fmt.Printf("Found opcode %s\n", opName)
	}

	if opName == "enter" {
		return parseEnterOpcode(indexName, words)
	}

	if len(words) < 3 { // opname shape operand
		return nil, fmt.Errorf("expected <opname> <shape> <operand> [<operand>], but %v found", words)
	}

	localRaw := strings.TrimSpace(indexName)
	local, localErr := strconv.ParseUint(localRaw, 10, 16)
	if localErr != nil {
		return nil, errors.New("Expected a valid local index, but found " + localRaw)
	}

	shape, knownShape := opcodes.ShapesByName[words[1]]
	if !knownShape {
		return nil, errors.New("Expected a valid shape, but found " + words[1])
	}

	switch opName {
	case "local":
		{
			refA, err := parseReference(words[2])
			if err != nil {
				return nil, err
			}
			return opcodes.MakeAssignment(uint16(local), refA, shape), nil
		}
	case "const":
		{
			fullInt, err := strconv.ParseInt(words[2], 10, 64)
			if err != nil {
				return nil, errors.New("Expected a valid integer constant, but found " + words[2])
			}
			return opcodes.MakeIConst(uint16(local), shape.Reshape(fullInt)), nil
		}
	default:
		{
			if fn, isBinary := operators.BinaryOperatorsNames[opName]; isBinary {
				if len(words) < 4 {
					return nil, errors.New("missing operand")
				}
				refA, errA := parseReference(words[2])
				if errA != nil {
					return nil, errA
				}
				refB, errB := parseReference(words[3])
				if errB != nil {
					return nil, errB
				}
				if verbose {
					fmt.Println("References: ", refA, refB)
				}
				return opcodes.MakeBinaryOp(uint16(local), opName, shape, refA, refB, fn), nil
			}

			if fn, isUnary := operators.UnaryOperatorsNames[opName]; isUnary {
				refA, err := parseReference(words[2])
				if err != nil {
					return nil, err
				}
				return opcodes.MakeUnaryOp(uint16(local), opName, shape, refA, fn), nil
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
	if strings.HasPrefix(line, "sigerr") {
		return ParseAssertOpcode(opcodes.MakeSigError, line[len("sigerr"):])
	}
	if strings.HasPrefix(line, "sigwarn") {
		return ParseAssertOpcode(opcodes.MakeSigWarning, line[len("sigwarn"):])
	}
	indexComment := strings.IndexByte(line, ';')
	if indexComment >= 0 {
		line = strings.TrimSpace(line[:indexComment])
	}
	if line == "" {
		return nil, errors.New("expected an opcode, but the line is empty")
	}
	line = strings.ReplaceAll(line, ",", " ")
	if line[0] == '%' || line[0] == '[' {
		indexEq := strings.IndexByte(line, '=')
		if indexEq > 2 {
			leftPart := line[1:indexEq]
			if verbose {
				fmt.Printf("Assign to reference %s\n", line[0:indexEq])
			}
			words := strings.Fields(line[indexEq+1:])
			return ParseValueOpcode(leftPart, words, verbose)
		}
		return nil, fmt.Errorf("expected <ref> '=' <rest-of-opcode>, but found %s instead", line)
	}
	words := strings.Fields(line)
	opcodeName := words[0]
	if verbose {
		fmt.Printf("Found opcode '%s'\n", opcodeName)
	}
	switch opcodeName {
	case "branch":
		{
			if len(words) < 4 {
				return nil, errors.New("expected 'branch %<ref> <value> <offset>', but found " + line)
			}
			ref, errRef := parseReference(words[1])
			if errRef != nil {
				return nil, errRef
			}
			val, errVal := strconv.ParseInt(words[2], 10, 64)
			if errVal != nil {
				return nil, errors.New("Expected a valid compare value, but found " + words[2])
			}
			offset, errOff := strconv.ParseInt(words[3], 10, 32)
			if errOff != nil {
				return nil, errors.New("Expected a valid branch offset, but found " + words[3])
			}
			return opcodes.MakeBranch(val, ref, int32(offset)), nil
		}
	case "goto":
		{
			if len(words) < 2 {
				return nil, errors.New("expected 'goto <offset>', but found " + line)
			}
			offset, err := strconv.ParseInt(words[1], 10, 32)
			if err != nil {
				return nil, errors.New("Expected a valid goto offset, but found " + words[1])
			}
			return opcodes.MakeGoto(int32(offset)), nil
		}
	case "leave":
		{
			refs := []uint16{}
			for _, r := range words[1:] {
				ref, err := parseReference(r)
				if err != nil {
					return nil, err
				}
				refs = append(refs, ref)
			}
			return opcodes.MakeLeave(refs...), nil
		}
	}
	return nil, errors.New("opcode is invalid")
}
