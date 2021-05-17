package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/operators"
	"github.com/aleferri/casmvm/pkg/vm"
	"github.com/aleferri/casmvm/pkg/vmio"
)

func neg(a int64) (int64, error) { return -a, nil }

func ParseLeftPart(localName string, words []string) (opcodes.Opcode, error) {
	if len(words) < 3 { // opname shape operand
		return nil, fmt.Errorf("Expected <opname> <shape> <operand> [<operand>], but %v found", words)
	}
	opName := words[0]
	local, _ := strconv.ParseUint(localName, 10, 16)
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

//ParseLineByLine provided source file and return the vm for the program
func ParseLineByLine(sourceFile string, debugMode bool) (*vm.NaiveVM, error) {
	var programfile, programErr = os.Open(sourceFile)
	if programErr != nil {
		wnd, _ := os.Getwd()
		return nil, fmt.Errorf("Error during opening of file %s from %s\n", sourceFile, wnd)
	}

	programCode := bufio.NewReader(programfile)

	var err error
	var line string
	var parseError error = nil

	listings := []opcodes.Opcode{}

	line, err = programCode.ReadString('\n')
	for err == nil && parseError == nil {
		eqIndex := strings.IndexByte(line, '=')
		if eqIndex > 0 {
			rawName := strings.TrimSpace(line[:eqIndex])
			localName := strings.TrimPrefix(rawName, "%")
			leftPart := line[eqIndex+1:]
			components := strings.Fields(leftPart)
			if debugMode {
				fmt.Printf("Recognized '%s'\n", components[0])
			}
			opcode, errOpcode := ParseLeftPart(localName, components)
			parseError = errOpcode
			listings = append(listings, opcode)
		} else {
			components := strings.Fields(line)
			opName := components[0]
			if debugMode {
				fmt.Printf("Recognized '%s'\n", opName)
			}
			switch opName {
			case "branch":
				{
					ref, _ := strconv.ParseUint(components[1], 10, 16)
					val, _ := strconv.ParseInt(components[2], 10, 64)
					offset, _ := strconv.ParseInt(components[3], 10, 64)
					listings = append(listings, opcodes.MakeBranch(val, uint16(ref), int32(offset)))
				}
			case "goto":
				{
					offset, _ := strconv.ParseInt(components[3], 10, 64)
					listings = append(listings, opcodes.MakeGoto(int32(offset)))
				}
			case "leave":
				{
					refs := []uint16{}
					for _, r := range components[1:] {
						ref, _ := strconv.ParseUint(r, 10, 16)
						refs = append(refs, uint16(ref))
					}
					listings = append(listings, opcodes.MakeLeave(refs...))
				}
			case "enter":
				{
					start, _ := strconv.ParseUint(components[1], 10, 16)
					frame, _ := strconv.ParseUint(components[2], 10, 32)
					refs := []uint16{}
					for _, r := range components[1:] {
						ref, _ := strconv.ParseUint(r, 10, 16)
						refs = append(refs, uint16(ref))
					}
					listings = append(listings, opcodes.MakeEnter(uint16(start), uint32(frame), refs))
				}
			case "sigwarn":
				{
					ref, _ := strconv.ParseUint(components[1], 10, 64)
					message := strings.Join(components[2:], ",")
					message = strings.TrimLeft(message, "\"")
					message = strings.TrimRight(message, "\"")
					listings = append(listings, opcodes.MakeSigWarning(message, uint16(ref)))
				}
			case "sigerr":
				{
					ref, _ := strconv.ParseUint(components[1], 10, 64)
					message := strings.Join(components[2:], ",")
					message = strings.TrimLeft(message, "\"")
					message = strings.TrimRight(message, "\"")
					listings = append(listings, opcodes.MakeSigError(message, uint16(ref)))
				}
			}
			parseError = nil
		}

		line, err = programCode.ReadString('\n')
	}

	callable := vm.MakeCallable()
	callable.Set(listings)

	return vm.MakeNaiveVM([]vm.Callable{callable}, vmio.MakeVMLoggerConsole(vmio.ALL), vm.MakeVMFrame()), parseError
}

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "-debug=true|false")
	flag.Parse()

	fileFound := false

	for _, f := range flag.Args() {
		if !strings.HasPrefix(f, "-") {
			fileFound = true

			vm, errAsm := ParseLineByLine(f, debugMode)

			if errAsm != nil {
				fmt.Println(errAsm.Error())
				break
			}

			results, execErr := vm.Enter(0)

			if execErr != nil {
				fmt.Println(execErr.Error(), execErr.OpcodeID())
				break
			}

			if results.Returns().IsEmpty() {
				fmt.Println("No result")
			} else {
				fmt.Println(vm.Frame().Returns())
			}
		} else if f == "--help" {
			fmt.Println("Usage: casmvm -debug=false|true filename.csm")
			fmt.Println(
				"Specified program will be parsed and executed on the fly, the last element of the stack after the computation will be printed in the console",
			)
			fmt.Println()
			fmt.Println("This program exists to debug dumps generated by casmeleon and to be used as a library by casmeleon v2")
		}
	}

	if !fileFound {
		fmt.Println("Usage: casmvm -debug=false|true filename.csm")
	}
}
