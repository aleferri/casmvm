package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/operators"
	"github.com/aleferri/casmvm/pkg/vm"
)

func neg(a int64) (int64, error) { return -a, nil }

func LineByLine(sourceFile string, debugMode bool) (*vm.NaiveVM, error) {
	var programfile, programErr = os.Open(sourceFile)
	if programErr != nil {
		wnd, _ := os.Getwd()
		return nil, fmt.Errorf("Error during opening of file %s from %s\n", sourceFile, wnd)
	}

	programCode := bufio.NewReader(programfile)

	var err error
	var line string

	listings := []opcodes.Opcode{}

	line, err = programCode.ReadString('\n')
	for err == nil {
		opcodeNameIndex := strings.IndexByte(line, ' ')
		if opcodeNameIndex < 0 {
			opcodeNameIndex = len(line)
		}
		opcodeName := strings.TrimSpace(line[0:opcodeNameIndex])
		if debugMode {
			fmt.Printf("Recognized '%s'\n", opcodeName)
		}
		left := strings.TrimSpace(strings.TrimPrefix(line, opcodeName))
		args := strings.Split(left, ",")
		switch opcodeName {
		case "rpush":
			{
				listings = append(listings, opcodes.MakeRPush())
			}
		case "iconst":
			{
				val, _ := strconv.ParseInt(args[0], 10, 64)
				listings = append(listings, opcodes.MakeIConst(val))
			}
		case "branch":
			{
				val, _ := strconv.ParseInt(args[0], 10, 64)
				offset, _ := strconv.ParseInt(args[1], 10, 64)
				listings = append(listings, opcodes.MakeBranch(val, int32(offset)))
			}
		case "call":
			{
				offset, _ := strconv.ParseInt(args[0], 10, 64)
				listings = append(listings, opcodes.MakeCall(int32(offset)))
			}
		case "goto":
			{
				offset, _ := strconv.ParseInt(args[0], 10, 64)
				listings = append(listings, opcodes.MakeGoto(int32(offset)))
			}
		case "return":
			{
				offset, _ := strconv.ParseUint(args[0], 10, 64)
				listings = append(listings, opcodes.MakeReturn(uint32(offset)))
			}
		case "neg":
			{
				listings = append(listings, opcodes.MakeUnaryOp(opcodeName, neg))
			}
		case "rstore":
			{
				offset, _ := strconv.ParseUint(args[0], 10, 64)
				listings = append(listings, opcodes.MakeRStore(uint32(offset)))
			}
		case "rload":
			{
				offset, _ := strconv.ParseUint(args[0], 10, 64)
				listings = append(listings, opcodes.MakeRLoad(uint32(offset)))
			}
		default:
			{
				for str, fn := range operators.BinaryOperatorsNames {
					if opcodeName == str {
						listings = append(listings, opcodes.MakeBinaryOp(str, fn))
					}
				}
			}
		}
		line, err = programCode.ReadString('\n')
	}

	return vm.MakeNaiveVM(listings), nil
}

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "-debug=true|false")
	flag.Parse()

	for _, f := range flag.Args() {
		if !strings.HasPrefix(f, "-") {

			vm, errAsm := LineByLine(f, debugMode)

			if errAsm != nil {
				fmt.Println(errAsm.Error())
				break
			}

			execErr := vm.Run(debugMode)

			if execErr != nil {
				fmt.Println(execErr.Error(), execErr.OpcodeID())
				break
			}

			if vm.EvalStack().Empty() {
				fmt.Println("No result")
			} else {
				fmt.Println(vm.EvalStack().Pop())
			}
		}
	}
}
